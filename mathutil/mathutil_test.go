package mathutil

import (
	"reflect"
	"testing"
)

func TestGCD(t *testing.T) {
	for _, tt := range []struct {
		a, b int
		want int
	}{
		{10, 9, 1},
		{12, 9, 3},
		{23, 13, 1},
		{16, 16, 16},
	} {
		if got := GCD(tt.a, tt.b); got != tt.want {
			t.Errorf("GCD(%d, %d): got %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestPhi(t *testing.T) {
	for _, tt := range []struct {
		n    int
		want int
	}{
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 2},
		{5, 4},
		{6, 2},
		{7, 6},
		{8, 4},
		{9, 6},
		{10, 4},
		{11, 10},
		{12, 4},
		{13, 12},
		{14, 6},
		{15, 8},
		{16, 8},
		{17, 16},
		{18, 6},
		{19, 18},
		{20, 8},
	} {
		if got := Phi(tt.n); got != tt.want {
			t.Errorf("Phi(%d): got %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestMu(t *testing.T) {
	for _, tt := range []struct {
		n    int
		want int
	}{
		{1, 1},
		{2, -1},
		{3, -1},
		{4, 0},
		{5, -1},
		{6, 1},
		{7, -1},
		{8, 0},
		{9, 0},
		{10, 1},
		{11, -1},
		{12, 0},
		{13, -1},
		{14, 1},
		{15, 1},
		{16, 0},
		{17, -1},
		{18, 0},
		{19, -1},
		{20, 0},
	} {
		if got := Mu(tt.n); got != tt.want {
			t.Errorf("Mu(%d): got %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestFactor(t *testing.T) {
	for _, tt := range []struct {
		n    int
		want []int
	}{
		{2, []int{2}},
		{3, []int{3}},
		{4, []int{2, 2}},
		{5, []int{5}},
		{12, []int{2, 2, 3}},
		{49, []int{7, 7}},
		{60, []int{2, 2, 3, 5}},
		{89, []int{89}},
	} {
		if got := Factor(tt.n); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Factor(%d): got %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestPrimes(t *testing.T) {
	got := Primes(100)
	want := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Primes(100): got %v, want %v", got, want)
	}
}
