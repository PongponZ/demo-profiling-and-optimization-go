package controller

import (
	"encoding/json"
	"log"

	"github.com/PongponZ/demo-profiling-and-optimization-go/super-worker/internal/usecase"
	"github.com/streadway/amqp"
)

type WorkerController struct {
	maxWorker      int
	pokemonUsecase *usecase.PokemonUsecase
	output         *amqp.Channel
}

func NewWorker(maxWorker int, pokemonUsecase *usecase.PokemonUsecase, output *amqp.Channel) *WorkerController {
	return &WorkerController{
		maxWorker:      maxWorker,
		pokemonUsecase: pokemonUsecase,
		output:         output,
	}
}

func (c *WorkerController) Start(messages <-chan amqp.Delivery) {
	for i := 0; i < c.maxWorker; i++ {
		go func() {
			for message := range messages {
				go c.processMessage(message)
			}
		}()
	}
}

func (c *WorkerController) processMessage(message amqp.Delivery) {
	log.Println("processing message ...")

	var job Job
	err := json.Unmarshal(message.Body, &job)
	if err != nil {
		log.Printf("error unmarshalling message: %v", err)
		return
	}

	pokemon := c.pokemonUsecase.GeneratePokemon(job.Name)
	log.Println("generating pokemon ...")

	data, err := json.Marshal(pokemon)
	if err != nil {
		log.Printf("error marshalling pokemon: %v", err)
		return
	}

	err = c.output.Publish(
		"",                  // exchange
		"pokemon_generated", // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		log.Printf("error publishing message: %v", err)
		return
	}

	log.Printf("published message ...")
}
