package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//JWT ..
var JWT string

//GetClient 。。
func GetClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}

}

//SetFixedHearder 设置相同的header
func SetFixedHearder(request *http.Request) {
	request.Header.Add("content-type", "application/json")
	// request.Header.Add("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJleHAiOjE1NDIxNjc1MDIsImlzX2FkbWluIjp0cnVlLCJpc3MiOiJsb3JhLWFwcC1zZXJ2ZXIiLCJuYmYiOjE1NDIwODExMDIsInN1YiI6InVzZXIiLCJ1c2VybmFtZSI6ImFkbWluIn0.0YhZ7ZTUmE_USA8hkDXCnzlSLSbVdLI__0oQVo_mJc0")
	request.Header.Add("grpc-metadata-authorization", JWT)
}

func loginError(err error) {
	log.Fatal("login error :", err)
}

//GetBody 发送请求并且获取请求回复body
func GetBody(url, method string, postBody []byte) ([]byte, error) {
	client := GetClient()
	request, err := http.NewRequest(method, url, strings.NewReader(string(postBody)))
	if err != nil {
		return nil, nil
	}
	SetFixedHearder(request)
	response, err := client.Do(request)
	if err != nil {
		return nil, nil
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil
	}
	return body, nil
}

//Login 登陆
func Login(username, password string) {
	postBody := map[string]string{"password": password, "username": username}
	post, err := json.Marshal(postBody)
	if err != nil {
		log.Fatal(err)
	}
	body, err := GetBody("https://127.0.0.1:8080/api/internal/login", "POST", post)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &postBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(postBody["jwt"])
	JWT = postBody["jwt"]
	// return postBody["jwt"]
}

//ApplicationsResult ..应用信息
type ApplicationsResult struct {
	ID                 string `json:"id"`          //应用id
	Name               string `json:"name"`        //应用名称
	Description        string `json:"description"` //应用描述
	OrganizationID     string `json:"organizationID"`
	ServiceProfileID   string `json:"serviceProfileID"`   //配置id
	ServiceProfileName string `json:"serviceProfileName"` //配置名称
	// Id, Name, Description, OrganizationID, ServiceProfileID, ServiceProfileName string
}

// ApplicationsResponse 。。主体
type ApplicationsResponse struct {
	TotalCount string               `json:"totalCount"` //总计数
	Rsesult    []ApplicationsResult `json:"result"`     //信息
}

//GetApplications 列表列出了可用的应用程序ID
func GetApplications() []string {
	body, err := GetBody("https://127.0.0.1:8080/api/applications?limit=999", "GET", nil)
	if err != nil || body == nil {
		return nil
	}

	var responseMap ApplicationsResponse
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var IDArray []string
	for _, info := range responseMap.Rsesult {
		IDArray = append(IDArray, info.ID)
	}
	return IDArray
}

//DevicesResult ..应用信息
type DevicesResult struct {
	DevEUI                              string `json:"devEUI"`
	Name                                string `json:"name"`
	ApplicationID                       string `json:"applicationID"`
	Description                         string `json:"description"`
	DeviceProfileID                     string `json:"deviceProfileID"`
	DeviceProfileName                   string `json:"deviceProfileName"`
	DeviceStatusBattery                 int    `json:"deviceStatusBattery"`
	DeviceStatusMargin                  int    `json:"deviceStatusMargin"`
	DeviceStatusExternalPowerSource     bool   `json:"deviceStatusExternalPowerSource"`
	DeviceStatusBatteryLevelUnavailable bool   `json:"deviceStatusBatteryLevelUnavailable"`
	DeviceStatusBatteryLevel            int    `json:"deviceStatusBatteryLevel"`
	LastSeenAt                          string `json:"lastSeenAt"`
}

// DevicesResponse 。。主体
type DevicesResponse struct {
	TotalCount string          `json:"totalCount"` //总计数
	Rsesult    []DevicesResult `json:"result"`     //信息
}

