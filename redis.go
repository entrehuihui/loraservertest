package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type measurement struct {
	Name   string
	Tags   map[string]string
	Values map[string]interface{}
}

// ns u ms s
func main() {
	ti := time.Now().UnixNano()
	fmt.Println(ti)
	return
	a := measurement{
		Name: "device_uplink",
		Tags: map[string]string{
			"application_name": "testApp",
			"dev_eui":          "1000000000000100",
			"device_name":      "testDev100",
			"dr":               "2",
			"frequency":        "470500000",
		},
		Values: map[string]interface{}{
			"f_cnt": 1,
			"rssi":  20,
			"snr":   45.000000,
			"value": 1,
		},
	}
	fmt.Println(a.String())
}

// device_uplink,application_name=testApp,dev_eui=1000000000000100,device_name=testDev100,dr=2,frequency=470500000 f_cnt=1i,rssi=20i,snr=45.000000,value=1i
// device_uplink,application_name=testApp,dev_eui=1000000000000100,device_name=testDev100,dr=2,frequency=470500000 f_cnt=1i,rssi=20i,snr=45.000000,value=1i
func (m measurement) String() string {
	var tags []string
	var values []string

	for k, v := range m.Tags {
		tags = append(tags, fmt.Sprintf("%s=%v", k, formatInfluxValue(v, false)))
	}

	for k, v := range m.Values {
		values = append(values, fmt.Sprintf("%s=%v", k, formatInfluxValue(v, true)))
	}

	// as maps are unsorted the order of tags and values is random.
	// this is not an issue for influxdb, but makes testing more complex.
	sort.Strings(tags)
	sort.Strings(values)

	return fmt.Sprintf("%s,%s %s", m.Name, strings.Join(tags, ","), strings.Join(values, ","))
}

func formatInfluxValue(v interface{}, quote bool) string {
	switch v := v.(type) {
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case int, uint, uint8, int8, uint16, int16, uint32, int32, uint64, int64:
		return fmt.Sprintf("%di", v)
	case string:
		if quote {
			return strconv.Quote(v)
		}
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
