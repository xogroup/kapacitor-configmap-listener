package kapacitor

import (
	"time"
)

// Watch polls Kapacitor to check if the instance has been replaced and renews desired state when necessary
func Watch(handler func()) chan struct{} {
	ticker := time.NewTicker(5 * time.Second)

	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				handler()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return stop
}
