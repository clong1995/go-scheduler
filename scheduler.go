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

// NewScheduler 执行器
// nextExecution 函数一定要写好不然很危险，或者使用提供的安全函数 Every EveryHourly EveryDaily EveryWeekly
func NewScheduler(process func(), nextExecution func() time.Time) *Scheduler {
	return &Scheduler{
		stopCh:        make(chan struct{}),
		process:       process,
		nextExecution: nextExecution,
	}
}

func (s *Scheduler) Start() {
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

func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

/*func (s *Scheduler) nextExecution() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location()).Add(-1 * time.Nanosecond)
}*/

// EveryDaily 安全的执行每天的某时刻的任务
func EveryDaily(hour, minute, second int) func() time.Time {
	return func() time.Time {
		now := time.Now()

		next := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			hour,
			minute,
			second,
			0,
			now.Location(),
		)

		if !next.After(now) {
			next = next.AddDate(0, 0, 1)
		}

		return next
	}
}

func Every(interval time.Duration) func() time.Time {
	if interval < 0 {
		if interval == time.Duration(-1<<63) {
			interval = time.Duration(1<<63 - 1)
		} else {
			interval = -interval
		}
	}
	return func() time.Time {
		return time.Now().Add(interval)
	}
}

// EveryHourly 每小时的 30 分 00 秒执行
// scheduler.EveryHourly(30, 0)
func EveryHourly(minute, second int) func() time.Time {
	return func() time.Time {
		now := time.Now()

		next := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			minute,
			second,
			0,
			now.Location(),
		)

		if !next.After(now) {
			next = next.Add(time.Hour)
		}

		return next
	}
}

// EveryWeekly 每周的周几的某时刻
func EveryWeekly(weekday time.Weekday, hour, minute, second int) func() time.Time {
	return func() time.Time {
		now := time.Now()

		// 当前是星期几
		daysUntil := (int(weekday) - int(now.Weekday()) + 7) % 7

		next := time.Date(
			now.Year(),
			now.Month(),
			now.Day()+daysUntil,
			hour,
			minute,
			second,
			0,
			now.Location(),
		)

		// 同一天但时间已经过去
		if !next.After(now) {
			next = next.AddDate(0, 0, 7)
		}

		return next
	}
}
