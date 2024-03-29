#!/bin/bash

tmux-workspace "TasksNTimes" "editor" -c "nvim && zsh"\
    -w "terminal" -c  "make serve && zsh" -c "cd ./pkg/serve/views && tailwindcss -i ./main.tw.css -o main.css --watch && zsh"
