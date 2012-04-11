#!/usr/bin/env bash
set -e
eval $($1/dist env -p)
bash run.bash --no-rebuild --banner
