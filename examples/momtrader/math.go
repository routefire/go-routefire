package main

import "math"

func mean(xs []float64) float64 {
	s := 0.0
	for _, x := range xs {
		s += x
	}
	return s / float64(len(xs))
}

func stdDev(xs []float64) float64 {
	s := 0.0
	mu := mean(xs)
	for _, x := range xs {
		y := x - mu
		z := y*y
		s += z
	}
	return math.Sqrt(s/float64(len(xs)))
}