package main

import "fmt"

func formatInKilometer(raw float64) string {
	return fmt.Sprintf("%.2f", raw/1000)
}

func formatStravaTime(t uint) string {
	sec := (t - uint(t/60)*60)
	min := (t - sec) / 60
	hour := uint(min / 60)
	min -= hour * 60
	return fmt.Sprintf("%d:%02d:%02d", hour, min, sec)
}
