package cli

import (
	"github.com/exapsy/capcli/internal/cli/commands"
	"github.com/spf13/cobra"
)

type App struct {
	rootCommand *cobra.Command
}

func NewApp() *App {
	rootCommand := commands.RootCommand

	app := &App{
		rootCommand: rootCommand,
	}

	return app
}

func (a *App) Run() error {
	return a.rootCommand.Execute()
}
