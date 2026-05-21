package cli

import (
	"fmt"
)

// Runner 执行具体子命令，便于用公共接口测试 CLI 行为。
type Runner func(command string) error

type App struct {
	run Runner
}

func NewApp(run Runner) *App {
	return &App{run: run}
}

func (app *App) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("unsupported arguments")
	}

	command := args[0]
	if command != "sync" && command != "dry-run" {
		return fmt.Errorf("unknown command: %s", command)
	}

	return app.run(command)
}
