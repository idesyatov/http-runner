# VHS tape for the HTTPRunner demo GIF.
# Render (Docker, no local deps):
#   docker run --rm -v "$PWD":/vhs ghcr.io/charmbracelet/vhs assets/demo.tape
# Requires bin/http-runner to be built first (make build-docker).

Output assets/demo.gif

Set Shell bash
Set FontSize 18
Set Width 1100
Set Height 700
Set Theme "Dracula"
Set Padding 20
Set TypingSpeed 35ms

Hide
Type "export PATH=$PWD/bin:$PATH && clear"
Enter
Show

Type "http-runner -url https://example.com"
Sleep 800ms
Enter
Sleep 2.5s

Type "http-runner -url https://example.com -count 100 -concurrency 20"
Sleep 800ms
Enter
Sleep 3.5s
