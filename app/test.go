package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"../handlers"
)

//TestAccess 入网请求测试
func TestAccess(gatewayNum, devNum int) {
	type useTimeInfo struct {
		useTime            int64
		devNum, gatewayNum int
	}
	useTimes := make([]useTimeInfo, 0)
	for times := 0; times < 1; times++ {
		func() {
			// gatewayNum := 15 //测试网关数量
			// devNum := 3000   //测试设备数量
			Macs := make([]string, 0)
			j := gatewayNum + 10000
			for i := 10000; i < j; i++ {
				Macs = append(Macs, fmt.Sprintf("fffb02426f9%d", i))
			}
			useTime, devNum, gatewayNum := TestAccessNetwork("testApp", "testDev", "127.0.0.1:1700", Macs, devNum, gatewayNum)
			useTimes = append(useTimes, useTimeInfo{
				useTime:    useTime,
				devNum:     devNum,
				gatewayNum: gatewayNum,
			})
		}()
		time.Sleep(time.Second * 3)
	}
	fmt.Println(useTimes)
}

//TestAccessNetwork 入网请求测试
//
func TestAccessNetwork(applicationName, devTestName, Addr string, Macs []string, devNum, gatewayNum int) (int64, int, int) {
	if len(Macs) < gatewayNum {
		log.Fatal("测试失败, 传入网关数量小于需求网关数")
	}
	//获取服务器配置列表
	serviceProfilesResult := GetServiceProfiles()
	if len(serviceProfilesResult) == 0 {
		log.Fatal("测试失败,没有服务器配置文件")
	}

	//CreateApplication 创建应用
LOOP:
	applicationsID := CreateApplication(applicationName, serviceProfilesResult[0])
	if applicationsID == "" {
		applicationsResults := GetApplications()
		if len(applicationsResults) == 0 {
			log.Fatal(" not have Application")
		}
		for _, applicationsResult := range applicationsResults {
			if applicationsResult.Name == applicationName {
				err := DelApplication(applicationsResult.ID)
				if err != nil {
					log.Fatal(err)
				}
				goto LOOP
			}
		}
	}
	fmt.Println("创建应用成功, ID:", applicationsID)
	// 批量创建设备
	func() {
		name := 0
		dev := devTestName
		eui := 1000000000000000
		for index := 0; index < devNum; {
			name++
			eui++
			err := PostDevice(fmt.Sprintf("%s%d", dev, name), "????", strconv.Itoa(eui), applicationsID, "ceec88af-86e9-4d3a-a372-693440be67d4", true)
			if err != nil {
				continue
			}
			index++
		}
	}()
	//获取应用内设备列表
	devicesinfos := GetDevices(applicationsID, "", 9999)
	if devNum != len(devicesinfos) {
		fmt.Println("应用内设备数量不足,删除应用重新创建!")
		goto LOOP
	}
	activationDeviceInfos := make([]ActivationDeviceInfo, 0)
	for _, devicesinfo := range devicesinfos {
		deviceKeysResult := DeviceKeysResult{
			DevEUI: devicesinfo.DevEUI,
			NwkKey: "11111111111111111111111111111111",
			AppKey: "11111111111111111111111111111111",
		}
		// 创建设备keys
		err := PostDeviceKeys(deviceKeysResult, "POST")
		if err != nil {
			fmt.Println("应用内设备key数量不足,删除应用重新创建!")
			goto LOOP
		}
		activationDeviceInfo := ActivationDeviceInfo{
			DevEUI: devicesinfo.DevEUI,
			NwkKey: deviceKeysResult.NwkKey,
		}
		activationDeviceInfos = append(activationDeviceInfos, activationDeviceInfo)
	}
	type ClientChan struct {
		client *handlers.GatewayClient
		rxchan chan int
	}
	startPort := 10000
LOOPCREATECLIENT:
	clientChans := make([]ClientChan, 0)
	for i := 0; i < gatewayNum; i++ {
		startPort++
		client, err := handlers.NewClient(fmt.Sprintf(":%d", startPort), Macs[i], Addr)
		if err != nil {
			fmt.Println(err)
			goto LOOPCREATECLIENT
		}
		clientChan := ClientChan{
			client: client,
			rxchan: make(chan int, 1),
		}
		clientChans = append(clientChans, clientChan)
	}
	if gatewayNum != len(clientChans) {
		goto LOOPCREATECLIENT
	}
	// go func() {
	// 	for i, clientChan := range clientChans {
	// 		go func(clientChan ClientChan, i int) {
	// 			for {
	// 				redb := make([]byte, 1024)
	// 				_, err := clientChan.client.Conn.Read(redb)
	// 				if err != nil {
	// 					fmt.Println("read error:", err)
	// 					return
	// 				}
	// 				fmt.Println("read from ", i, " server:", string(redb))
	// 			}
	// 		}(clientChan, i)
	// 	}
	// }()
	//查询数据库是否激活成功
	dbns, err := sql.Open("postgres", "host=127.0.0.1 user=postgres password=loraserver_ns dbname=loraserver_ns sslmode=disable")
	defer dbns.Close()
	if err != nil {
		log.Fatal("Open:", err)
	}
	fmt.Println("=----------------------------------------------------=")
	fmt.Println("will activation device number:", devNum)
	fmt.Println(" use gatewayNum:", gatewayNum)
	time.Sleep(time.Second * 3)
	num := devNum / gatewayNum
	timewg := sync.WaitGroup{}
	timewg.Add(1)
	activationDeviceSuccessNum := 0
	for i, clientChan := range clientChans {
		clientChan.rxchan <- 1
		go func(newClientChan ClientChan, iIdex int) {
			fmt.Println(iIdex)
			timewg.Wait()
			for _, activationDeviceInfo := range activationDeviceInfos[iIdex*num : (iIdex+1)*num] {
				newClientChan.client.JoinRequest(activationDeviceInfo.DevEUI, activationDeviceInfo.NwkKey)
				// time.Sleep(time.Millisecond * 100)
			}
			<-newClientChan.rxchan
		}(clientChan, i)
	}
	timeStart := time.Now().UnixNano() / 1e6
	//等待接受完成
	timewg.Done()
	//等待所有协程结束
	for _, clientChan := range clientChans {
		clientChan.rxchan <- 1
	}
	timeEnd := time.Now().UnixNano() / 1e6
	fmt.Println("=----------------------------------------------------=")
	fmt.Println("end time: ", timeEnd, "ms")
	fmt.Println("wait 3S to check activation infomation")
	time.Sleep(time.Second * 3)
	for _, activationDeviceInfo := range activationDeviceInfos {
		addr, nwkkey, err := GetDeviceActivationKey(dbns, activationDeviceInfo.DevEUI)
		if err != nil {
			fmt.Println("get device keys fail! DevEUI:", activationDeviceInfo.DevEUI)
			continue
		}
		activationDeviceSuccessNum++
		//发送一条数据 正式激活
		clientChans[0].client.SendData(addr, nwkkey, nwkkey, "Activation", "string", 1)
	}
	fmt.Println("=----------------------------------------------------=")
	fmt.Println(" use gatewayNum:", gatewayNum)
	fmt.Println("start time: ", timeStart, "ms")
	fmt.Println("end time: ", timeEnd, "ms")
	fmt.Println("use time : ", timeEnd-timeStart, "ms")
	fmt.Println("activation Device Success Number :", activationDeviceSuccessNum)
	return timeEnd - timeStart, activationDeviceSuccessNum, gatewayNum
}

