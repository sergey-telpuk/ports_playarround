package app

import (
	api "github.com/client-api/internal/api/grpc/port"
	"github.com/client-api/internal/file/reader"
	"github.com/client-api/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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
		Use:   "portreader",
		Short: "Run read port from json file",
		Run:   runPortReader,
	}

	rootCmd.AddCommand(cmds)

	return rootCmd
}

func runPortReader(cmd *cobra.Command, args []string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	rootCmd.Flags().String("grpcport", "", "The GRPC server port")
	_ = viper.BindPFlag("grpcport", rootCmd.Flags().Lookup("grpcport"))
	_ = viper.BindEnv("grpcport", "GRPC_PORT")

	rootCmd.Flags().String("pathjson", "", "The path of json")
	_ = viper.BindPFlag("pathjson", rootCmd.Flags().Lookup("pathjson"))
	_ = viper.BindEnv("pathjson", "PATH_JSON")

	grpcport := viper.GetString("grpcport")
	pathjson := viper.GetString("pathjson")

	file, err := os.Open(pathjson) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	conn, err := grpc.Dial(grpcport, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}
	defer conn.Close()

	srv := services.NewPortService(
		reader.NewReader(file),
		api.NewClient(conn),
	)

	go func() {
		<-c
		srv.Close()
		os.Exit(1)
	}()

	if err := srv.ReadFromJsonWorkerPool(5); err != nil {
		log.Fatal(err)
	}

	port, err := srv.GetPortByName("Ajman")
	log.Println(port)

	if err != nil {
		log.Fatal(err)
	}
}
