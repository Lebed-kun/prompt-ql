package callcmd

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/core"
)

type TCmdCallableFunction func(
	args interpreter.TFunctionInputChannel,
) interface{}
