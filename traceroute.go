package main

import (
	"fmt"
	"log"
	"math/rand"
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

	// Average RTT across successful probes
	AvgRTT time.Duration

	// Smallest RTT observed for this hop
	MinRTT time.Duration

	// Largest RTT observed for this hop
	MaxRTT time.Duration

	// Total number of probes sent
	SentProbes int

	// Number of probes that received a reply
	SuccessfulProbes int

	// Percentage of probes that were lost
	PacketLoss float64

	// Variation between RTT measurements
	Jitter time.Duration

	// RTT increase relative to the previous hop
	RTTIncrease time.Duration

	// Combined score used to identify bottlenecks
	BottleneckScore int

	// Overall severity classification
	Severity string

	//Performance Issue Flags
	RTTIncreaseFlag bool

	HighPacketLossFlag bool

	HighJitterFlag bool

	HighAvgRTTFlag bool

	HighMaxRTTFlag bool
}

type Fault struct {
	// Enables or disables fault injection
	Enabled bool

	// Probability that a probe will be treated as lost
	DropRate float64

	// Maximum artificial delay applied to a probe
	Delay time.Duration
}

func RunTraceroute(host string, faults Fault) []Hop {

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

	//Used to stop TTL discovery once destination responds
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

					// Unique sequence number for identiying individual probes
					//Combines TTL and probe number to track delayed responses
					Seq: ttl*10 + i,

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

			//PROBE TRACKING
			//Record that a probe was sent for this TTL
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

			// TIMEOUT HANDLING
			// No reply was received before the read deadline expired
			// Treat the probe as lost and continue to the next probe
			if err != nil {
				fmt.Printf("\nHop %d Probe %d\n", ttl, i+1)
				fmt.Println("Request timed out.")
				continue
			}

			//FAULT INJECTION
			// Simulate packet loss and delay to create degraded routes
			if faults.Enabled {

				// Simulate packet loss
				if rand.Float64() < faults.DropRate {
					fmt.Printf("\nHop %d Probe %d\n", ttl, i+1)
					fmt.Println("Simulating packet loss.")
					fmt.Println("Request timed out.")

					continue
				}

				// Simulate network delay
				if faults.Delay > 0 {

					delay := time.Duration(rand.Int63n(int64(faults.Delay)))

					fmt.Printf("\nSimulating delay of %v for Hop %d Probe %d\n", delay, ttl, i+1)

					time.Sleep(delay)
				}
			}

			// RTT TIMING
			// Measure round-trip time
			elapsed := time.Since(start)

			fmt.Printf("\nHop %d Probe %d\n", ttl, i+1)

			fmt.Println("RTT:", elapsed)

			// Print raw packet bytes in hexadecimal
			//fmt.Printf("% X\n", buf[:numRead])

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

				// Reverse DNS lookup used for reporting
				//Timing has already been completed before hostname resolution
				resolvedHost := ResolveHost(peer.String())

				fmt.Printf("Received Time Exceeded message: TTL: %d| Reply From %s| Time Elapsed: %v|\n", ttl, resolvedHost, elapsed)

				// HOP LOOKUP
				// Check whether this TTL already exists in the route.
				// If it does, append another RTT measurement.
				found := false

				for k := range hops {
					if hops[k].TTL == ttl {
						hops[k].RTTs = append(hops[k].RTTs, elapsed)

						hops[k].Host = resolvedHost
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
						Host:             resolvedHost,
						IPAddress:        peer.String(),
						SuccessfulProbes: 1,
					})
				}

			//DESTINATION REPLY
			//Host receives Echo Request and replies with Echo Reply message
			case ipv4.ICMPTypeEchoReply:

				resolvedHost := ResolveHost(peer.String())

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
				if echoReply.Seq != ttl*10+i {
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

						hops[k].Host = resolvedHost
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
						Host:             resolvedHost,
						IPAddress:        peer.String(),
						SuccessfulProbes: 1,
					})
				}

				// ROUTE COMPLETE
				// Destination has responded successfully
				// Set flag so traceroute stops after finishing the current hop
				destinationReached = true
				fmt.Println("Destination reached!")
				fmt.Printf("Echo Reply ID: %d, Seq: %d, Data: %s\n", echoReply.ID, echoReply.Seq, string(echoReply.Data))
				fmt.Printf("TTL: %d\n", ttl)

			default:
				fmt.Printf("Received ICMP message of type %v\n", rm.Type)
			}

			time.Sleep(1 * time.Second)

		}
		// Stop increasing TTL once the destination has been reached
		if destinationReached {
			break
		}
	}

	println("\nTraceroute complete.")

	return hops
}
