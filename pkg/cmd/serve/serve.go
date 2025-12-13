package serve

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/binaryarc/watcher/internal/auth"
	grpc "github.com/binaryarc/watcher/internal/grpcserver"
	"github.com/binaryarc/watcher/internal/keystore"
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
	port            int
	host            string
	disableAuth     bool
	keystorePathArg string
)

func init() {
	ServeCmd.Flags().IntVarP(&port, "port", "p", 9090, "Port to listen on")
	ServeCmd.Flags().StringVar(&host, "host", "0.0.0.0", "Host to bind to")
	ServeCmd.Flags().BoolVar(&disableAuth, "disable-auth", false, "Disable authentication (use for testing only)")
	ServeCmd.Flags().StringVar(&keystorePathArg, "keystore", "", "Path to keystore file (default: ~/.watcher/server/keys.json)")
}

func runServe(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf("%s:%d", host, port)

	// Load keystore
	keystorePath := keystorePathArg
	if keystorePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("❌ Failed to get home directory: %v\n", err)
			return
		}
		keysDir := filepath.Join(homeDir, ".watcher", "server")
		if err := os.MkdirAll(keysDir, 0700); err != nil {
			fmt.Printf("❌ Failed to create keys directory: %v\n", err)
			return
		}
		keystorePath = filepath.Join(keysDir, "keys.json")
	}

	store, err := keystore.NewStore(keystorePath)
	if err != nil {
		fmt.Printf("❌ Failed to load keystore: %v\n", err)
		return
	}

	// Create gRPC server with authentication
	var grpcServer *grpcLib.Server
	if disableAuth {
		fmt.Println("⚠️  Authentication DISABLED - not recommended for production")
		grpcServer = grpcLib.NewServer()
	} else {
		if store.IsEmpty() {
			fmt.Println("⚠️  No API keys registered - authentication is effectively disabled")
			fmt.Println("   Add keys with: watcher-server key add <api-key> \"<description>\"")
		} else {
			keyCount := len(store.List())
			fmt.Printf("✅ Authentication enabled (%d key(s) registered)\n", keyCount)
		}

		grpcServer = grpcLib.NewServer(
			grpcLib.UnaryInterceptor(auth.UnaryServerInterceptor(store)),
			grpcLib.StreamInterceptor(auth.StreamServerInterceptor(store)),
		)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("❌ Failed to listen on %s: %v\n", addr, err)
		return
	}

	watcherServer := grpc.NewWatcherServer()
	proto.RegisterWatcherServiceServer(grpcServer, watcherServer)
	reflection.Register(grpcServer)

	fmt.Printf("👁️  Watcher server listening on %s...\n", addr)
	fmt.Println("Press Ctrl+C to stop")

	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("❌ Failed to serve: %v\n", err)
	}
}
