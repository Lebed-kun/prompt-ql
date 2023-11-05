package callcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"
)

type TCmdCallableFunction func(
	args interpreter.TFunctionInputChannel,
) interface{}
