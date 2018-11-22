package main

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"../handlers"
	_ "github.com/lib/pq"
)

//GetDeviceActivationKey 获取激活后的密钥
func GetDeviceActivationKey(db *sql.DB, devEUI string) (string, string, error) {
	q := "\\x" + devEUI
	// fmt.Println("select dev_eui,dev_addr,nwk_s_enc_key from device_activation where dev_eui = cast('" + q + "' as bytea);")
	rows, err := db.Query("select dev_eui,dev_addr,nwk_s_enc_key from device_activation where dev_eui = cast('" + q + "' as bytea) order by created_at desc limit 1;")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	for rows.Next() {
		var devEui []uint8
		var devAddr, nwkSEncKey []uint8
		err = rows.Scan(&devEui, &devAddr, &nwkSEncKey)
		if err != nil {
			log.Fatal(err)
		}
		// encodedStr := hex.EncodeToString(devEui)
		// DeviceKeysInfo[encodedStr].Askey = encodedStr
		// DeviceKeysInfo[encodedStr].DevAddr = hex.EncodeToString(devAddr)
		// DeviceKeysInfo[encodedStr].Nskey = hex.EncodeToString(nwkSEncKey)
		// fmt.Println(encodedStr)
		return hex.EncodeToString(devAddr), hex.EncodeToString(nwkSEncKey), nil
	}
	return "", "", errors.New("get fail")
}

//ActivationDeviceInfo 要激活的设备信息
type ActivationDeviceInfo struct {
	DevEUI string
	NwkKey string
	Mac    string
	Addr   string
}

//DeviceKeyInfo ..
type DeviceKeyInfo struct {
	DevEUI, NwkKey, AppSKey, DevAddr string
}

