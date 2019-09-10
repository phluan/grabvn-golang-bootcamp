package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
)

func main() {
	client, err := statsd.New("127.0.0.1:8125",
		statsd.WithNamespace("http-echo."),     // prefix every metric with the app name
		statsd.WithTags([]string{"localhost"}), // send the EC2 availability zone as a tag with every metric
		// add more options here...
	)
	if err != nil {
		fmt.Print("fail to connect to Datadog agent")
		log.Fatal(err)
	}

	// Create our server
	logger := log.New()
	server := Server{
		logger:          logger,
		metricCollector: client,
	}

	// Start the server
	server.ListenAndServe()
}

// Server represents our server.
type Server struct {
	logger          *log.Logger
	metricCollector *statsd.Client
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe() {
	s.logger.Info("echo server is starting on port 8080...")
	http.HandleFunc("/", s.echo)
	http.ListenAndServe(":8080", nil)
}

// Echo echos back the request as a response
func (s *Server) echo(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")

	// 30% chance of failure
	if rand.Intn(100) < 30 {
		s.metricCollector.Count("echo_failure", 1, []string{"env:localhost"}, 1)
		writer.WriteHeader(500)
		writer.Write([]byte("a chaos monkey broke your server"))
		return
	}

	// Happy path
	writer.WriteHeader(200)
	request.Write(writer)
}
