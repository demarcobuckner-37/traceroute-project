package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	// USER INPUT
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a host: ")

	// Read user input until Enter is pressed
	host, _ := reader.ReadString('\n')

	// Remove extra spaces/newline characters
	host = strings.TrimSpace(host)

	// RunLatencyTool(host)
	//RunPingTool(host)
	hops := RunTraceroute(host)
	RunAvgRTT(hops)

}
