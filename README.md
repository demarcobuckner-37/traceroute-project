# Traceroute Project

## Overview

This project is a networking toolkit written in Go designed to explore low-level networking concepts related to:

- DNS resolution
- TCP latency testing
- ICMP ping operations
- Traceroute path discovery
- RTT analysis
- Network path analysis

The project is being developed as part of a directed study focused on networking fundamentals and CCNA-related concepts.

---

# Features

## TCP Latency Tool

- Resolves hostnames using DNS
- Establishes TCP connections
- Measures connection latency
- Calculates round-trip timing information

---

## ICMP Ping Tool

- Sends ICMP Echo Requests
- Receives ICMP Echo Replies
- Measures RTT values
- Tracks:
  - average RTT
  - minimum RTT
  - maximum RTT
  - packet loss
- Validates:
  - process IDs
  - sequence numbers

---

## Traceroute Tool

- Uses TTL manipulation to discover network routes
- Sends multiple probes per hop
- Handles:
  - ICMP Time Exceeded messages
  - ICMP Echo Replies
- Discovers intermediate routers between source and destination
- Stores route information using a custom:

```go
type Hop struct {
	TTL       int
	RTTs      []time.Duration
	Host      string
	IPAddress string
}

## RTT Analysis System

- Separates route collection from analysis logic
- Calculates average RTT per hop
- Dynamically handles variable RTT sample counts
- Supports future analysis extensions such as:
  - bottleneck detection
  - jitter analysis
  - route comparison
  - instability detection

---

# Technologies Used

- Go (Golang)
- ICMP networking
- Raw sockets
- IPv4 packet manipulation
- DNS resolution
- TCP networking

---

# Key Networking Concepts Explored

- TTL (Time To Live)
- ICMP packet structure
- Echo Requests and Echo Replies
- ICMP Time Exceeded messages
- Packet loss
- RTT measurement
- Route discovery
- Raw socket networking
- Network path analysis

---

# Current Project Status

Current implemented components:

- DNS resolution
- TCP latency measurement
- ICMP ping tool
- Traceroute route discovery
- RTT statistics
- Average RTT analysis per hop
- Dynamic hop storage using Go structs

---

# Planned Features

- Bottleneck detection
- Jitter calculations
- Route comparison
- Long-term route monitoring
- Network instability analysis
- Exportable route summaries

---

# Running the Project

Run the full project using:

```bash
sudo go run .

Administrator privileges are required for raw ICMP socket operations.

---

# Example Targets

- google.com
- github.com
- cloudflare.com

---

# Project Goals

The primary goal of this project is to gain deeper understanding of:

- networking fundamentals
- systems programming
- Go networking libraries
- ICMP/TCP behavior
- traceroute implementation techniques
- network analysis concepts
