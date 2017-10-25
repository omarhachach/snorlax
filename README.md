# Snorlax
Snorlax is another modular Discord bot, with loads of built in modules.

## Menu
* [Modules](#modules)
* [Getting Started](#getting-started)
* [Commands](#commands)
* [Running from source](#running-from-source)
* [Contributing](#contributing)
* [License](#license)

## Modules
* Role Manager
* Ping
* More to come!

## Getting Started
To get the public bot, you can just click [here]().

To run it yourself, grab one of the [releases](https://github.com/omar-h/snorlax/releases).<br>
Once downloaded go to the location of the downloaded binary, and run it like so:
```Bash
$ ./snorlax -token=<your-bot-token>
```

## Commands
Visit [the website](https://www.snorlaxbot.com/commands) for a full list of commands.

## Running From Source
To run the bot from source, you need to have Go installed.<br>
You also need to define a [GOPATH](https://golang.org/doc/code.html#GOPATH).

Follow these instructions to compile and run Snorlax form source.
```Bash
$ go get github.com/omar-h/snorlax
$ go install github.com/omar-h/snorlax/cmd/snorlax
```
It will install it in the bin folder in your GOPATH ($GOPATH/bin).<br>
If you have $GOPATH/bin in your PATH variable, you will be able to run it like so:
```Bash
$ snorlax -token=<your-bot-token>
```
To run it in debug mode, just add the `-debug` flag:
```Bash
$ snorlax -token=<your-bot-token> -debug
```

## Contributing
If you've found a bug, or have a suggestion feel free to open an [issue](https://github.com/omar-h/snorlax/issues).

You can contact the author on Discord: Omar H.#6299 or via [email](mailto:contact@omarh.net).

## License
Snorlax is licensed under the [MIT License](https://github.com/omar-h/snorlax/blob/master/LICENSE).
