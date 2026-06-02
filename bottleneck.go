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
	for i := 1; i < len(hop); i++ {
		h := &hop[i]
		prev := &hop[i-1]

		//RTT INCREASE CALCULATION
		//Mearsures latency growth between adjanceent hops
		rttIncrease := h.AvgRTT - prev.AvgRTT
		h.RTTIncrease = rttIncrease

		score := 0

		//Significant RTT increase may indicate a bottleneck or congestion point
		if prev.SuccessfulProbes > 0 {

			if rttIncrease > 30*time.Millisecond {

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

		//JITER DETECTION
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

		//Display the most severe bottleneck found so far
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
