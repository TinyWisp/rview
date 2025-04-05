package rview

import (
	"testing"
)

func TestBasicReactivity(t *testing.T) {
	a := NewRef[int](1)
	b := NewRef[int](1)
	sum := NewRef[int](2)

	if a.Get() != 1 || b.Get() != 1 || sum.Get() != 2 {
		t.Fatalf("1.Reactivity doesn't work as expected. a=%d, b=%d, sum=%d", a.Get(), b.Get(), sum.Get())
	}

	NewTracker(func() {
		sum.Set(a.Get() + b.Get())
	})

	if sum.Get() != 2 {
		t.Fatalf("2.Reactivity doesn't work as expected. a=%d, b=%d, sum=%d", a.Get(), b.Get(), sum.Get())
	}

	a.Set(2)
	b.Set(3)
	if sum.Get() != 5 {
		t.Fatalf("3.Reactivity doesn't work as expected. a=%d, b=%d, sum=%d", a.Get(), b.Get(), sum.Get())
	}

	a.Set(0)
	b.Set(0)
	for i := 0; i < 100; i++ {
		a.Set(a.Get() + 1)
		b.Set(b.Get() + 1)
	}
	if sum.Get() != 200 {
		t.Fatalf("4.Reactivity doesn't work as expected. a=%d, b=%d, sum=%d", a.Get(), b.Get(), sum.Get())
	}
}

func TestEmbedReactivity(t *testing.T) {
	a := NewRef[int](1)
	b := NewRef[int](1)
	avg := NewRef[int](1)

	sum := func() int {
		return a.Get() + b.Get()
	}

	NewTracker(func() {
		avg.Set(sum() / 2)
	})

	a.Set(2)
	b.Set(4)
	if avg.Get() != 3 {
		t.Fatalf("1.Reactivity doesn't work as expected. a=%d, b=%d, avg=%d", a.Get(), b.Get(), avg.Get())
	}
}

func TestIsRef(t *testing.T) {
	a := NewRef[int](1)

	if !isRef(a) {
		t.Fatalf("isRef doesn't work as expected.")
	}
}
