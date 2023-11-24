package commands

import (
	"context"
	"webhook/internal/config"
	"webhook/internal/servers/webhook"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CobraCommand interface {
	Command() *cobra.Command
	FullName() string
}

func NewSendCmd(commands []CobraCommand, lc fx.Lifecycle, shutdowner fx.Shutdowner, sendServer *webhook.Server, log *zap.Logger, config config.Configuration) *cobra.Command {

	//
	// we got only one default command - "send"
	//
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send requests to target url",
		Long:  "Send N requests to target url with selected rps rate limit",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			go sendServer.Start()
		},
	}

	log.Debug("config loaded", zap.Any("config", config))

	//
	// using default values from config if there are no flags
	//
	sendCmd.Flags().StringVarP(&sendServer.Settings.TargetURL, "url", "u", config.URL, "Webhook URL")
	sendCmd.Flags().IntVarP(&sendServer.Settings.RequestAmount, "amount", "a", config.Requests.Amount, "Requests amount")
	sendCmd.Flags().IntVarP(&sendServer.Settings.WorkersNumber, "workers", "w", 1, "Workers number")
	sendCmd.Flags().IntVarP(&sendServer.Settings.RPS, "rps", "r", config.Requests.PerSecond, "Requests per second limit")

	lc.Append(fx.Hook{

		OnStart: func(ctx context.Context) error {

			err := sendCmd.Execute()
			if err != nil {
				log.Error("send command exec", zap.Error(err))
				err := shutdowner.Shutdown(fx.ExitCode(1))
				if err != nil {
					log.Error("shutdowner shutdown", zap.Error(err))
				}
			}

			//
			// we can`t use shutdowner here because it will stop server until it finishes its work
			// so, the only way is to stop App manually
			//
			// Using cobra-commands isn`t best solution for this case
			//

			return nil

		},

		OnStop: func(ctx context.Context) error {
			err := sendServer.Stop()
			if err != nil {
				log.Error("server stop by OnStop", zap.Error(err))
			}
			return nil
		},
	})

	return sendCmd
}
