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

//URL lora服务器地址
var URL = "https://127.0.0.1:8080"

//PostErrorResponse ..
type PostErrorResponse struct {
	Error string `json:"error"`
	ID    string `json:"id"`
	// Message string `json:"message"`
	// Code    int    `json:"code"`
}

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
	body, err := GetBody(URL+"/api/internal/login", "POST", post)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &postBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(postBody["jwt"])
	if postBody["jwt"] == "" {
		log.Fatal("login fail")
	}
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
func GetApplications() []ApplicationsResult {
	body, err := GetBody(URL+"/api/applications?limit=9999", "GET", nil)
	if err != nil || body == nil {
		return nil
	}

	var responseMap ApplicationsResponse
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return responseMap.Rsesult
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
//search 查询相关的设备 全部设置传""
//return []DevicesResult 返回查询到的应用设备全部信息
func GetDevices(id string, search string) []DevicesResult {
	body, err := GetBody(URL+"/api/devices?limit=9999&applicationID="+id+"&search="+search, "GET", nil)
	if err != nil || body == nil {
		return nil
	}

	var responseMap DevicesResponse
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(responseMap.TotalCount)
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
	body, err := GetBody(URL+"/api/devices", "POST", buf)
	if err != nil || body == nil {
		return err
	}
	var postResponseMap PostErrorResponse
	err = json.Unmarshal(body, &postResponseMap)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if postResponseMap.Error == "" {
		fmt.Println("create device success!\n", "device name:", post.Device["name"], "\ndevice devEUI:", post.Device["devEUI"])
		return nil
	}
	fmt.Println("create device fail:" + postResponseMap.Error)
	return errors.New("create device fail:" + postResponseMap.Error)
}

//DeviceKeysResult ..keys(OTAA)信息
type DeviceKeysResult struct {
	DevEUI string `json:"devEUI"`
	NwkKey string `json:"nwkKey"`
	AppKey string `json:"appKey"`
}

// DeviceKeysResponse 。。keys(OTAA)主体
type DeviceKeysResponse struct {
	DeviceKeys DeviceKeysResult `json:"deviceKeys"` //信息
}

//GetDeviceKeys 获取设备keys(OTAA)
//devEUI 设备devEUI
//return string Application key, if error return ""
func GetDeviceKeys(devEUI string) (DeviceKeysResult, error) {
	var responseMap DeviceKeysResponse
	body, err := GetBody(URL+"/api/devices/"+devEUI+"/keys", "GET", nil)
	if err != nil || body == nil {
		return DeviceKeysResult{}, err
	}

	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return DeviceKeysResult{}, err
	}
	if responseMap.DeviceKeys.NwkKey == "" {
		// fmt.Println("get device key fail")
		return DeviceKeysResult{}, errors.New("get device key fail")
	}
	fmt.Println(responseMap.DeviceKeys)
	return responseMap.DeviceKeys, nil
}

// //DeviceKeysinfo ..post信息
// type DeviceKeysinfo struct {
// 	DevEUI string `json:"devEUI"`
// 	NwkKey string `json:"nwkKey"`
// 	AppKey string `json:"appKey"`
// }

//PostDeviceKeys 设置设备keys(OTAA)
//deviceKeys keys信息
//medhod 请求方式 post创建 put修改
func PostDeviceKeys(deviceKeys DeviceKeysResult, medhod string) error {
	type PostDeviceKeys struct {
		DeviceKeys DeviceKeysResult `json:"deviceKeys"`
	}
	postDeviceKeys := PostDeviceKeys{DeviceKeys: deviceKeys}
	buf, err := json.Marshal(postDeviceKeys)
	body, err := GetBody(URL+"/api/devices/"+deviceKeys.DevEUI+"/keys", medhod, buf)
	if err != nil || body == nil {
		return err
	}
	responseMap := PostErrorResponse{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return err
	}
	if responseMap.Error != "" {
		return errors.New("create device keys fail:" + responseMap.Error)
	}
	// fmt.Println("create keys success")
	return nil
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
//devEUI 设备devEUI, if error return "",err
func GetDeviceActivation(devEUI string) (DeviceActivationResult, error) {
	var responseMap DeviceActivationResponse
	body, err := GetBody(URL+"/api/devices/"+devEUI+"/activation", "GET", nil)
	if err != nil || body == nil {
		return responseMap.DeviceActivation, errors.New("GetBody error")
	}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return responseMap.DeviceActivation, err
	}
	if responseMap.DeviceActivation.DevAddr == "" {
		fmt.Println("get device activation fail")
		return responseMap.DeviceActivation, errors.New("get device activation fail")
	}
	// fmt.Println(responseMap)
	return responseMap.DeviceActivation, err
}

//SendDataBatchInfo ..
type SendDataBatchInfo struct {
	Data        string
	Devicesinfo DeviceActivationResult
}

//ServiceProfilesResult ..配置信息
type ServiceProfilesResult struct {
	ID              string `json:"id"`   //配置id
	Name            string `json:"name"` //配置名称
	OrganizationID  string `json:"organizationID"`
	NetworkServerID string `json:"networkServerID"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// ServiceProfilesResponse 。。配置信息主体
type ServiceProfilesResponse struct {
	TotalCount string               `json:"totalCount"` //总计数
	Rsesult    []ApplicationsResult `json:"result"`     //信息
}

//GetServiceProfiles 获取服务器配置列表
func GetServiceProfiles() []ApplicationsResult {
	body, err := GetBody(URL+"/api/service-profiles?limit=999", "GET", nil)
	if err != nil || body == nil {
		return nil
	}
	var responseMap ServiceProfilesResponse
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return responseMap.Rsesult
}

//CreateApplication 创建应用
func CreateApplication(applicationName string, applicationsResult ApplicationsResult) string {
	type ApplicationsInfo struct {
		Description      string `json:"description"`
		Name             string `json:"name"`
		OrganizationID   string `json:"organizationID"`
		ServiceProfileID string `json:"serviceProfileID"`
	}
	type ApplicationsInfoPost struct {
		Application ApplicationsInfo `json:"application"`
	}
	applicationsInfoPost := ApplicationsInfoPost{}
	applicationsInfoPost.Application = ApplicationsInfo{
		Description:      applicationName,
		Name:             applicationName,
		OrganizationID:   applicationsResult.OrganizationID,
		ServiceProfileID: applicationsResult.ID,
	}
	buf, err := json.Marshal(applicationsInfoPost)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	body, err := GetBody(URL+"/api/applications", "POST", buf)
	if err != nil || body == nil {
		fmt.Println(err, string(body))
		return ""
	}
	var postResponseMap PostErrorResponse
	err = json.Unmarshal(body, &postResponseMap)
	if err != nil {
		fmt.Println("创建测试应用失败!!! error:", err, "returnError", postResponseMap.Error)
		return ""
	}
	if postResponseMap.Error != "" {
		fmt.Println("创建测试应用失败!!! returnError", postResponseMap.Error)
		return ""
	}
	return postResponseMap.ID
}

//DelApplication 删除应用
func DelApplication(ID string) error {
	body, err := GetBody(URL+"/api/applications/"+ID, "DELETE", nil)
	if err != nil || body == nil {
		return nil
	}
	var responseMap PostErrorResponse
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return err
	}
	if responseMap.Error != "" {
		return errors.New(responseMap.Error)
	}
	return nil
}
