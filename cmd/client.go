package cmd

import (
	"context"
	"log"
	"time"

	service "grpctest/gen/service/go"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientOptions = struct {
	Addr string
	Name string
}{}

var clientCmd = cobra.Command{
	Use:   "client",
	Short: "grpctest client",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(clientOptions.Addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			// grpc.WithHeader
		)
		if err != nil {
			return err
		}
		defer conn.Close()

		c := service.NewGreeterClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		r, err := c.SayHello(ctx, &service.HelloRequest{Name: clientOptions.Name})

		if err != nil {
			return err
		}
		log.Println(r.GetMessage())
		return nil
	},
}

func init() {
	clientCmd.PersistentFlags().StringVarP(&clientOptions.Addr, "addr", "a", "", "address of grpc server")
	clientCmd.PersistentFlags().StringVarP(&clientOptions.Name, "name", "n", "", "name of client")

	rootCmd.AddCommand(&clientCmd)
}
