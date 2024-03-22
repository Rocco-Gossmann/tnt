#!/bin/bash

tmux-workspace "TasksNTimes" "editor" -c "nvim && zsh"\
    -w "terminal" -c "cd ./views && tailwindcss -i ./main.tw.css -o main.css --watch && zsh" -c "make serve && zsh"
