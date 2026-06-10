package main

import (
	"time"
)

func RunAnalyzeHops(hop []Hop) []Hop {

	//Analyze each hp and calculate network performance metrics
	//Used for bottleneck detection and route analysis
	for i := range hop {

		h := &hop[i]

		//PACKET LOSS CALCULATION
		//Percentage of probes that did not receive a response
		if h.SentProbes > 0 {

			h.PacketLoss = float64(h.SentProbes-h.SuccessfulProbes) / float64(h.SentProbes) * 100

		} else {
			h.PacketLoss = 0
		}

		//Skip hops with no successful probes
		if len(h.RTTs) == 0 {

			continue
		}

		//RTT STATISTICS
		//Calculate average, min, and max RTT values
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

		//JITTER CALCULATION
		//Average variation between consecutive RTT samples
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

		//Average RTT Calculation
		//Mean latency across all successful probes
		h.AvgRTT = total / time.Duration(len(h.RTTs))

	}
	return hop
}
