# Tasks n' Times
A quick CLI tool to keep track of how long certain tasks have taken



<!-- vim-markdown-toc GFM -->

* [Goal of the Project](#goal-of-the-project)
* [installation](#installation)
* [usage](#usage)
    * [Switching to another task.](#switching-to-another-task)
    * [Getting a list of Tasks and Times](#getting-a-list-of-tasks-and-times)
    * [Full list of Times](#full-list-of-times)
* [Autocompletion:](#autocompletion)
* [Building from Source](#building-from-source)
* [Paq (Probabbly asked Questions)](#paq-probabbly-asked-questions)

<!-- vim-markdown-toc -->


## Goal of the Project

I had a lot of "Project-Hopping" to do lately, and can barely keep track of what toke how long.
So this tool is meant to make things easier by documenting when and how long certain tasks have been executed.


## installation
Just put one of the binaries from the Release section into a folder registered in your `$PATH` (Linux, Mac)  or `%PATH%` (Windows).

Then rename the Binary, `tnt`

if you don't trust the binaries further down is a section on how to create your own builds.


## usage
for now this is CLI (Command line interface) only

```bash
# first create some tasks 
tnt tasks add "important project nr 1"
tnt tasks add "very important project"
tnt tasks add "some side project"
tnt tasks add "learning new stuff"
   
# you can see what task you created via
tnt tasks ls

# To remove a task and all its Recorded Times use
tnt tasks rm "some side project";
```

to start tracking times, use:
```bash
tnt s "important project nr 1"

#or
tnt start "important project nr 1"

#or
tnt switch "important project nr 1"

#these 3 are aliasses for the the same function
```

now the timer for "important project nr 1" is running until you either start
another timer via `s`, `start` or `switch`

or, if you want to stop the timers (for example because your shift is over) call:

```bash
tnt stop
```

### Switching to another task.
switching tasks is as easy as calling
```bash
tnt s "what ever is your new task"
```
Should another timer be running currently, it will be stopped and a new timer for the new task is started.


### Getting a list of Tasks and Times
just run.

```bash
tnt times sum 
```
and you'll get a list of all tasks and how much time you've spend with them in total.

### Full list of Times
To find out, when you did what, you can use the 
```bash
tnt timers ls
```
command. The most current started timer is on top.
The oldest timer is on the bottom.



## Autocompletion:
This tool was build with auto completion im mind, as task names can become quite long.
Autocompletion make switching between tasks/timers very fritctionless.

on bash you can activate the autocomletion by putting the following in your .bashrc.
```bash
source <(tnt completion bash)
```

For ZSH put this into your .zshrc
```bash
source <(tnt completion zsh)
```

Check the following command for more autocompletion options
```bash
tnt completion 
```

## Building from Source
For that, you need to install Go 1.21.3 or later.
Clone this repo. 
Enter the Clones Folder
Run `go mod tidy` to download the missing packages
and build The Binaries `go build`


## Paq (Probabbly asked Questions)
Where does `tnt` store it's persistent Data?
> All it's data is stored in the $HOME/.local/share/tnt - Folder
