package embedifcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"
)

type TCondCallableFunction func(
	args interpreter.TFunctionInputChannel,
) bool
