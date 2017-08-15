package utils

import (
	"strings"
)

// ExtractUserIDFromMention takes in a mention (<@ID>) and spits out only the ID.
func ExtractUserIDFromMention(mention string) string {
	return strings.Replace(strings.Replace(strings.Replace(mention, "<", "", -1), ">", "", -1), "@", "", -1)
}
