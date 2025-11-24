#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

script_dir=$(cd "$(dirname "$0")" && pwd)
readonly script_dir

cp "${script_dir}/display-switch.ini" "$HOME/Library/Preferences/display-switch.ini"
cp "${script_dir}/dev.haim.display-switch.daemon.plist" "$HOME/Library/LaunchAgents/"
