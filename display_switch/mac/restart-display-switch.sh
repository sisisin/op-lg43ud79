#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

proc_id=$(pgrep display_switch)
if [ -n "$proc_id" ]; then
  echo "Killing existing display_switch process (PID: $proc_id)"
  kill "$proc_id"
  sleep 1
fi

echo "Starting display_switch..."
launchctl unload ~/Library/LaunchAgents/dev.haim.display-switch.daemon.plist
launchctl load ~/Library/LaunchAgents/dev.haim.display-switch.daemon.plist