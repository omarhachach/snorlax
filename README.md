# Snorlax
Snorlax is another modular Discord bot, with loads of built in modules.

**This is still in early development, therefore it is subject to frequent breaking changes and rewrites.**

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
* Gambling
* Eval
* Music
* More to come!

## Getting Started
To get the public bot, you can just click [here]() (not public yet).

To run it yourself, grab one of the [releases](https://github.com/omar-h/snorlax/releases). (Still in early development, see [running from source](#running-from-source)).<br>
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
$ go get github.com/omar-h/snorlax
$ go install github.com/omar-h/snorlax/cmd/snorlax
```
It will install it in the bin folder in your GOPATH ($GOPATH/bin).<br>
If you have $GOPATH/bin in your PATH variable, you will be able to run it like so:
```Bash
$ snorlax
```

Remember to [create a config](#getting-started) before running the bot.

## Contributing
If you've found a bug, or have a suggestion feel free to open an [issue](https://github.com/omar-h/snorlax/issues).

You can contact the author on Discord: Omar H.#6299 or via [email](mailto:contact@omarh.net).

## License
Snorlax is licensed under the [MIT License](https://github.com/omar-h/snorlax/blob/master/LICENSE).
