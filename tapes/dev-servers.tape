# VHS Tape for Development Servers Example
# Shows auto-restart and colored output

Output tapes/dev-servers.gif

Set FontSize 14
Set Width 1200
Set Height 800
Set TypingSpeed 50ms
Set Theme "Catppuccin Mocha"

Type "# Development servers with auto-restart"
Sleep 1s
Enter

Type "cat examples/dev-servers.yaml"
Sleep 500ms
Enter

Sleep 2s

Type "./tapes/run-dev-servers.sh"
Sleep 500ms
Enter

# Let services run for a bit
Sleep 9s
