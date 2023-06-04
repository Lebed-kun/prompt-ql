package callcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"
)

type TCmdCallableFunction func(
	args interpreter.TFunctionInputChannel,
) interface{}
