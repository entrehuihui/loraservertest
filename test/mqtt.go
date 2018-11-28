/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package main

import (
	"crypto/tls"
	"flag"
	"fmt" //"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var onMessageReceivedNum = 0

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	// fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	// onMessageReceivedNum++
	// fmt.Println("完成数量: ", onMessageReceivedNum)
	endTime = time.Now().UnixNano() / 1e6
	chanRX <- 1
}

func countNum() {
	for {
		<-chanGo
		<-chanRX
		onMessageReceivedNum++
		fmt.Println("完成数量: ", onMessageReceivedNum)
	}
}

func mqttserver() {
	go countNum()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	hostname, _ := os.Hostname()
	server := flag.String("server", "tcp://127.0.0.1:1883", "The full URL of the MQTT server to connect to")
	topic := flag.String("topic1", "pressure_test", "Topic to subscribe to1")
	qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "loraserver_ns", "A username to authenticate to the MQTT server")
	password := flag.String("password", "loraserver123", "Password to match username")
	flag.Parse()
	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *server)
	}
	<-c
}
