# Tasks n' Times
A quick CLI tool to keep track of how long certain tasks have taken


## Goal of the Project

I had a lot of "Project-Hopping" to do lately, and can barely keep track of what toke how long.
So this tool is meant to make things easier by documenting when and how long certain tasks have been executed.


## installation:
Just put one of the binaries from the Release section into a folder registered in your `$PATH` (Linux, Mac)  or `%PATH%` (Windows).

Then rename the Binary, `tnt`

if you don't trust the code, you can build your own version, if you have at least go 1.21.3 installed.


## usage:
for now this CLI (Command line interface only)

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

or, it you want to stop the timer (for example because your shift is over) call:

```bash
tnt stop
```

## Switching to another task.
switching tasks is as easy as calling 
```bash
tnt s "what ever is your new task"
```
Should another timer be running currently, it will be stopped and a new timer for the new task is started.


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
