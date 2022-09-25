package logger

import (
	"os"
	"strings"
	"sync"
	"testing"
)

// TestAtomicWrites ensures that all writes to the log are atomic.
func TestAtomicWrites(t *testing.T) {
	l, err := New(Opts{
		Module:   "Main",
		Filename: "/tmp/logger_test.log",
	})
	if err != nil {
		t.Fatalf("New(_) = %v, need nil error", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Infof("info")
		}()
	}
	wg.Wait()
	if err := l.Close(); err != nil {
		t.Fatalf("Close() = %v, need nil error", err)
	}

	contents, err := os.ReadFile("/tmp/logger_test.log")
	if err != nil {
		t.Fatalf("os.ReadFile(_) = %v, need nil error", err)
	}
	for _, line := range strings.Split(string(contents), "\n") {
		if line == "" {
			continue
		}
		if !strings.HasSuffix(line, "[Main INFO] info") {
			t.Errorf(`line %q: suffix "[Main INFO] info" expected`, line)
		}
	}
	os.Remove("/tmp/logger_test.log")
}

// TestSingleton ensures that a logger file can be only instantiated once.
func TestSingleton(t *testing.T) {
	l1, err := New(Opts{
		Module:   "Main",
		Filename: "/tmp/logger1_test.log",
	})
	if err != nil {
		t.Fatalf("First logger: New(_) = %v, need nil error", err)
	}
	_, err = New(Opts{
		Module:   "Main",
		Filename: "/tmp/logger2_test.log",
	})
	if err == nil {
		t.Fatalf("Second logger: New(_) = nil, need error")
	}
	l1.Close()
	os.Remove("/tmp/logger1_test.log")
}

// TestVerbose checks that debug message are (or not) sent.
func TestVerbose(t *testing.T) {
	var l *logger
	var err error

	setup := func(verbose bool) {
		l, err = New(Opts{
			Module:   "Main",
			Filename: "/tmp/logger_test.log",
			Verbose:  verbose,
			Append:   true,
		})
		if err != nil {
			t.Fatalf("New(_) = %v, need nil error", err)
		}
	}
	teardown := func() bool {
		l.Close()
		contents, err := os.ReadFile("/tmp/logger_test.log")
		if err != nil {
			t.Fatalf("os.ReadFile(_) = %v, need nil error", err)
		}
		found := false
		for _, line := range strings.Split(string(contents), "\n") {
			if strings.Contains(line, "DEBUG") {
				found = true
				break
			}
		}
		if err := os.Remove("/tmp/logger_test.log"); err != nil {
			t.Fatalf("os.Remove(_) = %v, need nil error", err)
		}
		return found
	}

	setup(false)
	for i := 0; i < 5; i++ {
		l.Debugf("test")
	}
	if found := teardown(); found {
		t.Errorf("Verbosity test: found unexpected DEBUG statements")
	}

	setup(true)
	for i := 0; i < 5; i++ {
		l.Debugf("test")
	}
	if found := teardown(); !found {
		t.Errorf("Verbosity test: failed to find DEBUG statements")
	}
	os.Remove("/tmp/logger_test.log")
}
