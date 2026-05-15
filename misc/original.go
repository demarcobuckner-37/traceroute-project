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

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a host: ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)

	fmt.Println("You entered:\n", host)
	fmt.Println("IP Adresses:\n")

	tm := time.Now()
	dnsTime := time.Now()

	ips, err := net.LookupIP(host)
	dnsElapsed := time.Since(dnsTime)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, ip := range ips {
		fmt.Println(ip)
	}

	fmt.Println("_________________")

	address := host + ":80"

	connTime := time.Now()

	conn, err := net.Dial("tcp", address)
	connElapsed := time.Since(connTime)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Connection successful to", conn, "on", host)
		defer conn.Close()
	}

	fmt.Println("Connection successful to:", address)

	fmt.Printf("Start time:%s\n", tm)
	elapsed := time.Since(tm)
	tm1 := time.Now()
	fmt.Printf("End time:%s\n", tm1)
	fmt.Println("Total Time:", elapsed)
	fmt.Println("DNS lookup time:", dnsElapsed)
	fmt.Println("Connection time:", connElapsed)

}
