package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	// USER INPUT
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a host: ")

	// Read user input until Enter is pressed
	host, _ := reader.ReadString('\n')

	// Remove extra spaces/newline characters
	host = strings.TrimSpace(host)

	// ROUTE STORAGE
	// Store traceroute results for comparison and export
	var routeA []Hop
	var routeB []Hop

	// TOOL RESULTS
	// Store latest Ping and TCP Latency results for export
	var pingResult PingResult
	var tcpResult TCPResult

	//Continuosly display the menu and pocess user selections
	//until exit is selected
	for {

		// MENU DISPLAY
		fmt.Printf("\nCurrent Host: %s\n", host)

		fmt.Println("1. Run Traceroute")
		fmt.Println("2. Run Traceroute with Fault Injection")
		fmt.Println("3. Hop Analytics")
		fmt.Println("4. Export Route A JSON")
		fmt.Println("5. Export Route B JSON")
		fmt.Println("6. Change Host")
		fmt.Println("7. Run Ping")
		fmt.Println("8. Run TCP Latency Tool")
		fmt.Println("9. Export Ping JSON")
		fmt.Println("10. Export TCP JSON")
		fmt.Println("11. Exit")

		fmt.Print("\nSelect an option: \n")

		//Reads and validates menu input
		var choice int
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		_, err := fmt.Sscanf(input, "%d", &choice)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}

		// MENU HANDLER
		// Execute selected tool based on user choice
		switch choice {

		case 1:

			// BASELINE TRACEROUTE
			// Run normal route analysis without faults
			routeA = RunTraceroute(host, Fault{Enabled: false})
			routeA = RunAnalyzeHops(routeA)

		case 2:

			// FAULT-INJECTED TRACEROUTE
			// Simulate packet loss and delay for comparison testing
			routeB = RunTraceroute(host, Fault{
				Enabled:  true,
				DropRate: 0.3,
				Delay:    100 * time.Millisecond,
			})

			routeB = RunAnalyzeHops(routeB)

		case 3:

			// HOP ANALYTICS SUBMENU
			// Provides tools for summarizing, analyzing, and
			// comparing collected traceroute data.
			analyticsRunning := true

			for analyticsRunning {

				fmt.Println("\nHop Analytics:")
				fmt.Println("1. Summarize Route A")
				fmt.Println("2. Summarize Route B")
				fmt.Println("3. Analyze Route A")
				fmt.Println("4. Analyze Route B")
				fmt.Println("5. Compare Routes")
				fmt.Println("6. Exit Analytics")
				fmt.Println()

				// Read and validate submenu input
				var analyticsChoice int
				fmt.Print("Select an option: ")

				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				_, err := fmt.Sscanf(input, "%d", &analyticsChoice)
				if err != nil {
					fmt.Println("Please enter a valid number.")
					continue
				}

				// ANALYTICS MENU HANDLER
				switch analyticsChoice {

				case 1:

					// Display a high-level summary of Route A
					if len(routeA) == 0 {
						fmt.Println("\nNo Route A data available.")
						continue
					}

					RouteSummary(routeA)

				case 2:

					// Display a high-level summary of Route B
					if len(routeB) == 0 {
						fmt.Println("\nNo Route B data available.")
						continue
					}

					RouteSummary(routeB)

				case 3:

					// Run bottleneck detection and detailed hop
					// analysis on the baseline traceroute
					if len(routeA) == 0 {
						fmt.Println("\nNo Route A data available.")
						continue
					}

					RunBottleneckDetection(routeA)
					PrintHopAnalysis(routeA)

				case 4:

					// Run bottleneck detection and detailed hop
					// analysis on the fault-injected traceroute
					if len(routeB) == 0 {
						fmt.Println("\nNo Route B data available.")
						continue
					}

					RunBottleneckDetection(routeB)
					PrintHopAnalysis(routeB)

				case 5:

					// ROUTE COMPARISON
					// Compare baseline and fault-injected traceroute results
					if len(routeA) == 0 || len(routeB) == 0 {
						fmt.Println("Run Route A and Route B first.")
						continue
					}

					RunRouteComparison(routeA, routeB)

				case 6:

					// Return to the main application menu
					analyticsRunning = false

				default:
					fmt.Println("\nInvalid Output...")
				}

			}

		case 4:

			// EXPORT ROUTE A
			// Save the baseline traceroute results as JSON
			if len(routeA) == 0 {
				fmt.Println("No Route A data available.")
				break
			}

			err := ExportJSON("routeA.json", routeA)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Route A results exported to routeA.json")
			}

		case 5:

			// EXPORT ROUTE B
			// Save the fault-injected traceroute results as JSON
			if len(routeB) == 0 {
				fmt.Println("No Route B data available.")
				break
			}

			err := ExportJSON("routeB.json", routeB)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Route B results exported to routeB.json")
			}

		case 6:

			// HOST CHANGE
			// Update target host and clear previous results
			fmt.Print("Enter a new host: ")

			// Read user input until Enter is pressed
			host, _ = reader.ReadString('\n')

			// Remove extra spaces/newline characters
			host = strings.TrimSpace(host)

			fmt.Printf("Host changed to %s\n", host)

			// Clear stored results from previous host
			routeA = nil
			routeB = nil

			pingResult = PingResult{}
			tcpResult = TCPResult{}

		case 7:

			// PING TOOL
			// Run ICMP latency testing and display summary metrics
			pingResult = RunPingTool(host)

			fmt.Printf("\nPing Summary:\n")
			fmt.Printf("Host: %s\n", pingResult.Host)
			fmt.Printf("Average RTT: %v\n", pingResult.AvgRTT)
			fmt.Printf("Minimum RTT: %v\n", pingResult.MinRTT)
			fmt.Printf("Maximum RTT: %v\n", pingResult.MaxRTT)
			fmt.Printf("Jitter: %v\n", pingResult.Jitter)
			fmt.Printf("Packet Loss: %.2f%%\n", pingResult.PacketLoss)

		case 8:

			// TCP LATENCY TOOL
			// Measure DNS lookup and TCP connection performance
			tcpResult = RunLatencyTool(host)

			fmt.Printf("\nTCP Summary:\n")
			fmt.Printf("Host: %s\n", tcpResult.Host)
			fmt.Printf("DNS Lookup Time: %v\n", tcpResult.DNSLookupTime)
			fmt.Printf("Average Connection Time: %v\n", tcpResult.AvgConnectionTime)
			fmt.Printf("Minimum Connection Time: %v\n", tcpResult.MinConnectionTime)
			fmt.Printf("Maximum Connection Time: %v\n", tcpResult.MaxConnectionTime)
			fmt.Printf("Successful Connections: %d\n", tcpResult.SuccessfulConnections)
			fmt.Printf("Packet Loss: %.2f%%\n", tcpResult.PacketLoss)

		case 9:

			// EXPORT PING RESULTS
			// Save most recent ping statistics to JSON
			if pingResult.Host == "" {
				fmt.Println("No Ping data available.")
				break
			}

			err := ExportJSON("ping.json", pingResult)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Ping results exported to ping.json")
			}

		case 10:

			// EXPORT TCP RESULTS
			// Save most recent TCP latency statistics to JSON
			if tcpResult.Host == "" {
				fmt.Println("No TCP data available.")
				break
			}

			err := ExportJSON("tcp.json", tcpResult)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("TCP Latency results exported to tcp.json")
			}

		case 11:

			// PROGRAM EXIT
			// Terminate application
			fmt.Println("Exit..")
			return

		default:
			fmt.Println("Invalid option")
		}
		fmt.Println()
	}

}