//GetDevices 获取设备信息
//id 要查询的应用id
//return []DevicesResult 返回查询到的应用设备全部信息
func GetDevices(id string) []DevicesResult {
	body, err := GetBody("https://127.0.0.1:8080/api/devices?limit=999&applicationID="+id, "GET", nil)
	if err != nil || body == nil {
		return nil
	}

	var responseMap DevicesResponse
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(responseMap)
	return responseMap.Rsesult
}

//PostBody Post数据 包括创建设备的基本信息
//deviceInfo map[string]string
//deviceInfo  name 设备名称
//deviceInfo  description 设备说明
//deviceInfo  devEUI 设备EUI
//deviceInfo  applicationID 设备所属应用
//deviceInfo  deviceProfileID 配置文件ID
//deviceInfo  skipFCntCheck 是否禁用帧计数器验证
type PostBody struct {
	Device map[string]interface{} `json:"device"`
}

//PostDevice 创建给定的设备
//name 设备名称
//description 设备说明
//devEUI 设备EUI
//applicationID 设备所属应用
//deviceProfileID 配置文件ID
//skipFCntCheck 是否禁用帧计数器验证
//return bool
func PostDevice(name, description, devEUI, applicationID, deviceProfileID string, skipFCntCheck bool) error {
	var post PostBody
	post.Device = make(map[string]interface{})
	post.Device["name"] = name
	post.Device["description"] = description
	post.Device["devEUI"] = devEUI
	post.Device["applicationID"] = applicationID
	post.Device["deviceProfileID"] = deviceProfileID
	post.Device["skipFCntCheck"] = skipFCntCheck
	buf, err := json.Marshal(post)
	if err != nil {
		return errors.New("device info error")
	}
	body, err := GetBody("https://127.0.0.1:8080/api/devices", "POST", buf)
	if err != nil || body == nil {
		return err
	}
	fmt.Println(string(body))
	var responseMap DevicesResponse
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(responseMap)
	return nil
}

//DeviceKeysResult ..keys(OTAA)信息
type DeviceKeysResult struct {
	// DevEUI string `json:"devEUI"`
	NwkKey string `json:"nwkKey"`
	// AppKey string `json:"appKey"`
}

// DeviceKeysResponse 。。keys(OTAA)主体
type DeviceKeysResponse struct {
	DeviceKeys DeviceKeysResult `json:"deviceKeys"` //信息
}

//GetDeviceKeys 获取设备keys(OTAA)
//devEUI 设备devEUI
//return string Application key
func GetDeviceKeys(devEUI string) string {
	var responseMap DeviceKeysResponse
	body, err := GetBody("https://127.0.0.1:8080/api/devices/"+devEUI+"/keys", "GET", nil)
	if err != nil || body == nil {
		return ""
	}

	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(responseMap)
	return responseMap.DeviceKeys.NwkKey
}

//DeviceActivationResult ..Activation信息
type DeviceActivationResult struct {
	// DevEUI      string `json:"devEUI"`
	DevAddr    string `json:"devAddr"`
	AppSKey    string `json:"appSKey"`
	NwkSEncKey string `json:"nwkSEncKey"`
	// SNwkSIntKey string `json:"sNwkSIntKey"`
	// FNwkSIntKey string `json:"fNwkSIntKey"`
	// FCntUp      int    `json:"fCntUp"`
	// NFCntDown   int    `json:"nFCntDown"`
	// AFCntDown   int    `json:"aFCntDown"`
}

// DeviceActivationResponse 。。Activation主体
type DeviceActivationResponse struct {
	DeviceActivation DeviceActivationResult `json:"deviceActivation"` //信息
}

//GetDeviceActivation 获取设备keys(OTAA)
//devEUI 设备devEUI
func GetDeviceActivation(devEUI string) (DeviceActivationResult, error) {
	var responseMap DeviceActivationResponse
	body, err := GetBody("https://127.0.0.1:8080/api/devices/"+devEUI+"/activation", "GET", nil)
	if err != nil || body == nil {
		return responseMap.DeviceActivation, errors.New("GetBody error")
	}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return responseMap.DeviceActivation, err
	}
	fmt.Println(responseMap)
	return responseMap.DeviceActivation, err
}
