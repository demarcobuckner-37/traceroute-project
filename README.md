# Traceroute & Network Diagnostics Toolkit

## Overview

This project is a command-line network diagnostics toolkit written in **Go (Golang)**. It was developed as part of a directed study to gain hands-on experience with low-level networking concepts, systems programming, and network analysis techniques commonly encountered in network engineering and cybersecurity.

The toolkit combines several core networking utilities into a single application, including:

- DNS resolution
- TCP latency testing
- ICMP ping operations
- Traceroute path discovery
- Hop-by-hop route analysis
- Network bottleneck detection
- Fault injection and resilience testing
- Route comparison
- JSON export

The project emphasizes understanding how network traffic flows across the Internet while providing practical experience working with raw sockets, packet structures, and network performance metrics.

---

# Features

## Traceroute Engine

- Discovers network paths using **TTL manipulation**
- Sends multiple probes per hop
- Handles:
  - ICMP Time Exceeded messages
  - ICMP Echo Replies
- Records RTT measurements for each probe
- Tracks successful and failed responses
- Dynamically stores route data for later analysis

---

## Hop Analytics

The analytics engine processes collected traceroute data and calculates:

- Average RTT
- Minimum RTT
- Maximum RTT
- Packet Loss
- Jitter

Route analysis is separated from data collection, allowing metrics to be recalculated or expanded independently.

---

## Bottleneck Detection System

Analyzes each hop and identifies potential network bottlenecks using:

- RTT increases between adjacent hops
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

The system also identifies the highest-scoring bottleneck along the route.

---

## Fault Injection System

Simulates degraded network conditions for testing and analysis by introducing:

- Artificial packet loss
- Artificial network latency

This feature provides a simple form of chaos testing, allowing route behavior and analytics to be observed under adverse network conditions.

---

## Route Comparison

Compares a baseline traceroute against a fault-injected traceroute to identify:

- Latency differences
- Hop-by-hop performance changes
- Potential bottleneck shifts
- Overall route behavior under simulated faults

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
  - Jitter
- Validates:
  - Process IDs
  - Sequence Numbers

---

## TCP Latency Tool

- Resolves hostnames using DNS
- Measures DNS lookup time
- Establishes TCP connections
- Calculates:
  - Average connection time
  - Minimum connection time
  - Maximum connection time
  - Connection success rate
  - Packet loss percentage

---

## JSON Export

Export diagnostic data for future analysis or reporting:

- Route A (Baseline Traceroute)
- Route B (Fault-Injected Traceroute)
- Ping Results
- TCP Latency Results

---

# Technologies Used

- Go (Golang)
- ICMP Networking
- Raw Socket Programming
- IPv4 Packet Manipulation
- DNS Resolution
- TCP Networking
- JSON Serialization
- Command-Line Interface (CLI) Development

---

# Key Networking Concepts Explored

- Time To Live (TTL)
- ICMP Packet Structure
- Echo Requests and Echo Replies
- ICMP Time Exceeded Messages
- DNS Resolution
- TCP Connection Establishment
- Round-Trip Time (RTT)
- Packet Loss
- Jitter
- Route Discovery
- Network Bottleneck Detection
- Network Path Analysis
- Fault Injection
- Basic Chaos Engineering Concepts

---

# Project Architecture

The project is divided into independent components to simplify testing and future expansion.

### `RunTraceroute()`

Responsible for:

- Route discovery
- ICMP packet transmission
- ICMP response handling
- Hop data collection

### `RunAnalyzeHops()`

Responsible for:

- Average RTT calculation
- Minimum RTT calculation
- Maximum RTT calculation
- Packet loss calculation
- Jitter calculation

### `RunBottleneckDetection()`

Responsible for:

- RTT increase analysis
- Bottleneck scoring
- Severity classification
- Bottleneck flag generation

### `RouteSummary()`

Responsible for:

- Route statistics generation
- High-level route reporting

### `RunRouteComparison()`

Responsible for:

- Comparing baseline and fault-injected routes
- Identifying latency differences
- Highlighting route behavior changes

### `PrintHopAnalysis()`

Responsible for:

- Detailed hop reporting
- Bottleneck reporting
- Analytics presentation

### `RunPingTool()`

Responsible for:

- ICMP ping diagnostics
- RTT and packet loss measurement

### `RunLatencyTool()`

Responsible for:

- DNS timing measurements
- TCP connection latency analysis

### `ExportJSON()`

Responsible for:

- Exporting traceroute, ping, and TCP latency results to JSON files

---

# Current Project Status

## Implemented

- DNS Resolution
- TCP Latency Measurement
- ICMP Ping Tool
- Traceroute Route Discovery
- Hop Analytics
- RTT Statistics
- Packet Loss Analysis
- Jitter Analysis
- Bottleneck Detection
- Fault Injection
- Route Comparison
- Route Summaries
- JSON Export
- Dynamic Hop Storage Using Go Structs
- Interactive CLI Menu System

---

# Planned Features

- Long-Term Route Monitoring
- Historical Route Tracking
- Network Instability Trend Analysis
- Exportable Report Generation
- Enhanced Visualization and Reporting

---

# Running the Project

Run the project using:

```bash
sudo go run .
```

> **Note:** Administrator (root) privileges are required for raw ICMP socket operations.

---

# Example Targets

- google.com
- github.com
- cloudflare.com
- espn.com

---

# Project Goals

The primary goals of this project are to gain practical experience with:

- Network Troubleshooting
- Traceroute Implementation
- ICMP and TCP Networking
- Raw Socket Programming
- Network Performance Analysis
- Bottleneck Detection Techniques
- Fault Injection and Resilience Testing
- Go Networking Libraries
- Systems Programming
- CCNA and Network Engineering Concepts

---

# Future Direction

This project serves as a learning platform for exploring advanced networking and systems concepts. Future iterations may include long-term route monitoring, historical route analysis, enhanced reporting capabilities, and additional diagnostic utilities to expand the toolkit into a more comprehensive network analysis framework.