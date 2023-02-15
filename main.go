package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/MathPeixoto/go-financial-system/api"
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	_ "github.com/MathPeixoto/go-financial-system/doc/statik"
	"github.com/MathPeixoto/go-financial-system/gapi"
	"github.com/MathPeixoto/go-financial-system/pb"
	"github.com/MathPeixoto/go-financial-system/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration from file "app.env"
	config, err := util.LoadConfig("app.env")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	// If the environment is set to "development", set the logger to console output
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Connect to the database using the specified database driver and source
	var conn *sql.DB
	conn, err = sql.Open(config.DatabaseDriver, config.DatabaseSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	// run db migrations
	runDBMigration(config.MigrationURL, config.DatabaseSource)

	// Create a new store using the database connection
	store := db.NewStore(conn)

	// Start the gateway server in a new goroutine
	go runGatewayServer(config, store)
	// Start the gRPC server
	runGrpcServer(config, store)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration")
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run migration")
	}

	log.Info().Msg("migration completed")
}

// runGrpcServer starts a gRPC server and listens for incoming requests
func runGrpcServer(config util.Config, store db.Store) {
	// Create a new gapi server
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	// Register the bank server to the gRPC server
	pb.RegisterBankServer(grpcServer, server)
	// Register the gRPC server to use reflection
	reflection.Register(grpcServer)
	// Listen for incoming requests at the specified address
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Printf("Starting gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

// runGatewayServer starts the HTTP gateway server for the bank service with the given configuration and database store.
func runGatewayServer(config util.Config, store db.Store) {
	// Create a new server using the provided configuration and store.
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	// Define JSON options for the gRPC-JSON transcoder.
	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			// Use protobuf field names in the JSON output.
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			// Discard unknown fields in the JSON input.
			DiscardUnknown: true,
		},
	})

	// Create a gRPC-JSON transcoder serve mux.
	grpcMux := runtime.NewServeMux(jsonOptions)

	// Create a context with cancel function.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register the handler server to the gRPC-JSON transcoder serve mux.
	err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	// Create a new HTTP serve mux.
	mux := http.NewServeMux()
	// Mount the gRPC-JSON transcoder serve mux to the root path.
	mux.Handle("/", grpcMux)

	// Create a new file system using Statik.
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	// Create a handler for serving Swagger documentation from the Statik file system.
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// Mount the Swagger documentation handler to the /swagger/ path.
	mux.Handle("/swagger/", swaggerHandler)

	// Create a listener on the HTTP server address.
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	// Start serving HTTP requests using the created listener and HTTP serve mux.
	log.Printf("Starting HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}
}

// runGinServer starts the Gin HTTP server with the provided config and store.
func runGinServer(config util.Config, store db.Store) {
	// Creates a new server instance with the provided config and store.
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	// Starts the server and listens on the address specified in the config.
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
