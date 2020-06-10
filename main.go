package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
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
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_1;SharedAccessKey=TFIWmO9qj3aBiJhTffhJbEz5vCOgC5qG+bYosS2J5Zo=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_2;SharedAccessKey=l4mnChokjqjW+NmR6VWMAWBzhK2jo/X4zYB45tNzl7c=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_3;SharedAccessKey=cPvQelbknUqwXXksdumBZca75hHF6a+ObJqwC5eaQAo=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_4;SharedAccessKey=Vpfzc6YT4bkpnWPidNv2ZA84mizrwpQHbyGzyzcfSBI=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_5;SharedAccessKey=rN8zgh1/admTPpvEQlNrKALvGaj3QIZsYnufT9IPtoQ=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_6;SharedAccessKey=JvxNrcCLQGJVUNWrq+Gk53ZHiaueOhvdc83wkD+NJ1E=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_7;SharedAccessKey=4ksCHG7I0r3yBbfLUpfgJvLuplkT5sSJlI741ItJln8=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_8;SharedAccessKey=4uvUHGRAWAjBAqNITTfcMX21iGRnlG5L1NBwEydbuuE=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_9;SharedAccessKey=wHTT8Fs4Gqlg/da6M3dcLS342QRFeh5BSEVIsfOhdck=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_10;SharedAccessKey=xrDIYDi9cJwlbw409nUW075jnUgdzy3Qs5WHkbpcJak=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_11;SharedAccessKey=ljrnqFMsoofm2umLhMf66nw0qNzrUat14swv2/ai1VQ=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_12;SharedAccessKey=F8AdPjpiEjgbW7CeKIGGRRafa3RdKi51arvqJoWxx+k=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_13;SharedAccessKey=q6TxvSk5aXVVwgyGwMVgInKJTJVu6G/jgiYFLYpqY4Y=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_14;SharedAccessKey=0lMTgFN3cza35GaRZRMBJAx6jLnlvu+ylte8tu8ZPBk=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_15;SharedAccessKey=FUKAJPlmzijxup+go2Z39iGnxbOuhp1d1+kXo+Xuo3U=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_16;SharedAccessKey=mh4r1khaxLpUxKMTl4QHQPnbZS+MWC8rvigtDBBAVGw=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_17;SharedAccessKey=UzzF58OJgAtzsbTQzBNejc+ihZyfkKwebzBgIA+Os/E=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_18;SharedAccessKey=NIYaVvSjp5n+PTHFOeUK+b/UzNo83wxcUWv0bDIz9d0=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_19;SharedAccessKey=eVpMLwby41Gs8S82xJ+HoZ9c4lEOsucoIfyCYqd9IP0=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_20;SharedAccessKey=SDWCRC/99yLw6OSK9jJpyF6koONcYfY/xs0QQgiVSUY=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_21;SharedAccessKey=ByOkFeTfXk44owlSbFZh1iebIeSK32e2QRTtgMEjZ3E=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_22;SharedAccessKey=5xGfvvyTK6FIXfH0RU5Xq9+HqeF+7OjGFL72neE0t8Q=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_23;SharedAccessKey=u03hDshv7Ks3+OoJN2WHh0AosHQqcoacIMt23B+YSaw=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_24;SharedAccessKey=Ce5MynfeIKQrQfHqYBgZ69VCwDnN+ynBOCE74KpG1po=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_25;SharedAccessKey=h6MvwyrcgTWdxWwSSzPIVdS58DK1XQNbeh2FqiWOI2g=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_26;SharedAccessKey=qUbb4xX30NCEPs+VY6/G2PMnoJi8XonVLYyAY4PEzzc=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_27;SharedAccessKey=OytCK2DXOTksy4SXZhMNl1DGjctHm+IUcj2VTMZBC9Q=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_28;SharedAccessKey=HLY44ilQ75Ilf7T//kqTTm1kYMvNtYhsTXViQvsfhKw=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_29;SharedAccessKey=1kV8tan04eySHFx+0s0GNmieF1or6sD1LrTD1wbTNfg=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_30;SharedAccessKey=4F4pW/zvXvTu+GvRVNmn7E4+mQxJamAw152oDD2nZko=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_31;SharedAccessKey=ORnwOteJ3pSr1ImRR8CM9YWdceBlWzfEjMYqqd71d3E=",
	"HostName=jio-iot-hub.azure-devices.net;DeviceId=jio_device_32;SharedAccessKey=46fKtbdhQeeX35ue0Ic3LZmmNkI9ockrdGeRvPpPu7g=",
}

func main() {
	router = gin.Default()
	router.POST("/callback/", scefHandler)
	router.GET("/callback/", scefHandlerGet)
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

	if index >= 31 {
		index = 0
	}
	fmt.Println(index)
	fmt.Println(connectinStrings[index])

	go send_data(connectinStrings[index], string(sDec))
	index = index + 1
	c.Status(http.StatusNoContent)

}

func send_data(connectionString string, data string) {

	var client *iotdevice.Client
	var err error
	client, err = iotdevice.NewFromConnectionString(
		iotmqtt.New(), connectionString,
	)
	fmt.Printf("%T\n", client)
	fmt.Printf("%T\n", err)
	fmt.Println(client)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Error in creating new connection")
	}

	// connect to the iothub
	if err = client.Connect(context.Background()); err != nil {
		fmt.Println("debug 2")
		//log.Fatal(err)
	}
	fmt.Println(data)
	// send a device-to-cloud message
	if err = client.SendEvent(context.Background(), []byte(data)); err != nil {
		fmt.Println("debug 3")
		//log.Fatal(err)
	}
}

func scefHandlerGet(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "It Worked",
	})
}
