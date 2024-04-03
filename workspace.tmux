#!/bin/bash

tmux-workspace "TasksNTimes" "editor" -c "nvim && zsh"\
    -w "terminal" -c "make serve && zsh"
