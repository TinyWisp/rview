package rview

import (
	"testing"
)

func TestWatchRef(t *testing.T) {
	// int,  immediate=false
	func() {
		a := NewRef[int](1)
		newValCopy := 1
		oldValCopy := 1
		changedCounter := 0
		prefix := "WatchRef[int](..., ..., false)"

		unwatch := WatchRef[int](a, func(newVal int, oldVal int) {
			newValCopy = newVal
			oldValCopy = oldVal
			changedCounter += 1
		}, false)

		a.Set(1)
		if changedCounter != 0 {
			t.Fatalf("%s #1", prefix)
		}

		a.Set(2)
		if changedCounter != 1 {
			t.Fatalf("%s #2", prefix)
		}
		if newValCopy != 2 || oldValCopy != 1 {
			t.Fatalf("%s #2.1", prefix)
		}

		a.Set(3)
		if changedCounter != 2 {
			t.Fatalf("%s #3", prefix)
		}
		if newValCopy != 3 || oldValCopy != 2 {
			t.Fatalf("%s #3.1", prefix)
		}

		unwatch()
		a.Set(4)
		if changedCounter != 2 {
			t.Fatalf("%s #4", prefix)
		}
		if newValCopy != 3 || oldValCopy != 2 {
			t.Fatalf("%s #4.1", prefix)
		}
	}()

	// int,  immediate=true
	func() {
		a := NewRef[int](1)
		newValCopy := 1
		oldValCopy := 1
		changedCounter := 0
		prefix := "WatchRef[int](..., ..., true)"

		unwatch := WatchRef[int](a, func(newVal int, oldVal int) {
			newValCopy = newVal
			oldValCopy = oldVal
			changedCounter += 1
		}, true)

		a.Set(1)
		if changedCounter != 1 {
			t.Fatalf("%s #1", prefix)
		}
		if newValCopy != 1 || oldValCopy != 1 {
			t.Fatalf("%s #1.1", prefix)
		}

		a.Set(2)
		if changedCounter != 2 {
			t.Fatalf("%s #2", prefix)
		}
		if newValCopy != 2 || oldValCopy != 1 {
			t.Fatalf("%s #2.1", prefix)
		}

		a.Set(3)
		if changedCounter != 3 {
			t.Fatalf("%s #3", prefix)
		}
		if newValCopy != 3 || oldValCopy != 2 {
			t.Fatalf("%s #3.1", prefix)
		}

		unwatch()
		a.Set(4)
		if changedCounter != 3 {
			t.Fatalf("%s #4", prefix)
		}
		if newValCopy != 3 || oldValCopy != 2 {
			t.Fatalf("%s #4.1", prefix)
		}
	}()
}

func TestWatchRefs(t *testing.T) {
	func() {
		a := NewRef[int](1)
		b := NewRef[string]("hello")
		c := NewRef[bool](false)
		changedCounter := 0
		prefix := "WatchRefs([int, string, bool], ...)"

		unwatch := WatchRefs([]Watchable{a, b, c}, func() {
			changedCounter += 1
		}, false)

		if changedCounter != 0 {
			t.Fatalf("%s #1", prefix)
		}

		a.Set(1)
		if changedCounter != 0 {
			t.Fatalf("%s #2", prefix)
		}
		a.Set(2)
		if changedCounter != 1 {
			t.Fatalf("%s #2.1", prefix)
		}

		b.Set("hello")
		if changedCounter != 1 {
			t.Fatalf("%s #3", prefix)
		}
		b.Set("world")
		if changedCounter != 2 {
			t.Fatalf("%s #3.1", prefix)
		}

		c.Set(false)
		if changedCounter != 2 {
			t.Fatalf("%s #4", prefix)
		}
		c.Set(true)
		if changedCounter != 3 {
			t.Fatalf("%s #4.1", prefix)
		}

		unwatch()

		a.Set(3)
		b.Set("smart")
		c.Set(false)
		if changedCounter != 3 {
			t.Fatalf("%s #5", prefix)
		}
	}()
}

