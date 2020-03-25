package strings

import (
	"regexp"
	"strings"
)

//Capitalize return a capitalize text
func Capitalize(text string) string {
	wordsCapitalized := []string{}
	for _, word := range strings.Split(text, " ") {
		regex, _ := regexp.Compile("[a-zA-Z0-9]+")
		matchedString := regex.FindString(word)

		if len(matchedString) > 0 {
			wordsCapitalized = append(wordsCapitalized, strings.ToUpper(matchedString[0:1])+strings.ToLower(matchedString[1:]))
		}
	}

	return strings.Join(wordsCapitalized, " ")
}
