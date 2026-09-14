package anomaly

import (
	"testing"
)

func TestEWMAAnomalyDetector_NormalStream(t *testing.T) {
	detector := NewEWMAAnomalyDetector(0.15, 3.5)

	// Stream baseline normal metrics around 50.0
	for i := 0; i < 50; i++ {
		isAnomaly, zScore, _, _ := detector.Process(50.0 + float64(i%3))
		if isAnomaly {
			t.Fatalf("Unexpected anomaly flagged during normal baseline warmup at step %d: zScore=%.2f", i, zScore)
		}
	}
}

func TestEWMAAnomalyDetector_DetectsSpike(t *testing.T) {
	detector := NewEWMAAnomalyDetector(0.15, 3.5)

	// Warmup
	for i := 0; i < 30; i++ {
		detector.Process(40.0)
	}

	// Inject massive spike
	isAnomaly, zScore, _, _ := detector.Process(800.0)
	if !isAnomaly {
		t.Fatalf("Failed to detect 800.0 spike: zScore=%.2f (threshold=3.5)", zScore)
	}
}
