package healthz

import (
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestNewCheckEmptyLive(t *testing.T) {
	t.Parallel()

	h := NewCheck("", "", "")
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res1 := w.Result()
	if res1.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res1.Status)
	}

	r = httptest.NewRequest(http.MethodGet, "/ready", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.ready(w, r)
	res2 := w.Result()
	if res2.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res2.Status)
	}
}

func TestNewCheckValues(t *testing.T) {
	t.Parallel()

	h := NewCheck("livez", "readyz", "8081")
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.ready(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}
}

func TestNewCheckPrefixes(t *testing.T) {
	t.Parallel()

	h := NewCheck("/livez", "/readyz", ":8082")
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}
}

func TestLiveness(t *testing.T) {
	t.Parallel()

	h := NewCheck("livez", "", "8086")

	r := httptest.NewRequest(http.MethodGet, "/livez", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodPost, "/livez", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("handler returned unexpected status code: got %v want 405",
			res.Status)
	}
}

func TestReadiness(t *testing.T) {
	t.Parallel()

	h := NewCheck("", "readyz", "8087")

	// test ready
	h.Ready()

	r := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v want 200",
			res.Status)
	}

	r = httptest.NewRequest(http.MethodPost, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("handler returned unexpected status code: got %v want 405",
			res.Status)
	}

	// test notready
	h.NotReady()

	r = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w = httptest.NewRecorder()
	h.router().ServeHTTP(w, r)
	h.live(w, r)
	res = w.Result()
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("handler returned unexpected status code: got %v want 503",
			res.Status)
	}
}

func TestTerminating(t *testing.T) {
	h := NewCheck("", "", "")
	var term bool
	go func() {
		term = h.Terminating()
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	proc.Signal(syscall.SIGINT)

	if term != <-done {
		t.Errorf("termination return: got %v want true",
			term)
	}
}