//SendTADA 测试发送数据
//给应用内所有应用发送信息
// applicationsID 应用ID
func SendTADA(applicationsID string, gatewayNum, deviceNum int) {
	// //获取应用内设备列表
	devicesinfos := GetDevices(applicationsID, "", deviceNum)
	if len(devicesinfos) == 0 {
		return
	}
	if len(devicesinfos) < deviceNum {
		log.Fatal("测试失败,应用内设备不足")
	}
	temperature := 38.0
	humidity := 47.0
	r1 := rand.New(rand.NewSource(time.Now().UnixNano() + 3))
	r2 := rand.New(rand.NewSource(time.Now().UnixNano() + 5))
	ra1 := 0.47
	ra2 := 0.47
	sendDataBatchInfos := make([]SendDataBatchInfo, 0)
	for _, devicesinfo := range devicesinfos[:deviceNum] {
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
		sendDataBatchInfos = append(sendDataBatchInfos, sendDataBatchInfo)
	}
	Macs := make([]string, 0)
	j := gatewayNum + 10000
	for i := 10000; i < j; i++ {
		Macs = append(Macs, fmt.Sprintf("fffb02426f9%d", i))
	}
	startTime, endTime, txNum := SendDataConcurrent(Macs, "127.0.0.1:1701", sendDataBatchInfos, 4356, gatewayNum, 1)
	fmt.Println("startTime:", startTime, "endTime:", endTime, "useTime:", endTime-startTime, "txNum:", txNum)
}

