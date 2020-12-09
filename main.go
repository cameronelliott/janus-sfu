

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	

	"os"
	"strconv"
	"time"
	// "net/url"
	// "log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {	


	//MQTT.DEBUG = log.New(os.Stdout, "", 0)
	//MQTT.ERROR = log.New(os.Stdout, "", 0)
	hostname, _ := os.Hostname()
	


	urlflag := flag.String("url", "tcp://127.0.0.1:1883", "url of server")
	// topic := flag.String("topic", hostname, "topic to publish to")
	// qos := flag.Int("qos", 0, "publish message qos flag")
	// retained := flag.Bool("retained", false, "are messages retained")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "clientid")
	usernameflag := flag.String("username", "", "username")
	passwordflag  := flag.String("password", "", "password")

	
	flag.Parse()

	user:=*usernameflag
	pass:=*passwordflag


	// u, err := url.Parse(*urlflag)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(u.Hostname())





	connOpts := MQTT.NewClientOptions().AddBroker(*urlflag).SetClientID(*clientid).SetCleanSession(true)
	if user != "" {
		connOpts.SetUsername(user)
		if pass != "" {
			connOpts.SetPassword(pass)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)


	choke := make(chan [2]string)

	connOpts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println("ouch...")
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})



	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
	fmt.Printf("connected to %s\n", *urlflag)

	m := make(map[string]byte)
	
	
	//{"offline": "foo/fc8df0c80c1bbf0d"}
	//{"online": "foo/2fc95e1469902ea0"}
	m["status"] = 0		// online/offline messages	, qos=0
	m["from-janus"] = 0	// responses to commands  , qos = 0

		//not yet  m["from-janus-admin"] = 13




	if token := client.SubscribeMultiple(m,nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}


	
	// for {
	// 	// message, err := stdin.ReadString('\n')
	// 	// if err == io.EOF {
	// 	// 	os.Exit(0)
	// 	// }
	// 	// client.Publish(*topic, byte(*qos), *retained, message)
	// }

	




	for  {
		fmt.Println("waiting...")
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
	}



}