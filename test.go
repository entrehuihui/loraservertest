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
	"time"
)

func main() {
	gates, num, datarx := rxdata()
	if datarx == nil {
		log.Fatal("文件有误")
	}
	gates = 1
	// num = 40
	for index := 0; index < gates; index++ {
		rx(datarx[index*num : (index+1)*num])
		time.Sleep(time.Millisecond * 10)
	}
}

func rx(chanData []string) {
	num11 := 0
	conn, err := net.Dial("udp", "127.0.0.1:1700")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for _, data := range chanData {
		// fmt.Println(string(data))
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println(err)
			continue
		}
		data := make([]byte, 10000)
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// fmt.Println(data[:i])
		num11++
	}
	fmt.Println(num11)
}

func rxdata() (int, int, []string) {
	fileName := "rxdata_10_10.txt"
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
