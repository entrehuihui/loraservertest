package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJleHAiOjE1NDIyNjg3MTcsImlzcyI6ImxvcmEtYXBwLXNlcnZlciIsIm5iZiI6MTU0MjE4MjMxNywic3ViIjoidXNlciIsInVzZXJuYW1lIjoiYWRtaW4ifQ.SfKocoilLvg1dzQXJpIOEzDy3DfFVTfqUaZ02SaeFws"
	//获取JWT登陆信息
	Login("admin", "admin") // 获取应用ID
	applicationsResult := GetApplications()
	if len(applicationsResult) == 0 {
		log.Fatal(" not have Application")
	}
	//发送数据测试
	SendTADA(applicationsResult[0].ID, 1, 2000)

	//TestAccess 入网请求测试 网关数,设备数
	// TestAccess(10,100)
	// return

	//给应用内所有应用发送信息
	// sendDataAll(applicationsResult[0].ID)
	// return

	// 批量创建设备
	// CreateALL(applicationsResult[0].ID)
	//激活应用内所有应用
	// ActivationDevices(applicationsResult[0].ID, "fffb02426f96c917", "127.0.0.1:1700")
}

//CreateALL 批量创建设备
func CreateALL(applicationsResultID string) {
	name := 0
	dev := "dev"
	eui := 1111111111111110
	createBatchDevicesInfos := make([]CreateBatchDevicesInfo, 0)
	for index := 0; index < 1000; {
		createBatchDevicesInfo := CreateBatchDevicesInfo{
			PplicationsID: applicationsResultID,
			DeviceName:    fmt.Sprintf("%s%d", dev, name),
			DeviceeEUI:    strconv.Itoa(eui),
			Synopsis:      "????",
			DeviceProfile: "ceec88af-86e9-4d3a-a372-693440be67d4",
			DisValidation: true,
		}
		createBatchDevicesInfos = append(createBatchDevicesInfos, createBatchDevicesInfo)
		name++
		eui++
		index++
	}
	// CreateBatchDevices 批量创建设备
	CreateBatchDevices(10, createBatchDevicesInfos)
}

//给应用内所有应用发送信息
func sendDataAll(applicationsResultID string) {
	//获取应用内设备列表
	devicesinfos := GetDevices(applicationsResultID, "", 9999)
	if len(devicesinfos) == 0 {
		return
	}
	// fmt.Println(len(devicesinfos))
	temperature := 38.0
	humidity := 47.0
	r1 := rand.New(rand.NewSource(time.Now().UnixNano() + 3))
	r2 := rand.New(rand.NewSource(time.Now().UnixNano() + 5))
	ra1 := 0.47
	ra2 := 0.47
	chanSendDataBatchInfo := make(chan SendDataBatchInfo, 100)
	wg := sync.WaitGroup{}
	go SendDataBatch("fffb02426f96c917", "127.0.0.1:1700", chanSendDataBatchInfo, &wg, 4432, 10)
	wg.Add(1)
	for i := 0; i < 10; i++ {
		for _, devicesinfo := range devicesinfos[:] {
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
			// fmt.Println(data)
			deviceActivation, err := GetDeviceActivation(devicesinfo.DevEUI)
			if err != nil {
				continue
			}
			sendDataBatchInfo := SendDataBatchInfo{
				Data:        data,
				Devicesinfo: deviceActivation,
			}
			chanSendDataBatchInfo <- sendDataBatchInfo
		}
		// time.Sleep(time.Millisecond * 100)
	}
	close(chanSendDataBatchInfo)
	wg.Wait()
}

func mainold(daddr, appkey, nkey string) {
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
				client.SendData(daddr, appkey, nkey, data, "string", 0)
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
