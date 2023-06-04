package interpreterimpl

import (
	api "gitlab.com/jbyte777/prompt-ql/api"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"

	callcmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/call"
	getcmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/get"
	listenquerycmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/listenquery"
	openquerycmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/openquery"
	setcmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/set"
	wrappercmds "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/wrapper-cmds"
)

func makeCmdTable(
	gptApi *api.GptApi,
) interpreter.TExecutedFunctionTable {
	return interpreter.TExecutedFunctionTable{
		"open_query": openquerycmd.MakeOpenQueryCmd(gptApi),
		"listen_query": listenquerycmd.MakeListenQueryCmd(gptApi),
		"call": callcmd.CallCmd,
		"get": getcmd.GetCmd,
		"set": setcmd.SetCmd,
		"user": wrappercmds.MakeWrapperCmd("user"),
		"assistant": wrappercmds.MakeWrapperCmd("assistant"),
		"system": wrappercmds.MakeWrapperCmd("system"),
		"data": wrappercmds.MakeWrapperCmd("data"),
		"error": wrappercmds.MakeWrapperCmd("error"),
	}
}
