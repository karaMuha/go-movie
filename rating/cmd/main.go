package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/karaMuha/go-movie/pb"
	"github.com/karaMuha/go-movie/pkg/discovery"
	consul "github.com/karaMuha/go-movie/pkg/discovery/consul"
	"github.com/karaMuha/go-movie/rating/internal/core"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	grpchandler "github.com/karaMuha/go-movie/rating/internal/endpoint/grpc"
	"github.com/karaMuha/go-movie/rating/internal/endpoint/rest/v1"
	"github.com/karaMuha/go-movie/rating/internal/queue/consumer"
	"github.com/karaMuha/go-movie/rating/internal/repository/memory"
	"google.golang.org/grpc"
)

const serviceName = "rating"
const domain = "localhost"

func main() {
	log.Println("Starting movie rating service")
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the rating service on port %d", port)

	registry, err := consul.NewConsulRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	ratingRepo := memory.New()
	app := core.New(&ratingRepo)

	consumer := consumer.NewQueueConsumer(&app, "localhost:9092", "ratings", "rating")
	defer consumer.Reader.Close()
	go consumer.StartReading()

	// startRest(&app, port)
	startGrpc(&app, port)

}

func setupRestEndpoints(mux *http.ServeMux, ratingHandlerV1 rest.RatingHandlerV1) {
	ratingV1 := http.NewServeMux()
	ratingV1.HandleFunc("GET /get-rating", ratingHandlerV1.HandleGetRating)
	ratingV1.HandleFunc("POST /submit-rating", ratingHandlerV1.HandleSubmitRating)

	mux.Handle("/v1/", http.StripPrefix("/v1", ratingV1))
}

func startRest(app driving.IApplication, port int) {
	ratingHandlerV1 := rest.NewRatingHandlerV1(app)

	mux := http.NewServeMux()
	setupRestEndpoints(mux, ratingHandlerV1)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func startGrpc(app driving.IApplication, port int) {
	ratingHandler := grpchandler.NewRatingHandler(app)
	address := fmt.Sprintf("%s:%d", domain, port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterRatingServiceServer(server, &ratingHandler)
	server.Serve(listener)
}
