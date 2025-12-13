package run

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/auth"
	"github.com/binaryarc/watcher/internal/grpcserver"
	"github.com/binaryarc/watcher/internal/keystore"
	"github.com/binaryarc/watcher/proto"
	"github.com/spf13/cobra"
	grpcLib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Start Watcher gRPC server",
	Long:  `Start the Watcher server to accept remote observation requests`,
	Run:   runServer,
}

var (
	port            int
	host            string
	disableAuth     bool
	keystorePathArg string
)

func init() {
	Cmd.Flags().IntVarP(&port, "port", "p", 9090, "Port to listen on")
	Cmd.Flags().StringVar(&host, "host", "0.0.0.0", "Host to bind to")
	Cmd.Flags().BoolVar(&disableAuth, "disable-auth", false, "Disable authentication (use for testing only)")
	Cmd.Flags().StringVar(&keystorePathArg, "keystore", "", "Path to keystore file (default: ~/.watcher/server/keys.json)")
}

func runServer(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf("%s:%d", host, port)

	keystorePath := keystorePathArg
	if keystorePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Failed to get home directory: %v\n", err)
			return
		}
		keysDir := filepath.Join(homeDir, ".watcher", "server")
		if err := os.MkdirAll(keysDir, 0700); err != nil {
			fmt.Printf("Failed to create keys directory: %v\n", err)
			return
		}
		keystorePath = filepath.Join(keysDir, "keys.json")
	}

	store, err := keystore.NewStore(keystorePath)
	if err != nil {
		fmt.Printf("Failed to load keystore: %v\n", err)
		return
	}

	var grpcServer *grpcLib.Server
	if disableAuth {
		fmt.Println("WARNING: Authentication DISABLED - not recommended for production")
		grpcServer = grpcLib.NewServer()
	} else {
		if store.IsEmpty() {
			fmt.Println("WARNING: No API keys registered - all requests will be rejected")
			fmt.Println("Add keys with: watcher-server key add <api-key> \"<description>\"")
		} else {
			keyCount := len(store.List())
			fmt.Printf("Authentication enabled (%d key(s) registered)\n", keyCount)
		}

		grpcServer = grpcLib.NewServer(
			grpcLib.UnaryInterceptor(auth.UnaryServerInterceptor(store)),
			grpcLib.StreamInterceptor(auth.StreamServerInterceptor(store)),
		)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v\n", addr, err)
		return
	}

	watcherServer := grpcserver.NewWatcherServer()
	proto.RegisterWatcherServiceServer(grpcServer, watcherServer)
	reflection.Register(grpcServer)

	fmt.Printf("Watcher server listening on %s...\n", addr)
	fmt.Println("Press Ctrl+C to stop")

	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}
}
