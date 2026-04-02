package output

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Spinner displays an animated progress indicator in the terminal.
type Spinner struct {
	message string
	stop    chan struct{}
	done    sync.WaitGroup
}

// NewSpinner creates a spinner with the given message.
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan struct{}),
	}
}

// Start begins the spinner animation in a background goroutine.
func (s *Spinner) Start() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	s.done.Add(1)
	go func() {
		defer s.done.Done()
		i := 0
		for {
			select {
			case <-s.stop:
				// Clear the spinner line
				fmt.Fprintf(os.Stderr, "\r\033[K")
				return
			default:
				fmt.Fprintf(os.Stderr, "\r%s %s", frames[i%len(frames)], s.message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop halts the spinner and waits for cleanup.
func (s *Spinner) Stop() {
	close(s.stop)
	s.done.Wait()
}
