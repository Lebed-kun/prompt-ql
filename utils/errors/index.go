package errorsutils

import (
	"fmt"

	timeutils "gitlab.com/jbyte777/prompt-ql/v3/utils/time"
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
