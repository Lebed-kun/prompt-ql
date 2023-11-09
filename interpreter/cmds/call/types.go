package callcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"
)

type TCmdCallableFunction func(
	args interpreter.TFunctionInputChannel,
) interface{}
