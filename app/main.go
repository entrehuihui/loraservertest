package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"../common"
	"../handlers"
)

const addr = "127.0.0.1:1700"
const laddr = ":5453"

func main() {
	JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJleHAiOjE1NDIxNzkxMjEsImlzcyI6ImxvcmEtYXBwLXNlcnZlciIsIm5iZiI6MTU0MjA5MjcyMSwic3ViIjoidXNlciIsInVzZXJuYW1lIjoiYWRtaW4ifQ.nHZZYdPuLFuW4ei7muQQFmVEE8LU-WXPRNbb-6rZZQg"
	//获取JWT登陆信息
	// Login("admin", "admin")
	// if len(JWT) == 0 {
	// 	log.Fatal("login fail")
	// }
	// 获取应用ID
	// IDArray := GetApplications()
	// if len(IDArray) == 0 {
	// 	log.Fatal(" not have Application")
	// }
	//获取应用内设备列表
	// GetDevices("3")
	//获取设备keys
	// GetDeviceKeys("24c5d9e63257f58c")
	//获取设备Activation
	// GetDeviceActivation("24c5d9e63257f58c")

	PostDevice("114", "大大萨达萨达", "1111111111111114", "3", "ceec88af-86e9-4d3a-a372-693440be67d4", true)
	return
	rand.Seed(time.Now().UnixNano())
	timeStart := time.Now().UnixNano() / 1e6
	fmt.Println("start on:", timeStart)
	port := 65453
	wg := &sync.WaitGroup{}
	for i, mac := range common.Macs {
		laddr := fmt.Sprintf(":%d", port+i)
		client, err := handlers.NewClient(laddr, mac, addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		wg.Add(1)
		go func(c *handlers.GatewayClient) {
			defer wg.Done()
			c.Ping()
		}(client)

		who := false
		// who = true
		if who {
			if i == 0 {
				fmt.Println("======================")
				err = client.JoinRequest("2b7e151628aed2a6", "11111111111111111111111111111111")
				if err != nil {
					fmt.Println("sendData error:", err)
				}
				time.Sleep(5 * time.Second)
				return
			}
		}

		go func() {
			temperature := 38.0
			humidity := 47.0
			r1 := rand.New(rand.NewSource(time.Now().UnixNano() + 3))
			r2 := rand.New(rand.NewSource(time.Now().UnixNano() + 5))
			ra1 := 0.47
			ra2 := 0.47
			// for i := 0; i < 10000; i++ {
			for {
				// for {
				var count uint32 = 1
				temperature = temperature - (r1.Float64()-ra1)/10
				humidity = humidity - (r2.Float64()-ra2)/10
				if temperature <= 0 {
					ra1 += 0.01
				} else if temperature >= 40 {
					ra1 -= 0.01
				}
				if humidity <= 30 {
					ra2 += 0.01
				} else if humidity >= 80 {
					ra2 -= 0.01
				}
				deviceInfo := make(map[string]interface{})
				deviceInfo["temperature"] = strconv.FormatFloat(temperature, 'f', 2, 64)
				deviceInfo["humidity"] = strconv.FormatFloat(humidity, 'f', 2, 64)
				deviceInfo["time_device"] = time.Now().Unix()
				b, err := json.Marshal(deviceInfo)
				if err != nil {
					continue
				}
				data := string(b)
				fmt.Println(data)
				// client.SendData("01cb8a27", "e310aff65e7c9b2a01c4ffa3432e9896", "a781a5f436705dbcce9028925479c244", data, "string", count) //127--120
				// client.SendData("003b914d", "d4a7105edc0c4ffd918f3d1139c76626", "f2479e86f91eb46debc0793f04855edf", data, "string", count) //127--119
				// client.SendData("01f734a6", "c9790bbf6a4fd96a541914e32405d1f2", "8ee9c3dbdbbdc99f6bb18e01d5da78c5", data, "string", count) //127--111
				client.SendData("01d35d38", "9fc2687ce9b3755ce09d514a7e1969fa", "34d7aa1cf1e9332243a40bce3c10a814", data, "string", count) //127--110
				// client.SendData("019e2550", "47766f6056ad38127a717d593d01c93c", "2076c1a6da6f1b332365510fad6d4271", data, "string", count) //---120 dve1
				// client.SendData("01587d6f", "bf231e9619d9d2bbbefe4e53f4be23ea", "924100f401d983c170f2fc38bb4c56ba", data, "string", count) //120 --dev2
				time.Sleep(time.Millisecond * 5000)
				count++
			}
			timeEnd := time.Now().UnixNano() / 1e6
			fmt.Println("time end:", timeEnd)
			ioutil.WriteFile("time.txt", []byte(strconv.FormatInt(timeStart, 10)+":"+strconv.FormatInt(timeEnd, 10)), 0644)
		}()

		wg.Add(1)
		go func(client *handlers.GatewayClient) {
			defer wg.Done()
			for {
				redb := make([]byte, 1024)
				n, err := client.Conn.Read(redb)
				if err != nil {
					fmt.Println("read error:", err)
					return
				}
				if n > 12 {
					fmt.Printf("read from server:%x\n", redb[12:n])
				} else {
					fmt.Printf("read from server:%x\n", redb[:n])
				}
			}
		}(client)
	}
	wg.Wait()
}
