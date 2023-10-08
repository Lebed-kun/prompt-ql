package interpreterimpl

import (
	api "gitlab.com/jbyte777/prompt-ql/v2/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v2/custom-apis"
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/core"

	callcmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/call"
	getcmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/get"
	listenquerycmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/listenquery"
	openquerycmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/openquery"
	setcmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/set"
	wrappercmds "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/wrapper-cmds"
	hellocmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/hello"
	opensessioncmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/open-session"
	closesessioncmd "gitlab.com/jbyte777/prompt-ql/v2/interpreter/cmds/close-session"
)

func makeCmdTable(
	gptApi *api.GptApi,
	customApis *customapis.CustomModelsApis,
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
		"session_begin": opensessioncmd.OpenSessionCmd,
		"session_end": closesessioncmd.CloseSessionCmd,
	}
}
