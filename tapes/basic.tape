# VHS Tape for Basic Example
# Quick demo of basic concurrent execution

Output tapes/basic.gif

Set FontSize 14
Set Width 1200
Set Height 600
Set Theme "Catppuccin Mocha"

Type "# Simple concurrent commands with goncurrently"
Sleep 1s
Enter
Sleep 500ms

Type "cat examples/basic.yaml"
Sleep 500ms
Enter
Sleep 3s

Type "# Execute the example"
Sleep 500ms
Enter
Sleep 500ms

Type "./tapes/run-basic.sh"
Enter

# Wait for commands to run and show output
Sleep 5s
