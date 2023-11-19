package user

import (
  "fmt"
  "log"
  "regexp"
)

func IsUser(input string) bool {
	matched, err := regexp.MatchString(`<@!?[0-9]*>`, input)
	if err != nil {
		log.Println("Error checking for valid user:", err)
	}
	return matched
}

func UserToString(input string) string {
	r := regexp.MustCompile(`<@!?([0-9]*)>`)
	return r.FindStringSubmatch(input)[1]
}

func StringToUser(input string) string {
	return fmt.Sprintf("<@%s>", input)
}
