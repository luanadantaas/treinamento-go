package main

import (
	"challenge/internal/cache"
	"challenge/internal/database"
	gorm "challenge/internal/database/gorm"
	pb "challenge/internal/grpc"
	"challenge/internal/grpcRouter"
	"challenge/internal/logger"
	"challenge/internal/repository"
	"context"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const (
	LocalPort = "localhost:8081"
	Port = ":8081"
	HttpP = ":8082"
	DbType = "DATABASE_TYPE"
	Gorm = "gorm"
)


func main() {

	dBIMP, ok := os.LookupEnv(DbType)   
	var err error
	var repo repository.Repository
	if ok && dBIMP == Gorm {
		repo, err = gorm.New()
		if err != nil {
			logger.Log().Fatal(err)
		}

	} else {
		repo, err = database.New()
		if err != nil {
			logger.Log().Fatal(err)
		}
	}

	cache, err := cache.NewClient()
	if err != nil{
		logger.Log().Fatal(err)
	}

	lis, err := net.Listen("tcp", Port)
	if err != nil {
		logger.Log().Fatalf("fail to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, grpcRouter.NewService(repo, cache))

	go func() {
		err = s.Serve(lis)
		if err != nil {
			logger.Log().Fatalf("error running grpc server: %v", err)
		}
	}()

	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux()// setting up a dail up for gRPC service by specifying endpoint/target url
	err = pb.RegisterTaskServiceHandlerFromEndpoint(context.Background(), mux, LocalPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Log().Fatal(err)
	}
	
	// Creating a normal HTTP server
	server := http.Server{Handler: mux,}// creating a listener for server
	l, err := net.Listen("tcp", HttpP)
	if err!=nil {
		logger.Log().Fatal(err)
	}
	// start server
	err = server.Serve(l)
	if err != nil {
		logger.Log().Fatal(err)
	}
}
