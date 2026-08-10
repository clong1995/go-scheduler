package scheduler

import (
	"fmt"
	"sync"
	"time"

	pcolor "github.com/clong1995/go-ansi-color"
)

type Scheduler struct {
	stopCh        chan struct{}
	wg            sync.WaitGroup
	process       func()
	nextExecution func() time.Time
}

func NewScheduler(process func(), nextExecution func() time.Time) *Scheduler {
	return &Scheduler{
		stopCh:        make(chan struct{}),
		process:       process,
		nextExecution: nextExecution,
	}
}

func (s *Scheduler) start() {
	for {
		next := s.nextExecution()
		timer := time.NewTimer(time.Until(next))
		select {
		case <-s.stopCh:
			timer.Stop()
			return
		case <-timer.C:
			s.wg.Go(func() {
				defer func() {
					if r := recover(); r != nil {
						pcolor.PrintError("任务 panic: %v", fmt.Errorf("%v", r))
					}
				}()
				s.process()
				pcolor.Print("任务结束")
			})
		}
	}
}

func (s *Scheduler) stop() {
	close(s.stopCh)
	s.wg.Wait()
}

/*func (s *Scheduler) nextExecution() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location()).Add(-1 * time.Nanosecond)
}*/
