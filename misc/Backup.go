package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	// USER INPUT
	// Create a reader so the user can type a hostname
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a host: ")

	// Read user input until Enter is pressed
	host, _ := reader.ReadString('\n')

	// Remove extra spaces/newline characters
	host = strings.TrimSpace(host)

	// OUTPUT FORMATTING
	// Print basic information to the terminal
	fmt.Println("You entered:\n", host)
	fmt.Println()
	fmt.Println("IP Adresses:")

	//LATENCY TIMING
	// Start overall program timer
	tm := time.Now()
	// Start DNS timing
	dnsTime := time.Now()

	// DNS RESOLUTION
	// Resolve hostname into IP addresses
	ips, err := net.LookupIP(host)
	dnsElapsed := time.Since(dnsTime)

	// ERROR HANDLING
	// Stop program if DNS lookup fails
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// OUTPUT FORMATTING
	// Print all resolved IP addresses
	for _, ip := range ips {
		fmt.Println(ip)
	}

	fmt.Println("\nTCP Connectins:")

	// TCP CONNECTIONS
	// Add port 80 (HTTP) to hostname
	address := host + ":80"

	var totalConnTime time.Duration
	var AvgConnTime time.Duration

	// Perform multiple TCP probes to measure latency
	for i := 0; i < 3; i++ {

		// Start connection timing
		connTime := time.Now()

		// Attempt TCP connection to destination with a timeout of 3 seconds
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)

		// Measure connection latency
		connElapsed := time.Since(connTime)
		fmt.Println("Connection time:", connElapsed)

		// ERROR HANDLING
		// Check if TCP connection succeeded
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			// OUTPUT FORMATTING
			// Print successful connection info
			fmt.Println("Connection successful to", conn, "on", host)

			// Close connection after each probe
			conn.Close()
		}

		// Update running total and average latency
		totalConnTime += connElapsed
		AvgConnTime = (totalConnTime) / time.Duration(i+1)
	}

	fmt.Println("\nConnection successful to:", address)

	// LATENCY TIMING
	// Measure total runtime
	elapsed := time.Since(tm)

	// Record ending time
	tm1 := time.Now()

	// OUTPUT FORMATTING
	// Print timing information
	fmt.Printf("\nTiming Information:\n")

	fmt.Printf("Start time:%s\n", tm)
	fmt.Printf("End time:%s\n", tm1)

	fmt.Printf("Total Runtime:%s\n", elapsed)

	fmt.Printf("DNS Lookup time:%s\n", dnsElapsed)

	fmt.Printf("Average Connection Time: %s\n", AvgConnTime)

	fmt.Printf("Total Probe time: %s\n", totalConnTime)

}

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a host: ")

	// Read user input until Enter is pressed
	host, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	// Remove extra spaces/newline characters
	host = strings.TrimSpace(host)

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

	var totalRTT time.Duration

	var successfulProbes int
	var minTime time.Duration
	var maxTime time.Duration

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

		// Read incoming ICMP reply packet

		err = c.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			fmt.Println(err)
			return
		}
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

		fmt.Printf("\nProbe %d\n", i)

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
			return
		}

		fmt.Println("ICMP Type:", rm.Type)

		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:

			echoReply := rm.Body.(*icmp.Echo)

			if echoReply.ID != os.Getpid()&0xffff {
				fmt.Println("Received reply for another process")
				continue
			}
			if echoReply.Seq != i {
				fmt.Printf("Received out-of-order Echo Reply: Seq %d\n", echoReply.Seq)
				continue
			}
			fmt.Println("Received Echo Reply")

			fmt.Printf("Echo Reply ID: %d, Seq: %d, Data: %s\n", echoReply.ID, echoReply.Seq, string(echoReply.Data))

		default:
			fmt.Printf("Received ICMP message of type %v\n", rm.Type)
		}

		time.Sleep(1 * time.Second)
		totalRTT += elapsed
	}
	if successfulProbes > 0 {
		avgRTT := totalRTT / time.Duration(successfulProbes)
		fmt.Printf("\nAverage RTT over %d probes: %s\n", successfulProbes, avgRTT)
		fmt.Printf("Minimum RTT: %s\n", minTime)
		fmt.Printf("Maximum RTT: %s\n", maxTime)
	}
	packetLoss := float64(3-successfulProbes) / 3 * 100
	fmt.Printf("Packet Loss: %.2f%%\n", packetLoss)

}
