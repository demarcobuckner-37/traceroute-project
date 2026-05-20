package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func RunPingTool(host string) {

	fmt.Print("\nRunning Ping Tool...\n")

	// DNS RESOLUTION
	// Resolve hostname into an IPv4 address
	dstAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ICMP CONNECTION
	// Open raw ICMP packet listener on all interfaces
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close connection when program exits
	defer c.Close()

	//RTT TIMING
	//Track Total, Min, and Max Latency
	var totalRTT time.Duration
	var minTime time.Duration
	var maxTime time.Duration

	var successfulProbes int

	// ICMP PROBES
	// Send multiple Echo Requests to destination host
	for i := 1; i <= 3; i++ {

		// ICMP ECHO REQUEST
		// Build ICMP Echo Request packet
		m := icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				// Process ID helps identify our packet
				ID: os.Getpid() & 0xffff,

				// Sequence number for tracking requests
				Seq: i,

				// Packet payload data
				Data: []byte("WADDUP"),
			},
		}

		// Convert ICMP message into raw bytes
		b, err := m.Marshal(nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		// RTT TIMING
		// Start round-trip timer before sending packet
		start := time.Now()

		// SEND PACKET
		// Send ICMP Echo Request to destination
		_, err = c.WriteTo(b, dstAddr)
		if err != nil {
			fmt.Println(err)
			return
		}

		// RECEIVE REPLY
		// Create buffer for incoming packet data
		buf := make([]byte, 1024)

		//TIMEOUT HANDLING
		err = c.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			fmt.Println(err)
			return
		}
		// Read incoming ICMP reply packet
		numRead, peer, err := c.ReadFrom(buf)
		if err != nil {
			fmt.Printf("\nProbe %d\n", i)
			fmt.Println("Request timed out.")
			continue
		}
		successfulProbes++

		// RTT TIMING
		// Measure round-trip time
		elapsed := time.Since(start)

		//Initialize min and max times on first successful probe
		if successfulProbes == 1 {
			minTime = elapsed
			maxTime = elapsed

		} else {

			//Update min RTT
			if elapsed < minTime {
				minTime = elapsed
			}

			//Update max RTT
			if elapsed > maxTime {
				maxTime = elapsed
			}
		}

		fmt.Printf("\nProbe %d\n", i)

		fmt.Println("RTT:", elapsed)

		// Print responding host
		fmt.Printf("Received reply from %v\n", peer)

		// PACKET PARSING
		// Parse raw ICMP packet into structured message
		rm, err := icmp.ParseMessage(1, buf[:numRead])
		if err != nil {
			fmt.Println(err)
			return
		}

		//Print Parsed ICMP Type For Debugging
		fmt.Println("ICMP Type:", rm.Type)

		//ICMP MESSAGE HANDLING
		//Process different ICMP message types accordingly
		switch rm.Type {

		case ipv4.ICMPTypeEchoReply:

			//Convert message body to Echo type for access to ID, Seq, and Data
			echoReply := rm.Body.(*icmp.Echo)

			//Process Validation: Check if reply matches our process ID and sequence number
			if echoReply.ID != os.Getpid()&0xffff {
				fmt.Println("Received reply for another process")
				continue
			}

			//Sequence Validation: Check if reply is in expected order
			if echoReply.Seq != i {
				fmt.Printf("Received out-of-order Echo Reply: Seq %d\n", echoReply.Seq)

				continue
			}

			//Output Formatting: Print Echo Reply information, including ID, Sequence Number, and Payload Data
			fmt.Println("Received Echo Reply")

			fmt.Printf("Echo Reply ID: %d, Seq: %d, Data: %s\n", echoReply.ID, echoReply.Seq, string(echoReply.Data))

		default:

			fmt.Printf("Received ICMP message of type %v\n", rm.Type)
		}

		// Add delay between probes to avoid flooding the network
		time.Sleep(1 * time.Second)

		//Add Current RTT to Total RTT for average calculation later
		totalRTT += elapsed
	}

	//Calculate and Print Average RTT, Min RTT, Max RTT, and Packet Loss Percentage after all probes are complete
	if successfulProbes > 0 {
		avgRTT := totalRTT / time.Duration(successfulProbes)
		fmt.Printf("\nAverage RTT over %d probes: %s\n", successfulProbes, avgRTT)
		fmt.Printf("Minimum RTT: %s\n", minTime)
		fmt.Printf("Maximum RTT: %s\n", maxTime)
	}

	//Packet Loss Calculation: Calculate percentage of lost packets based on number of successful probes out of total probes sent
	packetLoss := float64(3-successfulProbes) / 3 * 100
	fmt.Printf("Packet Loss: %.2f%%\n", packetLoss)

}
