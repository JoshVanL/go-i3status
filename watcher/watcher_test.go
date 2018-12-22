package watcher

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_Watcher(t *testing.T) {
	errNotExp := errors.New("did not expect event")

	w, err := New()
	must(t, err)

	f1, ch1 := newFile(t, w)
	defer os.RemoveAll(f1.Name())

	select {
	case <-ch1:
		must(t, errNotExp)
	default:
	}

	_, err = f1.Write([]byte{' '})
	must(t, err)
	time.Sleep(time.Second / 2)

	select {
	case <-ch1:
	default:
		must(t, errNotExp)
	}

	f2, ch2 := newFile(t, w)
	defer os.RemoveAll(f2.Name())

	select {
	case <-ch1:
		must(t, errNotExp)
	case <-ch2:
		must(t, errNotExp)
	default:
	}

	_, err = f2.Write([]byte{' '})
	must(t, err)
	time.Sleep(time.Second / 2)

	select {
	case <-ch1:
		must(t, errNotExp)
	case <-ch2:
	default:
		must(t, errNotExp)
	}
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func newFile(t *testing.T, w *Watcher) (*os.File, chan struct{}) {
	f, err := ioutil.TempFile("", "watcher")
	must(t, err)

	ch := make(chan struct{})
	must(t, w.AddFile(f.Name(), ch))

	return f, ch
}
