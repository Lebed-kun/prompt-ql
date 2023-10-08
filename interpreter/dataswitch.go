package interpreterimpl

import (
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/core"

	openqueryswitch "gitlab.com/jbyte777/prompt-ql/v2/interpreter/dataswitches/openquery"
	defaultswitch "gitlab.com/jbyte777/prompt-ql/v2/interpreter/dataswitches/default"
	wrapperswitch "gitlab.com/jbyte777/prompt-ql/v2/interpreter/dataswitches/wrapper"
	rootswitch "gitlab.com/jbyte777/prompt-ql/v2/interpreter/dataswitches/root"
)

type TRootSwitchTable map[string]interpreter.TDataSwitchFunction

var rootSwitchTable TRootSwitchTable = map[string]interpreter.TDataSwitchFunction{
	"root": rootswitch.RootSwitch,
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
