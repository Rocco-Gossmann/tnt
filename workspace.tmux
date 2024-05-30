#!/bin/bash

tmux-workspace "TasksNTimes" "GIT" -c "lg && zsh"\
    -w "terminal" -c "make serve && zsh"
