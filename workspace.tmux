#!/bin/bash

tmux-workspace "TasksNTimes" "editor" -c "nvim && zsh"\
    -w "terminal" -c "zsh" -c "make dev && zsh"
