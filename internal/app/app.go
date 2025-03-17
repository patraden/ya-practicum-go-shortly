package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/handler"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/logger"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/memento"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/repository"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/server"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/urlgenerator"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/utils/postgres"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/version"
)

// App returns main app function as fx.App.
func App(appCfg *config.Config, logLevel zerolog.Level) *fx.App {
	appLogger := logger.NewLogger(logLevel)

	return fx.New(
		fx.StartTimeout(time.Minute),
		fx.StopTimeout(time.Minute),
		fx.Supply(appLogger),
		fx.Supply(appCfg),
		fx.Provide(func(l *logger.Logger) *zerolog.Logger { return l.GetLogger() }),
		fx.Provide(postgres.New),
		fx.Provide(fx.Annotate(urlgenerator.New, fx.As(new(urlgenerator.URLGenerator)))),
		fx.Provide(
			func(db *postgres.Database, l *zerolog.Logger, c *config.Config) (repository.URLRepository, error) {
				if c.DatabaseDSN != `` {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := db.Init(ctx); err != nil {
						return nil, err
					}

					return repository.NewDBURLRepository(db.ConnPool, l), nil
				}

				return repository.NewInMemoryURLRepository(), nil
			}),
		fx.Provide(
			shortener.NewInsistentShortener,
			remover.NewBatchRemover,
			func(r *remover.BatchRemover) remover.URLRemover { return r },
		),
		fx.Provide(
			fx.Annotate(handler.NewPingHandler, fx.As(new(handler.Handler)), fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(handler.NewDeleteHandler, fx.As(new(handler.Handler)), fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(handler.InsistentShortenerHandler, fx.As(new(handler.Handler)), fx.ResultTags(`group:"handlers"`)),
			fx.Annotate(handler.NewRouter, fx.ParamTags(`group:"handlers"`)),
		),
		fx.Provide(
			func(r repository.URLRepository) memento.Originator { return r },
			memento.NewStateManager,
		),
		fx.Provide(server.NewServer),
		fx.Invoke(fxAppInvoke),
		fx.WithLogger(appLogger.GetFxLogger()),
	)
}

func fxAppInvoke(
	lc fx.Lifecycle,
	log *zerolog.Logger,
	config *config.Config,
	remover *remover.BatchRemover,
	stateManager *memento.StateManager,
	server *server.Server,
	shutdowner fx.Shutdowner,
) {
	ctxRemover, removerCancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			// load state from disc
			if !config.ForceEmptyRepo {
				err := stateManager.RestoreFromFile()
				if err != nil && !errors.Is(err, e.ErrStateNotmplemented) {
					log.Error().Err(err).Msg("State restoration error")
				}
			}

			appHandleSignals(shutdowner, log)
			appServerStart(shutdowner, server, log)
			remover.Start(ctxRemover)
			version := version.NewVersion(log)
			version.Log()
			logStart(log, config)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			removerCancel()
			remover.Stop(ctx)

			err := server.Shutdown(ctx)
			if err != nil {
				return err
			}

			// preserve state to disc
			if !config.ForceEmptyRepo {
				err := stateManager.StoreToFile()
				if err != nil && !errors.Is(err, e.ErrStateNotmplemented) {
					log.Error().Err(err).Msg("State preservation error")
				}
			}

			logStop(log, config)

			return nil
		},
	})
}

func appHandleSignals(shutdowner fx.Shutdowner, log *zerolog.Logger) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-stopChan
		log.Info().
			Str("Signal", sig.String()).
			Msg("Shutdown signal received")

		err := shutdowner.Shutdown()
		if err != nil {
			log.Error().Err(err).
				Str("Signal", sig.String()).
				Msg("Failed to shutdown")
		}
	}()
}

func appServerStart(shutdowner fx.Shutdowner, server *server.Server, log *zerolog.Logger) {
	go func() {
		err := server.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).
				Msg("Stopping app due to http server error")

			err := shutdowner.Shutdown()
			if err != nil {
				log.Error().Err(err).
					Msg("Failed to shutdown")
			}
		}
	}()
}

func logStart(log *zerolog.Logger, config *config.Config) {
	log.Info().
		Str("SERVER_ADDRESS", config.ServerAddr).
		Str("BASE_URL", config.BaseURL).
		Bool("ENABLE_HTTPS", config.EnableHTTPS).
		Bool("FORCE_EMPTY", config.ForceEmptyRepo).
		Str("FILE_STORAGE_PATH", config.FileStoragePath).
		Msg("App started")
}

func logStop(log *zerolog.Logger, config *config.Config) {
	log.Info().
		Str("SERVER_ADDRESS", config.ServerAddr).
		Str("BASE_URL", config.BaseURL).
		Bool("ENABLE_HTTPS", config.EnableHTTPS).
		Bool("FORCE_EMPTY", config.ForceEmptyRepo).
		Str("FILE_STORAGE_PATH", config.FileStoragePath).
		Msg("App stopped")
}
