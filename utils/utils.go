package utils

import (
	"strings"
)

// ExtractUserIDFromMention takes in a mention (<@ID>) and spits out only the ID.
func ExtractUserIDFromMention(mention string) string {
	return strings.Replace(strings.Replace(strings.Replace(mention, "<", "", -1), ">", "", -1), "@", "", -1)
}

// GetStringFromQuotes finds a "string which spans multiple spaces" in a split message.
// Then takes that and replaces the Quote string with a single string value of the quote contents.
func GetStringFromQuotes(parts []string) []string {
	var (
		// found is for controlling the loop-
		found bool
		// str is the string we're searching for in quotes.
		str string
		// startQuote holds the location of the quote
		startQuote int
	)

	for k, val := range parts {
		if !found {
			if str == "" && val[0] == '"' {
				// startQuote location is equal the index of the current value.
				startQuote = k
				if val[len(val)-1] == '"' {
					// The string has been found, mark this so it stops looping.
					found = true
					// Since the quoted message is in one location in the slice, we'll just remove the quotes and replace it.
					parts[k] = val[:len(val)-1][1:]
				} else {
					str = str + val[1:] + " "
				}
				// If string isn't started just skip.
			} else if str != "" {
				// If last char of val is a quote, finish the string search.
				if val[len(val)-1] == '"' {
					found = true
					// Explanation: append the result of append(parts[:startQuote], str+" "+val[:len(val)-1]) with parts[k+1:]...
					// First in append(parts[:startQuote], str+" "+val[:len(val)-1]) we're appending parts[:startQuote]
					// parts[:startQuote] is contents of parts which came before where the first quote started.
					// Then we're appending the quoted message, which we've extracted.
					// We're taking this new slice (parts[:startQuote] + the quoted message) and appending parts[k+1:]...
					// parts[k+1:]... is the contents of parts which came before the last quote.
					// The last quote location is in the current index (stored in var k).
					parts = append(append(parts[:startQuote], str+val[:len(val)-1]), parts[k+1:]...)
				} else {
					str = str + val + " "
				}
			}
		}
	}

	return parts
}
