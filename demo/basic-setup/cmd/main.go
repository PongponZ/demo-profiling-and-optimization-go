package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"net/http"
	_ "net/http/pprof"

	"github.com/PongponZ/demo-profiling-and-optimization-go/basic-setup/internal/handler"
	"github.com/PongponZ/demo-profiling-and-optimization-go/libs"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
	"github.com/xyproto/randomstring"
)

type Job struct {
	Name string `json:"name"`
}

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

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	rabbitMQQueue := os.Getenv("RABBITMQ_QUEUE")

	rmq := libs.NewRabbitMQClient(rabbitMQURL, 100)
	defer rmq.Close()

	rmq.QueueDeclare(rabbitMQQueue)

	handler := handler.NewLeakHandler()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Demo profiling and optimization in Go")
	})

	app.Get("/goleak", handler.GoroutineLeak)
	app.Get("/block", handler.Block)      // Route that causes blocking (mutex contention)
	app.Get("/alloc", handler.Alloc)      // Route that causes heavy allocations
	app.Get("/cpu", handler.CPUIntensive) // Route that causes high CPU usage

	app.Get("/publish/:number", func(c *fiber.Ctx) error {
		number := c.Params("number")
		numberInt, err := strconv.Atoi(number)
		if err != nil {
			return c.SendString("invalid number")
		}

		for range numberInt {
			j := Job{
				Name: randomstring.HumanFriendlyString(7),
			}

			data, err := json.Marshal(j)
			if err != nil {
				return c.SendString("error marshalling job")
			}

			rmq.Channel().Publish(
				"",
				rabbitMQQueue,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        data,
				})
		}

		return c.SendString(fmt.Sprintf("published %d messages", numberInt))
	})

	app.Listen(":3010")
}
