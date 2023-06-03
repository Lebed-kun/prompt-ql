package errorsutils

import (
	"fmt"

	timeutils "gitlab.com/jbyte777/prompt-ql/utils/time"
)

func LogError(
	moduleName string,
	methodName string,
	err error,
) error {
	return fmt.Errorf(
		"[%v] [%v / %v] => %v",
		timeutils.NowTimestamp(),
		moduleName,
		methodName,
		err.Error(),
	)
}
