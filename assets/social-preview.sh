#!/bin/sh
# Generate the GitHub social preview card (1280x640) for HTTPRunner.
# Upload the result under repo Settings -> Social preview.
# Render (Docker, no local deps):
#   docker run --rm -v "$PWD":/img -w /img --entrypoint sh dpokidov/imagemagick assets/social-preview.sh
set -e
F=/usr/share/fonts/truetype/dejavu
magick -size 1280x640 xc:'#16161e' \
  -fill '#7aa2f7' -draw 'rectangle 0,0 14,640' \
  -gravity NorthWest \
  -font "$F/DejaVuSans-Bold.ttf" -pointsize 96 -fill '#c0caf5' -annotate +80+140 'HTTPRunner' \
  -font "$F/DejaVuSans.ttf"      -pointsize 34 -fill '#a9b1d6' -annotate +84+275 'Load testing for HTTP requests — written in Go' \
  -font "$F/DejaVuSansMono.ttf"  -pointsize 25 -fill '#9ece6a' -annotate +84+405 '$ http-runner -url https://example.com -count 100 -concurrency 20' \
  -font "$F/DejaVuSans.ttf"      -pointsize 23 -fill '#565f89' -annotate +84+525 'github.com/idesyatov/http-runner     ·     MIT License     ·     Go 1.23' \
  assets/social-preview.png
