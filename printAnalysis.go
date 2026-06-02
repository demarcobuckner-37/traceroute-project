package main

import (
	"fmt"
)

func PrintHopAnalysis(hops []Hop) {

	for _, h := range hops {

		host := h.Host
		ip := h.IPAddress

		if host == "" {
			host = "*"
		}

		if ip == "" {
			ip = "*"
		}

		fmt.Println("Hop Analysis:")

		fmt.Printf("Hop: %d\n", h.TTL)
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

	}
}
