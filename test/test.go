package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var fileName = "../rxdata_1_1.txt"
var chanRX = make(chan int, 1000)
var wgTime = sync.WaitGroup{}
var endTime, startTime int64
var gates, num int
var chanGo = make(chan int, 50)

func main() {
	go mqttserver()
	gates, num, datarx := rxdata()
	if datarx == nil {
		log.Fatal("文件有误")
	}
	// fmt.Println(len(datarx[0]))
	// fmt.Println(datarx[0])
	// return
	// gates = 1
	// num = 1
	wg := sync.WaitGroup{}
	wgTime.Add(1)
	for i := 0; i < gates; i++ {
		wg.Add(1)
		go func(index int) {
			rx(datarx[index*num : (index+1)*num])
			// time.Sleep(time.Millisecond * 10)
			wg.Done()
		}(i)
	}
	startTime = time.Now().UnixNano() / 1e6
	wgTime.Done()
	wg.Wait()
	// endTime = time.Now().UnixNano() / 1e6
	// fmt.Printf("send data number : %d\nsend success number :%d\nsend data start time : %d(ms)\nsend data end time : %d(ms)\nAll use time:%d(ms)", gates*num, onMessageReceivedNum, startTime, endTime, endTime-startTime)
	func() {
		for {
			var buffer [512]byte
			_, err := os.Stdin.Read(buffer[:])
			if err != nil {
				fmt.Println("read error:", err)
			}
			a := strconv.FormatFloat(float64(onMessageReceivedNum)/float64(endTime-startTime)*1000.00, 'f', 2, 64)
			fmt.Printf("send data number : %d\nsend success number :%d\nsend data start time : %d(ms)\nsend data end time : %d(ms)\nAll use time: %d (ms)\n AVG(S)= %s\n", gates*num, onMessageReceivedNum, startTime, endTime, endTime-startTime, a)
			os.Exit(1)
		}
	}()
}

func stdin() {
	for {
		var buffer [512]byte
		_, err := os.Stdin.Read(buffer[:])
		if err != nil {
			fmt.Println("read error:", err)
		}
		fmt.Println(gates, num)
		fmt.Printf("send data number : %d\nsend success number :%d\nsend data start time : %d(ms)\nsend data end time : %d(ms)\nAll use time:%d(ms)", gates*num, onMessageReceivedNum, startTime, endTime, endTime-startTime)
		os.Exit(1)
	}
}
func rx(chanData []string) {
	num11 := 0
	conn, err := net.Dial("udp", "127.0.0.1:1700")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	wgTime.Wait()
	data1 := make([]byte, 512)
	for _, data := range chanData {
		// fmt.Println(string(data))
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println(err)
			continue
		}
		conn.Read(data1)
		// _, err := conn.Read(data)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }
		chanGo <- 1
		num11++
	}
	fmt.Println(num11)
}

func rxdata() (int, int, []string) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return 0, 0, nil
	}
	defer file.Close()
	data := make([]string, 0)
	buf := bufio.NewReader(file)
	line, err := buf.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			log.Fatal("File read ok!")
		} else {
			log.Fatal("Read file error!", err)
		}
	}
	line = strings.TrimSpace(line)
	lines := strings.Split(line, "&&")
	if len(lines) != 2 {
		log.Fatal("文件有误")
	}
	gates, err := strconv.Atoi(lines[0])
	if err != nil {
		log.Fatal("文件有误")
	}
	num, err := strconv.Atoi(lines[1])
	if err != nil {
		log.Fatal("文件有误")
	}
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return 0, 0, nil
			}
		}
		line = strings.TrimSpace(line)
		data = append(data, line)
		// time.Sleep(time.Millisecond * 100)
	}
	return gates, num, data
}
