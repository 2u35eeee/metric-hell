package game

import "math"

func CompletionPercent(s State, totalNodes int) float64 {
	if totalNodes <= 0 {
		return 0
	}
	value := float64(len(s.CompletedNodes)) / float64(totalNodes) * 100
	if value > 99.9 {
		return 99.9
	}
	return math.Round(value*10) / 10
}