func TestComputed(t *testing.T) {
	func() {
		a := NewRef[int](1)
		b := NewRef[int](3)
		sum := Computed[int](func() int {
			return a.Get() + b.Get()
		})
		prefix := "Computed[int](...)"

		if sum.Get() != 4 {
			t.Fatalf("%s #1", prefix)
		}

		a.Set(2)
		if sum.Get() != 5 {
			t.Fatalf("%s #2", prefix)
		}

		b.Set(4)
		if sum.Get() != 6 {
			t.Fatalf("%s #3", prefix)
		}
	}()

	// with conditional branch statements inside
	func() {
		a := NewRef[int](1)
		b := NewRef[int](2)
		c := NewRef[int](4)
		d := NewRef[int](6)
		e := NewRef[int](8)
		f := NewRef[int](9)
		g := NewRef[int](10)
		h := NewRef[int](5)
		calcType := NewRef[string]("sum")
		prefix := "Computed[int](...) switch"

		res := Computed[int](func() int {
			switch calcType.Get() {
			case "sum":
				return a.Get() + b.Get()

			case "avg":
				return (c.Get() + d.Get()) / 2

			case "multiply":
				return e.Get() * f.Get()

			case "divide":
				return g.Get() / h.Get()
			}

			return 0
		})

		if res.Get() != 3 {
			t.Fatalf("%s #1", prefix)
		}

		calcType.Set("avg")
		if res.Get() != 5 {
			t.Fatalf("%s #2", prefix)
		}

		calcType.Set("multiply")
		if res.Get() != 72 {
			t.Fatalf("%s #3", prefix)
		}

		calcType.Set("divide")
		if res.Get() != 2 {
			t.Fatalf("%s #4", prefix)
		}

		a.Set(8)
		b.Set(16)
		calcType.Set("sum")
		if res.Get() != 24 {
			t.Fatalf("%s #5", prefix)
		}
	}()

	// watch computed result
	func() {
		a := NewRef[int](1)
		b := NewRef[int](2)
		sum := Computed(func() int {
			return a.Get() + b.Get()
		})

		newValCopy := 3
		oldValCopy := 3
		changedCounter := 0
		prefix := "computed[int](...) watch"
		WatchRef[int](sum, func(newVal int, oldVal int) {
			newValCopy = newVal
			oldValCopy = oldVal
			changedCounter += 1
		}, false)

		a.Set(2)
		if changedCounter != 1 {
			t.Fatalf("%s #1", prefix)
		}
		if newValCopy != 4 || oldValCopy != 3 {
			t.Fatalf("%s #1.1", prefix)
		}

		b.Set(3)
		if changedCounter != 2 {
			t.Fatalf("%s #2", prefix)
		}
		if newValCopy != 5 || oldValCopy != 4 {
			t.Fatalf("%s #2.1", prefix)
		}
	}()
}

func TestRunReactively(t *testing.T) {
	func() {
		a := NewRef[int](1)
		b := NewRef[int](2)
		c := NewRef[int](4)
		d := NewRef[int](6)
		e := NewRef[int](8)
		f := NewRef[int](9)
		g := NewRef[int](10)
		h := NewRef[int](5)
		res := NewRef[int](0)
		calcType := NewRef[string]("sum")
		prefix := "RunReactively(...)"

		unwatch := RunReactively(func() {
			switch calcType.Get() {
			case "sum":
				res.Set(a.Get() + b.Get())

			case "avg":
				res.Set((c.Get() + d.Get()) / 2)

			case "multiply":
				res.Set(e.Get() * f.Get())

			case "divide":
				res.Set(g.Get() / h.Get())
			}
		})

		if res.Get() != 3 {
			t.Fatalf("%s #1", prefix)
		}

		calcType.Set("avg")
		if res.Get() != 5 {
			t.Fatalf("%s #2", prefix)
		}

		calcType.Set("multiply")
		if res.Get() != 72 {
			t.Fatalf("%s #3", prefix)
		}

		calcType.Set("divide")
		if res.Get() != 2 {
			t.Fatalf("%s #4", prefix)
		}

		a.Set(8)
		b.Set(16)
		calcType.Set("sum")
		if res.Get() != 24 {
			t.Fatalf("%s #5", prefix)
		}

		unwatch()

		a.Set(10)
		b.Set(10)
		if res.Get() != 24 {
			t.Fatalf("%s #6", prefix)
		}
	}()
}

