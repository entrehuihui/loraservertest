package main

import (
	"crypto/aes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"encoding/base64"

	"../lorawan"
	"github.com/jacobsa/crypto/cmac"
	log "github.com/sirupsen/logrus"
)

func decryptFRMPayload(key string, uplink bool, devAddr string, fcnt uint32, data string) {
	var aesKey lorawan.AES128Key
	var deviceAddr lorawan.DevAddr
	err := aesKey.UnmarshalText([]byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = deviceAddr.UnmarshalText([]byte(devAddr))
	if err != nil {
		fmt.Println(err)
		return
	}
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	decryptData, err := encryptFRMPayload(aesKey, uplink, deviceAddr, fcnt, dataBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%X\n", decryptData)
}

func main() {
	// if err := unmarshalJoinrequest("EF88B34A06DF195E3F249C321140DE1B", "AABAGSQBxSYsWG0AJAF3SgAnB2drIY4="); err != nil {
	// 	log.Println("unmarshalJoinrequest error:", err)
	// 	return
	// }
	unmarshalPhypayload("gH7JP2aAdgACoNP9UoOss3wJ")
	//decryptFRMPayload("fa2832248c0e8f9c72c6eb989f8d2fd8", false, "663fc97e", 0, "QBEAAADAwhIJ03P9wqXFEux2JXA=")
}

func unmarshalJoinrequest(appKey string, b64Data string) error {
	fmt.Println("入网请求原始base64数据:", b64Data)
	fmt.Println("appkey:", appKey)
	var appKey_bytes lorawan.AES128Key
	err := appKey_bytes.UnmarshalText([]byte(appKey))
	if err != nil {
		return err
	}

	var phyPayload lorawan.PHYPayload
	err = phyPayload.UnmarshalText([]byte(b64Data))
	if err != nil {
		return err
	}

	_, ok := phyPayload.MACPayload.(*lorawan.JoinRequestPayload)
	if !ok {
		return errors.New("not join request")
	}
	js, err := json.MarshalIndent(phyPayload, " ", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(js))
	mic, err := calculateJoinRequestMIC(appKey_bytes, phyPayload)
	if err != nil {
		return nil
	}
	fmt.Printf("计算出来的MIC:%x\n", mic)
	return nil
}

func unmarshalMacpayload(appSessionKey string, b64Data string) {
	var appKey lorawan.AES128Key

	err := appKey.UnmarshalText([]byte(appSessionKey))
	if err != nil {
		log.Fatal("unmarshal test error:", err)
	}

	var phyPayload lorawan.PHYPayload
	err = phyPayload.UnmarshalText([]byte(b64Data))
	if err != nil {
		log.Fatal("UnmarshalBinary data error:", err)
	}

	macPayload, ok := phyPayload.MACPayload.(*lorawan.MACPayload)
	if !ok {
		log.Println("not mac paylaod!")
		return
	}
	js, err := json.MarshalIndent(macPayload, " ", " ")
	if err == nil {
		fmt.Println(string(js))
	}
	for _, payload := range macPayload.FRMPayload {
		if dataPayload, ok := payload.(*lorawan.DataPayload); ok {

			decryB, err := encryptFRMPayload(appKey, false, macPayload.FHDR.DevAddr, macPayload.FHDR.FCnt, dataPayload.Bytes)
			if err != nil {
				log.Error("decrypt data error:", err)
				break
			}
			fmt.Printf("%x\n", decryB)
		}
	}
}

func encryptFRMPayload(key lorawan.AES128Key, uplink bool, devAddr lorawan.DevAddr, fCnt uint32, data []byte) ([]byte, error) {
	pLen := len(data)
	if pLen%16 != 0 {
		// append with empty bytes so that len(data) is a multiple of 16
		data = append(data, make([]byte, 16-(pLen%16))...)
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	if block.BlockSize() != 16 {
		return nil, errors.New("lorawan: block size of 16 was expected")
	}

	s := make([]byte, 16)
	a := make([]byte, 16)
	a[0] = 0x01
	if !uplink {
		a[5] = 0x01
	}

	b, err := devAddr.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(a[6:10], b)
	binary.LittleEndian.PutUint32(a[10:14], uint32(fCnt))

	for i := 0; i < len(data)/16; i++ {
		a[15] = byte(i + 1)
		block.Encrypt(s, a)
		for j := 0; j < len(s); j++ {
			data[i*16+j] = data[i*16+j] ^ s[j]
		}
	}

	return data[0:pLen], nil
}

func calculateJoinRequestMIC(key lorawan.AES128Key, p lorawan.PHYPayload) ([]byte, error) {
	if p.MACPayload == nil {
		return []byte{}, errors.New("lorawan: MACPayload should not be empty")
	}
	jrPayload, ok := p.MACPayload.(*lorawan.JoinRequestPayload)
	if !ok {
		return []byte{}, errors.New("lorawan: MACPayload should be of type *JoinRequestPayload")
	}

	micBytes := make([]byte, 0, 19)

	b, err := p.MHDR.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}
	micBytes = append(micBytes, b...)

	b, err = jrPayload.MarshalBinary()
	if err != nil {
		return nil, err
	}
	micBytes = append(micBytes, b...)

	hash, err := cmac.New(key[:])
	if err != nil {
		return []byte{}, err
	}
	if _, err = hash.Write(micBytes); err != nil {
		return nil, err
	}
	hb := hash.Sum([]byte{})
	if len(hb) < 4 {
		return []byte{}, errors.New("lorawan: the hash returned less than 4 bytes")
	}
	return hb[0:4], nil
}

func unmarshalPhypayload(data string) {
	var payload lorawan.PHYPayload
	err := payload.UnmarshalText([]byte(data))
	if err != nil {
		log.Println("unmarshal phypayload error:", err)
		return
	}
	jsData, err := json.MarshalIndent(payload, " ", " ")
	if err != nil {
		log.Println("encode json error:", err)
		return
	}
	fmt.Println(string(jsData))
}
