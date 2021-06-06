package app

import (
	api "github.com/port-domain/internal/api/port"
	domain "github.com/port-domain/internal/domain/port"
	"github.com/port-domain/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
)

func init() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})
}

var (
	// Root command
	rootCmd = &cobra.Command{
		TraverseChildren: true,
	}
)

func New() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "grpcserver",
		Short: "Run grpc server",
		Run:   runGRPCServer,
	}

	rootCmd.AddCommand(cmds)

	return rootCmd
}

func runGRPCServer(cmd *cobra.Command, args []string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	rootCmd.Flags().String("grpcport", "", "The GRPC server port")
	_ = viper.BindPFlag("grpcport", rootCmd.Flags().Lookup("grpcport"))
	_ = viper.BindEnv("grpcport", "GRPC_PORT")

	grpcport := viper.GetString("grpcport")

	apiServer := api.NewApiService(domain.NewRepository())

	srv := service.NewPortDomainService(apiServer)

	go func() {
		<-c
		srv.Close()
		os.Exit(1)
	}()

	if err := srv.Run(grpcport); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
