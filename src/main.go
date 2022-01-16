package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const CONTROL_MESSAGE_REGISTER_DEVICE = 0
const CONTROL_MESSAGE_SHARE_USERS = 3

const CONTROL_USER_MESSAGE_REQUEST_CHAT = 1
const CONTROL_USER_MESSAGE_ACCEPT_CHAT = 2
const CONTROL_TOPIC = "topic/control"

type User struct {
	Id     string `json:"id"`
	Status bool   `json:"status"`
}

type ControlMessage struct {
	Id     string `json:"id"`
	Status bool   `json:"status"`
	Type   int    `json:"type"`
	Topic  string `json:"topic"`
	Users  []User `json:"users"`
}

type Chat struct {
	TargetId string        `json:"targetId"`
	Topic    string        `json:"topic"`
	Message  []MessageChat `json:"messageChat"`
}

type MessageChat struct {
	Message  string `json:"message"`
	Time     string `json:"time"`
	SenderId string `json:"senderId"`
}

var usersList = make([]User, 0)
var usersChatList = make([]Chat, 0)
var requestMessageList = make([]User, 0)

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

// go run ./src/main.go -s
// go run ./src/main.go -id=1

func main() {
	userId := flag.String("id", "-1", "user id")
	createPersistentMessage := flag.Bool("s", false, "Iniciar aplicação")
	flag.Parse()

	opts := createNewClient(userId)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if *createPersistentMessage {
		sendCurrentUsersListToControl(client)
		fmt.Println("Configurações iniciadas com sucesso")
		client.Disconnect(250)
		os.Exit(0)
	}

	// register user on vector. if it already exists change it status to online

	userControlSub(client, userId)

	controlTopicSub(client, *userId)

	go func() {
		for {
			time.Sleep(8 * time.Second)
			fmt.Printf("%+v\n", usersList)
		}
	}()

	for {
	}

	// use it to break
	// client.Disconnect(250)
}

func registerUserToList(client mqtt.Client, userId string) {
	var user *User = nil

	for _, currentUser := range usersList {
		if currentUser.Id == userId {
			user = &currentUser
			break
		}
	}

	if user != nil && user.Status {
		fmt.Println("User is already nline")
		os.Exit(3)
	}

	if user != nil {
		user.Status = true
	} else {
		usersList = append(usersList, User{Id: userId, Status: true})
	}
}

func handleChatMessages(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	topic := msg.Topic()

	var message MessageChat

	json.Unmarshal([]byte(payload), &message)

	var currentChat Chat

	for _, chat := range usersChatList {
		if chat.Topic == topic {
			currentChat = chat
			break
		}
	}

	currentChat.Message = append(currentChat.Message, message)
}

func sendMessageToUser(client mqtt.Client, topic string, senderId string, messageString string) {
	client.Publish(topic, 0, false, MessageChat{
		Message:  messageString,
		Time:     dateNowAsString(),
		SenderId: senderId,
	})
}

func dateNowAsString() string {
	return time.Now().Format(time.RFC3339)
}

func acceptRequestToUserChat(client mqtt.Client, userId string, targetId string) {
	topic := createTopicString(targetId + "_" + userId + "_" + dateNowAsString())

	client.Subscribe(topic, 1, handleChatMessages)
	// subscrever ao topico x_y_timestamp

	// enviar mensagem ao target_solicitante com o tópico
	topicControlRequester := createControlTopicString(targetId)

	payload := stringify(&ControlMessage{
		Topic: topic,
		Id:    userId,
		Type:  CONTROL_USER_MESSAGE_ACCEPT_CHAT,
	})

	client.Publish(topicControlRequester, 0, false, payload)

	// save to chats list
	usersChatList = append(usersChatList, Chat{
		TargetId: targetId,
		Topic:    topic,
		Message:  make([]MessageChat, 100),
	})
}

func sendRequestToUser(client mqtt.Client, userId string, targetId string) {
	topic := createControlTopicString(targetId)

	message := &ControlMessage{
		Id:   userId,
		Type: CONTROL_USER_MESSAGE_REQUEST_CHAT,
	}

	payload := stringify(message)

	stringMessage := stringify(payload)

	client.Publish(topic, 0, false, stringMessage)
}

func createControlTopicString(userId string) string {
	return createTopicString(userId + "_control")
}

func createTopicString(name string) string {
	return "topic/" + name
}

func controlTopicSub(client mqtt.Client, userId string) {
	firstSharedUsers := true

	var messagePubHandlerControl mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		var message ControlMessage

		json.Unmarshal([]byte(payload), &message)
		fmt.Printf("%+v\n", message)

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
		case CONTROL_MESSAGE_SHARE_USERS:
			{
				fmt.Println("lksklqs")
				usersList = message.Users

				if firstSharedUsers {
					registerUserToList(client, userId)

					// inform others you are online
					sendCurrentUsersListToControl(client)

					firstSharedUsers = false
				}

				break
			}
		}
	}

	// Descobrir pra que serve o 1
	token := client.Subscribe(CONTROL_TOPIC, 0, messagePubHandlerControl)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", CONTROL_TOPIC)
}

func sendCurrentUsersListToControl(client mqtt.Client) {
	allUSers := &ControlMessage{
		Type:  CONTROL_MESSAGE_SHARE_USERS,
		Users: usersList,
	}

	fmt.Println(stringify(allUSers))

	client.Publish(CONTROL_TOPIC, 1, true, stringify(allUSers))
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
		fmt.Println(payload)
		// mudar isso de control pra control account

		switch message.Type {
		case CONTROL_USER_MESSAGE_REQUEST_CHAT:
			{
				// add message to REQUEST_MESSAGE_LIST;
				requestMessageList = append(requestMessageList, User{
					Id: message.Id,
				})

				break
			}
		case CONTROL_USER_MESSAGE_ACCEPT_CHAT:
			{
				usersChatList = append(usersChatList, Chat{
					TargetId: message.Id,
					Topic:    message.Topic,
				})
			}
		}
	}

	topic := createControlTopicString(*user)

	token := client.Subscribe(topic, 1, userPubHandlerControl)
	token.Wait()
}
