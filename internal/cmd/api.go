package cmd

import (
	"context"

	"tribe-payments-wallet-golang-interview-assignment/internal/api"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
	"tribe-payments-wallet-golang-interview-assignment/internal/http"
	"tribe-payments-wallet-golang-interview-assignment/internal/wallet"

	"database/sql"

	"github.com/go-chi/chi/v5"
	_ "github.com/microsoft/go-mssqldb" // Example for MSSQL Server, adjust if using another driver.
	"github.com/spf13/cobra"
	"github.com/sumup-oss/go-pkgs/errors"
	"github.com/sumup-oss/go-pkgs/logger"
	"github.com/sumup-oss/go-pkgs/os"
	"github.com/sumup-oss/go-pkgs/task"
	"moul.io/chizap"
)

//nolint:gocognit
func NewApiCmd(osExecutor os.OsExecutor) *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Run application server",
		Long:  "Run application server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			cfg, err := config.NewServerConfig()
			if err != nil {
				return errors.Wrap(err, "failed to create runtime config")
			}

			log, err := logger.NewZapLogger(
				logger.Configuration{
					Level:         cfg.Log.Level,
					Encoding:      logger.EncodingJSON,
					StdoutEnabled: cfg.Log.StdoutEnabled,
				},
			)
			if err != nil {
				return errors.Wrap(err, "failed to create logger")
			}

			defer log.Sync() //nolint:errcheck

			shutdownTask := newShutdownTask(
				log,
				osExecutor,
				cfg.GracefulShutdownTimeout,
			)

			db, err := sql.Open("sqlserver", "server=localhost\\SQLEXPRESS;database=wallet_db;trusted_connection=yes;")
			if err != nil {
				return errors.Wrap(err, "failed to connect to database")
			}
			defer db.Close()

			// Initialise Wallet Repository
			walletRepo := wallet.NewRepository(db)

			// Initialise Wallet Service with the Repository
			walletService := wallet.NewService(walletRepo)

			mux := chi.NewRouter()
			mux.Use(
				http.Recovery(
					log,
					api.WritePanicResponse(log),
				),
				chizap.New(
					log.Logger,
					&chizap.Opts{
						WithReferer:   true,
						WithUserAgent: true,
					},
				),
			)

			// Pass walletService to RegisterRoutes
			api.RegisterRoutes(mux, log, walletService)

			httpServer := http.NewServer(
				log,
				cfg.ListenAddress,
				mux,
				http.WithMaxHeaderBytes(cfg.MaxHeaderBytes),
				http.WithReadHeaderTimeout(cfg.ReadHeaderTimeout),
				http.WithReadTimeout(cfg.ReadTimeout),
				http.WithServerShutdownTimeout(cfg.GracefulShutdownTimeout),
				http.WithWriteTimeout(cfg.WriteTimeout),
			)

			taskGroup := task.NewGroup()
			taskGroup.Go(
				shutdownTask.Run,
				httpServer.Run,
			)

			err = taskGroup.Wait(ctx)
			if err != nil {
				return errors.Wrap(err, "task group exits with error")
			}

			log.Info("taskGroup successfully shutdown")

			return nil
		},
	}
}