//SendDataConcurrent 批量并发发送信息
//mac 发送网关
//addr 网关地址
// SendDataBatchInfo 信息集合
// wg 同步发送锁
// port 发送起始端口
// goNum 协程数量--模拟同时发送信息网关数数
// infoNum 重复发送信息数
func SendDataConcurrent(mac []string, addr string, chanSendDataBatchInfo []SendDataBatchInfo, port, goNum, infoNum int) (int64, int64, int) {
	if goNum == 0 {
		goNum = 1
	}
	type ClientInfo struct {
		client *handlers.GatewayClient
		rxchan chan int
		txNum  int
	}
	clientInfos := make([]ClientInfo, 0)
	for i := 0; i < goNum; {
		client, err := handlers.NewClient(fmt.Sprintf(":%d", port), mac[i], addr)
		if err != nil {
			fmt.Println(err)
			port++
			continue
		}
		clientInfo := ClientInfo{
			client: client,
			rxchan: make(chan int, 1),
			txNum:  0,
		}
		clientInfos = append(clientInfos, clientInfo)
		i++
		port++
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	num := len(chanSendDataBatchInfo) / goNum
	txNumAll := 0
	fd, err := os.OpenFile(fmt.Sprintf("../rxdata_%d_%d.txt", goNum, num), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("log error", err)
	}
	fd.WriteString(fmt.Sprintf("%d&&%d\n", goNum, num))
	for i := range clientInfos {
		clientInfos[i].rxchan <- 1
		func(iIdex int) {
			// wg.Wait()
			for index := 0; index < infoNum; index++ {
				for _, activationDeviceInfo := range chanSendDataBatchInfo[iIdex*num : (iIdex+1)*num] {
					data, err := clientInfos[iIdex].client.RetSendData(activationDeviceInfo.Devicesinfo.DevAddr, activationDeviceInfo.Devicesinfo.AppSKey, activationDeviceInfo.Devicesinfo.NwkSEncKey, activationDeviceInfo.Data, "string", 0)
					if err != nil {
						continue
					}
					_, err = fd.WriteString(string(data))
					if err != nil {
						continue
					}
					fd.WriteString("\n")
					clientInfos[iIdex].txNum++
					// fmt.Println((buf[:i]))
				}
			}
			<-clientInfos[iIdex].rxchan
		}(i)
	}
	fmt.Println("start to sendData ready:")
	time.Sleep(time.Second * 1)
	timeStart := time.Now().UnixNano() / 1e6
	wg.Done()
	//等待所有协程结束
	for i := range clientInfos {
		clientInfos[i].rxchan <- 1
		txNumAll = txNumAll + clientInfos[i].txNum
		fmt.Println("txNumAll ", i, ": ", clientInfos[i].txNum)
	}
	timeEnd := time.Now().UnixNano() / 1e6
	return timeStart, timeEnd, txNumAll
}
