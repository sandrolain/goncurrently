# VHS Tape for Microservices Example
# This generates a demo GIF showing the microservices example

Output tapes/microservices.gif

Set FontSize 14
Set Width 1200
Set Height 800
Set TypingSpeed 50ms
Set Theme "Catppuccin Mocha"

Type "# Running microservices example with goncurrently"
Sleep 1s
Enter

Type "cat examples/microservices.yaml"
Sleep 500ms
Enter

Sleep 2s

Type "./tapes/run-microservices.sh"
Sleep 500ms
Enter

# Wait for setup and let the services run
Sleep 11s
