package main

import (
	"fmt"
	"time"
)

func RunAnalyzeHops(hop []Hop) []Hop {

	fmt.Printf("\nCalculating Average RTT for each hop...\n")

	// Loop through each hop and calculate average RTT

	for i := range hop {

		h := &hop[i]

		if h.SentProbes > 0 {

			h.PacketLoss = float64(h.SentProbes-h.SuccessfulProbes) / float64(h.SentProbes) * 100

		} else {
			h.PacketLoss = 0
		}

		if len(h.RTTs) == 0 {

			continue
		}

		fmt.Printf("Hop: %d | RTTs: %v\n", h.TTL, h.RTTs)

		total := time.Duration(0)

		for _, rtt := range h.RTTs {

			total += rtt

			if rtt < h.MinRTT || h.MinRTT == 0 {
				h.MinRTT = rtt
			}
			if rtt > h.MaxRTT {
				h.MaxRTT = rtt
			}

		}

		var totalDiff time.Duration

		for i := 1; i < len(h.RTTs); i++ {

			diff := h.RTTs[i] - h.RTTs[i-1]
			if diff < 0 {
				diff = -diff
			}
			totalDiff += diff
		}
		if len(h.RTTs) > 1 {
			h.Jitter = totalDiff / time.Duration(len(h.RTTs)-1)
		} else {
			h.Jitter = 0
		}

		h.AvgRTT = total / time.Duration(len(h.RTTs))

		fmt.Printf("Average RTT: %v\n", h.AvgRTT)
		fmt.Printf("Min RTT: %v\n", h.MinRTT)
		fmt.Printf("Max RTT: %v\n", h.MaxRTT)
		fmt.Printf("Packet Loss: %.2f%%\n", h.PacketLoss)
		fmt.Printf("Jitter: %v\n", h.Jitter)

		fmt.Print("\n")

	}
	return hop

}
