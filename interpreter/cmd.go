package interpreterimpl

import (
	api "gitlab.com/jbyte777/prompt-ql/v3/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v3/custom-apis"
	interpreter "gitlab.com/jbyte777/prompt-ql/v3/core"

	callcmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/call"
	getcmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/get"
	listenquerycmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/listenquery"
	openquerycmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/openquery"
	setcmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/set"
	wrappercmds "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/wrapper-cmds"
	hellocmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/hello"
	opensessioncmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/open-session"
	closesessioncmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/close-session"
	headercmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/header"
	embedifcmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/embed-if"
	nopcmd "gitlab.com/jbyte777/prompt-ql/v3/interpreter/cmds/nop"
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
		"header": headercmd.HeaderCmd,
		"embed_if": embedifcmd.EmbedIfCmd,
		"nop": nopcmd.NopCmd,
	}
}
