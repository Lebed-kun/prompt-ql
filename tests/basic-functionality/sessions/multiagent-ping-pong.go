package sessionstests

import (
	"fmt"
	"sync"
	"time"

	interpretercore "gitlab.com/jbyte777/prompt-ql/v2/core"
	interpreter "gitlab.com/jbyte777/prompt-ql/v2/interpreter"
)

func setupFirstAgent() *interpreter.TPromptQL {
	defaultGlobals := interpretercore.TGlobalVariablesTable{
		"myRef111": "@myVar111",
		"myVar111": "Hello, PromptQL!",
		"myFunc111": func(args []interface{}) interface{} {
			return nil
		},
	}
	agent := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
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
	agent := interpreter.New(
		interpreter.TPromptQLOptions{
			DefaultExternalGlobals: defaultGlobals,
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
	)

	return agent
}

// 07-10-2023: Works +++
func MultiagentSessionPingPongTest() {
	agent1 := setupFirstAgent()
	agent2 := setupSecondAgent()
	agent1CmdBox := make(chan string)
	agent2CmdBox := make(chan string)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	waitGroup.Add(1)

	// Alice agent goroutine
	go func() {
		sessCount := 0

		if agent1.Instance.IsSessionClosed() {
			agent1.Instance.Execute(
				`{~session_begin /}`,
			)
		}
		
		agent2CmdBox <- `
			{~session_begin /}
			From Alice
		`

		for cmd := range agent1CmdBox {
			time.Sleep(time.Second)
			result := agent1.Instance.Execute(cmd)
			resultStr, _ := result.ResultDataStr()

			fmt.Printf(
				`
				  ==== Alice =====
				  Alice: %v
				`,
				resultStr,
			)

			if sessCount >= 5 {
				agent2CmdBox <- fmt.Sprintf(
					`
						{~get from="myVar" /}
						{~set to="myVar"}%v{/set}
						{~session_end /}
					`,
					fmt.Sprintf(
						"Bob_%v",
						sessCount,
					),
				)

				if !agent1.Instance.IsSessionClosed() {
					agent1.Instance.Execute(
						`{~session_end /}`,
					)
				}

				break
			} else {
				agent2CmdBox <- fmt.Sprintf(
					`
						{~get from="myVar" /}
						{~set to="myVar"}%v{/set}
					`,
					fmt.Sprintf(
						"Bob_%v",
						sessCount,
					),
				)
			}
			sessCount++
		}

		waitGroup.Done()
	}()

	// Bob agent goroutine
	go func() {
		sessCount := 0

		for cmd := range agent2CmdBox {
			time.Sleep(time.Second)
			result := agent2.Instance.Execute(cmd)
			resultStr, _ := result.ResultDataStr()

			fmt.Printf(
				`
				  ==== Bob =====
				  Bob: %v
				`,
				resultStr,
			)

			if agent2.Instance.IsSessionClosed() {
				break
			}

			agent1CmdBox <- fmt.Sprintf(
				`
					{~get from="myVar" /}
					{~set to="myVar"}%v{/set}
				`,
				fmt.Sprintf(
					"Alice_%v",
					sessCount,
				),
			)
			sessCount++
		}

		waitGroup.Done()
	}()

	waitGroup.Wait()
}
