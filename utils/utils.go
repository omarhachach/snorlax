package utils

import (
	"bytes"
	"strings"
)

// ExtractUserIDFromMention takes in a mention (<@ID>) and spits out only the ID.
func ExtractUserIDFromMention(mention string) string {
	return strings.Replace(strings.Replace(strings.Replace(mention, "<", "", -1), ">", "", -1), "@", "", -1)
}

// GetStringFromQuotes finds a "string which spans multiple spaces" in a split message.
// Then takes that and replaces the Quote string with a single string value of the quote contents.
func GetStringFromQuotes(parts []string) []string {
	found := false
	var buffer bytes.Buffer
	var newParts []string

	for _, val := range parts {
		if !found {
			if val[0] == '"' {
				if val[len(val)-1] == '"' {
					found = true
					buffer.WriteString(val[:len(val)-1][1:])
					newParts = append(newParts, buffer.String())
				} else {
					buffer.WriteString(val[1:])
				}
			} else if buffer.Len() != 0 {
				if val[len(val)-1] == '"' {
					found = true
					buffer.WriteString(" " + val[:len(val)-1])
					newParts = append(newParts, buffer.String())
				} else {
					buffer.WriteString(" " + val)
				}
			} else {
				newParts = append(newParts, val)
			}
		} else {
			newParts = append(newParts, val)
		}
	}

	return newParts
}
