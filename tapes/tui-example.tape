# VHS Tape for TUI Example
# This generates a demo GIF showing the TUI mode

Output tapes/tui-example.gif

Set FontSize 13
Set Width 1400
Set Height 1000
Set Theme "Catppuccin Mocha"

Type "# TUI Mode - Terminal User Interface"
Sleep 1s
Enter
Sleep 500ms

Type "cat examples/tui-example.yaml"
Sleep 500ms
Enter
Sleep 3s

Type "# Starting TUI mode..."
Sleep 500ms
Enter
Sleep 500ms

Type "./tapes/run-tui.sh"
Enter

# Wait for TUI to fully initialize and render
Sleep 2s

# Let it run to show activity in all panels
Sleep 10s

# Send Ctrl+C to exit gracefully
Ctrl+C
Sleep 1s
