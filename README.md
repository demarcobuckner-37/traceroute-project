# Traceroute Project

## Overview

This project is a networking toolkit written in Go designed to explore low-level networking concepts related to:

- DNS resolution
- TCP latency testing
- ICMP ping operations
- Traceroute path discovery
- RTT analysis
- Network path analysis
- Bottleneck detection

The project is being developed as part of a directed study focused on networking fundamentals, network analysis, and CCNA-related concepts.

---

# Features

## TCP Latency Tool

- Resolves hostnames using DNS
- Establishes TCP connections
- Measures connection latency
- Calculates connection timing statistics

---

## ICMP Ping Tool

- Sends ICMP Echo Requests
- Receives ICMP Echo Replies
- Measures RTT values
- Calculates:
  - Average RTT
  - Minimum RTT
  - Maximum RTT
  - Packet Loss
- Validates:
  - Process IDs
  - Sequence Numbers

---

## Traceroute Tool

- Uses TTL manipulation to discover network routes
- Sends multiple probes per hop
- Handles:
  - ICMP Time Exceeded messages
  - ICMP Echo Replies
- Discovers intermediate routers between source and destination
- Records RTT measurements for every hop
- Tracks successful and failed probes

---

## Hop Analysis System

Calculates performance metrics for every discovered hop:

- Average RTT
- Minimum RTT
- Maximum RTT
- Packet Loss
- Jitter

The analysis engine is separated from route collection to allow future expansion and testing.

---

## Bottleneck Detection System

Analyzes each hop and identifies potential network bottlenecks using:

- RTT increases between hops
- Packet loss thresholds
- Jitter thresholds
- Average RTT thresholds
- Maximum RTT thresholds

Each hop receives:

- Bottleneck Score
- Severity Classification

### Severity Levels

- None
- Minor
- Moderate
- Severe

The system also identifies the highest-scoring bottleneck within the route.

---

# Technologies Used

- Go (Golang)
- ICMP Networking
- Raw Sockets
- IPv4 Packet Manipulation
- DNS Resolution
- TCP Networking

---

# Key Networking Concepts Explored

- TTL (Time To Live)
- ICMP Packet Structure
- Echo Requests and Echo Replies
- ICMP Time Exceeded Messages
- Packet Loss
- Jitter
- RTT Measurement
- Route Discovery
- Raw Socket Networking
- Network Bottleneck Detection
- Network Path Analysis

---

# Project Architecture

The project is divided into independent components:

### RunTraceroute()

Responsible for:

- Route discovery
- ICMP packet transmission
- ICMP response handling
- Hop data collection

### RunAnalyzeHops()

Responsible for:

- Average RTT calculation
- Minimum RTT calculation
- Maximum RTT calculation
- Packet loss calculation
- Jitter calculation

### RunBottleneckDetection()

Responsible for:

- RTT increase analysis
- Bottleneck scoring
- Severity classification
- Bottleneck flag generation

### PrintHopAnalysis()

Responsible for:

- Reporting
- Route presentation
- Bottleneck reporting

---

# Current Project Status

## Implemented

- DNS Resolution
- TCP Latency Measurement
- ICMP Ping Tool
- Traceroute Route Discovery
- RTT Statistics
- Average RTT Analysis
- Minimum RTT Analysis
- Maximum RTT Analysis
- Packet Loss Analysis
- Jitter Analysis
- Bottleneck Detection
- Severity Classification
- Dynamic Hop Storage Using Go Structs

---

# Planned Features

- Route Summary Generation
- Route Comparison
- Long-Term Route Monitoring
- Network Instability Analysis
- Exportable Reports
- Historical Route Tracking

---

# Running the Project

Run the project using:

```bash
sudo go run .
```

Administrator privileges are required for raw ICMP socket operations.

---

# Example Targets

- google.com
- github.com
- cloudflare.com
- espn.com

---

# Project Goals

The primary goals of this project are to gain practical experience with:

- Network Troubleshooting Concepts
- Traceroute Implementation
- ICMP and TCP Networking
- Network Performance Analysis
- Bottleneck Detection Techniques
- Go Networking Libraries
- Systems Programming
- CCNA-Related Networking Concepts