//SendActivationDevice 发送激活申请
//activationDeviceInfos 激活信息
//startPort 发送端口
func SendActivationDevice(activationDeviceInfos []ActivationDeviceInfo, startPort int) []DeviceKeyInfo {
	deviceKeyInfos := make([]DeviceKeyInfo, 0)
	dbns, err := sql.Open("postgres", "host=127.0.0.1 user=postgres password=loraserver_ns dbname=loraserver_ns sslmode=disable")
	if err != nil {
		log.Fatal("Open:", err)
	}
	defer dbns.Close()
	for _, activationDeviceInfo := range activationDeviceInfos {
		client, err := handlers.NewClient(fmt.Sprintf(":%d", startPort), activationDeviceInfo.Mac, activationDeviceInfo.Addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//发送激活请求
		err = client.JoinRequest(activationDeviceInfo.DevEUI, activationDeviceInfo.NwkKey)
		if err != nil {
			fmt.Println("ActivationDevice fail:", activationDeviceInfo)
		}
		time.Sleep(time.Second * 1)
		fmt.Println("get keys", activationDeviceInfo.DevEUI)
		addr, nwkkey, err := GetDeviceActivationKey(dbns, activationDeviceInfo.DevEUI)
		if err != nil {
			fmt.Println("get device keys fail! DevEUI:", activationDeviceInfo.DevEUI)
			continue
		}
		deviceKeyInfo := DeviceKeyInfo{
			DevEUI:  activationDeviceInfo.DevEUI,
			NwkKey:  nwkkey,
			DevAddr: addr,
		}
		fmt.Println(addr, nwkkey)
		//发送一条数据 正式激活
		client.SendData(addr, nwkkey, nwkkey, "Activation", "string", 1)
		deviceKeyInfos = append(deviceKeyInfos, deviceKeyInfo)
	}
	return deviceKeyInfos
}

//ActivationDevices  激活应用内所有应用
//ApplicationsID 要激活的应用ID
func ActivationDevices(ApplicationsID, Mac, Addr string) {
	//批量激活
	activationDeviceInfos := make([]ActivationDeviceInfo, 0)
	//获取应用内设备列表
	devicesinfos := GetDevices(ApplicationsID, "", 9999)
	// fmt.Println(devicesinfos)
	activationDeviceInfoNum := 10
	chanActivationDeviceInfo := make(chan ActivationDeviceInfo, activationDeviceInfoNum)
	//异步创建keys
	go func() {
		numfail := 0
		for _, devicesinfo := range devicesinfos {
			// fmt.Println(devicesinfo.DevEUI)
			//获取设备keys 如果已经设置则跳过设置
			deviceKeysResult, err := GetDeviceKeys(devicesinfo.DevEUI)
			if err != nil {
				numfail++
				deviceKeysResult = DeviceKeysResult{
					DevEUI: devicesinfo.DevEUI,
					NwkKey: "11111111111111111111111111111111",
					AppKey: "11111111111111111111111111111111",
				}
				//创建设备keys
				PostDeviceKeys(deviceKeysResult, "POST")
			}
			activationDeviceInfo := ActivationDeviceInfo{
				DevEUI: devicesinfo.DevEUI,
				NwkKey: deviceKeysResult.NwkKey,
			}
			chanActivationDeviceInfo <- activationDeviceInfo
		}
		//创建完成 关闭管道
		close(chanActivationDeviceInfo)
	}()
	rangActivationDeviceInfoNum := 1
	activationDevicePort := 5432
	wg := sync.WaitGroup{}
	for activationDeviceInfo := range chanActivationDeviceInfo {
		activationDeviceInfos = append(activationDeviceInfos, activationDeviceInfo)
		if rangActivationDeviceInfoNum%101 == 0 {
			go func(rangActivationDeviceInfoNum int) {
				wg.Add(1)
				SendActivationDeviceSameApplications(activationDeviceInfos[rangActivationDeviceInfoNum-101:rangActivationDeviceInfoNum], activationDevicePort, Mac, Addr)
				wg.Done()
			}(rangActivationDeviceInfoNum)
			activationDevicePort++
		}
		rangActivationDeviceInfoNum++
	}
	//管道被关闭, 激活剩余部分
	SendActivationDeviceSameApplications(activationDeviceInfos[rangActivationDeviceInfoNum-rangActivationDeviceInfoNum%101:], activationDevicePort, Mac, Addr)
	//等待所有协程结束
	wg.Wait()
}

//SendActivationDeviceSameApplications 同一网关批量激活
//activationDeviceInfos 激活信息
//startPort 发送端口
func SendActivationDeviceSameApplications(activationDeviceInfos []ActivationDeviceInfo, startPort int, Mac, Addr string) []DeviceKeyInfo {
	deviceKeyInfos := make([]DeviceKeyInfo, 0)
	client, err := handlers.NewClient(fmt.Sprintf(":%d", startPort), Mac, Addr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	dbns, err := sql.Open("postgres", "host=127.0.0.1 user=postgres password=loraserver_ns dbname=loraserver_ns sslmode=disable")
	if err != nil {
		log.Fatal("Open:", err)
	}
	defer dbns.Close()
	for _, activationDeviceInfo := range activationDeviceInfos {
		//发送激活请求
		err = client.JoinRequest(activationDeviceInfo.DevEUI, activationDeviceInfo.NwkKey)
		if err != nil {
			fmt.Println("ActivationDevice fail:", activationDeviceInfo)
		}
		time.Sleep(time.Second * 1)
		fmt.Println("get keys", activationDeviceInfo.DevEUI)
		addr, nwkkey, err := GetDeviceActivationKey(dbns, activationDeviceInfo.DevEUI)
		if err != nil {
			fmt.Println("get device keys fail! DevEUI:", activationDeviceInfo.DevEUI)
			continue
		}
		deviceKeyInfo := DeviceKeyInfo{
			DevEUI:  activationDeviceInfo.DevEUI,
			NwkKey:  nwkkey,
			DevAddr: addr,
		}
		fmt.Println(addr, nwkkey)
		//发送一条数据 正式激活
		client.SendData(addr, nwkkey, nwkkey, "Activation", "string", 1)
		deviceKeyInfos = append(deviceKeyInfos, deviceKeyInfo)
	}
	return deviceKeyInfos
}

//CreateBatchDevicesInfo 批量创建的设备信息
type CreateBatchDevicesInfo struct {
	PplicationsID, DeviceName, DeviceeEUI, Synopsis, DeviceProfile string
	DisValidation                                                  bool
}

//CreateBatchDevices 批量创建设备
//chanNum 协程数
//createBatchDevicesInfos 设备信息
func CreateBatchDevices(chanNum int, createBatchDevicesInfos []CreateBatchDevicesInfo) []CreateBatchDevicesInfo {
	chanRES := make(chan int, chanNum)
	wg := sync.WaitGroup{}
	createBatchDevicesInfoError := make([]CreateBatchDevicesInfo, 0)
	for _, createBatchDevicesInfo := range createBatchDevicesInfos {
		chanRES <- 1
		wg.Add(1)
		go func(createBatchDevicesInfo CreateBatchDevicesInfo) {
			err := PostDevice(createBatchDevicesInfo.DeviceName, createBatchDevicesInfo.Synopsis, createBatchDevicesInfo.DeviceeEUI, createBatchDevicesInfo.PplicationsID, createBatchDevicesInfo.DeviceProfile, createBatchDevicesInfo.DisValidation)
			if err != nil {
				createBatchDevicesInfoError = append(createBatchDevicesInfoError, createBatchDevicesInfo)
			}
			<-chanRES
			wg.Done()
		}(createBatchDevicesInfo)
	}
	wg.Wait()
	return createBatchDevicesInfoError
}

//SendDataBatch 批量发送信息
//mac 发送网关
//addr 网关地址
// chanSendDataBatchInfo 信息管道
// wg 等待锁
// port 发送起始端口
// goNum 协程数量--模拟同时发送信息设备数
func SendDataBatch(mac, addr string, chanSendDataBatchInfo chan SendDataBatchInfo, wg *sync.WaitGroup, port, goNum int) {
	if goNum == 0 {
		goNum = 1
	}
	type ClientInfo struct {
		client *handlers.GatewayClient
		rxchan chan int
	}
	clientInfos := make([]ClientInfo, 0)
	for i := 0; i < goNum; i++ {
		client, err := handlers.NewClient(fmt.Sprintf(":%d", port), mac, addr)
		if err != nil {
			fmt.Println(err)
			port++
			continue
		}
		clientInfo := ClientInfo{
			client: client,
			rxchan: make(chan int, 1),
		}
		clientInfos = append(clientInfos, clientInfo)
		port++
	}
	num := 0
	nums := 0
	numss := 0
	for sendDataBatchInfo := range chanSendDataBatchInfo {
		nums++
		for {
			select {
			case clientInfos[num].rxchan <- 1:
				numss++
				fmt.Println("!!!", num)
				go func(clientInfo ClientInfo) {
					clientInfo.client.SendData(sendDataBatchInfo.Devicesinfo.DevAddr, sendDataBatchInfo.Devicesinfo.AppSKey, sendDataBatchInfo.Devicesinfo.NwkSEncKey, sendDataBatchInfo.Data, "string", 0)
					// time.Sleep(time.Second * 1)
					<-clientInfo.rxchan
				}(clientInfos[num])
				num++
				if num == goNum {
					num = 0
				}
			default:
				num++
				if num == goNum {
					num = 0
				}
			}
		}
	}
	for _, sclientInfo := range clientInfos {
		sclientInfo.rxchan <- 1
	}
	time.Sleep(time.Second * 3)
	fmt.Println("send data numbers : ", nums, "  numss : ", numss)
	wg.Done()
	// num := len(chanSendDataBatchInfo) / goNum
	// for i, clientChan := range clientInfos {
	// 	clientChan.rxchan <- 1
	// 	go func(newClientChan ClientInfo, iIdex int) {
	// 		fmt.Println(iIdex)
	// 		wg.Wait()
	// 		for _, activationDeviceInfo := range chanSendDataBatchInfo[iIdex*num : (iIdex+1)*num] {
	// 			newClientChan.client.SendData(activationDeviceInfo.Devicesinfo.DevAddr, activationDeviceInfo.Devicesinfo.AppSKey, activationDeviceInfo.Devicesinfo.NwkSEncKey, activationDeviceInfo.Data, "string", 0)
	// 			// time.Sleep(time.Millisecond * 100)
	// 		}
	// 		<-newClientChan.rxchan
	// 	}(clientChan, i)
	// }
	// //等待所有协程结束
	// for _, clientChan := range clientInfos {
	// 	clientChan.rxchan <- 1
	// }
}
