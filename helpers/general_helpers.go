package helpers

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Display a warning if there was an error
// Return true if a warning was displayed
// Return false otherwise
func Warn(err error, message string) bool {
	if err != nil {
		println(fmt.Sprintf("Warning: %s", message))
		return true
	}

	return false
}

func WarnWithPrompt(err error, message string, prompt string, cb func(resp string)) error {
	if !Warn(err, message) {
		return nil
	}

	print(fmt.Sprintf("%s : ", prompt))
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return err
	}

	defer cb(response)

	return nil
}

// Look for the search value in the provided haystack.
// If found, set the value of needle to the first matching value in the haystack
// Return true or false respectively if the needle was in the haystack or not
func SearchFor[N comparable, H []N](search N, haystack H, needle *N) bool {
	// s for straw. Clever right?
	for _, s := range haystack {
		if s == search {
			if needle != nil {
				*needle = s
			}
			return true
		}
	}

	return false
}

// Return true if the needle is in the haystack
func IsIn[N comparable, H []N](needle N, haystack H) bool {
	return SearchFor(needle, haystack, nil)
}

// Convert ascii to integer in one line by swallowing errors.
// If you are not absolutely certain that s is a valid convertable string, don't use this method.
// Prefer strconv when possible.
func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Return true if the file being executed is a go test file
func Testing() bool {
	// allow tests to overwrite the default behavior
	if os.Getenv("GO_ENVIRON") == "PROD" {
		return false
	}
	teststr := os.Args[len(os.Args)-1]
	if matched, _ := regexp.Match(`^\^Test\w+\$$`, []byte(teststr)); matched {
		println("NOTICE: RUNNING IN TEST ENVIRONMENT. EXPECT INCONSISTENCIES")
		return true
	}

	return false
}

func HandleError(err error, msg ...string) {
	if err != nil {
		replacer := strings.NewReplacer("{ERROR}", err.Error())
		for i, s := range msg {
			msg[i] = replacer.Replace(s)
		}
		log.Fatal(msg)
	}
}
