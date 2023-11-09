﻿# PromptQL: Flexible language for prompting agents based on ML models

<img src="./readme-content/simple-dialog.png" />
<img src="./readme-content/simple-dialog2.png" />
<img src="./readme-content/simple-dialog3.png" />
<img src="./readme-content/simple-dialog4.png" />

It's a zero-dependencies library for orchestrating agents based on ML models like `gpt3.5-turbo` . The default ML model API is based on the OpenAI API: https://platform.openai.com/docs/api-reference . Full list of supported models is here: https://platform.openai.com/docs/models/model-endpoint-compatibility 


## Getting started
```
go get -u gitlab.com/jbyte777/prompt-ql/vX - for >= v2.x release
go get -u gitlab.com/jbyte777/prompt-ql - for v1.x release
```

Making a basic query is just like writing plain HTML or another template:
```
func BasicQueryTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			OpenAiBaseUrl: openAiBaseUrl,
			OpenAiKey: openAiKey,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~open_query to="query1" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to machine learning
			{/open_query}
			{~listen_query from="query1" /}
		`,
	)

	// ...
}
```

Then you can extract its result like this:
```
resultStr, _ := result.ResultDataStr()
errStr, _ := result.ResultErrorStr()
```


## Opening a query doesn't block execution of code

You can easily batch multiple queries without waiting for completion of previously sent query:

```
func NonBlockingQueriesTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"logtime": testutils.LogTimeForProgram,
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			OpenAiBaseUrl: openAiBaseUrl,
			OpenAiKey: openAiKey,
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~open_query to="query1" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to learn statistics step by step.
			{/open_query}
			{~call fn=@logtime }
				open query1
			{/call}
			=======================
			{~open_query to="query2" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to make a solar panel step by step.
			{/open_query}
			{~call fn=@logtime }
				open query2
			{/call}
			=======================
			Answer1: {~listen_query from="query1" /}
			{~call fn=@logtime }
				listen query1
			{/call}
			=======================
			Answer2: {~listen_query from="query2" /}
			{~call fn=@logtime }
				listen query2
			{/call}
			=======================
		`,
	)

	// ...
}
```

This prints a log with timestamp after execution of each PromptQL command (remember that user-defined function can do anything):
<img src="./readme-content/non-blocking-requests/logs.png" />


## You can define your own ML models APIs

This allows to extend default set of PromptQL models beyond OpenAI capabilities. For example, you can define an API to your local Llama model like this:

```
func makeLlamaDoQuery(
	pathToLlamaCommand string,
	pathToLlamaModel string,
) customapis.TDoQueryFunc {
	return func(
		model string,
		temperature float64,
		inputs interpretercore.TFunctionInputChannelTable,
		execInfo interpretercore.TExecutionInfo,
	) (string, error) {
		prompt := llamaComposePrompt(inputs)

		cmd := exec.Command(
			pathToLlamaCommand,
			"-m",
			pathToLlamaModel,
			"--temp",
			fmt.Sprintf("%.1f", temperature),
			"-p",
			fmt.Sprintf("\"%v\"", prompt),
		)

		res, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf(
				"ERROR (line=%v, charpos=%v): %v",
				execInfo.Line,
				execInfo.CharPos,
				err.Error(),
			)
		}

		return string(res), nil
	}
}
```


Then bind it to PromptQL:

```
llamaDoQuery := makeLlamaDoQuery(pathToLlamaCommand, pathToLlamaModel)
interpreterInst.CustomApis.RegisterModelApi(
	"llama",
	llamaDoQuery,
)
```

And finally execute your query. Provide additional `user` flag to the `open_query` command. With it PromptQL knows it's user defined model, not OpenAI's:

```
result := interpreterInst.Instance.Execute(
		`
			{~open_query user to="query1" model="llama"}
				{~system}
					You are a helpful assistant.
				{/system}
				I want a response to the following question:
				Write a guide to cook pasta
			{/open_query}
			{~listen_query from="query1" /}
		`,
	)
```

<img src="./readme-content/custom-llama-model/result.png" />


## Post-process answer from ML model with user defined functions
You can define your own functions for query program. This allows you to prettify ML model output for example:

```
func QueryWithPostprocessFunctionTest(
	openAiBaseUrl string,
	openAiKey string,
) {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"postprocess": postProcessFunctionTest,
	}
	interpreterInst := interpreter.New(
		interpreter.TPromptQLOptions{
			OpenAiBaseUrl: openAiBaseUrl,
			OpenAiKey: openAiKey,
			DefaultExternalGlobals: defaultGlobals,
		},
	)

	result := interpreterInst.Instance.Execute(
		`
			{~open_query to="query1" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to machine learning step by step
			{/open_query}
			{~set to="queryres"}
				{~listen_query from="query1" /}
			{/set}
			Raw result is:
			{~get from="queryres" /}

			JSON result is:
			{~call fn=@postprocess }
				{~get from="queryres" /}
			{/call}
		`,
	)

	// ...
}

```
This gives you this output for example:
<img src="./readme-content/query-with-fn1.jpg" />

<img src="./readme-content/query-with-fn2.jpg" />

<img src="./readme-content/query-with-fn3.jpg" />


## Execute queries partially (and continously)

PromptQL v2.x+ is designed as an imperative protocol for ML-based agents. So partial execution of code is achievable at language level, with the `{~session_begin /}` and `{~session_end /}` commands. However, API of PromptQL v2.x for Golang still supports the `ExecutePartial` method: it executes PromptQL program just like session is always opened during it.

Session is like a thread of execution: all internal variables keep their state, execution context stack keeps its state during it.

You can view examples for this in a test located by a path of `/tests/basic-functionality/sessions` 


## External vs internal variables

PromptQL v2.x makes difference between external and internal variables as it's message-oriented protocol. These kinds of variables differ in few points:

* Accessing undefined external variables resolves into a run-time error. Accessing undefined internal variables resolves into `nil` value;
* Direct writes to external variables resolve into a run-time error. Internal variables can be both read and written; 
* External variables values are preserved between sessions. Internal variables values are lost after completing a PromptQL session;

For more examples, you can visit the `/tests/basic-functionality/external-vs-internal-variables` tests.


## Use references in your commands

References are names of variables in some interpreter table. They are prefixed with `$` sign if it refernces to internal variable. Or with `$@` if it references to external variable. You can use them for non-string values, variadic commands etc.

```
result := interpreterInst.Execute(
		`
			{~open_query to="query1" model="gpt-3.5-turbo-16k"}
				{~system}
					You are a helpful and terse assistant.
				{/system}
				I want a response to the following question:
				Write a comprehensive guide to machine learning
			{/open_query}
			{~$cmd $cmdarg=$cmdval /}
		`,
		interpretercore.TGlobalVariablesTable{
			"cmd": "listen_query",
			"cmdarg": "from",
			"cmdval": "query1",
		},
	)
```


## Supported PromptQL commands v4.x
### 1. Basic query commands
They are basic building blocks of a query to ML models. They have been introduced since the v1.x version

 - `{~open_query user to="X" model="Y" temperature="Z"}<execution_text>{/open_query}` - sends prompt request for given ML model that's defined by `<execution_text>` . It doesn't block execution of query. The command doesn't return any data. `<execution_text>` defines an input data for the command as follows:
```
 - "!user <text>", "!data <text>" -> USER input channel;
 - "!assistant <text>" -> ASSISTANT input channel;
 - "!system <text>" -> SYSTEM input channel;
 - "!error <text>" -> ERROR input channel;
```

Static arguments for the command are:
```
 - "user" - is a flag. If it's set, it **forces** a `model` to match user-defined model. If it's not set, then `model` is assumed as an OpenAI model **by default**. In this case if `model` is not supported by OpenAI, then `model` is assumed as a user-defined model;
 - "sync" - is a flag. If it's set, then query is executed in a blocking manner and returns a text response like the `listen_query` command do. If it's not set, then query is executed in parallel and result is stored in the `to` handle;
 - "to" - is a name of variable to store a query handle. Variable name prefixed with the "@" sign is a name of external variable, otherwise it's internal variable name. "to" is a required parameter if query is **asynchronous** (i.e. no `sync` flag). It's not required if query is **synchronous** (i.e. `sync` flag is set);
 - "model" - is a name of chosen ML model. Default value is "gpt-3.5-turbo";
 - "temperature" - is a temperature of chosen ML model. Default value is 1.0;
```

 - `{~listen_query from="X" /}` - waits for OpenAI ML model query from "X" variable to complete. The command doesn't receive any additional inputs. It returns a text with the `!assistant` tag if succeed, otherwise it returns an error with the `!error` tag;
 
Static arguments for the command are:
```
 - "from" - is a name of variable from which result is fetched. Variable name prefixed with the "@" sign is a name of external variable, otherwise it's internal variable name. "from" is a required parameter;
```

 - `{~call fn="F"}<execution_text>{/call}` - calls function from `fn` variable. Command returns error if `fn` variable doesn't exist or the variable doesn't contain function with the type `func([]interface{}) interface{}` . Otherwise the command returns a data from the execution of `fn`. `<execution_text>` defines an input data for the command as follows:
```
 - "!user <text>", "!assistant <text>", "!system <text>", "!data <text>" -> DATA channel;
 - "!error <text>" -> ERROR channel;
 - error -> ERROR channel;
 - text without a tag -> DATA channel;
 - any non-string and non-error value -> DATA channel;
```
`DATA` channel contains array of arguments for function

Static arguments for the command are:
```
 - "fn" - is a name of variable where called function is stored. Variable name prefixed with the "@" sign is a name of external variable, otherwise it's internal variable name. "fn" is a required parameter;
```

 - `{~get from="X" /}` - gets data from the `from` variable. The command doesn't receive any additional data;

Static arguments for the command are:
```
 - "from" - is a name of variable from which data is retrieved. Variable name prefixed with the "@" sign is a name of external variable, otherwise it's internal variable name. "from" is a required parameter;
```

 - `{~set to="X"}<execution_text>{/set}` -  stores data defined by `<execution_text>` in the `X` variable. The command doesn't return any value. `<execution_text>` defines an input data for the command as follows:
```
 - "!user <text>", "!assistant <text>", "!system <text>", "!data <text>" -> DATA channel;
 - "!error <text>" -> ERROR channel;
 - error -> ERROR channel;
 - text without a tag -> DATA channel;
 - any non-string and non-error value -> DATA channel;
```

Static arguments for the command are:
```
 - "to" - is a name of variable to which data is stored. Variable name prefixed with the "@" sign is a name of external variable, otherwise it's internal variable name. "to" is a required parameter;
```

 - Wrapper commands. They wrap a text with corresponding prompt tag: `!user`, `!assistant`, `!system`, `!data` or `!error`. This is useful for separating roles of ML model query texts, for specific error handling etc. They are defined like this:
```
{~user}<execution_text>{/user}
{~assistant}<execution_text>{/assistant}
{~system}<execution_text>{/system}
{~data}<execution_text>{/data}
{~error}<execution_text>{/error}
```
They receive all input data in the `DATA` channel;

### 2. Agent messaging commands
They are useful for communicating API of some called agent to calling agent. They have been introduced since the v2.x version

 - `{~hello /}` - returns a set of ML models, external variables and code embeddings. They're defined in given PromptQL instance. It's useful for acknowleding user or other automatic agent of given PromptQL agent capabilities. The command returns a JSON string with following structure:
 ```
   {
      "myModels": {
				"gpt-4": "description of gpt-4",
				...
				"myModel": "description of myModel",
			},
      "myVariables": {
				"myVar1": "description of myVar1",
				...
			},
			"myEmbeddings": {
				"myEmbedding1": "description of myEmbedding1",
				...
			},
	 }
 ```
 - `{~header from="Sender agent" to="Receiver agent" /}` - returns a message header formatted in JSON. It's useful for dynamic routing of PromptQL message to arbitrary known agent. The command returns a JSON string with following structure:
 ```
   {
      "fromAgent": "Sender agent id/name",
			"toAgent": "Receiver agent id/name",
	 }
 ```

### 3. Execution life-cycle commands
They are useful for controlling agent's execution flow at language level. They have been introduced since the v2.x version

 - `{~session_begin /}` - opens a current execution session. After opening a session and execution of PromptQL chunk, a state of interpreter is saved (except its cursor pointing to program text). The command brings a basic management of execution flow to protocol/language level;
 - `{~session_end /}` - closes a current execution session. After closing a session and execution of PromptQL chunk, a full state of interpreter is lost. The command brings a basic management of execution flow to protocol/language level;

### 4. Code embedding commands
They are useful for embedding PromptQL code as data in messages. It's useful for later code execution: by forwarded agent, for separation of interfacing and implementation etc. They have been introduced since the v3.x version.

Embeddable can contain placeholders defined as a special literal with the `%` sign like `%embd_arg`. They are resolved on code expansion.

 - `{~embed_if cond=@conditionFunc}<arg1><arg2>...<yes_branch><no_branch>{/embed_if}` - checks condition `cond` with arguments `<arg1><arg2>...`. If it's `true`, then `<yes_branch>` is returned. Otherwise, `<no_branch>` is returned. `cond` should have a type of `func([]interface{}) bool`. The command is useful for conditionally embedding PromptQL code on executing agent side;
 - `{~embed_def name="embed_name" desc="embed_description"}<PromptQL code as text>{/embed_def}` - registers a `<PromptQL code as text>` as an expandable chunk of code for later expansion and execution by `embed_name`. `embed_description` can be optionally provided for communicating layout of embedding in the `{~hello /}` command;
 - `{~embed_exp name="embed_name"}<arg1=val1>...<argN=valN>{/embed_exp}` - expands a PromptQL code defined as embedding `embed_name`. `<arg1=val1>...<argN=valN>` can be optionally provided to pass placeholders to embedding. Placeholder is defined as a special literal like: `%arg1`, ..., `%argN`;

### 5. Miscellaneous commands
- `{~nop /}` - returns empty characters sequence `\x00` . It's a phantom command that can serve as an argument filler for other PromptQL commands;


## Additional features
 - References to entries in some global variables table are supported. You can use them by prefixing a name with the `$` sign like:
```
For references to internal variables:
{~$command $arg=$val /}

For references to external variables:
{~$@command $@arg=$@val /}
```
 - Defining custom ML model APIs. It can be obtained with the `RegisterModelApi` method (see below)
 - Code embedding. You can include the PromptQL code in the brackets like this: `<% {~open_query sync}<other PromptQL code>{/open_query} %>`. And the code inside brackets won't be executed. Instead it will be returned just like plain string. This is useful for late PromptQL code processing or forwarding it to other agents.


## Interpreter API
 - `func New(options TPromptQLOptions) *TPromptQL` - creates an instance of PromptQL with default "closed" state of session. 

 The function receives parameters from the `TPromptQLOptions` structure that contains:
 ```
 - "OpenAiBaseUrl" - is an URL to OpenAI API. For example, "https://api.openai.com". It's a required parameter for OpenAI models use-cases. Otherwise it can be omitted;
 - "OpenAiKey" - is your OpenAI API key. You can set up it on "https://platform.openai.com/account/api-keys". It's a required parameter for OpenAI models use-cases. Otherwise it can be omitted;
 - "OpenAiListenQueryTimeoutSec" - is a timeout for listening prompting query from an OpenAI model. Default value is 30 seconds;
 - "CustomApisListenQueryTimeoutSec" - is a timeout for listening prompting query from a user-defined ML model. Default value is 30 seconds;
 - "DefaultExternalGlobals" - is a table of predefined external global variables. You can provide there custom functions, constants, services etc. Default value is *nil*;
 ```

 - PromptQL.Instance methods:
	- Basic API:
		- `func (self *Interpreter) Execute(program string) *TInterpreterResult` - executes query as a part of **current** session. I.e. if session is closed, then state of interpreter is completely reset after execution. Otherwise only interpreter cursor is reset.

		The method receives parameters:
		```
		- "program" - is an executed PromptQL program;
		```

		The method returns `*TInterpreterResult` which consists of:
		```
		- "Result" - is a collection of input channels for root context (which represents a final result). It contains "data" and "error" channels;
		- "Error" - is a parsing error;
		- "Complete" - is a flag for completeness of execution of PromptQL program. Execution of PromptQL chunk is complete in 3 cases:
			1. When parsing error occurs (i.e. `*TInterpreterResult.Error != nil`);
			2. When all PromptQL commands are executed in current chunk. Only root context is left;
			3. When runtime/execution error occurs (i.e. when `*TInterpreterResult.Result` contains `error` data);
		```
		For nice formatting of `Result`, you can use methods `func (self *TInterpreterResult) ResultDataStr() (string, bool)` and `func (self *TInterpreterResult) ResultErrorStr() (string, bool)`.

		Notice that these methods formats a result accumulated as a text on all root input channel entries. For getting the latest clean result you can use the `func (self *TInterpreterResult) ResultLatestData(chanName string) interface{}`


		- `func (self *Interpreter) ExecutePartial(program string) *TInterpreterResult` - executes query as a part of **open** session. Only interpreter cursor is reset.

		The method receives parameters:
		```
		- "program" - is an executed PromptQL program;
		```

		- `func (self *PromptQL) Instance.Reset()` - for manually resetting all interpreter state. You can use it for very specific use-cases when standard PromptQL execution flow is not suitable; 
		- `func (self *PromptQL) Instance.IsDirty() bool` - determines if interpreter is in process of execution PromptQL session. It's `false` after execution of closed session and after calling the `Reset` method;

	- Globals API:
		- `func (self *Interpreter) Instance.SetExternalGlobals(globals TGlobalVariablesTable)` - for late setup of default external variables table;
		- `func (self *Interpreter) Instance.SetExternalGlobalVar(name string, val interface{}, description string)` - for late setup of some external variable (with optional description);
		- `func (self *Interpreter) Instance.GetExternalGlobalsList() map[string]string` - for getting list of external globals with descriptions. It's primarily used by the `{~hello /}` command, but you can use it in other scenarios on your own;

	- Sessions API:
		- `func (self *Interpreter) Instance.OpenSession()` - for opening PromptQL execution session. It's primarily used by the `{~session_begin /}` command, but you can use it in other scenarios on your own (carefully as it modifies an instance state!);
		- `func (self *Interpreter) Instance.CloseSession()` - for closing PromptQL execution session. It's primarily used by the `{~session_end /}` command, but you can use it in other scenarios on your own (carefully as it modifies an instance state!);
		- `func (self *Interpreter) Instance.IsSessionClosed() bool` - returns a flag of **current** session state;
	
	- Embeddings API:
		- `func (self *Interpreter) Instance.GetEmbeddingsList() map[string]string` - for getting list of embeddings with descriptions. It's primarily used by the `{~hello /}` command, but you can use it in other scenarios on your own;
		- `func (self *Interpreter) Instance.RegisterEmbedding(name string, code string, description string)` - for registering some PromptQL code chunk for later expansion or execution (with optional description), It's primarily used by the `{~embed_def}{/embed_def}` command, but you can use it in other scenarios on your own;
		- `func (self *Interpreter) Instance.ExpandEmbedding(name string, args TEmbeddingArgsTable) (string, error)` - for expanding PromptQL code chunk (with optional `args`). It's primarily used by the `{~embed_exp}{/embed_exp}` command, but you can use it in other scenarios on your own;


 - PromptQL.CustomApis methods:

 - `func (self *PromptQL) CustomApis.RegisterModelApi(name string, doQuery TDoQueryFunc, description string)` - defines ML model API with its own unique name and function for processing queries. And optional description if provided. The `doQuery` function is defined by this convention:

	```
	  func(
	    model string,
	    temperature float64,
	    inputs interpreter.TFunctionInputChannelTable,
	    execInfo interpreter.TExecutionInfo,
    ) (string, error)
	```

	This function should block if it contains some blocking requests to IO, DB, network etc. As it executes in separate goroutine that pushes result to query handle;

- `func (self *CustomModelsApis) GetAllModelsList() map[string]string` - returns a list of user-defined ML models with their descriptions. It's primarily used by the `{~hello /}` command, but you can use it for other scenarios on your own;


## Architecture

Interpreter has simple stack-based architecture like this:
<img src="./readme-content/prompt-ql-exec-context-architecture.png" />

<img src="./readme-content/prompt-ql-interpreter-state-diagram.png" />

Each stack entry consists of **execution context**. It defines executed command with static arguments (defined with `<arg>=<val>`) and input channels (this data is filled after execution of inner commands). A context can also be in 4 states:

 - `StackFrameStateExpectCmd` - expecting a command name for "opening" command;
 - `StackFrameStateExpectArg` - expecting a current argument name;
 - `StackFrameStateExpectVal` - expecting a current argument value;
 - `StackFrameStateIsClosing` - current top context stack frame is about to leave the stack and be executed. This is done after the command mode (defined with `{}` brackets) is switched back to the plain text mode of interpreter;
 - `StackFrameStateFullfilled` - state that's set after filling all command info (command name and static arguments). It's introduced for better distinguising "opening" and "closing" commands;
 - `StackFrameStateExpectCmdAfterFullfill` - expecting a command name for "closing" command. Ot's introduced for handling errors of mismatching command tags (ex. `{~open_query}<some_text>{/call}`);

The overall state diagram of context states is:
<img src="./readme-content/promptql-exec-context-state-diagram.png" />
