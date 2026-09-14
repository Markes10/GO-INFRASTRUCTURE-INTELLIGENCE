package anomaly

import (
	"math"
	"sync"
)

type EWMAAnomalyDetector struct {
	mu           sync.Mutex
	alpha        float64
	ewmaMean     float64
	ewmaVariance float64
	threshold    float64
	initialized  bool
}

func NewEWMAAnomalyDetector(alpha float64, thresholdSigmas float64) *EWMAAnomalyDetector {
	return &EWMAAnomalyDetector{
		alpha:     alpha,
		threshold: thresholdSigmas,
	}
}

func (d *EWMAAnomalyDetector) Process(val float64) (isAnomaly bool, zScore float64, mean float64, stdDev float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.initialized {
		d.ewmaMean = val
		d.ewmaVariance = 1.0
		d.initialized = true
		return false, 0, val, 1.0
	}

	diff := val - d.ewmaMean
	variance := d.ewmaVariance
	if variance < 1e-6 {
		variance = 1e-6
	}
	stdDev = math.Sqrt(variance)
	zScore = math.Abs(diff) / stdDev

	isAnomaly = zScore > d.threshold

	// Update EWMA moments
	d.ewmaMean = d.alpha*val + (1.0-d.alpha)*d.ewmaMean
	d.ewmaVariance = (1.0-d.alpha)*(d.ewmaVariance + d.alpha*diff*diff)

	return isAnomaly, zScore, d.ewmaMean, stdDev
}
