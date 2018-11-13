// package main

// import (
// 	"log"
// 	"math/rand"
// 	"time"

// 	"github.com/influxdata/influxdb/client/v2"
// )

// const (
// 	MyDB     = "entre"
// 	username = "root"
// 	password = ""
// )

// func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
// 	q := client.Query{
// 		Command:  cmd,
// 		Database: MyDB,
// 	}
// 	if response, err := clnt.Query(q); err == nil {
// 		if response.Error() != nil {
// 			return res, response.Error()
// 		}
// 		res = response.Results
// 	} else {
// 		return res, err
// 	}
// 	return res, nil
// }

// func writePoints(clnt client.Client) {
// 	sampleSize := 1 * 10000
// 	rand.Seed(42)
// 	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
// 		Database:  MyDB,
// 		Precision: "us",
// 	})

// 	for i := 0; i < sampleSize; i++ {
// 		tags := map[string]string{
// 			"system_name": "dev001",
// 		}
// 		fields := map[string]interface{}{
// 			"value": i,
// 		}
// 		pt, err := client.NewPoint("test", tags, fields, time.Now())
// 		if err != nil {
// 			log.Fatalln("Error: ", err)
// 		}
// 		bp.AddPoint(pt)
// 	}

// 	err := clnt.Write(bp)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	//fmt.Printf("%d task done\n",num)
// }

// func main() {
// 	// Make client
// 	c, err := client.NewHTTPClient(client.HTTPConfig{
// 		Addr:     "http://localhost:8086",
// 		Username: username,
// 		Password: password,
// 	})

// 	if err != nil {
// 		log.Fatalln("Error: ", err)
// 	}
// 	// _, err = queryDB(c, fmt.Sprintf("CREATE DATABASE %s", MyDB))
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	defer writePoints(c)

// }

package main

import (
	"fmt"
)

func main() {
	b := [8]byte{36, 197, 217, 230, 50, 87, 245, 140}
	// b, _ := hex.DecodeString("24 c5 d9 e6 32 57 f5 8c")
	fmt.Println(b)

}
