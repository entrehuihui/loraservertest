package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"../lorawan"
	"../models/gateway"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//"QAUAEXYAvwECHflfQwumoqY1LtPxLQUhWwkeNNx2RWgVpm0+AvFtARr30xk+LIGrSnSsuwiCxcBXT7PUju/avOUU/bkY4FkQAVs4CA5g2CmKZcVFsFuJ9fsPeNhi"
func encode(devaddr lorawan.DevAddr, askey, nskey lorawan.AES128Key, data []byte, count uint32) ([]byte, error) {
	var port uint8
	port = 2
	payload := lorawan.PHYPayload{
		MHDR: lorawan.MHDR{
			MType: lorawan.UnconfirmedDataUp,
		},
		MACPayload: &lorawan.MACPayload{
			FHDR: lorawan.FHDR{
				DevAddr: devaddr,
				FCtrl: lorawan.FCtrl{
					ADR:       false,
					ADRACKReq: false,
					ACK:       false,
					FPending:  false,
					ClassB:    false,
				},
				FCnt:  count,
				FOpts: nil,
			},
			FPort: &port,
			FRMPayload: []lorawan.Payload{&lorawan.DataPayload{
				Bytes: data,
			}},
		},
	}

	if err := payload.EncryptFRMPayload(askey); err != nil {
		return nil, err
	}
	payload.SetMIC(nskey)

	return payload.MarshalBinary()
}

type GatewayClient struct {
	Conn net.Conn
	Mac  lorawan.EUI64
}

func NewClient(laddr, mac, addr string) (*GatewayClient, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	var gatewaymac lorawan.EUI64
	if err := gatewaymac.UnmarshalText([]byte(mac)); err != nil {
		return nil, err
	}
	return &GatewayClient{
		Conn: conn,
		Mac:  gatewaymac,
	}, nil
}

func (c *GatewayClient) Ping() {
	for {
		//fmt.Println("Test ping strat...")
		pullData := gateway.PullDataPacket{
			ProtocolVersion: 2,
			RandomToken:     1,
			GatewayMAC:      [8]byte(c.Mac),
		}
		wData, _ := pullData.MarshalBinary()
		fmt.Printf("ping data:%x\n", wData)
		_, err := c.Conn.Write(wData)
		if err != nil {
			log.Error("gateway client ping error:", err)
		}
		time.Sleep(time.Second * 30)
	}
}

func (c *GatewayClient) SendData(devAddr, askey, nskey, str string, dataType string, count uint32) error {
	var devaddr lorawan.DevAddr
	err := devaddr.UnmarshalText([]byte(devAddr))
	if err != nil {
		return errors.Wrap(err, "Device address error")
	}
	var asKey lorawan.AES128Key
	var nsKey lorawan.AES128Key
	err = asKey.UnmarshalText([]byte(askey))
	if err != nil {
		return errors.Wrap(err, "application session key format error")
	}
	err = nsKey.UnmarshalText([]byte(nskey))
	if err != nil {
		return errors.Wrap(err, "network session key format error")
	}

	var b []byte
	if strings.ToLower(dataType) == "hex" {
		b, err = hex.DecodeString(str)
		if err != nil {
			return err
		}
	} else if strings.ToLower(dataType) == "string" {
		b = []byte(str)
	} else {
		return errors.New("data_type invalid,valid(hex/string)")
	}

	data, err := encode(devaddr, asKey, nsKey, b, count)
	if err != nil {
		return errors.Wrap(err, "encode phypayload error")
	}
	now := time.Now()
	var compactTime gateway.CompactTime
	compactTime = gateway.CompactTime(now)
	bs64Data := base64.StdEncoding.EncodeToString(data)
	// md5str := md5.New()
	// md5str.Write([]byte("044816c60202feff"))
	// md5str.Write([]byte(fmt.Sprintf("%d", now.Unix())))
	// md5str.Write([]byte("KCogiaertpoqij4nr"))
	// md5res := hex.EncodeToString(md5str.Sum(nil))
	rand.Seed(time.Now().Unix())
	pushdata := gateway.PushDataPacket{
		ProtocolVersion: 1,
		RandomToken:     1,
		GatewayMAC:      c.Mac,
		Payload: gateway.PushDataPayload{
			RXPK: []gateway.RXPK{
				gateway.RXPK{
					Tmst: uint32(now.Unix()),
					Time: &compactTime,
					Chan: 1,
					RFCh: 1,
					Freq: 470.5,
					Stat: 1,
					Modu: "LORA",
					DatR: gateway.DatR{
						LoRa: "SF10BW125",
					},
					CodR: "4/5",
					LSNR: -18,
					RSSI: -5,
					Size: uint16(len(data)),
					Data: bs64Data,
					RSig: []gateway.RSig{
						gateway.RSig{
							Ant:   1,
							Chan:  3,
							LSNR:  45.0,
							RSSIC: 20,
						},
						gateway.RSig{
							Ant:   2,
							Chan:  5,
							LSNR:  45.0,
							RSSIC: 20,
						},
					},
				},
			},
			// Stat: &gateway.Stat{RXFW: 3, RXOK: 3, RXNb: 3, Time: gateway.ExpandedTime(now),
			// 	TXNb:  30,
			// 	CPU:   uint32(rand.Intn(100)),
			// 	Memf:  uint32(rand.Intn(100)),
			// 	Memd:  uint32(rand.Intn(100)),
			// 	Discf: uint32(rand.Intn(100)),
			// 	Discd: uint32(rand.Intn(100)),
			// 	MD5:   md5res,
			// },
		},
	}

	wData, err := pushdata.MarshalBinary()
	// fmt.Println(wData)
	if err != nil {
		return errors.Wrap(err, "pushdata marshal error")
	}
	// send(wData)
	// fmt.Println("jsonUP:", string(wData[12:]))
	_, err = c.Conn.Write(wData)

	return err
}

