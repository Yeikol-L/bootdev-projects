package video

import (
	"testing"
)

func TestSimpli(t *testing.T) {
	tests := []struct {
		n, d         int
		wantN, wantD int
	}{
		{1920, 1080, 16, 9},
		{1280, 720, 16, 9},
		{2560, 1080, 64, 27},
		{1024, 768, 4, 3},
		{-4, 8, -1, 2},
		{0, 5, 0, 1},
	}

	for _, tt := range tests {
		n, d := simpli(tt.n, tt.d)
		t.Logf("From %d/%d to %d/%d", tt.n, tt.d, n, d)
		if n != tt.wantN || d != tt.wantD {
			t.Errorf("simpli(%d, %d) = %d/%d, want %d/%d",
				tt.n, tt.d, n, d, tt.wantN, tt.wantD)
		}
	}
}
