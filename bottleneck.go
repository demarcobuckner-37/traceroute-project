package main

import (
	"fmt"
	"time"
)

func RunBottleneckDetection(hop []Hop) []Hop {

	fmt.Print("\nRunning Bottleneck Detection...\n")

	highestScore := 0
	highestHop := 0

	for i := 1; i < len(hop); i++ {
		h := &hop[i]
		prev := &hop[i-1]

		rttIncrease := h.AvgRTT - prev.AvgRTT
		h.RTTIncrease = rttIncrease

		score := 0

		if prev.SuccessfulProbes > 0 {

			if rttIncrease > 30*time.Millisecond {

				fmt.Printf("Significant RTT increase detected between Hop %d and Hop %d: %v\n", prev.TTL, h.TTL, rttIncrease)

				score += 2
			}
		}

		if h.PacketLoss >= 99.9 {

			fmt.Printf("Hop %d is unresponsive with 100%% packet loss\n", h.TTL)

			score += 4

		} else if h.PacketLoss > 20 {

			fmt.Printf("High packet loss detected at Hop %d: %.2f%%\n", h.TTL, h.PacketLoss)

			score += 2

		}

		if h.Jitter > 30*time.Millisecond {

			fmt.Printf("High jitter detected at Hop %d: %v\n", h.TTL, h.Jitter)
			score++

		}

		if h.AvgRTT > 100*time.Millisecond {

			fmt.Printf("High average RTT detected at Hop %d: %v\n", h.TTL, h.AvgRTT)
			score += 2

		}

		if h.MaxRTT > 200*time.Millisecond {

			fmt.Printf("High max RTT detected at Hop %d: %v\n", h.TTL, h.MaxRTT)
			score++

		}

		h.BottleneckScore = score

		if score > 0 {

			switch {
			case score >= 4:
				h.Severity = "Severe"
			case score >= 2:
				h.Severity = "Moderate"
			case score > 0:
				h.Severity = "Minor"

			}

		}

		if score > highestScore {
			highestScore = score
			highestHop = h.TTL
		}

	}

	if highestScore > 0 {

		fmt.Printf("\nHighest bottleneck score: %d at Hop %d\n", highestScore, highestHop)

	} else {

		fmt.Printf("\nNo significant bottlenecks detected.\n")
	}

	return hop

}
