package main

import (
	"container/list"
	"sync"
	"time"
)

// OPCTimer is used to signal read and write operation after a defined time interval
type OPCTimer struct {
	interval      time.Duration // Duration of an interval
	intervalForce time.Duration // Duration of force read and write
	trackers      *list.List    // List of all values to be tracked
}

// AddTracker adds a tracker to tracker list
func (opcTimer *OPCTimer) AddTracker(tracker *Tracker) {
	opcTimer.trackers.PushBack(tracker)
}

// Run starts the timer interval.
func (opcTimer *OPCTimer) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		forceTimerDuration := opcTimer.intervalForce
		timerDuration := opcTimer.interval
		var forceTimer <-chan time.Time
		if forceTimerDuration > 0 {
			forceTimer = time.NewTimer(forceTimerDuration).C
		}
		var timer <-chan time.Time
		if timerDuration > 0 {
			timer = time.NewTimer(timerDuration).C
		}
		for {
			select {
			case <-forceTimer:
				forceTimer = time.NewTimer(forceTimerDuration).C
				printf("**** FORCE!!! ****\n")
				for e := opcTimer.trackers.Front(); e != nil; e = e.Next() {
					if e.Value.(*Tracker).force == false {
						continue
					}
					success, data := e.Value.(*Tracker).ReadValue()
					if success {
						e.Value.(*Tracker).WriteValue(data)
					}
				}
				printf("**** **** ****\n")
			case <-timer:
				timer = time.NewTimer(timerDuration).C
				printf("**** TIMER ****\n")
				for e := opcTimer.trackers.Front(); e != nil; e = e.Next() {
					success, data := e.Value.(*Tracker).ReadValue()
					if success {
						if e.Value.(*Tracker).lastValue != data.Value {
							e.Value.(*Tracker).force = false
							e.Value.(*Tracker).WriteValue(data)
							e.Value.(*Tracker).lastValue = data.Value
						} else {
							printf("%s not changed\n", data.Name)
							e.Value.(*Tracker).force = true
						}
					} else {
						e.Value.(*Tracker).force = true
					}
				}
				printf("**** **** ****\n")
			}
		}
	}()
}
