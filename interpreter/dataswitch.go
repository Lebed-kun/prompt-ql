package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/core"

	openqueryswitch "gitlab.com/jbyte777/prompt-ql/interpreter/dataswitches/openquery"
	defaultswitch "gitlab.com/jbyte777/prompt-ql/interpreter/dataswitches/default"
	wrapperswitch "gitlab.com/jbyte777/prompt-ql/interpreter/dataswitches/wrapper"
)

type TRootSwitchTable map[string]interpreter.TDataSwitchFunction

var rootSwitchTable TRootSwitchTable = map[string]interpreter.TDataSwitchFunction{
	"root": defaultswitch.DefaultSwitch,
	"open_query": openqueryswitch.OpenQuerySwitch,
	"call": defaultswitch.DefaultSwitch,
	"set": defaultswitch.DefaultSwitch,
	"user": wrapperswitch.WrapperSwitch,
	"assistant": wrapperswitch.WrapperSwitch,
	"system": wrapperswitch.WrapperSwitch,
	"data": wrapperswitch.WrapperSwitch,
	"error": wrapperswitch.WrapperSwitch,
}

func rootDataSwitch(
	topCtx *interpreter.TExecutionStackFrame,
	rawData interface{},
) {
	if rawData == nil {
		return
	}

	switchFn, hasSwitchFn := rootSwitchTable[topCtx.FnName]
	if !hasSwitchFn {
		return
	}

	switchFn(topCtx, rawData)
}
