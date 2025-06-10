#!/bin/bash

INIT_SCRIPT="$HOME/.config/tmux-init/init.sh"
if [ -f "$INIT_SCRIPT" ]; then
    source "$INIT_SCRIPT"
fi
