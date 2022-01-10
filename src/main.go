package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/twinj/uuid"
)

type User struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func main() {
	userId := flag.String("userId", "-1", "id do usuário")
	flag.Parse()

	jsonFile, _ := os.Open("users.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []User
	json.Unmarshal([]byte(byteValue), &result)

	fmt.Println(result)
	fmt.Printf("%+v\n", result)

	addNewToList := true

	for _, value := range result {
		if value.Id == *userId {
			addNewToList = false
			value.Status = "online"

			file, _ := json.MarshalIndent(result, "", " ")

			_ = ioutil.WriteFile("users.json", file, 0644)
			break
		}
	}

	if addNewToList {
		var idToCreate string

		if *userId == "-1" {
			idToCreate = uuid.NewV4().String()
		} else {
			idToCreate = *userId
		}

		newUser := User{
			Id: idToCreate,
		}

		result = append(result, newUser)
		file, _ := json.MarshalIndent(result, "", " ")
		_ = ioutil.WriteFile("users.json", file, 0644)
	}

	for {
		var input string

		fmt.Println("1 - Listar Usuários")
		fmt.Println("2 - Criar novo grupo")
		fmt.Println("3 - Listar grupos")
		fmt.Println("4 - Nova conversa")
		fmt.Println("5 - Listagem do histórico de solicitação recebidas")
		fmt.Println("6 - Listagem das confirmações de aceitação da solicitação de batepapo")

		fmt.Scanln(&input)

		switch input {
		case "1":
			fmt.Println("one")
		case "2":
			fmt.Println("two")
		case "3":
			fmt.Println("three")
		default:
			fmt.Println("Nennum")
		}

		fmt.Println(input)
	}
}

// var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// 	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
// }

// var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
// 	fmt.Println("Connected")
// }

// var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
// 	fmt.Printf("Connect lost: %v", err)
// }

// func createNewClient() *mqtt.ClientOptions {
// 	var broker = "localhost"
// 	var port = 1883

// 	opts := mqtt.NewClientOptions()
// 	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
// 	opts.SetClientID("go_mqtt_client")
// 	opts.SetUsername("emqx")
// 	opts.SetPassword("public")
// 	opts.SetDefaultPublishHandler(messagePubHandler)

// 	opts.OnConnect = connectHandler
// 	opts.OnConnectionLost = connectLostHandler

// 	return opts
// }

// func main() {
// 	opts := createNewClient()

// 	client := mqtt.NewClient(opts)
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		panic(token.Error())
// 	}

// 	sub(client)
// 	publish(client)

// 	client.Disconnect(250)
// }

// func publish(client mqtt.Client) {
// 	num := 10
// 	for i := 0; i < num; i++ {
// 		text := fmt.Sprintf("Message %d", i)
// 		token := client.Publish("topic/test", 0, false, text)
// 		token.Wait()
// 		time.Sleep(time.Second)
// 	}
// }

// func sub(client mqtt.Client) {
// 	topic := "topic/test"
// 	token := client.Subscribe(topic, 1, nil)
// 	token.Wait()
// 	fmt.Printf("Subscribed to topic: %s", topic)
// }
