package testutils

import (
	"fmt"
	"io/ioutil"
	timeutils "gitlab.com/jbyte777/prompt-ql/v5/utils/time"
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

func SaveToFile(args []interface{}) interface{} {
	if len(args) < 2 {
		return fmt.Errorf("path and blob must be specified")
	}

	path, isPathStr := args[0].(string)
	if !isPathStr {
		return fmt.Errorf("path is not a string")
	}

	blob, isBlobBytes := args[1].([]byte)
	if !isBlobBytes {
		return fmt.Errorf("blob is not a byte slice")
	}
	
	err := ioutil.WriteFile(
		path,
		blob,
		0666,
	)
	return err
}
