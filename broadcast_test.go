package chanx

import (
	"testing"
)

func TestBroadcast(t *testing.T) {
	// Test with integer channels
	t.Run("broadcast integers", func(t *testing.T) {
		src := make(chan int)
		dst1 := make(chan int)
		dst2 := make(chan int)

		err := Broadcast(src, dst1, dst2)
		if err != nil {
			t.Fatalf("Failed to setup broadcast: %v", err)
		}

		// Send values
		go func() {
			src <- 1
			src <- 2
			src <- 3
			close(src)
		}()

		// Receive from first destination
		go func() {
			expected := []int{1, 2, 3}
			for i, want := range expected {
				got := <-dst1
				if got != want {
					t.Errorf("dst1[%d] = %v; want %v", i, got, want)
				}
			}
		}()

		// Receive from second destination
		expected := []int{1, 2, 3}
		for i, want := range expected {
			got := <-dst2
			if got != want {
				t.Errorf("dst2[%d] = %v; want %v", i, got, want)
			}
		}
	})

	// Test with string channels
	t.Run("broadcast strings", func(t *testing.T) {
		src := make(chan string)
		dst1 := make(chan string)
		dst2 := make(chan string)

		err := Broadcast(src, dst1, dst2)
		if err != nil {
			t.Fatalf("Failed to setup broadcast: %v", err)
		}

		go func() {
			src <- "hello"
			src <- "world"
			close(src)
		}()

		// Check both destinations receive the same values
		for _, want := range []string{"hello", "world"} {
			got1 := <-dst1
			got2 := <-dst2
			if got1 != want || got2 != want {
				t.Errorf("got %q, %q; want %q, %q", got1, got2, want, want)
			}
		}
	})

	// Test error cases
	t.Run("error cases", func(t *testing.T) {
		cases := []struct {
			name string
			src  interface{}
			dst  []interface{}
			want string
		}{
			{
				name: "non-channel source",
				src:  42,
				dst:  []interface{}{make(chan int)},
				want: "src must be a channel",
			},
			{
				name: "send-only source",
				src:  make(chan<- int),
				dst:  []interface{}{make(chan int)},
				want: "src channel must be a readable:[<-chan] or [chan]",
			},
			{
				name: "type mismatch",
				src:  make(chan int),
				dst:  []interface{}{make(chan string)},
				want: "dst channel element type must match src channel",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := Broadcast(tc.src, tc.dst...)
				if err == nil || err.Error() != tc.want {
					t.Errorf("Broadcast() error = %v; want %v", err, tc.want)
				}
			})
		}
	})
}
