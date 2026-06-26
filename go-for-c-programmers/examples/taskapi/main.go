// Command taskapi is a small in-memory REST service for "tasks". It shows the
// Go 1.22 net/http router (method and path patterns like "GET /tasks/{id}"),
// JSON encoding and decoding, a logging middleware, a mutex-protected store
// (handlers run concurrently, one goroutine per request), server timeouts, and
// graceful shutdown driven by an OS signal.
//
// Try it:
//
//	go run ./taskapi
//	curl -s localhost:8080/tasks
//	curl -s -XPOST localhost:8080/tasks -d '{"title":"buy milk"}'
//	curl -s localhost:8080/tasks/1
//	curl -s -XDELETE localhost:8080/tasks/1 -i
package main

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Task is one to-do item. The struct tags control the JSON field names; only
// exported (capitalized) fields are encoded, so every field we want in the
// output must start with an uppercase letter.
type Task struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// store is an in-memory set of tasks guarded by a mutex. Every HTTP handler
// runs in its own goroutine, so two requests can hit the store at the same
// time. The mutex makes each operation safe: only one goroutine touches the
// map at once.
type store struct {
	mu     sync.Mutex
	nextID int64
	tasks  map[int64]Task
}

func newStore() *store {
	return &store{tasks: make(map[int64]Task), nextID: 1}
}

// list returns all tasks sorted by ID. Map iteration order is random in Go, so
// we sort to give callers a stable, predictable order.
func (s *store) list() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	slices.SortFunc(out, func(a, b Task) int { return cmp.Compare(a.ID, b.ID) })
	return out
}

func (s *store) create(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Task{ID: s.nextID, Title: title}
	s.tasks[t.ID] = t
	s.nextID++
	return t
}

func (s *store) get(id int64) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *store) delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return false
	}
	delete(s.tasks, id)
	return true
}

// app bundles the dependencies the handlers need. Methods on *app become our
// HTTP handlers, which is a clean way to give them access to the store without
// global variables.
type app struct {
	store *store
}

// writeJSON sets the content type, writes the status line, and encodes v. Once
// WriteHeader has run the status is already sent, so an encode failure can only
// be logged, not turned into a different HTTP error.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json: %v", err)
	}
}

// writeError sends a JSON object like {"error":"..."} with the given status.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// parseID reads the {id} path segment and converts it to an int64.
func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func (a *app) listTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.list())
}

func (a *app) createTask(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
	}
	// MaxBytesReader caps the body so a huge request cannot exhaust memory.
	// DisallowUnknownFields rejects bodies with stray fields.
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	t := a.store.create(in.Title)
	w.Header().Set("Location", fmt.Sprintf("/tasks/%d", t.ID))
	writeJSON(w, http.StatusCreated, t)
}

func (a *app) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return
	}
	t, ok := a.store.get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "no task with that id")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (a *app) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return
	}
	if !a.store.delete(id) {
		writeError(w, http.StatusNotFound, "no task with that id")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// statusWriter wraps http.ResponseWriter so the logging middleware can see the
// status code that a handler chose. ResponseWriter has no getter for it, so we
// record it as it passes through WriteHeader.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// logging is middleware: it takes a handler and returns a new handler that logs
// one line per request after the inner handler runs.
func logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(sw, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, sw.status, time.Since(start))
	})
}

func main() {
	app := &app{store: newStore()}

	// The patterns include the HTTP method and a named wildcard {id}. The
	// router (Go 1.22+) matches on both, so we do not parse the method by hand.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", app.listTasks)
	mux.HandleFunc("POST /tasks", app.createTask)
	mux.HandleFunc("GET /tasks/{id}", app.getTask)
	mux.HandleFunc("DELETE /tasks/{id}", app.deleteTask)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: logging(mux),
		// Timeouts protect the server from slow or stuck clients. Without them
		// a single client could hold a connection open forever.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// NotifyContext returns a context that is canceled when the process gets
	// SIGINT (Ctrl-C) or SIGTERM. That cancellation is our shutdown trigger.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run the server in its own goroutine so main can wait for the signal. The
	// channel is buffered so this goroutine never blocks even if nobody reads.
	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		// ListenAndServe failed to even start (for example, port in use).
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	// Give in-flight requests up to 10 seconds to finish, then stop. Shutdown
	// makes ListenAndServe above return http.ErrServerClosed, the clean case.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		return
	}
	log.Println("server stopped cleanly")
}
