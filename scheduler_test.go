package scheduler

import (
	"testing"
	"time"
)

func TestTes(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sched := NewScheduler(
				func() {
					now := time.Now()
					println(now.Format(time.DateTime))
				},
				func() time.Time {
					now := time.Now()
					return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+5, 0, now.Location())
				},
			)
			go sched.start()
		})
	}

	time.Sleep(30 * time.Second)
}
