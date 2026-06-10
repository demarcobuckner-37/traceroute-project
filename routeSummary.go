package main

import (
	"fmt"
	"time"
)

func RouteSummary(hops []Hop) {

	fmt.Println("\nRoute Summary:")

	// ROUTE STATISTICS
	// Track overall route performance metrics
	totalHops := len(hops)

	var totalRTT time.Duration
	var totalLoss float64

	// Track the most severe bottleneck detected
	highestScore := 0
	worstHop := 0

	// AGGREGATE ROUTE DATA
	// Combine hop statistics to produce route-wide averages.
	for _, h := range hops {

		totalRTT += h.AvgRTT
		totalLoss += h.PacketLoss

		if h.BottleneckScore > highestScore {
			highestScore = h.BottleneckScore
			worstHop = h.TTL
		}
	}

	// ROUTE AVERAGES
	// Calculate average latency and packet loss
	avgRTT := totalRTT / time.Duration(totalHops)
	avgLoss := totalLoss / float64(totalHops)

	// SUMMARY OUTPUT
	// Display overall route health metrics.
	fmt.Printf("Total Hops: %d\n", totalHops)
	fmt.Printf("Average Route RTT: %v\n", avgRTT)
	fmt.Printf("Average Packet Loss: %.2f%%\n", avgLoss)
	fmt.Printf("Highest Bottleneck Score: %d\n", highestScore)
	fmt.Printf("Worst Hop: %d\n", worstHop)
}
