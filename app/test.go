package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"../handlers"
)

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
	fmt.Println(applicationsID)
	// 批量创建设备
	func() {
		name := 0
		dev := devTestName
		eui := 1000000000000000
		createBatchDevicesInfos := make([]CreateBatchDevicesInfo, 0)
		for index := 0; index < devNum; index++ {
			// err := PostDevice(strconv.Itoa(name), "大大萨达萨达", strconv.Itoa(eui), IDArray[0], "ceec88af-86e9-4d3a-a372-693440be67d4", true)
			createBatchDevicesInfo := CreateBatchDevicesInfo{
				PplicationsID: applicationsID,
				DeviceName:    fmt.Sprintf("%s%d", dev, name),
				DeviceeEUI:    strconv.Itoa(eui),
				Synopsis:      "????",
				DeviceProfile: "ceec88af-86e9-4d3a-a372-693440be67d4",
				DisValidation: true,
			}
			createBatchDevicesInfos = append(createBatchDevicesInfos, createBatchDevicesInfo)
			name++
			eui++
		}
		// CreateBatchDevices 批量创建设备
		CreateBatchDevices(gatewayNum, createBatchDevicesInfos)
	}()
	activationDeviceInfos := make([]ActivationDeviceInfo, 0)
	//获取应用内设备列表
	devicesinfos := GetDevices(applicationsID, "")
	for _, devicesinfo := range devicesinfos {
		deviceKeysResult := DeviceKeysResult{
			DevEUI: devicesinfo.DevEUI,
			NwkKey: "11111111111111111111111111111111",
			AppKey: "11111111111111111111111111111111",
		}
		//创建设备keys
		err := PostDeviceKeys(deviceKeysResult, "POST")
		if err != nil {
			continue
		}
		activationDeviceInfo := ActivationDeviceInfo{
			DevEUI: devicesinfo.DevEUI,
			NwkKey: deviceKeysResult.NwkKey,
		}
		activationDeviceInfos = append(activationDeviceInfos, activationDeviceInfo)
	}
	devNum = len(activationDeviceInfos)
	type ClientChan struct {
		client *handlers.GatewayClient
		rxchan chan int
	}
	startPort := 10000
	clientChans := make([]ClientChan, 0)
	for i := 0; i < gatewayNum; i++ {
		startPort++
		client, err := handlers.NewClient(fmt.Sprintf(":%d", startPort), Macs[i], Addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		clientChan := ClientChan{
			client: client,
			rxchan: make(chan int, 1),
		}
		clientChans = append(clientChans, clientChan)
	}
	clientChanNum := 0
	gatewayNum = len(clientChans)
	go func() {
		for _, clientChan := range clientChans {
			go func(clientChan ClientChan) {
				redb := make([]byte, 1024)
				n, err := clientChan.client.Conn.Read(redb)
				if err != nil {
					fmt.Println("read error:", err)
					return
				}
				if n > 12 {
					fmt.Printf("read from server:%x\n", redb[12:n])
				} else {
					fmt.Printf("read from server:%x\n", redb[:n])
				}
			}(clientChan)
		}
	}()
	time.Sleep(time.Second * 3)
	fmt.Println("=----------------------------------------------------=")
	fmt.Println("will activation device number:", devNum)
	fmt.Println(" use gatewayNum:", gatewayNum)
	timeStart := time.Now().UnixNano() / 1e6
	for _, activationDeviceInfo := range activationDeviceInfos {
		select {
		case clientChans[clientChanNum].rxchan <- 1:
			go func(num int, newActivationDeviceInfo ActivationDeviceInfo) {
				//发送激活请求
				err := clientChans[num].client.JoinRequest(newActivationDeviceInfo.DevEUI, newActivationDeviceInfo.NwkKey)
				if err != nil {
					fmt.Println("ActivationDevice fail:", newActivationDeviceInfo)
				}
				<-clientChans[num].rxchan
			}(clientChanNum, activationDeviceInfo)
			clientChanNum++
			if clientChanNum == gatewayNum {
				clientChanNum = 0
			}
		}
	}
	//等待所有协程结束
	for _, clientChan := range clientChans {
		clientChan.rxchan <- 1
	}
	timeEnd := time.Now().UnixNano() / 1e6
	fmt.Println("=----------------------------------------------------=")
	fmt.Println("will activation device number:", devNum)
	fmt.Println(" use gatewayNum:", gatewayNum)
	fmt.Println("end time: ", timeEnd, "ms")
	fmt.Println("wait 3S to check activation infomation")
	time.Sleep(time.Second * 15)
	//查询数据库是否激活成功
	dbns, err := sql.Open("postgres", "host=127.0.0.1 user=postgres password=loraserver_ns dbname=loraserver_ns sslmode=disable")
	if err != nil {
		log.Fatal("Open:", err)
	}
	activationDeviceSuccessNum := 0
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
