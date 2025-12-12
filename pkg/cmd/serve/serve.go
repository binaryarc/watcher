package serve

import (
	"fmt"
	"net"

	grpc "github.com/binaryarc/watcher/internal/grpcserver"
	"github.com/binaryarc/watcher/proto"
	"github.com/spf13/cobra"
	grpcLib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Watcher gRPC server",
	Long:  `Start the Watcher server to accept remote observation requests`,
	Run:   runServe,
}

var (
	port int
	host string
)

func init() {
	ServeCmd.Flags().IntVarP(&port, "port", "p", 9090, "Port to listen on")
	ServeCmd.Flags().StringVar(&host, "host", "0.0.0.0", "Host to bind to")
}

func runServe(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf("%s:%d", host, port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("❌ Failed to listen on %s: %v\n", addr, err)
		return
	}

	grpcServer := grpcLib.NewServer()
	watcherServer := grpc.NewWatcherServer()

	proto.RegisterWatcherServiceServer(grpcServer, watcherServer)
	reflection.Register(grpcServer)
	fmt.Printf("👁️  Watcher server listening on %s...\n", addr)
	fmt.Println("Press Ctrl+C to stop")

	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("❌ Failed to serve: %v\n", err)
	}
}
