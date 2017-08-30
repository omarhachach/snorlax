package utils

import (
	"bytes"
	"strings"
)

// ExtractUserIDFromMention takes in a mention (<@ID>) and spits out only the ID.
func ExtractUserIDFromMention(mention string) string {
	return strings.Replace(strings.Replace(strings.Replace(mention, "<", "", -1), ">", "", -1), "@", "", -1)
}

// GetStringFromParts finds a "String which spans multiple spaces" in a split message.
func GetStringFromParts(parts []string) (string, []string) {
	quotesFound := false
	finish := false
	var buffer bytes.Buffer
	var newParts []string

	for _, val := range parts {
		if !finish {
			if !quotesFound && val[0] == '"' {
				quotesFound = true
				if val[len(val)-1] == '"' {
					finish = true
					val = val[:len(val)-1]
				}
				buffer.WriteString(val[1:])
			} else if quotesFound {
				if val[len(val)-1] == '"' {
					finish = true
					val = val[:len(val)-1]
				}
				buffer.WriteString(" " + val)
			} else {
				newParts = append(newParts, val)
			}
		} else {
			newParts = append(newParts, val)
		}
	}
	return buffer.String(), newParts
}
