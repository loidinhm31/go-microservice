package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/loidinhm31/go-microservice/common"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

const logTopic = "logs_topic"

var tools common.Tools

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	randomQueue, err := declareRandomQueue(channel)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		channel.QueueBind(
			randomQueue.Name,
			topic,
			logTopic,
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := channel.Consume(
		randomQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for msg := range messages {
			var payload Payload
			_ = json.Unmarshal(msg.Body, &payload)

			go handlePayload(payload)
		}
	}()

	log.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", randomQueue.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever it gets
		err := logEvent(payload)
		if err != nil {
			log.Println("Error when logging event:", err)
		}
	case "auth":
		// authenticate

	// you can have as many cases as you want, as long as you write the logic

	default:
		log.Println("Cannot find action to handle")
	}
}

func logEvent(entry Payload) error {

	jsonData, _ := json.Marshal(entry)

	logServiceURL := fmt.Sprintf("http://logger-service:%s/log", common.LoggerPort)

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
