package run

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/kirychukyurii/wasker-directory/internal/config"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
)

func init() {
	pf := Command.PersistentFlags()

	pf.StringVarP(&config.Path, "config", "c", "config/config.yml",
		"this parameter is used to start the service application")
}

const ServiceName = "wasker-directory"

var Command = &cobra.Command{
	Use:          "run",
	Short:        "",
	Example:      "",
	SilenceUsage: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(config.Path)
		if err := viper.ReadInConfig(); err != nil {
			panic(errors.Wrap(err, "read config"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		app := fx.New(
			Module,
			fx.WithLogger(
				func(log logger.Logger) fxevent.Logger {
					return &logger.FxLogger{
						Logger: &log.Logger,
					}
				},
			),
		)

		app.Run()
	},
}
