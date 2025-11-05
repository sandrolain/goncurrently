package main

import (
	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
)

func defaultCommandColors() []*color.Color {
	return []*color.Color{
		color.New(color.FgCyan),
		color.New(color.FgGreen),
		color.New(color.FgMagenta),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgRed),
	}
}

func defaultPanelStyles(commands []CommandConfig) map[string]panelAppearance {
	panelColors := []tcell.Color{
		tcell.GetColor("aqua"),
		tcell.GetColor("springgreen"),
		tcell.GetColor("fuchsia"),
		tcell.GetColor("yellow"),
		tcell.GetColor("dodgerblue"),
		tcell.GetColor("indianred"),
	}

	styles := make(map[string]panelAppearance, len(commands)+1)
	for i, c := range commands {
		panelColor := panelColors[i%len(panelColors)]
		styles[c.Name] = panelAppearance{
			BorderColor:     panelColor,
			TitleColor:      panelColor,
			BackgroundColor: tcell.ColorDefault,
		}
	}
	if _, ok := styles[basePanelName]; !ok {
		styles[basePanelName] = panelAppearance{
			BorderColor:     tcell.ColorDarkGray,
			TitleColor:      tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
		}
	}
	return styles
}
