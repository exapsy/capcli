package commands

import (
	"github.com/spf13/cobra"
)

const (
	CliName             = "capcli"
	CliDescriptionShort = "capcli is a CLI for interacting with the Capital API"
	CliDescriptionLong  = "capcli is a CLI for interacting with the Capital API"
)

var (
	RootCommand = &cobra.Command{
		Use:   CliName,
		Short: CliDescriptionShort,
		Long:  CliDescriptionLong,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				panic(err)
				return
			}
		},
	}
)

func init() {
	RootCommand.AddCommand(
		inspectCommand,
	)
}
