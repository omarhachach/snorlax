# Snorlax
[![GitHub tag](https://img.shields.io/github/release/omarhachach/snorlax.svg?style=flat-square)](https://github.com/omarhachach/snorlax/releases)
[![Report Card](https://img.shields.io/badge/report%20card-a%2B-c0392b.svg?style=flat-square)](https://goreportcard.com/report/github.com/omarhachach/snorlax)
[![Documentation](https://img.shields.io/badge/documentation-godoc-1abc9c.svg?style=flat-square)](https://godoc.org/github.com/omarhachach/snorlax)
[![Powered By](https://img.shields.io/badge/powered%20by-go-blue.svg?style=flat-square)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT%20License-1abc9c.svg?style=flat-square)](https://github.com/omarhachach/snorlax/blob/master/LICENSE)

Snorlax is yet another modular Discord bot, with loads of premade modules and other features.

**This is still in early development, therefore it is subject to frequent breaking changes and rewrites.**

## Menu
* [Features](#features)
* [Modules](#modules)
* [Getting Started](#getting-started)
* [Commands](#commands)
* [Running from source](#running-from-source)
* [Contributing](#contributing)
* [License](#license)

## [Premade Modules](https://github.com/omarhachach/snorlax/tree/master/modules)
* Administration
* Eval
* Gambling (WIP)
* Moderation
* Music
* Role Manager
* More to come!

## Getting Started
To get the public bot, you can just click [here]() (invalid link, not public yet).

To run it yourself, grab one of the [releases](https://github.com/omarhachach/snorlax/releases). (Still in early development, see [running from source](#running-from-source)).<br>
Once downloaded go to the location of the downloaded binary, create a config.

This a sample config:
```JSON
{
    "autoDelete": true,
    "dbPath": "./snorlax.db",
    "debug": false,
    "displayAuthor": false,
    "token": "bot-token",
    "owners": [
        "140254342170148864"
    ]
}
```
When you've created a config, run the program:
```Bash
$ ./snorlax
```

## Commands
Visit [the website](https://www.snorlaxbot.com/commands) (not finished, see help command) for a full list of commands.

## Running From Source
To run the bot from source, you need to have Go installed.<br>
You also need to define a [GOPATH](https://golang.org/doc/code.html#GOPATH).

Follow these instructions to compile and run Snorlax form source.
```Bash
$ go get github.com/omarhachach/snorlax
$ go install github.com/omarhachach/snorlax/cmd/snorlax
```
It will install it in the bin folder in your GOPATH ($GOPATH/bin).<br>
If you have $GOPATH/bin in your PATH variable, you will be able to run it like so:
```Bash
$ snorlax
```

Remember to [create a config](#getting-started) before running the bot.

## Contributing
If you've found a bug, or have a suggestion feel free to open an [issue](https://github.com/omarhachach/snorlax/issues).

You can contact me (the author) on Discord: Omar H.#6299 or via [email](mailto:contact@omarh.net).

## License
Snorlax is licensed under the [MIT License](https://github.com/omarhachach/snorlax/blob/master/LICENSE).
