package sample

import "math"

func BankRound(x float64, precision int) float64 {
	scale := math.Pow10(precision)
	xScaled := x * scale
	rounded := math.Round(xScaled)
	diff := math.Abs(xScaled - rounded)
	if diff == 0.5 {
		if int64(rounded)%2 != 0 {
			rounded -= math.Copysign(1, rounded)
		}
	}
	return rounded / scale
}
