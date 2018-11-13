package main

import (
	"fmt"

	"../lorawan/band"
)

// Channel defines a gateay channel.
type Channel struct {
	Modulation       band.Modulation `json:"modulation"`
	Frequency        int             `json:"frequency"`
	Bandwidth        int             `json:"bandwidth"`
	Bitrate          int             `json:"bitrate,omitempty"`          // FSK modulation only
	SpreadingFactors []int           `json:"spreadingFactors,omitempty"` // LoRa modulation only
}

func main() {
	fmt.Println(470300000)
	fmt.Println(470300000 + 200000*7)
	chans := []Channel{
		Channel{
			band.LoRaModulation,
			923200000,
			125,
			0,
			[]int{10, 11, 12},
		},
		Channel{
			band.LoRaModulation,
			923400000,
			125,
			0,
			[]int{10, 11, 12},
		},
	}

	calc(chans)
}

var radioBandwidthPerChannelBandwidth = map[int]int{
	500000: 1100000, // 500kHz channel
	250000: 1000000, // 250kHz channel
	125000: 925000,  // 125kHz channel
}

const defaultRadioBandwidth = 925000

func calc(channels []Channel) {
	for _, c := range channels {
		channelBandwidth := c.Bandwidth * 1000
		//channelMax := c.Frequency + (channelBandwidth / 2)
		radioBandwidth, ok := radioBandwidthPerChannelBandwidth[channelBandwidth]
		if !ok {
			radioBandwidth = defaultRadioBandwidth
		}
		minRadioCenterFreq := c.Frequency - (channelBandwidth / 2) + (radioBandwidth / 2)
		fmt.Println("minRadioCenterFreq:", minRadioCenterFreq)
	}
}
