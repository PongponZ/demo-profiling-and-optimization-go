package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/PongponZ/demo-profiling-and-optimization-go/libs"
	"github.com/PongponZ/demo-profiling-and-optimization-go/super-worker/internal/controller"
	"github.com/PongponZ/demo-profiling-and-optimization-go/super-worker/internal/repo"
	"github.com/PongponZ/demo-profiling-and-optimization-go/super-worker/internal/usecase"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xyproto/randomstring"
)

func main() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	go func() {
		// pprof listening on port 6060
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	// Prometheus metrics endpoint on port 2112
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics server started on :2112")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// ไม่ต้องสนใจสิ่งนี้
	pokemonServer := simulateHttp()
	defer pokemonServer.Close()

	config := readConfig()

	rmq := libs.NewRabbitMQClient(config.RabbitMQURL, 1000)
	defer rmq.Close()

	rmq.QueueDeclare(config.RabbitMQQueue)
	rmq.QueueDeclare("pokemon_generated")
	msgs := rmq.Consume(config.RabbitMQQueue, "worker")

	pokemonRepo := repo.NewPokemon(pokemonServer.URL)
	pokemonUsecase := usecase.NewPokemonUsecase(pokemonRepo)
	worker := controller.NewWorker(config.MaxWorkers, pokemonUsecase, rmq.Channel())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("worker started ...")

	go worker.Start(msgs)

	<-signalChan
}

type config struct {
	MaxWorkers    int
	RabbitMQURL   string
	RabbitMQQueue string
}

func readConfig() *config {
	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		log.Fatalf("failed to parse MAX_WORKERS: %v", err)
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	rabbitMQQueue := os.Getenv("RABBITMQ_QUEUE")
	return &config{
		MaxWorkers:    maxWorkers,
		RabbitMQURL:   rabbitMQURL,
		RabbitMQQueue: rabbitMQQueue,
	}
}

func simulateHttp() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			limitInt = 3
		}

		abilities := map[string]int{}
		for i := 0; i < limitInt; i++ {
			abilities[randomstring.HumanFriendlyString(7)] = rand.Intn(100)
		}

		data, err := json.Marshal(abilities)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))

	return server
}
