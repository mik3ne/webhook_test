package app

import (
	"webhook/internal/commands"
	"webhook/internal/config"
	"webhook/internal/servers/webhook"
	"webhook/internal/services"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func CreateApp(cfgLoader func() (config.Configuration, error)) *fx.App {
	return fx.New(
		fx.Provide(
			webhook.NewServer,
			zap.NewDevelopment,
			services.NewRequestSender,
			services.NewTaskGenerator,
			services.NewRateLimiter,
			cfgLoader,
			fx.Annotate(
				commands.NewSendCmd,
				fx.ParamTags(`group:"commands"`),
			),
		),
		fx.Invoke(func(*cobra.Command) {}),
	)
}
