package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/siamak4mo/wpool"
)

func do(ID int) bool {
	// this funtion runs in each job
	// return true on success and false on failure
	time.Sleep(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)

	return true
}

func main() {
	var wg sync.WaitGroup

	const (
		JOB_C    = 12 // number of jobs
		WORKER_C = 2  // number of workers
		MRTC     = 3  // max retry number for each job
	)

	q := make(chan wpool.Job)

	go wpool.GOJOB(q, WORKER_C)
	wg.Add(JOB_C)

	for i := 0; i < JOB_C; i++ { // add new job
		q <- wpool.Job{
			i,     // job unique ID
			100,   // disable ttl
			false, // is not done
			do,    // job function struct
			MRTC,  // max retry number
			&wg}   // wait group
	}

	wg.Wait()
}
