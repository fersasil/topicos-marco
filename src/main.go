package main

import (
	"encoding/json"
	"flag"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const CONTROL_MESSAGE_REGISTER_DEVICE = 0
const CONTROL_TOPIC = "topic/control"

type User struct {
	Id     string `json:"id"`
	Status bool   `json:"status"`
}

type ControlMessage struct {
	Id     string `json:"id"`
	Status bool   `json:"status"`
	Type   int    `json:"type"`
}

var usersList = make([]User, 100)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func createNewClient(userId *string) *mqtt.ClientOptions {
	var broker = "localhost"
	var port = 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(string(*userId))

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	return opts
}

// go run ./src/main.go -id=1

func main() {
	userId := flag.String("id", "-1", "user id")
	flag.Parse()

	opts := createNewClient(userId)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	userControlSub(client, userId)

	controlTopicSub(client)

	message := ControlMessage{
		Id:     *userId,
		Status: true,
		Type:   CONTROL_MESSAGE_REGISTER_DEVICE,
	}

	controlTopicPub(client, &message, false)

	for {
	}

	// use it to break
	// client.Disconnect(250)
}

func controlTopicSub(client mqtt.Client) {
	var messagePubHandlerControl mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		var message ControlMessage

		json.Unmarshal([]byte(payload), &message)

		switch message.Type {
		case CONTROL_MESSAGE_REGISTER_DEVICE:
			{
				user := &User{
					Id:     message.Id,
					Status: message.Status,
				}

				usersList = append(usersList, *user)
				fmt.Println("Usuário adicionado- " + user.Id)
				break
			}
		}
	}

	// Descobrir pra que serve o 1
	token := client.Subscribe(CONTROL_TOPIC, 1, messagePubHandlerControl)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", CONTROL_TOPIC)
}

func controlTopicPub(client mqtt.Client, message *ControlMessage, retain bool) {
	stringMessage := stringify(message)

	client.Publish(CONTROL_TOPIC, 0, retain, stringMessage)
}

func stringify(item interface{}) string {
	b, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
	}

	return string(b)
}

func userControlSub(client mqtt.Client, user *string) {
	var userPubHandlerControl mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		var message ControlMessage

		json.Unmarshal([]byte(payload), &message)

		// switch message.Type {
		// case CONTROL_MESSAGE_REGISTER_DEVICE:
		// 	{
		// 		user := &User{
		// 			Id:     message.Id,
		// 			Status: message.Status,
		// 		}

		// 		usersList = append(usersList, *user)
		// 		fmt.Println("Usuário adicionado- " + user.Id)
		// 		break
		// 	}
		// }
	}

	topic := "topic/" + *user + "_control"

	token := client.Subscribe(topic, 1, userPubHandlerControl)
	token.Wait()
}
