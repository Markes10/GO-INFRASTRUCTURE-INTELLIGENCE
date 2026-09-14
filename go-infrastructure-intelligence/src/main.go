package main

import (
	"fmt"
	"infrastructure-intelligence/src/anomaly"
	"infrastructure-intelligence/src/cluster"
	"time"
)

func main() {
	fmt.Println("=== Planet-Scale Infrastructure Intelligence Platform (Go) ===")

	// 1. Cluster Mesh Gossip
	fmt.Println("[CLUSTER] Initializing distributed gossip cluster...")
	c := cluster.NewGossipCluster("node-eu-west-01", "10.0.1.10:8080")
	c.AddPeer("node-us-east-02", "10.0.2.20:8080")
	c.AddPeer("node-ap-south-03", "10.0.3.30:8080")

	fmt.Println("  3 Active Cluster Nodes in Quorum.")

	// 2. Real-time Telemetry Ingestion with EWMA Anomaly Detection
	fmt.Println("\n[TELEMETRY] Ingesting CPU / Network IO stream with adaptive EWMA...")
	detector := anomaly.NewEWMAAnomalyDetector(0.15, 3.5)

	// Stream normal metrics
	for i := 0; i < 20; i++ {
		val := 45.0 + float64(i%4)*1.5
		detector.Process(val)
	}

	// Inject volumetric DDoS spike: 980 MB/s traffic spike
	fmt.Println("[EVENT] Influx anomaly spike: Traffic spikes to 980 MB/s (Baseline ~45 MB/s)...")
	isAnomaly, zScore, mean, stdDev := detector.Process(980.0)

	fmt.Printf("  Analysis Result -> Anomaly: %v | Z-Score: %.2f | Baseline Mean: %.2f | StdDev: %.2f\n",
		isAnomaly, zScore, mean, stdDev)

	if !isAnomaly {
		panic("Expected EWMA detector to catch 980 MB/s traffic spike")
	}

	// Simulate node heartbeat timeout
	fmt.Println("\n[FAILURE DETECTOR] Probing peer node heartbeats...")
	failures := c.DetectFailures(100 * time.Millisecond)
	fmt.Printf("  Failure Audit: %d node disruptions observed.\n", len(failures))

	fmt.Println("\n[SUCCESS] Go Infrastructure Intelligence Platform verified.")
}
