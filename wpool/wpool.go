package wpool

import (
	"log"
	"sync"
	"time"
)

type Job struct {
	Jid  int            // unique ID
	TTL  int            // time to leave (in milisecond) [set -1 to disable timeout]
	Done bool           // result of the job (False is failure)
	Jfun func(int) bool // struct of function and its params
	Mrtc int            // max retrying count [set -1 to disable it]
	Wg   *sync.WaitGroup
}

type worker struct {
	wid  int      // worker ID
	jobs chan Job // jobs queue
}

func (j *Job) DO() {
	if j.TTL < 0 {
		// without timeout
		j.Done = j.Jfun(j.Jid)
	} else if j.TTL >= 0 {
		// with timeout
		res := make(chan bool, 1)

		go func() {
			defer close(res)
			res <- j.Jfun(j.Jid)
		}()

		select {
		case <-time.After(time.Duration(j.TTL) * time.Millisecond):
			return
		case r := <-res:
			j.Done = r
			return
		}
	}
}

func (w *worker) work() {
	for {
		select {
		case j := <-w.jobs:
			log.Printf("__worker%d__  running: jobID=%d, ttl=%d\n", w.wid, j.Jid, j.TTL)
			rtc := 0
			for !j.Done {
				j.DO()
				if j.Done {
					log.Printf("job %d is done.\n", j.Jid)
					j.Wg.Done()
				} else {
					if j.Mrtc == -1 || rtc < j.Mrtc {
						rtc += 1
						log.Printf("error jobID=%d retrying, rtc=%d\n", j.Jid, rtc)
					} else {
						j.Done = true
						j.Wg.Done()
						log.Printf("** ERROR ** jobID=%d reached max retry number and terminated.\n", j.Jid)
					}
				}
			}
		}
	}
}

func GOJOB(queue chan Job, workerc int) {
	log.SetFlags(2)
	workers := make([]worker, workerc)

	for i := 0; i < workerc; i++ {
		workers[i] = worker{i, queue}
		go workers[i].work()
	}
}
