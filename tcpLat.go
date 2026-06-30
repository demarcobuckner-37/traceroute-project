package main

import (
	"fmt"
	"net"
	"time"
)

type TCPResult struct {
	//Target host that was tested
	Host string

	StartTime time.Time

	EndTime time.Time

	TotalRunTime time.Duration

	// Time required to resolve hostname to IP address
	DNSLookupTime time.Duration

	// Average TCP connection latency across successful probes
	AvgConnectionTime time.Duration

	// Fastest successful TCP connection
	MinConnectionTime time.Duration

	// Slowest successful TCP connection
	MaxConnectionTime time.Duration

	// Number of successful TCP connections established
	SuccessfulConnections int

	// Percentage of failed connection attempts
	PacketLoss float64
}

func RunLatencyTool(host string) TCPResult {

	// Display all IP addresses returned by DNS lookup
	fmt.Println("IP Addresses:")

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
		return TCPResult{}
	}
	// OUTPUT FORMATTING
	// Print all resolved IP addresses
	for _, ip := range ips {
		fmt.Println(ip)
	}

	fmt.Println("\nTCP Connections:")

	// DESTINATION ADDRESS
	// Combine hostname and HTTPS port into a valid network address
	address := net.JoinHostPort(host, "443")

	var totalConnTime time.Duration
	var avgConnTime time.Duration

	var minConnTime time.Duration
	var maxConnTime time.Duration

	var successfulConnections int

	// Perform multiple TCP probes to measure latency
	for i := 0; i < 3; i++ {

		// Start connection timing
		connTime := time.Now()

		// Attempt TCP connection to destination with a timeout of 3 seconds
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)

		// Measure connection latency
		connElapsed := time.Since(connTime)

		fmt.Println("Connection time:", connElapsed)

		//CONNECTION STATISTICS
		if err == nil {

			if successfulConnections == 0 {
				//First succesful probe is baseline
				minConnTime = connElapsed
				maxConnTime = connElapsed
			} else {

				//Update minimum connection latency
				if connElapsed < minConnTime {
					minConnTime = connElapsed
				}

				//Update maximum connection latency
				if connElapsed > maxConnTime {
					maxConnTime = connElapsed
				}
			}

			//Track successful connections and total latency
			successfulConnections++
			totalConnTime += connElapsed

			fmt.Println("Connection successful to", conn, "on", host)

			conn.Close()

		} else {
			fmt.Println("Error:", err)
		}

	}

	// AVERAGE CONNECTION TIME
	// Calculate mean TCP connection latency using only successful probes
	if successfulConnections > 0 {
		avgConnTime = totalConnTime / time.Duration(successfulConnections)
	}

	// LATENCY TIMING
	// Measure total runtime
	elapsed := time.Since(tm)

	// Record ending time
	tm1 := time.Now()

	// PACKET LOSS CALCULATION
	// Percentage of connection attempts that failed
	packetLoss := float64(3-successfulConnections) / 3 * 100

	return TCPResult{
		Host:                  host,
		StartTime:             tm,
		EndTime:               tm1,
		TotalRunTime:          elapsed,
		DNSLookupTime:         dnsElapsed,
		AvgConnectionTime:     avgConnTime,
		MinConnectionTime:     minConnTime,
		MaxConnectionTime:     maxConnTime,
		SuccessfulConnections: successfulConnections,
		PacketLoss:            packetLoss,
	}

}
