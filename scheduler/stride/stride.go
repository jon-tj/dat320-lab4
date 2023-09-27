package stride

import (
	"dat320/lab4/scheduler/cpu"
	"dat320/lab4/scheduler/job"
	"time"
)

type stride struct {
	queue job.Jobs
	cpu   *cpu.CPU
}

func New(cpus []*cpu.CPU, quantum time.Duration) *stride {
	if len(cpus) != 1 {
		panic("take the L")
	}
	return &stride{
		queue: make(job.Jobs, 0),
		cpu:   cpus[0],
	}
}

func (s *stride) Add(job *job.Job) {
	s.queue = append(s.queue, job)
}

// Tick runs the scheduled jobs for the system time, and returns
// the number of jobs finished in this tick. Depending on scheduler requirements,
// the Tick method may assign new jobs to the CPU before returning.
func (s *stride) Tick(systemTime time.Duration) int {
	jobsFinished := 0
	if s.cpu.IsRunning() {
		if s.cpu.Tick() {
			jobsFinished += 1
			s.reassign()
		}
	} else {
		s.reassign()
	}
	return jobsFinished
}

// reassign assigns a job to the cpu
func (s *stride) reassign() {
	job := s.getNewJob()
	s.cpu.Assign(job)
}

// getNewJob finds a new job to run on the CPU, removes the job from the queue and returns the job
func (s *stride) getNewJob() *job.Job {

	if len(s.queue) == 0 {
		return nil
	}
	job := s.queue[MinPass(s.queue)]
	return job
}

// minPass returns the index of the job with the lowest pass value.
func MinPass(theJobs job.Jobs) int {
	lowest := 0
	val := theJobs[0].Pass
	for i, job := range theJobs {
		if job.Pass < val {
			val = job.Pass
			lowest = i
		}
	}
	return lowest
}
