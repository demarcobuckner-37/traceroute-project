package main

import (
	"fmt"
)

func PrintHopAnalysis(hops []Hop) {

	fmt.Println("\nHop Analysis:")

	highestScore := 0
	highestHop := 0

	// HOP ANALYSIS
	// Display detailed performance metrics for each hop.
	for _, h := range hops {

		// Replace missing responses with placeholder values.
		ip := h.IPAddress
		host := h.Host

		if host == "" {
			host = "*"
		}

		if ip == "" {
			ip = "*"
		}

		// HOP METRICS
		// Display latency, packet loss, jitter,
		// and bottleneck analysis results.
		fmt.Printf("\nHop %d:\n", h.TTL)
		fmt.Printf("Host: %s\n", host)
		fmt.Printf("IP Address: %s\n", ip)
		fmt.Printf("Average RTT: %v\n", h.AvgRTT)
		fmt.Printf("Packet Loss: %.2f%%\n", h.PacketLoss)
		fmt.Printf("Jitter: %v\n", h.Jitter)
		fmt.Printf("RTT Increase: %v\n", h.RTTIncrease)
		fmt.Printf("Max RTT: %v\n", h.MaxRTT)
		fmt.Printf("Min RTT: %v\n", h.MinRTT)
		fmt.Printf("Bottleneck Score: %d\n", h.BottleneckScore)
		fmt.Printf("Severity: %s\n", h.Severity)
		fmt.Print("\n")

		// ISSUE REPORTING
		// Display specific conditions that contributed
		// to the hop's bottleneck score
		if h.Severity != "None" {
			fmt.Println("Issues Detected:")

			if h.RTTIncreaseFlag {
				fmt.Printf("- Significant RTT Increase\n")
			}
			if h.HighPacketLossFlag {
				fmt.Printf("- High Packet Loss\n")
			}
			if h.HighJitterFlag {
				fmt.Printf("- High Jitter\n")
			}
			if h.HighAvgRTTFlag {
				fmt.Printf("- High Average RTT\n")
			}
			if h.HighMaxRTTFlag {
				fmt.Printf("- High Max RTT\n")
			}

		}

		if h.BottleneckScore > highestScore {
			highestScore = h.BottleneckScore
			highestHop = h.TTL
		}

		fmt.Print("\n-----------------------------\n\n")
	}
	if highestScore > 0 {
		fmt.Printf("\nHighest bottleneck score: %d at Hop %d\n", highestScore, highestHop)
	} else {
		fmt.Println("\nNo significant bottlenecks detected.")
	}

}
