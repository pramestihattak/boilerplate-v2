package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"boilerplate-v2/service"
	"boilerplate-v2/service/auth"
	"boilerplate-v2/service/feed"
	"boilerplate-v2/storage/postgres"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	port   = flag.Int("port", 50051, "The server port")
	logger *logrus.Logger
	config *viper.Viper
)

func init() {
	config = viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("backend/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	if config.GetString("google.clientID") == "" || config.GetString("google.clientSecret") == "" || config.GetString("google.callbackURL") == "" {
		logger.Fatal("error google credential missing")
	}

	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
}

func main() {
	flag.Parse()

	storage, err := postgres.NewStorage(logger, config)
	if err != nil {
		logger.Fatal("error initializing postgres storage", err.Error())
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)

	authService := auth.NewAuthService(logger, storage)
	feedService := feed.NewFeedService(logger, storage)

	service.RegisterServices(s,
		auth.RegisterService(authService),
		feed.RegisterService(feedService),
	)

	// start gRPC server
	go func() {
		log.Printf("server listening at %v", lis.Addr())
		log.Fatalln(s.Serve(lis))
		s.GracefulStop()
	}()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := service.RegisterGateways(gwmux, lis.Addr().String(), opts,
		auth.RegisterGateway(),
		feed.RegisterGateway(),
	); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
