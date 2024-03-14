# Tasks n' Times
A CLI tool to keep track of how long certain tasks have taken

<!-- vim-markdown-toc GFM -->

* [Goal of the Project](#goal-of-the-project)
* [installation](#installation)
* [usage](#usage)
    * [Switching to another task.](#switching-to-another-task)
    * [Getting a list of Tasks and how much time was spend in them](#getting-a-list-of-tasks-and-how-much-time-was-spend-in-them)
    * [Getting a full list of Times](#getting-a-full-list-of-times)
    * [Filtering the times](#filtering-the-times)
* [Autocompletion:](#autocompletion)
* [Probabbly asked Questions](#probabbly-asked-questions)
* [Build it yourself](#build-it-yourself)

<!-- vim-markdown-toc -->


## Goal of the Project

I had a lot of "Project-Hopping" to do lately, and can barely keep track of what toke how long.
This tool is meant to make things easier by documenting when and how long certain tasks have been executed.

## installation
- grab one of the Binaries from the [Releases - Section](https://github.com/Rocco-Gossmann/tnt/releases/latest)
- rename it to `tnt` (for Linux/Unix) or `tnt.exe` (for Windows)
- put it into a directory that is listed in your $PATH or %PATH%
- setup the Autocompletion as described further down in [Autocompletion](#autocompletion)
Done

Should your system not be covered by the precompiled binaries, please use the
build guide further down [Build it yourself](#build-it-yourself)


## usage
For now this is CLI (Command line interface) tool only
Easily usable via Windows Powershell/CMD or Unix Bash, ZSH, etc. 

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


### Getting a list of Tasks and how much time was spend in them
just run.

```bash
tnt times sum 
```
and you'll get a list of all tasks and how much time you've spend with them in total.

### Getting a full list of Times
To find out, when you did what, you can use the 
```bash
tnt times ls
```
command. The most current started timer is on top.
The oldest timer is on the bottom.

### Filtering the times 
both of the `tnt times` commands can take a `-t` or `--task` flag, to filter the results for a given task.
The value for this flag is the Taskname give during `tnt task add` or used during `tnt s | start | switch`.


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


## Probabbly asked Questions
 **Q**: Where does `tnt` store it's persistent Data?
> All it's data is stored in the `$HOME/.local/share/tnt` directory.
> On Windows that should be `%USERPROFILE%\.local\share\tnt`.


## Build it yourself

- Install GO (https://go.dev/doc/install) 
- clone this repo
```bash
git clone https://github.com/Rocco-Gossmann/tnt.git
```
- Enter the directory
```bash
cd tnt
```

- Build the project
```bash
CGO_ENABLED=1 go build -ldflags="-s -w -X main.Version=`git describe --tags --abbrev=0`"
```

- Copy the `tnt` file into one of the folders listed by your environment vars $PATH or %PATH% var
- or extend $PATH 
```bash
echo "export PATH=\"`pwd`:\$PATH\"" >> ~/.bashrc

# or for ZSH

echo "export PATH=\"`pwd`:\$PATH\"" >> ~/.zshrc
```


