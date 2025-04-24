package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"ecom/global"
	"ecom/internal/database"
	"ecom/internal/service"
	"ecom/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumeMessage struct {
	rabbitMQManager *rabbitmq.QueueManager
	testService     service.ITestService
}

func NewConsumeMessage(
	testService service.ITestService,
) *ConsumeMessage {
	return &ConsumeMessage{
		rabbitMQManager: global.RabbitMQManager,
		testService:     testService,
	}
}

func (c *ConsumeMessage) marshalBody(response rabbitmq.QueueResponse) []byte {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v\n", err)
		return []byte("{}")
	}
	return jsonResponse
}

type BodyMessage struct {
	Data   interface{} `json:"data"`
	Action string      `json:"action"`
	UserID string      `json:"user_id"`
}

func (c *ConsumeMessage) RegisterConsumers() {
	fmt.Println("RegisterConsumers")
	// Order Queue Consumer
	number := strings.Split(global.Config.Queue.Test, ":")[1]
	name := strings.Split(global.Config.Queue.Test, ":")[0]
	numberInt, err := strconv.Atoi(number)
	if err != nil {
		log.Printf("Failed to convert number to int: %v\n", err)
		return
	}
	for i := 0; i < numberInt; i++ {
		queue := fmt.Sprintf("%s:%d", name, i)
		fmt.Println("queue", queue)
		err := c.rabbitMQManager.Consume(queue, func(msg amqp.Delivery) {
			fmt.Println("Received message:", string(msg.Body))
			body := BodyMessage{}
			err := json.Unmarshal(msg.Body, &body)
			if err != nil {
				log.Printf("Failed to unmarshal body: %v\n", err)
				c.sendResponse(msg, rabbitmq.QueueResponse{
					CodeResult: http.StatusBadRequest,
					Data:       nil,
					Error:      "Failed to parse message body",
				})
				return
			}
			fmt.Println("Body:", body)
			var response rabbitmq.QueueResponse

			switch body.Action {
			case "create":
				// Convert the data map to JSON bytes
				dataBytes, err := json.Marshal(body.Data)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				var req database.CreateTestParams
				err = json.Unmarshal(dataBytes, &req)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				test, err := c.testService.CreateTest(&req)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				jsonResponse, err := json.Marshal(test)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				response.CodeResult = http.StatusOK
				response.Data = &jsonResponse
				c.sendResponse(msg, response)
				break
			case "update":
				// Convert the data map to JSON bytes
				dataBytes, err := json.Marshal(body.Data)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				var req database.UpdateTestParams
				err = json.Unmarshal(dataBytes, &req)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				test, err := c.testService.UpdateTest(&req)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				jsonResponse, err := json.Marshal(test)
				if err != nil {
					response.CodeResult = http.StatusBadRequest
					response.Data = nil
					response.Error = err.Error()
					c.sendResponse(msg, response)
					return
				}

				response.CodeResult = http.StatusOK
				response.Data = &jsonResponse
				c.sendResponse(msg, response)
				break
			default:
				log.Printf("Unknown action in routing key: %s\n", body.Action)
				response.CodeResult = http.StatusBadRequest
				response.Data = nil
				response.Error = "Unknown action"
				c.sendResponse(msg, response)
			}
		})
		if err != nil {
			fmt.Println("Failed to consume from log_queue:", err)
		}
	}
}

// Helper function to send response
func (c *ConsumeMessage) sendResponse(msg amqp.Delivery, response rabbitmq.QueueResponse) {
	if msg.ReplyTo != "" {
		c.rabbitMQManager.Channel.Publish(
			"",
			msg.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: msg.CorrelationId,
				Body:          c.marshalBody(response),
			},
		)
	}
}
