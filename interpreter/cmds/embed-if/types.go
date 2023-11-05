package embedifcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"
)

type TCondCallableFunction func(
	args interpreter.TFunctionInputChannel,
) bool
