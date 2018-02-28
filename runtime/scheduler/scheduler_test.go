package scheduler

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Many tests schedule a job for every second, and then wait at most a second
// for it to run.  This amount is just slightly larger than 1 second to
// compensate for a few milliseconds of runtime.
const OneSecond = 1*time.Second + 10*time.Millisecond

type testJob struct {
	wg   *sync.WaitGroup
	name string
}

func (t testJob) Run() {
	fmt.Println("running ", t.name)
	t.wg.Done()
}

// Simple test using Runnables.
func TestJob(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	Init()
	AddJob("0 0 0 30 Feb ?", testJob{wg, "job0"})
	AddJob("0 0 0 1 1 ?", testJob{wg, "job1"})
	AddJob("* * * * * ?", testJob{wg, "job2"})
	AddJob("1 0 0 1 1 ?", testJob{wg, "job3"})
	Schedule(5*time.Second+5*time.Nanosecond, testJob{wg, "job4"})
	Schedule(5*time.Minute, testJob{wg, "job5"})
	fmt.Println("#############################")

	select {
	case <-time.After(OneSecond):
		t.FailNow()
	case <-wait(wg):
	}

}

// Add a job, start cron, expect it runs.
func TestAddBeforeRunning(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	Init()
	Addfunc("* * * * * ?", func() { wg.Done(); fmt.Println("okok") })

	// Give cron 2 seconds to run our job (which is always activated).
	select {
	case <-time.After(OneSecond):
		t.Fatal("expected job runs")
	case <-wait(wg):
	}
}
func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}
