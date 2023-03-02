package cmd

import (
	"context"
	"fmt"
	service "grpctest/gen/service/go"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var serverOptions = struct {
	Port        int
	GatewayPort int
}{}

type server struct {
	service.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *service.HelloRequest) (*service.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("md: %+v, ok: %v", md, ok)

	log.Printf("Received: %v", in.GetName())
	return &service.HelloReply{Message: "Hello " + in.GetName()}, nil
}

var serverCmd = cobra.Command{
	Use:   "server",
	Short: "grpctest server",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetFlags(log.Flags() | log.Lshortfile)
		exitChan := make(chan error, 2)
		defer close(exitChan)

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverOptions.Port))
		if err != nil {
			return err
		}
		// lis = debugnet.NewDebugListener(lis)

		s := grpc.NewServer()

		reflection.Register(s)
		service.RegisterGreeterServer(s, &server{})

		go func() {
			if err := s.Serve(lis); err != nil {
				exitChan <- err
			} else {
				exitChan <- nil
			}
		}()

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		err = service.RegisterGreeterHandlerFromEndpoint(ctx, mux, fmt.Sprintf("127.0.0.1:%d", serverOptions.Port), opts)
		if err != nil {
			return err
		}

		err = http.ListenAndServe(fmt.Sprintf(":%d", serverOptions.GatewayPort), mux)
		if err != nil {
			return err
		}

		s.GracefulStop()

		<-exitChan
		return nil
	},
}

func init() {
	serverCmd.PersistentFlags().IntVarP(&serverOptions.Port, "port", "p", 0, "listen port for grpc server")
	serverCmd.PersistentFlags().IntVarP(&serverOptions.GatewayPort, "gateway-port", "", 0, "listen port for grpc gateway server")
	rootCmd.AddCommand(&serverCmd)
}
