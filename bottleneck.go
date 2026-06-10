package main

import (
	"fmt"
	"time"
)

func RunBottleneckDetection(hop []Hop) []Hop {

	fmt.Print("\nRunning Bottleneck Detection...\n")

	//Track the most severe bottleneck found in the route
	highestScore := 0
	highestHop := 0

	//Compare each hop against the previous
	for i := 0; i < len(hop); i++ {
		h := &hop[i]

		//Reset all detection flags before evaluating this hop
		h.RTTIncreaseFlag = false
		h.HighPacketLossFlag = false
		h.HighJitterFlag = false
		h.HighAvgRTTFlag = false
		h.HighMaxRTTFlag = false

		//Bottleneck score accumulates points based on detected issues
		score := 0

		if i > 0 {
			prev := &hop[i-1]

			//RTT INCREASE CALCULATION
			//Mearsures latency growth between adjancent hops
			rttIncrease := h.AvgRTT - prev.AvgRTT
			h.RTTIncrease = rttIncrease

			//Significant RTT increase may indicate a bottleneck or congestion point
			if prev.SuccessfulProbes > 0 && rttIncrease > 30*time.Millisecond {

				h.RTTIncreaseFlag = true
				score += 2
			}
		}

		//PACKET LOSS DETECTION
		//Heavily weighted because packet loss directly impacts network performance
		if h.PacketLoss >= 99.9 {

			h.HighPacketLossFlag = true
			score += 4

		} else if h.PacketLoss > 20 {

			h.HighPacketLossFlag = true
			score += 2

		}

		//JITTER DETECTION
		//High jitter can indicate instability and congestion
		if h.Jitter > 30*time.Millisecond {

			h.HighJitterFlag = true
			score++

		}

		// AVERAGE RTT DETECTION
		// Detects consistently high latency.
		if h.AvgRTT > 100*time.Millisecond {

			h.HighAvgRTTFlag = true
			score += 2

		}

		// MAXIMUM RTT DETECTION
		// Detects occasional latency spikes.
		if h.MaxRTT > 200*time.Millisecond {

			h.HighMaxRTTFlag = true
			score++

		}

		//Store overall bottleneck score
		h.BottleneckScore = score

		//Classify severity based on score
		switch {
		case score >= 4:
			h.Severity = "Severe"
		case score >= 2:
			h.Severity = "Moderate"
		case score > 0:
			h.Severity = "Minor"
		default:
			h.Severity = "None"

		}

		// Track the highest bottleneck score observed.
		if score > highestScore {
			highestScore = score
			highestHop = h.TTL
		}

	}

	//Display most significant bottleneck detected in the route
	if highestScore > 0 {

		fmt.Printf("\nHighest bottleneck score: %d at Hop %d\n", highestScore, highestHop)

	} else {

		fmt.Printf("\nNo significant bottlenecks detected.\n")
	}

	return hop

}
