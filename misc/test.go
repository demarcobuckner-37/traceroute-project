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