func TestRunAndWatch(t *testing.T) {
	func() {
		a := NewRef[int](1)
		b := NewRef[int](2)
		sum := 0
		runCounter := 0
		changedCounter := 0
		prefix := "RunAndWatch(..., ...) without calling watcher.RunAndWatch()"

		runWhat := func() {
			sum = a.Get() + b.Get()
			runCounter += 1
		}

		unwatch := RunAndWatch(runWhat, func(watcher *Watcher) {
			changedCounter += 1
		})

		if sum != 3 {
			t.Fatalf("%s #1", prefix)
		}
		if changedCounter != 0 {
			t.Fatalf("%s #1.1", prefix)
		}
		if runCounter != 1 {
			t.Fatalf("%s #1.2", prefix)
		}

		a.Set(2)
		if sum != 3 {
			t.Fatalf("%s #2", prefix)
		}
		if changedCounter != 1 {
			t.Fatalf("%s #2.1", prefix)
		}
		if runCounter != 1 {
			t.Fatalf("%s #2.2", prefix)
		}

		b.Set(3)
		if sum != 3 {
			t.Fatalf("%s #3", prefix)
		}
		if changedCounter != 2 {
			t.Fatalf("%s #3.1", prefix)
		}
		if runCounter != 1 {
			t.Fatalf("%s #3.2", prefix)
		}

		unwatch()

		a.Set(3)
		if sum != 3 {
			t.Fatalf("%s #4", prefix)
		}
		if changedCounter != 2 {
			t.Fatalf("%s #4.1", prefix)
		}
		if runCounter != 1 {
			t.Fatalf("%s #4.2", prefix)
		}
	}()

	func() {
		a := NewRef[int](1)
		b := NewRef[int](2)
		sum := 0
		runCounter := 0
		changedCounter := 0
		prefix := "RunAndWatch(..., ...) with calling watcher.RunAndWatch()"

		runWhat := func() {
			sum = a.Get() + b.Get()
			runCounter += 1
		}

		unwatch := RunAndWatch(runWhat, func(watcher *Watcher) {
			changedCounter += 1
			watcher.RunAndWatch()
		})

		if sum != 3 {
			t.Fatalf("%s #1", prefix)
		}
		if changedCounter != 0 {
			t.Fatalf("%s #1.1", prefix)
		}
		if runCounter != 1 {
			t.Fatalf("%s #1.2", prefix)
		}

		a.Set(2)
		if sum != 4 {
			t.Fatalf("%s #2", prefix)
		}
		if changedCounter != 1 {
			t.Fatalf("%s #2.1", prefix)
		}
		if runCounter != 2 {
			t.Fatalf("%s #2.2", prefix)
		}

		b.Set(3)
		if sum != 5 {
			t.Fatalf("%s #3", prefix)
		}
		if changedCounter != 2 {
			t.Fatalf("%s #3.1", prefix)
		}
		if runCounter != 3 {
			t.Fatalf("%s #3.2", prefix)
		}

		unwatch()

		a.Set(3)
		if sum != 5 {
			t.Fatalf("%s #4", prefix)
		}
		if changedCounter != 2 {
			t.Fatalf("%s #4.1", prefix)
		}
		if runCounter != 3 {
			t.Fatalf("%s #4.2", prefix)
		}
	}()

}

func TestIsRef(t *testing.T) {
	a := NewRef[int](1)

	if !isRef(a) {
		t.Fatalf("1. isRef doesn't work as expected.")
	}

	b := 1
	if isRef(b) {
		t.Fatalf("2. isRef doesn't work as expected.")
	}

	c := NewRef[string]("hello")
	if !isRef(c) {
		t.Fatalf("3. isRef doesn't work as expected.")
	}

	d := "hello"
	if isRef(d) {
		t.Fatalf("3. isRef doesn't work as expected.")
	}
}
