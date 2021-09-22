package hello

import "testing"

func TestHello(t *testing.T) {
	want := "안녕, 세상."
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
