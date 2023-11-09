package hellocommandtests

import (
	"fmt"
	"sync"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v4/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v4/interpreter"
)

func setupFirstAgent() *interpreter.TPromptQL {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myRef111": "@myVar111",
		"myVar111": "Hello, PromptQL!",
		"myFunc111": func(args []interface{}) interface{} {
			return nil
		},
	}
	defaultGlobalsMeta := interpretercore.TExternalGlobalsMetaTable{
		"myRef111": &interpretercore.TExternalGlobalMetaInfo{
			Description: "an Alice's reference to @myVar",
		},
		"myVar111": &interpretercore.TExternalGlobalMetaInfo{
			Description: "an Alice's example external variable",
		},
		"myFunc111": &interpretercore.TExternalGlobalMetaInfo{
			Description: "an Alice's example external function",
		},
	}
	agent := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
			DefaultExternalGlobalsMeta: defaultGlobalsMeta,
		},
	)
	agent.CustomApis.RegisterModelApi(
		"myMlModel111",
		func(
			model string,
			temperature float64,
			inputs interpretercore.TFunctionInputChannelTable,
			execInfo interpretercore.TExecutionInfo,
		) (string, error) {
			return "", nil
		},
		"my model 111",
	)

	return agent
}

func setupSecondAgent() *interpreter.TPromptQL {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myVar222": "Hello, PromptQL!",
		"myFunc222": func(args []interface{}) interface{} {
			return nil
		},
	}
	defaultGlobalsMeta := interpretercore.TExternalGlobalsMetaTable{
		"myVar222": &interpretercore.TExternalGlobalMetaInfo{
			Description: "a Bob's example external variable",
		},
		"myFunc222": &interpretercore.TExternalGlobalMetaInfo{
			Description: "a Bob's example external function",
		},
	}
	agent := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
			DefaultExternalGlobalsMeta: defaultGlobalsMeta,
		},
	)
	agent.CustomApis.RegisterModelApi(
		"myMlModel222",
		func(
			model string,
			temperature float64,
			inputs interpretercore.TFunctionInputChannelTable,
			execInfo interpretercore.TExecutionInfo,
		) (string, error) {
			return "", nil
		},
		"my model 222",
	)

	return agent
}

// 07-10-2023: Works +++
// 04-11-2023: Works with descriptions +++
func MultiagentPingPongTest() {
	agent1 := setupFirstAgent()
	agent2 := setupSecondAgent()
	agent1CmdBox := make(chan string)
	agent2CmdBox := make(chan string)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	waitGroup.Add(1)

	// Alice agent goroutine
	go func() {
		agent2CmdBox <- `
		  Bob's agent layout is:
			{~hello /}
		`
		cmdForAgent1 := <-agent1CmdBox
		result := agent1.Instance.Execute(cmdForAgent1)

		if result.Error != nil {
			waitGroup.Done()
			panic(result.Result)
		}

		fmt.Println("===================")
		resultStr, _ := result.ResultDataStr()
		errStr, _ := result.ResultErrorStr()
		fmt.Printf(
			"Alice response:\n%v\n",
			resultStr,
		)
		fmt.Printf(
			"Alice error:\n%v\n",
			errStr,
		)
		fmt.Println("===================")

		waitGroup.Done()
	}()

	// Bob agent goroutine
	go func() {
		cmdForAgent2 := <-agent2CmdBox
		agent1CmdBox <- `
			Alice's agent layout is:
			{~hello /}
		`
		result := agent2.Instance.Execute(cmdForAgent2)

		if result.Error != nil {
			waitGroup.Done()
			panic(result.Result)
		}

		fmt.Println("===================")
		resultStr, _ := result.ResultDataStr()
		errStr, _ := result.ResultErrorStr()
		fmt.Printf(
			"Bob response:\n%v\n",
			resultStr,
		)
		fmt.Printf(
			"Bob error:\n%v\n",
			errStr,
		)
		fmt.Println("===================")

		waitGroup.Done()
	}()

	waitGroup.Wait()
}
