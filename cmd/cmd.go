package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kirychukyurii/wasker-directory/cmd/run"
)

func init() {
	//Command.AddCommand(serve.Command)
	Command.AddCommand(run.Command)
	// Command.AddCommand(migrate.Command)
}

var (
	version    = "0.0.0"
	commit     = "hash"
	commitDate = "date"
)

var Command = &cobra.Command{
	Use:          "wasker",
	Short:        "wasker - a Task Tracker System for Support team",
	SilenceUsage: true,
	Long:         "",
	Version:      fmt.Sprintf("%s, commit %s, date %s", version, commit, commitDate),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(
				"requires at least one arg, " +
					"you can view the available parameters through `--help`",
			)
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run:               func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := Command.Execute(); err != nil {
		os.Exit(-1)
	}
}
