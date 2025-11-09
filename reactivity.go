package rview

import (
	"reflect"
	"sync"
)

type Watcher struct {
	refs        []Watchable
	mutex       *sync.Mutex
	watchWhat   func()
	doSomething func(Watchable, interface{}, interface{})
}

type ActiveWatcherMgr struct {
	watchers []*Watcher
}

type Ref[T any] struct {
	value    T
	oldValue T
	watchers []*Watcher
}

type Watchable interface {
	AddWatcher(watcher *Watcher)
	RemoveWatcher(watcher *Watcher)
	Trigger()
}

var activeWatcherMgr *ActiveWatcherMgr = &ActiveWatcherMgr{}

func (awm *ActiveWatcherMgr) Push(watcher *Watcher) {
	awm.watchers = append(awm.watchers, watcher)
}

func (awm *ActiveWatcherMgr) Pop() *Watcher {
	count := len(awm.watchers)
	if count == 0 {
		return nil
	}

	watcher := awm.watchers[count-1]
	awm.watchers = awm.watchers[:count-1]

	return watcher
}

func (awm *ActiveWatcherMgr) ActiveWatcher() *Watcher {
	count := len(awm.watchers)
	if count == 0 {
		return nil
	}

	return awm.watchers[count-1]
}

func (awm *ActiveWatcherMgr) Depth() int {
	return len(awm.watchers)
}

func NewWatcher(watchWhat func(), doSomething func(Watchable, interface{}, interface{})) *Watcher {
	watcher := &Watcher{
		refs:        make([]Watchable, 0),
		watchWhat:   watchWhat,
		mutex:       &sync.Mutex{},
		doSomething: doSomething,
	}
	watcher.RunAndWatch()
	return watcher
}

func (t *Watcher) RunAndWatch() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	activeWatcherMgr.Push(t)
	defer activeWatcherMgr.Pop()

	t.clean()
	t.watchWhat()
}

func (t *Watcher) AddRef(ref Watchable) {
	t.refs = append(t.refs, ref)
}

func (t *Watcher) clean() {
	if len(t.refs) == 0 {
		return
	}

	for _, ref := range t.refs {
		ref.RemoveWatcher(t)
	}
	t.refs = t.refs[:0]
}

func NewRef[T any](t T) *Ref[T] {
	return &Ref[T]{
		value:    t,
		oldValue: t,
	}
}

func (r *Ref[T]) Get() T {
	activeWatcher := activeWatcherMgr.ActiveWatcher()
	if activeWatcher != nil {
		activeWatcher.AddRef(r)
		r.AddWatcher(activeWatcher)
	}

	return r.value
}

func (r *Ref[T]) Set(val T) {
	r.oldValue = r.value
	r.value = val

	if !r.isEqual(r.oldValue, r.value) {
		r.Trigger()
	}
}

func (r *Ref[T]) isEqual(a T, b T) bool {
	aval := reflect.ValueOf(a)
	bval := reflect.ValueOf(b)

	if aval.Comparable() {
		return aval.Equal(bval)
	}

	return reflect.DeepEqual(a, b)
}

func (r *Ref[T]) isRef() bool {
	return true
}

func (r *Ref[T]) Type() reflect.Type {
	return reflect.TypeOf(r.value)
}

func (r *Ref[T]) AddWatcher(watcher *Watcher) {
	for _, owatcher := range r.watchers {
		if owatcher == watcher {
			return
		}
	}

	r.watchers = append(r.watchers, watcher)
}

func (r *Ref[T]) RemoveWatcher(watcher *Watcher) {
	for idx, owatcher := range r.watchers {
		if owatcher == watcher {
			r.watchers = append(r.watchers[:idx], r.watchers[idx+1:]...)
		}
	}
}

func (r *Ref[T]) Trigger() {
	for _, watcher := range r.watchers {
		watcher.doSomething(r, r.value, r.oldValue)
	}
}

func isRef(v interface{}) bool {
	_, ok := v.(interface{ isRef() bool })
	return ok
}

func WatchRef[T any](ref *Ref[T], fnOnChanged func(newVal, oldVal T), immediate bool) func() {
	watcher := NewWatcher(func() {
		ref.Get()
	}, func(ref Watchable, newValIntf interface{}, oldValIntf interface{}) {
		fnOnChanged(newValIntf.(T), oldValIntf.(T))
	})

	if immediate {
		ref.Trigger()
	}

	return func() {
		ref.RemoveWatcher(watcher)
	}
}

func WatchRefs(refs []Watchable, fnOnChanged func(), immediate bool) func() {
	watcher := NewWatcher(func() {
		for _, ref := range refs {
			refVal := reflect.ValueOf(ref)
			refVal.MethodByName("Get").Call([]reflect.Value{})
		}
	}, func(ref Watchable, newValIntf interface{}, oldValIntf interface{}) {
		fnOnChanged()
	})

	if immediate {
		fnOnChanged()
	}

	return func() {
		for _, ref := range refs {
			ref.RemoveWatcher(watcher)
		}
	}
}

func RunAndWatch(runWhat func(), fnOnChanged func(watcher *Watcher)) func() {
	stopped := false

	watcher := (*Watcher)(nil)
	watcher = NewWatcher(runWhat, func(ref Watchable, newVal interface{}, oldVal interface{}) {
		if !stopped {
			fnOnChanged(watcher)
		}
	})

	return func() {
		stopped = true
		watcher.clean()
	}
}

func RunReactively(runWhat func()) func() {
	stopped := false

	watcher := (*Watcher)(nil)
	watcher = NewWatcher(runWhat, func(ref Watchable, newVal interface{}, oldVal interface{}) {
		if !stopped {
			watcher.RunAndWatch()
		}
		if stopped {
			watcher.clean()
		}
	})

	return func() {
		stopped = true
		watcher.clean()
	}
}

func Computed[T any](fnCompute func() T) *Ref[T] {
	res := NewRef[T](*new(T))

	fnComputeWrapper := func() {
		res.Set(fnCompute())
	}

	watcher := (*Watcher)(nil)
	watcher = NewWatcher(fnComputeWrapper, func(ref Watchable, newVal interface{}, oldVal interface{}) {
		watcher.RunAndWatch()
	})

	return res
}
