package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"../common"
	"../handlers"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	timeStart := time.Now().UnixNano() / 1e6
	fmt.Println("start on:", timeStart)
	port := 65453
	wg := &sync.WaitGroup{}
	for i, mac := range common.Macs {
		laddr := fmt.Sprintf(":%d", port+i)
		client, err := handlers.NewClient(laddr, mac, "127.0.0.1:1700")
		if err != nil {
			fmt.Println(err)
			continue
		}
		type A struct {
			addr, ns, as string
		}
		a := []A{
			A{
				addr: "00b26c1d",
				ns:   "3befba8d46b5234b2dd2d573bd23c43a",
				as:   "e9d2b9af47631ca45ad5abfd55245bb4",
			},
			A{
				addr: "000fdc5d",
				ns:   "c682ed99287d2cc8414898d96f034b86",
				as:   "18610db1e7e410eed0c30ce5d3cce66b",
			},
			A{
				addr: "00153014",
				as:   "16c2095c65dd2c04083584ffe0f21008",
				ns:   "e26a0f094a6c7a20252e9c671be1c1ca",
			},
		}
		for _, a0 := range a {
			fmt.Println(1)
			go func(aa A) {
				temperature := 330
				humidity := 60
				elec := 50
				r1 := rand.New(rand.NewSource(time.Now().UnixNano() + 3))
				r2 := rand.New(rand.NewSource(time.Now().UnixNano() + 5))
				ra1 := 5
				ra2 := 1
				ra3 := 2
				sendTimes := 1 //发送频率
				for {
					temps := make([]uint16, 0)
					humis := make([]uint16, 0)
					elecs := make([]uint16, 0)
					var count uint32 = 1
					for i := 0; i < sendTimes; i++ {
						temperature = temperature - (r1.Intn(10) - ra1)
						humidity = humidity - (r2.Intn(2) - ra2)
						elec = elec - (r2.Intn(4) - ra3)
						if temperature <= 10 {
							ra1++
						} else if temperature >= 400 {
							ra1--
						}
						if humidity <= 30 {
							ra2++
						} else if humidity >= 90 {
							ra2--
						}
						if humidity <= 10 {
							ra3++
						} else if humidity >= 120 {
							ra3--
						}
						temps = append(temps, (uint16)(temperature))
						humis = append(humis, (uint16)(humidity))
						elecs = append(elecs, (uint16)(elec))
					}
					fmt.Println(temps, humis, elecs)
					data := get(temps, humis, elecs)
					client.SendData(aa.addr, aa.as, aa.ns, data, "string", count)
					// client.SendData("00b26c1d", "e9d2b9af47631ca45ad5abfd55245bb4", "3befba8d46b5234b2dd2d573bd23c43a", data, "string", count) //---120 dve1
					// client.SendData("01587d6f", "bf231e9619d9d2bbbefe4e53f4be23ea", "924100f401d983c170f2fc38bb4c56ba", data, "string", count) //120 --dev2
					// client.SendData("01587d6f", "bf231e9619d9d2bbbefe4e53f4be23ea", "924100f401d983c170f2fc38bb4c56ba", data, "string", count) //120 --dev2
					time.Sleep(time.Minute * time.Duration(sendTimes))
					count++
				}
			}(a0)
			time.Sleep(time.Second * 1)
		}

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
func get(t, h, e []uint16) string {
	//ff020558 01045be3da80 02050113fffff803040 030fff7 0404002ffff7050300ffa8ff00
	//ff020418 01045bfdef5b 02050113fffff803040 030fff7 0404002ffff7
	buf := []byte{0xff, 0x02, 0x05, 29}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, (uint32)(time.Now().Unix()))
	b = append([]byte{1, 4}, b...)
	// t := []uint16{275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275, 275}
	temp := encryption(2, t)
	b = append(b, temp...)
	// h := []uint16{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}
	humit := encryption1(3, h)
	b = append(b, humit...)
	// e := []uint16{47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47, 47}
	elec := encryption1(4, e)
	b = append(b, elec...)
	// buf = append(buf, uint8(len(b)))
	buf = append(buf, b...)
	// fmt.Println(buf)
	c := hex.EncodeToString(buf)
	// fmt.Println(c)
	return c + "050300ffa8ff00"
}

func encryption(num uint8, data []uint16) []byte {
	// fmt.Println(len(t))
	bold := make([]byte, 0)
	told := (uint16)(0)
	tlen := (uint16)(0)
	for _, t1 := range data {
		if told == t1 {
			if tlen == 0 {
				b := (byte)(0xf1)
				bold = append(bold, b)
			} else {
				bold[len(bold)-1] = bold[len(bold)-1] + 1
				if bold[len(bold)-1] == 0xff {
					tlen = 0
					continue
				}
			}
			tlen++
			continue
		}
		told = t1
		tlen = 0
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, t1)
		bold = append(bold, b...)
	}
	// fmt.Println("temp ", bold)
	b := make([]byte, 2)
	//分包号
	b[0] = num
	//保存数据长度
	b[1] = (uint8)(len(bold))
	bold = append(b, bold...)
	// fmt.Println(bold)
	// c := hex.EncodeToString(bold)
	// fmt.Println(c)
	return bold
}
func encryption1(num uint8, data []uint16) []byte {
	// fmt.Println(len(t))
	bold := make([]byte, 0)
	told := (uint16)(0)
	tlen := (uint16)(0)
	for _, t1 := range data {
		if told == t1 {
			if tlen == 0 {
				b := (byte)(0xf1)
				bold = append(bold, b)
			} else {
				bold[len(bold)-1] = bold[len(bold)-1] + 1
				if bold[len(bold)-1] == 0xff {
					tlen = 0
					continue
				}
			}
			tlen++
			continue
		}
		told = t1
		tlen = 0
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, t1)
		bold = append(bold, b[1])
	}
	// fmt.Println("temp ", bold)
	b := make([]byte, 2)
	//分包号
	b[0] = num
	//保存数据长度
	b[1] = (uint8)(len(bold))
	bold = append(b, bold...)
	// fmt.Println(bold)
	// c := hex.EncodeToString(bold)
	// fmt.Println(c)
	return bold
}

// ff02051d01045bfdff3d022f017e018101850187018b018af1018f019201960198019601970192018f01890182017e01790176016c0162015a01520318003c003d003ef3003ff20040f200410042f20043f30044f3043000360038003c003f00430047004a004e0051005200560059005c0060006400660069006a006c006d0070007100740077050300ffa8ff00
// ff02051d01045bfdff3d022f017e018101850187018b018af1018f019201960198019601970192018f01890182017e01790176016c0162015a01520318003c003d003ef3003ff20040f200410042f20043f30044f3043000360038003c003f00430047004a004e0051005200560059005c0060006400660069006a006c006d0070007100740077050300ffa8ff00
