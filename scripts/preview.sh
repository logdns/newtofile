#!/bin/bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
APP="$ROOT/build/NewToFile.app"
if [[ ! -x "$APP/Contents/MacOS/NewToFileBridge" ]]; then ./scripts/build.sh; fi
"$APP/Contents/MacOS/NewToFile" --preview "$ROOT/build/preview.png"
"$APP/Contents/MacOS/NewToFile" --preview "$ROOT/build/preview-dark.png" --dark
"$APP/Contents/MacOS/NewToFile" --preview "$ROOT/build/preview-settings.png" --settings
"$APP/Contents/MacOS/NewToFile" --preview "$ROOT/build/preview-forced-light.png" --settings --dark --appearance 1 --vibrancy
"$APP/Contents/MacOS/NewToFile" --preview "$ROOT/build/preview-forced-dark.png" --settings --appearance 2 --vibrancy
echo "Rendered MyGo previews in $ROOT/build/"
