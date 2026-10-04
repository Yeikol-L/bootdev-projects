package pokecache

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
		{
			key: "",
			val: []byte("empty key"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Fatalf("expected to find key %q", c.key)
			}
			if diff := cmp.Diff(c.val, val); diff != "" {
				t.Errorf("Get(%q) mismatch (-want +got):\n%s", c.key, diff)
			}
		})
	}
}

func TestGetMissingKey(t *testing.T) {
	cases := []struct {
		name string
		key  string
	}{
		{name: "empty cache", key: "https://example.com"},
		{name: "second miss does not deadlock", key: "https://example.com/other"},
	}

	cache := NewCache(5 * time.Second)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			done := make(chan struct{})
			go func() {
				val, ok := cache.Get(c.key)
				if ok {
					t.Errorf("expected key %q to be missing", c.key)
				}
				if diff := cmp.Diff([]byte(nil), val); diff != "" {
					t.Errorf("Get(%q) mismatch (-want +got):\n%s", c.key, diff)
				}
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("Get blocked: mutex probably not released")
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const interval = 10 * time.Millisecond
	cases := []struct {
		name      string
		wait      time.Duration
		wantFound bool
	}{
		{
			name:      "entry still present before interval",
			wait:      0,
			wantFound: true,
		},
		{
			name:      "entry reaped after interval",
			wait:      interval * 5,
			wantFound: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cache := NewCache(interval)
			const key = "https://example.com"
			cache.Add(key, []byte("testdata"))

			time.Sleep(c.wait)

			_, ok := cache.Get(key)
			if diff := cmp.Diff(c.wantFound, ok); diff != "" {
				t.Errorf("found mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
