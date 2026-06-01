package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Hop struct {
	// TTL value that discovered this hop
	TTL int
	// RTT measurements from multiple probes
	RTTs []time.Duration
	// Hostname of responding router
	Host string
	// IP address of responding router
	IPAddress string

	AvgRTT time.Duration

	MinRTT time.Duration

	MaxRTT time.Duration

	SentProbes int

	SuccessfulProbes int

	PacketLoss float64

	Jitter time.Duration

	RTTIncrease time.Duration

	BottleneckScore int

	Severity string
}

func RunTraceroute(host string) []Hop {

	// DNS RESOLUTION
	// Resolve hostname into an IPv4 address
	dstAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// ICMP CONNECTION
	// Open raw ICMP packet listener on all interfaces
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Close connection when program exits
	defer c.Close()

	//ROUTE STORAGE
	// Slice used to store discovered hops and RTT measurements.
	var hops []Hop

	var probesSent int
	var totalRTT time.Duration

	var successfulProbes int
	var minTime time.Duration
	var maxTime time.Duration

	destinationReached := false

	// Increase TTL one hop at a time to reveal the route
	for ttl := 1; ttl <= 64; ttl++ {

		//MULTIPLE PROBES
		// Send three probes per hop to gather RTT statistics.
		for i := 0; i < 3; i++ {

			// Set TTL for outgoing packets
			if err := c.IPv4PacketConn().SetTTL(ttl); err != nil {
				log.Fatalf("SetTTL failed: %s", err)
			}

			// ICMP ECHO REQUEST
			// Build ICMP Echo Request packet
			m := icmp.Message{
				Type: ipv4.ICMPTypeEcho, Code: 0,
				Body: &icmp.Echo{
					// Process ID helps identify our packet
					ID: os.Getpid() & 0xffff,

					// Sequence number for tracking requests
					Seq: ttl,

					// Packet payload data
					Data: []byte("WADDUP"),
				},
			}

			// Convert ICMP message into raw bytes
			b, err := m.Marshal(nil)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			// RTT TIMING
			// Start round-trip timer before sending packet
			start := time.Now()

			// SEND PACKET
			// Send ICMP Echo Request to destination
			_, err = c.WriteTo(b, dstAddr)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			found := false

			for k := range hops {
				if hops[k].TTL == ttl {
					hops[k].SentProbes++
					found = true
					break
				}
			}

			if !found {
				hops = append(hops, Hop{
					TTL:        ttl,
					SentProbes: 1,
				})
			}

			probesSent++

			// RECEIVE REPLY
			// Create buffer for incoming packet data
			buf := make([]byte, 1024)

			// Read incoming ICMP reply packet
			err = c.SetReadDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				fmt.Println(err)
				return nil
			}
			numRead, peer, err := c.ReadFrom(buf)
			if err != nil {
				fmt.Printf("\nHop %d Probe %d\n", ttl, i+1)
				fmt.Println("Request timed out.")
				continue
			}
			successfulProbes++

			// RTT TIMING
			// Measure round-trip time
			elapsed := time.Since(start)
			totalRTT += elapsed

			if successfulProbes == 1 {
				minTime = elapsed
				maxTime = elapsed
			} else {
				if elapsed < minTime {
					minTime = elapsed
				}
				if elapsed > maxTime {
					maxTime = elapsed
				}
			}

			fmt.Printf("\nHop %d Probe %d\n", ttl, i+1)

			fmt.Println("RTT:", elapsed)

			// Print raw packet bytes in hexadecimal
			//fmt.Printf("% X\n", buf[:numRead])

			// Print responding host
			fmt.Printf("Received reply from %v\n", peer)

			// PACKET PARSING

			// Parse raw ICMP packet into structured message
			rm, err := icmp.ParseMessage(1, buf[:numRead])
			if err != nil {
				fmt.Println(err)
				return nil
			}

			switch rm.Type {

			//INTERMEDIATE HOP REPLY
			//Router discards when TTL expires and sends Time Exceeded message back to sender
			case ipv4.ICMPTypeTimeExceeded:

				fmt.Printf("Received Time Exceeded message: TTL: %d| Reply From %v| Time Elasped: %v|\n", ttl, peer, elapsed)

				// HOP LOOKUP
				// Check whether this TTL already exists in the route.
				// If it does, append another RTT measurement.
				found := false

				for k := range hops {
					if hops[k].TTL == ttl {
						hops[k].RTTs = append(hops[k].RTTs, elapsed)

						hops[k].Host = peer.String()
						hops[k].IPAddress = peer.String()

						hops[k].SuccessfulProbes++

						found = true
						break
					}
				}

				// NEW HOP DISCOVERED
				// Create a new hop record if this TTL has not been seen before.
				if !found {
					hops = append(hops, Hop{
						TTL:              ttl,
						RTTs:             []time.Duration{elapsed},
						Host:             peer.String(),
						IPAddress:        peer.String(),
						SuccessfulProbes: 1,
					})
				}

			//DESTINATION REPLY
			//Host receives Echo Request and replies with Echo Reply message
			case ipv4.ICMPTypeEchoReply:

				echoReply, ok := rm.Body.(*icmp.Echo)
				if !ok {
					fmt.Println("Failed to parse ICMP Echo Reply body")
					continue
				}

				// VALIDATION CHECKS
				// Check that reply matches our process ID to avoid confusion with other traceroute instances
				if echoReply.ID != os.Getpid()&0xffff {
					fmt.Println("Received reply for another process")
					continue
				}

				// Check that sequence number matches expected TTL to ensure correct order
				if echoReply.Seq != ttl {
					fmt.Printf("Received out-of-order Echo Reply: Seq %d\n", echoReply.Seq)
					continue
				}

				// HOP LOOKUP
				// Check if this TTL already exists in the route data
				found := false

				for k := range hops {

					// RTT STORAGE
					// Add RTT measurement to existing hop
					if hops[k].TTL == ttl {
						hops[k].RTTs = append(hops[k].RTTs, elapsed)

						hops[k].Host = peer.String()
						hops[k].IPAddress = peer.String()

						hops[k].SuccessfulProbes++

						found = true
						break
					}
				}

				// NEW HOP RECORD
				// Create a new hop entry if this TTL has not been seen before
				if !found {
					hops = append(hops, Hop{
						TTL:              ttl,
						RTTs:             []time.Duration{elapsed},
						Host:             peer.String(),
						IPAddress:        peer.String(),
						SuccessfulProbes: 1,
					})
				}

				// ROUTE COMPLETE
				// Destination reached, stop further TTL discovery
				destinationReached = true
				fmt.Println("Destination reached!")
				fmt.Printf("Echo Reply ID: %d, Seq: %d, Data: %s\n", echoReply.ID, echoReply.Seq, string(echoReply.Data))
				fmt.Printf("TTL: %d\n", ttl)

			default:
				fmt.Printf("Received ICMP message of type %v\n", rm.Type)
			}

			time.Sleep(1 * time.Second)

		}
		if destinationReached {
			break
		}
	}

	// TRACEROUTE STATISTICS
	// Calculate overall RTT and packet loss metrics
	if successfulProbes > 0 {
		avgRTT := totalRTT / time.Duration(successfulProbes)
		fmt.Printf("\nAverage RTT over %d succesfull probes: %s\n", successfulProbes, avgRTT)
		fmt.Printf("Minimum RTT: %s\n", minTime)
		fmt.Printf("Maximum RTT: %s\n", maxTime)
	}
	packetLoss := float64(probesSent-successfulProbes) / float64(probesSent) * 100
	fmt.Printf("Packet Loss: %.2f%%\n", packetLoss)

	println("\nTraceroute complete.")

	// ROUTE OUTPUT
	// Print all discovered hops and their RTT measurements
	fmt.Printf("\nTraceroute Results:\n")
	for _, hop := range hops {
		fmt.Printf("TTL: %d | Host: %s | IP: %s\n", hop.TTL, hop.Host, hop.IPAddress)
	}

	return hops
}
