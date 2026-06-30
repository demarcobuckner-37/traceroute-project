package main

import (
	"fmt"
	"time"
)

func RunRouteComparison(hop1 []Hop, hop2 []Hop) {

	fmt.Println("\nRoute Comparison Analysis:")

	//Compare two traceroute results to identify changes in network performance over time
	//This can help detect emerging bottlenecks or improvements in the path

	// ROUTE LENGTH CHECK
	// Detect changes in the number of hops between routes
	if len(hop1) != len(hop2) {
		fmt.Printf("\nRoute length changed: Route 1 has %d hops, Route 2 has %d hops\n", len(hop1), len(hop2))
	}

	// HOP-BY-HOP COMPARISON
	for i := range hop1 {

		if i >= len(hop2) {
			break
		}

		h1 := hop1[i]
		h2 := hop2[i]

		if h1.TTL != h2.TTL {
			continue
		}
		oldIP := h1.IPAddress
		newIP := h2.IPAddress

		if oldIP == "" {
			oldIP = "*"
		}

		if newIP == "" {
			newIP = "*"
		}

		// ROUTE CHANGE DETECTION
		// Check whether traffic is taking a different path.
		if h1.IPAddress != h2.IPAddress {
			fmt.Printf("Route change detected at Hop %d: %s -> %s\n", h1.TTL, oldIP, newIP)
			fmt.Print("\n")
			continue
		}

		// PERFORMANCE DIFFERENCE CALCULATION
		// Measure changes in latency, packet loss
		// jitter, and bottleneck score
		rttDiff := h2.AvgRTT - h1.AvgRTT

		packetLossDiff := h2.PacketLoss - h1.PacketLoss

		jitterDiff := h2.Jitter - h1.Jitter

		bottleneckScoreDiff := h2.BottleneckScore - h1.BottleneckScore

		var changes []string

		//Detect significant changes in performance metrics
		if rttDiff > 30*time.Millisecond {
			changes = append(changes, fmt.Sprintf("RTT increased by %v", rttDiff))
		}
		if rttDiff < -30*time.Millisecond {
			changes = append(changes, fmt.Sprintf("RTT decreased by %v", -rttDiff))
		}
		if packetLossDiff > 20 {
			changes = append(changes, fmt.Sprintf("Packet Loss increased by %.2f%%", packetLossDiff))
		}
		if packetLossDiff < -20 {
			changes = append(changes, fmt.Sprintf("Packet Loss decreased by %.2f%%", -packetLossDiff))
		}

		if jitterDiff > 30*time.Millisecond {
			changes = append(changes, fmt.Sprintf("Jitter increased by %v", jitterDiff))
		}
		if jitterDiff < -30*time.Millisecond {
			changes = append(changes, fmt.Sprintf("Jitter decreased by %v", -jitterDiff))
		}

		if bottleneckScoreDiff > 0 {
			changes = append(changes, fmt.Sprintf("Bottleneck Score increased by %d", bottleneckScoreDiff))
		}
		if bottleneckScoreDiff < 0 {
			changes = append(changes, fmt.Sprintf("Bottleneck Score decreased by %d", -bottleneckScoreDiff))
		}

		if len(changes) > 0 {
			fmt.Printf("Significant performance change detected at Hop %d (%s)\n", h1.TTL, h1.IPAddress)
			for _, change := range changes {
				fmt.Printf("  - %s\n", change)
			}
		}

		// HOP SUMMARY
		// Display side-by-side metrics for both routes.
		fmt.Printf("\nComparing Hop %d:\n", h1.TTL)

		host := h1.Host
		ip := h1.IPAddress

		if host == "" {
			host = "*"
		}
		if ip == "" {
			ip = "*"
		}

		fmt.Printf("Host: %s | IP: %s\n", host, ip)
		fmt.Printf("Avg RTT: Route 1: %v -> Route 2: %v\n", h1.AvgRTT, h2.AvgRTT)
		fmt.Printf("Packet Loss: Route 1: %.2f%% -> Route 2: %.2f%%\n", h1.PacketLoss, h2.PacketLoss)
		fmt.Printf("Jitter: Route 1: %v -> Route 2: %v\n", h1.Jitter, h2.Jitter)
		fmt.Printf("RTT Increase: Route 1: %v -> Route 2: %v\n", h1.RTTIncrease, h2.RTTIncrease)
		fmt.Printf("Max RTT: Route 1: %v -> Route 2: %v\n", h1.MaxRTT, h2.MaxRTT)
		fmt.Printf("Min RTT: Route 1: %v -> Route 2: %v\n", h1.MinRTT, h2.MinRTT)
		fmt.Printf("Bottleneck Score: Route 1: %d -> Route 2: %d\n", h1.BottleneckScore, h2.BottleneckScore)
		fmt.Printf("Severity: Route 1: %s -> Route 2: %s\n", h1.Severity, h2.Severity)
		fmt.Print("\n-----------------------------\n\n")
	}

}
