package utils

import (
	"errors"
	"regexp"
	"strconv"
)

// ExtractUserIDFromMention takes in a mention (<@ID>) and spits out only the ID.
func ExtractUserIDFromMention(mention string) string {
	if mention[:2] == "<@" && mention[len(mention)-1:] == ">" {
		return mention[2:][:len(mention)-3] // -3 because we remove 2 chars.
	}

	return mention
}

// GetStringFromQuotes finds a "string which spans multiple spaces" in a split
// message. It then takes that and replaces the Quote string with a single
// string value of the quote contents.
func GetStringFromQuotes(parts []string) []string {
	var (
		// str is the string we're searching for in quotes.
		str string
		// startQuote holds the location of the quote
		startQuote int
	)

	// the length of the original parts
	length := len(parts)
	startQuote = -1
	val := ""
	for k := 0; k < length; k++ {
		val = parts[k]
		if len(val) == 0 {
			continue
		}

		switch {
		// If a startQuote hasn't been found, and the first byte is a quote,
		// then set the startQuote, and start the string.
		case val[0] == '"' && startQuote == -1:
			if val[len(val)-1] == '"' {
				parts[k] = val[:len(val)-1][1:]
			} else {
				startQuote = k
				str = val[1:] + " "
			}
		// If a startQuote has been found, and the last byte is a quote, then
		// remove the old parts and append the new ones.
		case val[len(val)-1] == '"' && startQuote >= 0:
			// Take the parts before startQuote, append the quote string (str)
			// to it, and take the parts after the current index and append
			// them to the new parts.
			parts = append(append(parts[:startQuote], str+val[:len(val)-1]), parts[k+1:]...)
			// Reset k to be at the index just after the current combined
			// string - so we don't check the combined string.
			newLen := len(parts)
			k = k - (length - newLen) + 1
			length = newLen
			startQuote = -1
		default:
			// If a start quote has been found, just add current value
			// to the current quote string.
			if startQuote >= 0 {
				str = str + val + " "
			}
		}
	}

	return parts
}

var hexColorRegex = regexp.MustCompile("^([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")

// This holds the errors that HexColorToInt will throw.
var (
	ErrColorInvalid = errors.New("hex code is invalid")
)

// HexColorToInt converts a hex color code to an int, which is usable in the
// Discord API. The hexCode is without the "#" prefix.
func HexColorToInt(hexCode string) (int, error) {
	if !hexColorRegex.MatchString(hexCode) {
		return 0, ErrColorInvalid
	}

	color, err := strconv.ParseInt(hexCode, 16, 32)
	if err != nil {
		return 0, err
	}

	return int(color), nil
}
