package main

import (
	"encoding/base64"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type ScefData struct {
	Data   string `json: data`
	Msisdn string `json: msisdn`
}

var router *gin.Engine

func main() {
	router = gin.Default()
	router.POST("/mo/callback/", scefHandler)
	router.Run(":7000")
}

func scefHandler(c *gin.Context) {

	var u ScefData
	c.BindJSON(&u)
	fmt.Println(u.Msisdn)
	sDec, err := base64.StdEncoding.DecodeString(u.Data)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}

	fmt.Println(string(sDec))
	go send_data(string(sDec))
	c.Status(http.StatusNoContent)

}

func send_data(data string) {
	//topic := flag.String("topic", "", "The topic name to/from which to publish/subscribe")
	//broker := flag.String("broker", "tcp://iot.eclipse.org:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	//password := flag.String("password", "", "The password (optional)")
	//user := flag.String("user", "", "The User (optional)")
	//id := flag.String("id", "testgoid", "The ClientID (optional)")
	//cleansess := flag.Bool("clean", false, "Set Clean Session (default false)")
	//qos := flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	//num := flag.Int("num", 1, "The number of messages to publish or subscribe (default 1)")
	//payload := flag.String("message", "", "The message text to publish (default empty)")
	//action := flag.String("action", "", "Action publish or subscribe (required)")
	//store := flag.String("store", ":memory:", "The Store Directory (default use memory store)")
	//flag.Parse()
	topic := "devices/jio_device_23/messages/events/"
	action := "pub"
	broker := "ssl://jio-iot-hub.azure-devices.net:8883"
	user := "jio-iot-hub.azure-devices.net/jio_device_23/?api-version=2018-06-30"
	password := "SharedAccessSignature sr=jio-iot-hub.azure-devices.net%2Fdevices%2Fjio_device_23&sig=WxAhEEqnaDwCauZEAOzzMMj5SO8lMN5ka8%2FvItP%2BhM8%3D&se=1591081761"
	id := "jio_device_23"
	cleansess := false
	qos := 0
	num := 1
	store := ":memory:"
	payload := data

	if action != "pub" && action != "sub" {
		fmt.Println("Invalid setting for -action, must be pub or sub")
		return
	}

	if topic == "" {
		fmt.Println("Invalid setting for -topic, must not be empty")
		return
	}

	fmt.Printf("Sample Info:\n")
	fmt.Printf("\taction:    %s\n", action)
	fmt.Printf("\tbroker:    %s\n", broker)
	fmt.Printf("\tclientid:  %s\n", id)
	fmt.Printf("\tuser:      %s\n", user)
	fmt.Printf("\tpassword:  %s\n", password)
	fmt.Printf("\ttopic:     %s\n", topic)
	fmt.Printf("\tmessage:   %s\n", payload)
	fmt.Printf("\tqos:       %d\n", qos)
	fmt.Printf("\tcleansess: %v\n", cleansess)
	fmt.Printf("\tnum:       %d\n", num)
	fmt.Printf("\tstore:     %s\n", store)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(cleansess)
	if store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(store))
	}

	if action == "pub" {
		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println("Creating connection")
			panic(token.Error())
		}
		fmt.Println("Sample Publisher Started")
		for i := 0; i < num; i++ {
			fmt.Println("---- doing publish ----")
			token := client.Publish(topic, byte(qos), false, payload)
			token.Wait()
		}

		client.Disconnect(250)
		fmt.Println("Sample Publisher Disconnected")
	} else {
		receiveCount := 0
		choke := make(chan [2]string)

		opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
			choke <- [2]string{msg.Topic(), string(msg.Payload())}
		})

		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := client.Subscribe(topic, byte(qos), nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		for receiveCount < num {
			incoming := <-choke
			fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			receiveCount++
		}

		client.Disconnect(10)
		fmt.Println("Sample Subscriber Disconnected")
	}
}
