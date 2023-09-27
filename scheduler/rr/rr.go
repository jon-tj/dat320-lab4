package rr

import (
	"dat320/lab4/scheduler/cpu"
	"dat320/lab4/scheduler/job"
	"time"
)

type roundRobin struct {
	queue       job.Jobs
	cpu         *cpu.CPU
	index       int
	quantum     time.Duration
	quantumTime time.Duration //brother i do not know the way
}

func New(cpus []*cpu.CPU, quantum time.Duration) *roundRobin {
	if len(cpus) != 1 {
		panic("rr scheduler supports only a single CPU")
	}
	return &roundRobin{
		cpu:     cpus[0],
		queue:   make(job.Jobs, 0),
		index:   0,
		quantum: quantum,
	}
}

func (rr *roundRobin) Add(job *job.Job) {
	rr.queue = append(rr.queue, job)
}

// Tick runs the scheduled jobs for the system time, and returns
// the number of jobs finished in this tick. Depending on scheduler requirements,
// the Tick method may assign new jobs to the CPU before returning.
func (rr *roundRobin) Tick(systemTime time.Duration) int {
	rr.quantumTime += systemTime
	if rr.quantumTime >= rr.quantum {
		rr.index += 1
	}
	jobsFinished := 0
	if rr.cpu.IsRunning() {
		if rr.cpu.Tick() {
			jobsFinished++
			rr.reassign()
		}
	} else {
		// CPU is idle, find new job in own queue
		rr.reassign()
	}
	return jobsFinished
}
func remove(slice job.Jobs, s int) job.Jobs { //bscly just skip the middle element
	return append(slice[:s], slice[s+1:]...) //"..." means we unpack the slice
}

// reassign assigns a job to the cpu
func (rr *roundRobin) reassign() {
	nxtJob := rr.getNewJob()
	rr.cpu.Assign(nxtJob)

}

// getNewJob finds a new job to run on the CPU, removes the job from the queue and returns the job
func (rr *roundRobin) getNewJob() *job.Job {

	if len(rr.queue) == 0 {
		return nil
	}
	rr.index = rr.index % len(rr.queue)
	job := rr.queue[rr.index]
	return job
}
