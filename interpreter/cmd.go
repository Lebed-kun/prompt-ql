package interpreterimpl

import (
	api "gitlab.com/jbyte777/prompt-ql/v5/api"
	customapis "gitlab.com/jbyte777/prompt-ql/v5/custom-apis"
	loggerapis "gitlab.com/jbyte777/prompt-ql/v5/logger-apis"
	interpreter "gitlab.com/jbyte777/prompt-ql/v5/core"

	// basic commands
	callcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/call"
	getcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/get"
	listenquerycmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/listenquery"
	openquerycmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/openquery"
	setcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/set"
	wrappercmds "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/wrapper-cmds"
	nopcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/nop"
	
	// communication commands
	hellocmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/hello"
	headercmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/header"
	msgbegincmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/msg-begin"
	msgendcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/msg-end"
	msgrestartchaincmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/msg-restart-chain"

	// execution life-cycle commands
	opensessioncmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/open-session"
	closesessioncmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/close-session"
	unsafeclearvarscmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/unsafe-clear-vars"
	unsafepreinitvarscmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/unsafe-preinit-vars"
	unsafeclearstackcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/unsafe-clear-stack"

	// code embedding commands
	embedifcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/embed-if"
	embeddefcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/embed-def"
	embedexpcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/embed-exp"

	// blob data commands
	blobreadfromfilecmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/blob-read-from-file"
	blobreadfromurlcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/blob-read-from-url"

	// debugging commands
	debugcmd "gitlab.com/jbyte777/prompt-ql/v5/interpreter/cmds/debug"
)

var cmdsMetaInfo = interpreter.TCommandMetaInfoTable{
	"debug": &interpreter.TCommandMetaInfo{
		IsErrorTolerant: true,
	},
	"session_begin": &interpreter.TCommandMetaInfo{
		IsErrorTolerant: true,
	},
	"session_end": &interpreter.TCommandMetaInfo{
		IsErrorTolerant: true,
	},
	"unsafe_clear_vars": &interpreter.TCommandMetaInfo{
		IsErrorTolerant: true,
	},
	"unsafe_preinit_vars": &interpreter.TCommandMetaInfo{
		IsErrorTolerant: true,
	},
	"unsafe_clear_stack": &interpreter.TCommandMetaInfo{
		IsErrorTolerant: true,
	},
}

func makeCmdTable(
	gptApi *api.GptApi,
	customApis *customapis.CustomModelsApis,
	loggerApis *loggerapis.LoggerApis,
	readFromFileTimeoutSec uint,
	readFromUrlTimeoutSec uint,
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
		"nop": nopcmd.NopCmd,

		// communication commands
		"hello": hellocmd.MakeHelloCmd(gptApi, customApis),
		"header": headercmd.HeaderCmd,
		"msg_begin": msgbegincmd.MsgBeginCmd,
		"msg_end": msgendcmd.MsgEndCmd,
		"msg_restart_chain": msgrestartchaincmd.MsgRestartChainCmd,

		// execution life-cycle commands
		"session_begin": opensessioncmd.OpenSessionCmd,
		"session_end": closesessioncmd.CloseSessionCmd,
		"unsafe_clear_vars": unsafeclearvarscmd.UnsafeClearVarsCmd,
		"unsafe_preinit_vars": unsafepreinitvarscmd.UnsafePreinitVarsCmd,
		"unsafe_clear_stack": unsafeclearstackcmd.UnsafeClearStackCmd,

		// code embedding commands
		"embed_if": embedifcmd.EmbedIfCmd,
		"embed_def": embeddefcmd.EmbedDefCmd,
		"embed_exp": embedexpcmd.EmbedExpCmd,

		// blob data commands
		"blob_read_from_file": blobreadfromfilecmd.MakeBlobReadFromFileCmd(readFromFileTimeoutSec),
		"blob_read_from_url": blobreadfromurlcmd.MakeBlobReadFromUrlCmd(readFromUrlTimeoutSec),

		// debugging commands
		"debug": debugcmd.MakeDebugCmd(loggerApis),
	}
}
