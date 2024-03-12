# Tasks n' Times
A quick CLI tool to keep track of how long certain tasks have taken



<!-- vim-markdown-toc GFM -->

* [Goal of the Project](#goal-of-the-project)
* [installation](#installation)
        * [Why no prebuild Binaries](#why-no-prebuild-binaries)
* [usage](#usage)
    * [Switching to another task.](#switching-to-another-task)
    * [Getting a list of Tasks and Times](#getting-a-list-of-tasks-and-times)
    * [Full list of Times](#full-list-of-times)
* [Autocompletion:](#autocompletion)
* [Probabbly asked Questions](#probabbly-asked-questions)

<!-- vim-markdown-toc -->


## Goal of the Project

I had a lot of "Project-Hopping" to do lately, and can barely keep track of what toke how long.
So this tool is meant to make things easier by documenting when and how long certain tasks have been executed.


## installation

For now building the project yourself is going to be your only option. 
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


#### Why no prebuild Binaries
Due to my lack of experience with cross-compilation.
Cross compiling Go projects, that don't have depencies is easy, yes.
However, this project made the mistake of using Sqlite3.
Sqlite3 requires CGO_ENABLED=1, which in return requires platform specific 
libs.
Maybe I'll try to find another way to handle persisten date, but for now,
this is it.



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
tnt times ls
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



## Probabbly asked Questions
 **Q**: Where does `tnt` store it's persistent Data?
> All it's data is stored in the `$HOME/.local/share/tnt` - Folder.
> On Windows that should be `%%USERPROFILE%%\.local\share\tnt`.
