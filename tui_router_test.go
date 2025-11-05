package main

import (
	"testing"
)

func TestCalculateGridDimensions(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		wantRows int
		wantCols int
	}{
		{
			name:     "zero panels",
			total:    0,
			wantRows: 1,
			wantCols: 1,
		},
		{
			name:     "one panel",
			total:    1,
			wantRows: 1,
			wantCols: 1,
		},
		{
			name:     "two panels",
			total:    2,
			wantRows: 1,
			wantCols: 2,
		},
		{
			name:     "three panels",
			total:    3,
			wantRows: 2,
			wantCols: 2,
		},
		{
			name:     "four panels",
			total:    4,
			wantRows: 2,
			wantCols: 2,
		},
		{
			name:     "five panels",
			total:    5,
			wantRows: 2,
			wantCols: 3,
		},
		{
			name:     "nine panels",
			total:    9,
			wantRows: 3,
			wantCols: 3,
		},
		{
			name:     "ten panels",
			total:    10,
			wantRows: 3,
			wantCols: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRows, gotCols := calculateGridDimensions(tt.total)
			if gotRows != tt.wantRows || gotCols != tt.wantCols {
				t.Errorf("calculateGridDimensions(%d) = (%d, %d), want (%d, %d)",
					tt.total, gotRows, gotCols, tt.wantRows, tt.wantCols)
			}
		})
	}
}

func TestCreatePanelView(t *testing.T) {
	styles := defaultPanelStyles([]CommandConfig{{Name: "test"}})
	style := styles["test"]

	view := createPanelView("test", style)
	if view == nil {
		t.Fatal("createPanelView() returned nil")
	}

	// Test with basePanelName
	baseStyle := styles[basePanelName]
	baseView := createPanelView(basePanelName, baseStyle)
	if baseView == nil {
		t.Fatal("createPanelView() for basePanelName returned nil")
	}
}
