package interpreterimpl

import (
	api "gitlab.com/jbyte777/prompt-ql/v4/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v4/custom-apis"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/core"

	// basic commands
	callcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/call"
	getcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/get"
	listenquerycmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/listenquery"
	openquerycmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/openquery"
	setcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/set"
	wrappercmds "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/wrapper-cmds"

	// communication commands
	hellocmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/hello"
	headercmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/header"

	// execution life-cycle commands
	opensessioncmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/open-session"
	closesessioncmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/close-session"
	unsafeclearvarscmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/unsafe-clear-vars"
	unsafeclearstackcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/unsafe-clear-stack"

	// code embedding commands
	embedifcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/embed-if"
	embeddefcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/embed-def"
	embedexpcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/embed-exp"

	// misc
	nopcmd "gitlab.com/jbyte777/prompt-ql/v4/interpreter/cmds/nop"
)

func makeCmdTable(
	gptApi *api.GptApi,
	customApis *customapis.CustomModelsApis,
) interpreter.TExecutedFunctionTable {
	return interpreter.TExecutedFunctionTable{
		// basic commands
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

		// communication commands
		"hello": hellocmd.MakeHelloCmd(gptApi, customApis),
		"header": headercmd.HeaderCmd,

		// execution life-cycle commands
		"session_begin": opensessioncmd.OpenSessionCmd,
		"session_end": closesessioncmd.CloseSessionCmd,
		"unsafe_clear_vars": unsafeclearvarscmd.UnsafeClearVarsCmd,
		"unsafe_clear_stack": unsafeclearstackcmd.UnsafeClearStackCmd,

		// code embedding commands
		"embed_if": embedifcmd.EmbedIfCmd,
		"embed_def": embeddefcmd.EmbedDefCmd,
		"embed_exp": embedexpcmd.EmbedExpCmd,

		// misc
		"nop": nopcmd.NopCmd,
	}
}
