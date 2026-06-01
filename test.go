package main

import (
	"fmt"
)

func main() {

	host := "www.google.com"
	hop := RunTraceroute(host)
	hop = RunAnalyzeHops(hop)
	RunBottleneckDetection(hop)

	for _, h := range hop {

		fmt.Printf("\nHop: %d | Sent Probes: %d | Successful Probes: %d\n", h.TTL, h.SentProbes, h.SuccessfulProbes)

		/*fmt.Printf("Hop: %d | Host: %s | IP: %s | Avg RTT: %v | Min RTT: %v | Max RTT: %v | Packet Loss: %.2f%% | Jitter: %v\n",
		h.TTL, h.Host, h.IPAddress, h.AvgRTT, h.MinRTT, h.MaxRTT, h.PacketLoss, h.Jitter) */

	}

}
