# Autofresh

[![Godoc](https://godoc.org/github.com/TerrenceHo/autofresh?status.svg)](http://godoc.org/github.com/TerrenceHo/autofresh)
[![GoReportCard](https://goreportcard.com/badge/github.com/TerrenceHo/autofresh)](https://goreportcard.com/report/github.com/TerrenceHo/autofresh)

        ___         __        ______               __
       /   | __  __/ /_____  / ____/_______  _____/ /_
      / /| |/ / / / __/ __ \/ /_  / ___/ _ \/ ___/ __ \
     / ___ / /_/ / /_/ /_/ / __/ / /  /  __(__  ) / / /
    /_/  |_\__,_/\__/\____/_/   /_/   \___/____/_/ /_/

![Sample](https://media.giphy.com/media/9JwRV7puvSTSGg3wEl/giphy.gif)

Autofresh is a simple live-reload development server that rebuilds your program
every time a file is saved, added, or deleted. Gone are the days where you need
to manually recompile a program to test it. Install Autofresh, run it in your
project directory, and start coding away. Once you save your file, your program
will automatically recompile, as seen above.

### Installation

You can checkout the [releases](https://github.com/TerrenceHo/autofresh/releases) 
page and download the latest version of autofresh.
Install the binary compiled for your operating system and architecture. The same 
binaries are also available under the directory `bin`.

If you have a Go runtime installed, you can `go get` this repository and
automatically install it. Make sure you put the GOBIN in your PATH.

```bash 
go get github.com/TerrenceHo/autofresh 
```

### Dependencies

##### Watchman

Watchman is hard dependency, because autofresh uses watchman to recursively
watch your project directory, rather use a custom file watcher.  Install
watchman [here](https://facebook.github.io/watchman/docs/install.html).

Note that because Watchman does not yet fully support Windows, Autofresh will
not yet work entirely correctly on Windows. Windows support should come in a few
months, and Autofresh will be kept up to date if/when that happens. As for now,
use Linux/MacOS with Autofresh.

### Usage

The easiest way to use Autofresh is to run `autofresh` in your command line
application/terminal. Create a autofresh-config file in your directory with the
appropriate build/run commands, and autofresh will automatically start compiling 
and refreshing your program.

### Configuration

Autofresh takes in three ways of configuration, flags, environment variables, a
configuration file autofresh-config.{json,yml/yaml,toml}, and finally the
defaults, with the each item taking precedence in that order (i.e. if both flags
and environment variables set watchman's path, the flag's path will be taken).
The easiest way to configure Autofresh is to use a autofresh-config file. A
sample autofresh-config file is shown below. Comments are added for clarity.

```json 
{
    "watchman": "/usr/local/bin/watchman",  --> path to watchman executable
    "build": "go build -o sample .",        --> command to build program
    "run": "./sample",                      --> command to run program 
    "suffixes": ["go", "txt"]               --> file extensions to watch
}
```

Additionally, Autofresh takes in flags and environment variables, with the same
name as the variables in the autofresh-config file. Environment variables should
be prefixed with AUTOFRESH\_, like AUTOFRESH\_BUILD or AUTOFRESH\_RUN. Flags and
their command descriptions can be seen with `autofresh --help`

### Live reload vs hot reload

While partially inspired by React Native's hot reloading development server,
Autofresh does not hot reload, only live reload. It is intended to work for most
languages/scripts, and simply rebuilds and restarts the program.

However, Autofresh works with all languages, as long as you provide a build
command and a run command.

### Contributing
You should fork this repository, and then clone your fork into your GOPATH. The 
Makefile is included for development convenience. Once you have your appropriate
changes, make a pull request! Most changes will be merged into the develop
branch first, and will keep master branch stable.
