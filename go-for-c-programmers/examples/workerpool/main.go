// Command workerpool runs many jobs through a fixed, bounded set of workers.
// Each job fetches one URL. To stay fully offline and deterministic, the URLs
// point at an in-process test server (net/http/httptest) started by the program
// itself, so there is no real network and the run is repeatable on any machine.
//
// It demonstrates a bounded worker pool, a context timeout that cancels all
// work at once, result aggregation, and clean termination with no leaked
// goroutines. Run it under the race detector to confirm it is data-race free:
//
//	go run -race ./workerpool
package main

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"time"
)

// Job is one unit of work: fetch one URL.
type Job struct {
	ID  int
	URL string
}

// Result is the outcome of one Job. Err is nil on success.
type Result struct {
	JobID  int
	Status int
	Bytes  int
	Err    error
}

// fetch performs one HTTP GET that respects the context. If the context is
// canceled or its deadline passes, the in-flight request is aborted and Do
// returns an error. It reports the status code and how many body bytes arrived.
func fetch(ctx context.Context, client *http.Client, job Job) Result {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		return Result{JobID: job.ID, Err: err}
	}
	resp, err := client.Do(req)
	if err != nil {
		return Result{JobID: job.ID, Err: err}
	}
	defer resp.Body.Close()
	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return Result{JobID: job.ID, Status: resp.StatusCode, Err: err}
	}
	return Result{JobID: job.ID, Status: resp.StatusCode, Bytes: int(n)}
}

// worker pulls jobs until the jobs channel is closed or the context is
// canceled, sends a Result for each, and then returns. Selecting on ctx.Done in
// both the receive and the send guarantees the worker never blocks forever, so
// it cannot leak.
func worker(ctx context.Context, client *http.Client, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return // jobs channel closed: no more work
			}
			res := fetch(ctx, client, job)
			select {
			case results <- res:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	// In-process server: the whole program is self-contained and offline. The
	// handler writes a short, fixed body so every run produces the same totals.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "result for %s\n", r.URL.Path)
	}))
	defer srv.Close()

	const (
		numJobs    = 20
		numWorkers = 4 // the bound: at most this many requests are ever in flight
	)

	jobList := make([]Job, numJobs)
	for i := range numJobs {
		jobList[i] = Job{ID: i, URL: fmt.Sprintf("%s/work/%d", srv.URL, i)}
	}

	// One deadline for the whole batch. If the work ran long, ctx would be
	// canceled and every worker and the producer would stop promptly.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// A shared client is safe for concurrent use by many goroutines.
	client := &http.Client{Timeout: 2 * time.Second}

	jobs := make(chan Job)
	results := make(chan Result)

	// Start a fixed number of workers. This is what makes the pool bounded:
	// adding more jobs does not add more concurrency.
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go worker(ctx, client, jobs, results, &wg)
	}

	// Producer: feed every job, then close the jobs channel. Selecting on
	// ctx.Done means a cancellation cannot leave this goroutine stuck sending.
	go func() {
		defer close(jobs)
		for _, job := range jobList {
			select {
			case jobs <- job:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Closer: when every worker has returned, close results so the range below
	// ends cleanly. This is the standard fan-in finish.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate. Results arrive in whatever order workers finish, so collect
	// them first and sort by job ID for stable, repeatable output.
	var collected []Result
	for res := range results {
		collected = append(collected, res)
	}
	slices.SortFunc(collected, func(a, b Result) int { return cmp.Compare(a.JobID, b.JobID) })

	var ok, failed, totalBytes int
	for _, res := range collected {
		if res.Err != nil {
			failed++
			log.Printf("job %d failed: %v", res.JobID, res.Err)
			continue
		}
		ok++
		totalBytes += res.Bytes
	}
	fmt.Printf("done: %d ok, %d failed, %d jobs total, %d bytes\n", ok, failed, len(collected), totalBytes)
	if err := ctx.Err(); err != nil {
		fmt.Printf("note: context ended early: %v\n", err)
	}
}
