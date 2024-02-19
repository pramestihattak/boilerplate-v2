package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"boilerplate-v2/service"
	"boilerplate-v2/service/auth"
	"boilerplate-v2/storage/postgres"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	port        = 50051
	portGateway = 8090
	logger      *logrus.Logger
	config      *viper.Viper
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
	storage, err := postgres.NewStorage(logger, config)
	if err != nil {
		logger.Fatal("error initializing postgres storage", err.Error())
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)

	service.RegisterServices(s,
		auth.RegisterService(auth.NewService(logger, storage)),
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
	); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", portGateway),
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http://0.0.0.0:%d", portGateway)
	log.Fatalln(gwServer.ListenAndServe())
}
