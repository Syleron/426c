package main

import "time"

/**
Function to schedule the execution every x time as time.Duration.
*/
func scheduler(method func() bool, delay time.Duration) {
	for _ = range time.Tick(delay) {
		end := method()
		if end {
			break
		}
	}
}
