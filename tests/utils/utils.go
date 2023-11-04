package testutils

import (
	"fmt"
	timeutils "gitlab.com/jbyte777/prompt-ql/v3/utils/time"
)

func LogTimeForProgram(args []interface{}) interface{} {
	if len(args) < 1 {
		return ""
	}

	log, isLogStr := args[0].(string)
	if !isLogStr {
		return ""
	}

	fmt.Printf(
		"[%v] %v",
		timeutils.NowTimestamp(),
		log,
	)

	return ""
}