func EncodeCommand() ([]byte, error) {
	var chmask lorawan.ChMask
	chmask.UnmarshalBinary([]byte{0xc0, 0x00})
	chmask_byte, _ := chmask.MarshalBinary()
	fmt.Printf("%08b\n", chmask_byte)
	cmds := []lorawan.MACCommand{
		lorawan.MACCommand{
			CID: lorawan.LinkADRReq,
			Payload: &lorawan.LinkADRReqPayload{
				DataRate: uint8(12),
				TXPower:  uint8(1),
				ChMask:   chmask,
				Redundancy: lorawan.Redundancy{
					ChMaskCntl: uint8(0),
					NbRep:      uint8(3),
				},
			},
		},
	}

	buf := bytes.NewBufferString("")
	for _, cmd := range cmds {
		b, _ := cmd.MarshalBinary()
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (c *GatewayClient) JoinRequest(deveui, appKey string) error {
	var (
		appEUI lorawan.EUI64
		devEUI lorawan.EUI64
		aesKey lorawan.AES128Key
	)
	err := appEUI.UnmarshalText([]byte("526973696e674847"))
	if err != nil {
		return err
	}
	err = devEUI.UnmarshalText([]byte(deveui))
	if err != nil {
		return err
	}

	err = aesKey.UnmarshalText([]byte(appKey))

	nonce := make([]byte, 2)
	rand.Seed(time.Now().UnixNano())
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}
	fmt.Printf("devNonce:%x\n", nonce)
	var devNonce lorawan.DevNonce
	copy(devNonce[:], nonce)
	d := &lorawan.JoinRequestPayload{
		AppEUI:   appEUI,
		DevEUI:   devEUI,
		DevNonce: devNonce,
	}

	py := &lorawan.PHYPayload{
		MHDR: lorawan.MHDR{
			MType: lorawan.JoinRequest,
			Major: lorawan.LoRaWANR1,
		},
		MACPayload: d,
	}
	py.SetMIC(aesKey)

	data, err := py.MarshalBinary()
	if err != nil {
		return err
	}
	bs64Data := base64.StdEncoding.EncodeToString(data)
	compactTime := gateway.CompactTime(time.Now())
	pushdata := gateway.PushDataPacket{
		ProtocolVersion: 1,
		RandomToken:     1,
		GatewayMAC:      c.Mac,
		Payload: gateway.PushDataPayload{
			RXPK: []gateway.RXPK{
				gateway.RXPK{
					Tmst: uint32(time.Now().Unix()),
					Time: &compactTime,
					Chan: 1,
					RFCh: 1,
					Freq: 472.3,
					Stat: 1,
					Modu: "LORA",
					DatR: gateway.DatR{
						LoRa: "SF12BW125",
					},
					CodR: "4/5",
					LSNR: -18,
					RSSI: -5,
					Size: uint16(len(data)),
					Data: bs64Data,
				},
			},
			Stat: &gateway.Stat{RXFW: 3, RXOK: 3, RXNb: 3, Time: gateway.ExpandedTime(time.Now()),
				TXNb: 30},
		},
	}

	wData, err := pushdata.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "pushdata marshal error")
	}

	fmt.Println(string(wData[:]))
	_, err = c.Conn.Write(wData)

	return err
}

func send(message []byte) {
	// serverAddr := "127.0.0.1" + ":" + "10006"
	// conn, err := net.Dial("udp", serverAddr)
	conn, err := net.Dial("udp", "127.0.0.1:1700")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("rx over !!!!!!!")
	// msg := make([]byte, 200)
	// conn.Read(msg)
	// fmt.Println(string(msg))
}
