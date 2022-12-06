package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

//go:embed tmpl/main.tmpl
var mainTemple string

//go:embed tmpl/agent.tmpl
var agentTemple string

//go:embed tmpl/agent_cmd.tmpl
var agentCmdTemple string

//go:embed tmpl/agent_func.tmpl
var agentFuncTemple string

//go:embed tmpl/command.tmpl
var commandTemple string

//go:embed tmpl/condition.tmpl
var conditionTemple string

//go:embed tmpl/condition_func.tmpl
var conditionFuncTemple string

//go:embed tmpl/cron.tmpl
var cronTemple string

//go:embed tmpl/cron_func.tmpl
var cronFuncTemple string

//go:embed tmpl/form.tmpl
var formTemple string

//go:embed tmpl/form_func.tmpl
var formFuncTemple string

//go:embed tmpl/group_func.tmpl
var groupFuncTemple string

//go:embed tmpl/input_func.tmpl
var inputFuncTemple string

const BasePath = "./server/extra/bots"

func main() {
	// args
	bot := flag.String("bot", "", "bot package")
	rule := flag.String("rule", "command", "rule type") // input,group,agent,command,condition,cron,form
	flag.Parse()
	if *bot == "" {
		panic("bot args error")
	}

	// schema
	data := schema{
		BotName:      *bot,
		HasInput:     false,
		HasGroup:     false,
		HasCommand:   false,
		HasAgent:     false,
		HasCondition: false,
		HasCron:      false,
		HasForm:      false,
	}
	parseRule(*rule, &data)

	// check dir
	_, err := os.Stat(BasePath)
	if os.IsNotExist(err) {
		panic("bots NotExist")
	}
	dir := fmt.Sprintf("%s/%s", BasePath, data.BotName)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(filePath(data.BotName, "bot.go"), parseTemplate(mainTemple, data), os.ModePerm)
		if err != nil {
			panic(err)
		}
		if data.HasAgent {
			cmdDir := fmt.Sprintf("%s/%s/cmd", BasePath, data.BotName)
			_, err = os.Stat(cmdDir)
			if os.IsNotExist(err) {
				err = os.Mkdir(cmdDir, os.ModePerm)
				if err != nil {
					panic(err)
				}
			}

			err = os.WriteFile(filePath(data.BotName, "agent.go"), parseTemplate(agentTemple, data), os.ModePerm)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile(filePath(data.BotName, "cmd/main.go"), parseTemplate(agentCmdTemple, data), os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		if data.HasCommand {
			err = os.WriteFile(filePath(data.BotName, "command.go"), parseTemplate(commandTemple, data), os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		if data.HasCondition {
			err = os.WriteFile(filePath(data.BotName, "condition.go"), parseTemplate(conditionTemple, data), os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		if data.HasCron {
			err = os.WriteFile(filePath(data.BotName, "cron.go"), parseTemplate(cronTemple, data), os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		if data.HasForm {
			err = os.WriteFile(filePath(data.BotName, "form.go"), parseTemplate(formTemple, data), os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	} else {
		if !fileExist(data.BotName, "bot.go") {
			panic("dir exist, but bot.go file not exist")
		}
		if data.HasInput {
			// append
			appendFileContent(filePath(data.BotName, "bot.go"), parseTemplate(inputFuncTemple, data))
		}
		if data.HasGroup {
			// append
			appendFileContent(filePath(data.BotName, "bot.go"), parseTemplate(groupFuncTemple, data))
		}
		if !fileExist(data.BotName, "agent.go") {
			if data.HasAgent {
				cmdDir := fmt.Sprintf("%s/%s/cmd", BasePath, data.BotName)
				_, err = os.Stat(cmdDir)
				if os.IsNotExist(err) {
					err = os.Mkdir(cmdDir, os.ModePerm)
					if err != nil {
						panic(err)
					}
				}

				err = os.WriteFile(filePath(data.BotName, "agent.go"), parseTemplate(agentTemple, data), os.ModePerm)
				if err != nil {
					panic(err)
				}
				err = os.WriteFile(filePath(data.BotName, "cmd/main.go"), parseTemplate(agentCmdTemple, data), os.ModePerm)
				if err != nil {
					panic(err)
				}

				// append
				appendFileContent(filePath(data.BotName, "bot.go"), parseTemplate(agentFuncTemple, data))
			}
		}
		if !fileExist(data.BotName, "condition.go") {
			if data.HasCondition {
				err = os.WriteFile(filePath(data.BotName, "condition.go"), parseTemplate(conditionTemple, data), os.ModePerm)
				if err != nil {
					panic(err)
				}
				// append
				appendFileContent(filePath(data.BotName, "bot.go"), parseTemplate(conditionFuncTemple, data))
			}
		}
		if !fileExist(data.BotName, "cron.go") {
			if data.HasCron {
				err = os.WriteFile(filePath(data.BotName, "cron.go"), parseTemplate(cronTemple, data), os.ModePerm)
				if err != nil {
					panic(err)
				}
				// append
				appendFileContent(filePath(data.BotName, "bot.go"), parseTemplate(cronFuncTemple, data))
			}
		}
		if !fileExist(data.BotName, "form.go") {
			if data.HasForm {
				err = os.WriteFile(filePath(data.BotName, "form.go"), parseTemplate(formTemple, data), os.ModePerm)
				if err != nil {
					panic(err)
				}
				// append
				appendFileContent(filePath(data.BotName, "bot.go"), parseTemplate(formFuncTemple, data))
			}
		}
	}

	fmt.Println("ok")
}

type schema struct {
	BotName      string
	HasInput     bool
	HasGroup     bool
	HasCommand   bool
	HasAgent     bool
	HasCondition bool
	HasCron      bool
	HasForm      bool
}

func filePath(botName, fileName string) string {
	return fmt.Sprintf("%s/%s/%s", BasePath, botName, fileName)
}

func fileExist(botName, fileName string) bool {
	_, err := os.Stat(filePath(botName, fileName))
	return !os.IsNotExist(err)
}

func parseTemplate(text string, data interface{}) []byte {
	buf := bytes.NewBufferString("")
	t, err := template.New("tmpl").Parse(text)
	if err != nil {
		panic(err)
	}
	err = t.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func parseRule(rule string, data *schema) {
	rules := strings.Split(rule, ",")
	for _, item := range rules {
		switch item {
		case "input":
			data.HasInput = true
		case "group":
			data.HasGroup = true
		case "agent":
			data.HasAgent = true
		case "command":
			data.HasCommand = true
		case "condition":
			data.HasCondition = true
		case "cron":
			data.HasCron = true
		case "form":
			data.HasForm = true
		}
	}
}

func appendFileContent(filePath string, content []byte) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(content)
	if err != nil {
		panic(err)
	}

	_ = file.Close()
}
