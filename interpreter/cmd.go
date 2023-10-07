package interpreterimpl

import (
	api "gitlab.com/jbyte777/prompt-ql/api"
	customapis "gitlab.com/jbyte777/prompt-ql/custom-apis"
	interpreter "gitlab.com/jbyte777/prompt-ql/core"

	callcmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/call"
	getcmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/get"
	listenquerycmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/listenquery"
	openquerycmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/openquery"
	setcmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/set"
	wrappercmds "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/wrapper-cmds"
	hellocmd "gitlab.com/jbyte777/prompt-ql/interpreter/cmds/hello"
)

func makeCmdTable(
	gptApi *api.GptApi,
	customApis *customapis.CustomLLMApis,
) interpreter.TExecutedFunctionTable {
	return interpreter.TExecutedFunctionTable{
		"open_query": openquerycmd.MakeOpenQueryCmd(gptApi, customApis),
		"listen_query": listenquerycmd.MakeListenQueryCmd(gptApi, customApis),
		"call": callcmd.CallCmd,
		"get": getcmd.GetCmd,
		"set": setcmd.SetCmd,
		"user": wrappercmds.MakeWrapperCmd("user"),
		"assistant": wrappercmds.MakeWrapperCmd("assistant"),
		"system": wrappercmds.MakeWrapperCmd("system"),
		"data": wrappercmds.MakeWrapperCmd("data"),
		"error": wrappercmds.MakeWrapperCmd("error"),
		"hello": hellocmd.MakeHelloCmd(gptApi, customApis),
	}
}
