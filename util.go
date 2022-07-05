package main

import "fmt"

func FormatMinutes(m int) string {
	s := fmt.Sprintf("%dh", m/60)
	if m%60 > 0 {
		s += fmt.Sprintf("%dm", m%60)
	}

	return s
}
