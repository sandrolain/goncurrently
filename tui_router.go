package main

import (
	"fmt"
	"io"
	"math"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tuiRouter struct {
	app         *tview.Application
	baseName    string
	views       map[string]*tview.TextView
	defaultView *tview.TextView
	stopOnce    sync.Once
	runDone     chan struct{}
	runErr      error
}

func newTUIRouter(baseName string, commandNames []string, styles map[string]panelAppearance) (*tuiRouter, error) {
	app := tview.NewApplication()
	sectionNames := make([]string, 0, len(commandNames)+1)
	sectionNames = append(sectionNames, baseName)
	sectionNames = append(sectionNames, commandNames...)
	if len(sectionNames) == 0 {
		sectionNames = []string{baseName}
	}
	seen := make(map[string]struct{}, len(sectionNames))
	unique := make([]string, 0, len(sectionNames))
	for _, name := range sectionNames {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		unique = append(unique, name)
	}
	sectionNames = unique
	layout, views := buildTUILayout(sectionNames, styles)

	defaultView := views[sectionNames[0]]
	if _, ok := views[baseName]; !ok {
		baseName = sectionNames[0]
	}

	t := &tuiRouter{
		app:         app,
		baseName:    baseName,
		views:       views,
		defaultView: defaultView,
		runDone:     make(chan struct{}),
	}

	app.SetRoot(layout, true)
	if focusView, ok := views[baseName]; ok {
		app.SetFocus(focusView)
	}

	go func() {
		t.runErr = app.Run()
		close(t.runDone)
	}()

	return t, nil
}

func buildTUILayout(sectionNames []string, styles map[string]panelAppearance) (*tview.Flex, map[string]*tview.TextView) {
	rows, cols := calculateGridDimensions(len(sectionNames))
	rootFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rowFlexes := make([]*tview.Flex, rows)
	for r := 0; r < rows; r++ {
		rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
		rootFlex.AddItem(rowFlex, 0, 1, false)
		rowFlexes[r] = rowFlex
	}

	views := make(map[string]*tview.TextView, len(sectionNames))
	for idx, name := range sectionNames {
		style := styles[name]
		view := createPanelView(name, style)
		rowFlexes[idx/cols].AddItem(view, 0, 1, idx == 0)
		views[name] = view
	}
	return rootFlex, views
}

func calculateGridDimensions(total int) (rows int, cols int) {
	if total <= 0 {
		return 1, 1
	}
	cols = int(math.Ceil(math.Sqrt(float64(total))))
	if cols < 1 {
		cols = 1
	}
	rows = int(math.Ceil(float64(total) / float64(cols)))
	if rows < 1 {
		rows = 1
	}
	return rows, cols
}

func createPanelView(name string, style panelAppearance) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetTitle(name)
	if style.BorderColor != tcell.ColorDefault {
		textView.SetBorderColor(style.BorderColor)
	}
	if style.TitleColor != tcell.ColorDefault {
		textView.SetTitleColor(style.TitleColor)
	}
	if style.BackgroundColor != tcell.ColorDefault {
		textView.SetBackgroundColor(style.BackgroundColor)
	}
	textView.SetScrollable(true)
	textView.SetWrap(true)
	return textView
}

func (t *tuiRouter) BaseWriter() io.Writer {
	view := t.views[t.baseName]
	if view == nil {
		view = t.defaultView
	}
	return &textViewWriter{
		app:         t.app,
		view:        view,
		prefix:      color.New(color.FgHiCyan).Sprint("[gonc] "),
		atLineStart: true,
	}
}

func (t *tuiRouter) LineWriter(name string, col *color.Color, prefix string) func(string) {
	view, ok := t.views[name]
	if !ok || view == nil {
		view = t.views[t.baseName]
	}
	if view == nil {
		view = t.defaultView
	}
	if view == nil {
		return func(string) {
			// no-op: no view available for rendering
		}
	}
	coloredPrefix := prefix
	if col != nil {
		coloredPrefix = col.Sprint(prefix)
	}
	return func(line string) {
		t.app.QueueUpdateDraw(func() {
			writer := tview.ANSIWriter(view)
			fmt.Fprintf(writer, lineJoinFormat, coloredPrefix, line) //nolint:errcheck
			view.ScrollToEnd()
		})
	}
}

func (t *tuiRouter) Stop() {
	t.stopOnce.Do(func() {
		if t.app != nil {
			t.app.Stop()
		}
	})
}

func (t *tuiRouter) Add() {

}

func (t *tuiRouter) Done() {

}

func (t *tuiRouter) Wait() {
	<-t.runDone
}

type textViewWriter struct {
	app         *tview.Application
	view        *tview.TextView
	prefix      string
	atLineStart bool
}

func (w *textViewWriter) Write(p []byte) (int, error) {
	if w == nil || w.view == nil {
		return len(p), nil
	}
	text := string(p)
	w.app.QueueUpdateDraw(func() {
		tvWriter := tview.ANSIWriter(w.view)
		remaining := text
		for len(remaining) > 0 {
			newlineIdx := strings.IndexByte(remaining, '\n')
			if w.atLineStart && w.prefix != "" {
				fmt.Fprint(tvWriter, w.prefix) //nolint:errcheck
			}
			if newlineIdx == -1 {
				fmt.Fprint(tvWriter, remaining) //nolint:errcheck
				w.atLineStart = false
				break
			}
			fmt.Fprint(tvWriter, remaining[:newlineIdx]) //nolint:errcheck
			fmt.Fprint(tvWriter, "\n")                   //nolint:errcheck
			w.atLineStart = true
			remaining = remaining[newlineIdx+1:]
		}
		w.view.ScrollToEnd()
	})
	return len(p), nil
}
