package errorhandler

import (
	"fmt"
)

func FailOnError(err error, msg string) {
	if err != nil {
		message := fmt.Sprintf("%s: %s", msg, err.Error())
		fmt.Println(message)

	}
}

func Fatal(err error, msg string) {
	if err != nil {
		message := fmt.Sprintf("%s: %s", msg, err.Error())
		fmt.Println(message)
	}
}

func LogInfo(message string) {
	fmt.Println(message)
}
