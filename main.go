package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
)

type ScefData struct {
	Data   string `json: data`
	Msisdn string `json: msisdn`
}

var index int = 0
var router *gin.Engine
var connectinStrings [32]string = [32]string{
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_01;SharedAccessKey=0IdXWVPsoS6ImzITBxt/4AbR352DwNG7wUJjeKtdsWg=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_02;SharedAccessKey=9mpAntX8Dnp/Mwrmdo1+yNYHBx0c/xgtpuYSYV+zuOg=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_03;SharedAccessKey=5wtz/SG/xjdq5HjobBT2AzrAeawI6PcXhIEhbrFlj8w=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_04;SharedAccessKey=fywSKWOPBIhd/h9Fx4/qcSUnNJQv3G7ceGWSf/XVisA=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_05;SharedAccessKey=y7kBsxv4YrgjrqAHcEGFCjBkYK+5DdYQkUOmLsQ0rSk=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_06;SharedAccessKey=gQn3tN5P6NwgNR7GLRGjDAExKQ0woEkHzM77tNdRTmw=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_07;SharedAccessKey=kCH37m5Ce6n4eX9nyw3TBLOXCD8ePSijQakMUnUXDTE=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_08;SharedAccessKey=8yE2toa0HhGvwOTPvnfxgsAjhMWF0GKSoNY3qUQEFv4=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_09;SharedAccessKey=/zGevtkDfu2vC1w+EkI6UmYn21u/vsz3Edar6bvvAmE=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_10;SharedAccessKey=zn4QrgJ9N7ZPDV5yrEXI4i/tZq73AkzRdGzPoZ3C4Ko=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_11;SharedAccessKey=4J7AUb7Qp1yfZ7X7Tn68LivxTjv9fWy3MwTmJxNSwbM=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_12;SharedAccessKey=WkBvEeFwBG3SzYeF1K1nIKX9dNzjBDhspgSHr61Yw9M=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_13;SharedAccessKey=jSBvlneg+c9Md5JqQ/XwvnAJirux93Tzow7ed+MYWHY=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_14;SharedAccessKey=IWw7PxjgL+oTRMRMa1xj470+Y0yobbVg4MlyV2EFs8M=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_15;SharedAccessKey=kbjPxlc7qV6R97uULdxmfhPSZKEXfv5QSoellVUGR3Y=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_16;SharedAccessKey=x0/vOFOU+cd4p+1QaWt4BNkR85IynDBS91R60s0cWk0=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_17;SharedAccessKey=VRUHG2pBLoQ4dOiITlOtqiNG0x9YENQw70rxp5fWq8M=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_18;SharedAccessKey=MSZ5oca8lxv9JS0GcyzAtda3Pv17/9qdRG+DDUcR508=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_19;SharedAccessKey=Mx61NTCx7C4PFmq9U8vQlWXylywam6D/y78myCj8XT8=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_20;SharedAccessKey=h8StTNr7tFNHOIBda/V28yVy47Y5IgHACkcOqVZo8Lw=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_21;SharedAccessKey=gpPdWD7UmtS4M+l9w343nsRhrJZgTzPO8u9jiUMThhw=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_22;SharedAccessKey=Hdo/DMYgWUMWhBLWEyegWbuIBxIDHf53k5y6i5mtPT0=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_23;SharedAccessKey=78NQlFGhCjiYc47pvU9xGFNzTIV2pL4mBNSf/QkJOhI=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_24;SharedAccessKey=NlN1eR34/6lHT2ofGjWw8P7R+jKfNA1JLlnXsmtFbQ0=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_25;SharedAccessKey=P512f/og+4dfLExJDt0n9j2VjdM0bW36mErD2bSFfHI=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_26;SharedAccessKey=CHK7xqMhTjHNs3Wx8aqQLwYV9yyLKxSfvvT3U6m05IA=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_27;SharedAccessKey=5FOQuKW0ZlctM03ac7L4XjGo3qlRm6t7MNMGzQ55gDg=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_28;SharedAccessKey=ak8y+M1FWzWal9nURSN5EkAsI20fUTJC6xZsPeNJluY=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_29;SharedAccessKey=XA4Fbn/0LWUP7jM2ytpTuadEuOl0GiSIRZpBZsNRGeo=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_30;SharedAccessKey=77tV3X7XarihOgSifmtiS7tLi/QUh3RYLMxs/2u/VIQ=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_31;SharedAccessKey=dCmNRZXO6dyOvjClWxQc/stcymGZuTEx3O1UYTFXJFQ=",
	"HostName=jio-iot-central.azure-devices.net;DeviceId=jio_device_32;SharedAccessKey=OvxAINrSxHviAjbl2QGN1+md7nf+b2ZbwU6yl2J5H58=",
}

var clients []*iotdevice.Client

func main() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	createConnection()
	router.POST("/callback/", scefHandler)
	router.GET("/callback/", scefHandlerGet)
	router.Run(":7000")
}

func createConnection() {
	var client *iotdevice.Client
	var err error

	for _, connectionString := range connectinStrings {
		client, err = iotdevice.NewFromConnectionString(
			iotmqtt.New(), connectionString,
		)
		fmt.Println(client)
		if err != nil {
			log.Fatal(err)
		}

		// connect to the iothub
		if err = client.Connect(context.Background()); err != nil {
			fmt.Println("debug 2")
			log.Fatal(err)
		}
		clients = append(clients, client)
	}

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

	if index >= 31 {
		index = 0
	}
	fmt.Println(index)
	fmt.Println(connectinStrings[index])
	go send_data(index, string(sDec))
	index = index + 1
	c.Status(http.StatusNoContent)

}

func send_data(index int, data string) {
	var err error
	fmt.Println(data)
	// send a device-to-cloud message
	//fmt.Println(clients[1])
	if err = clients[index].SendEvent(context.Background(), []byte(data)); err != nil {
		fmt.Println("debug 3")
		log.Fatal(err)
	}
}

func scefHandlerGet(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "It Worked",
	})
}
