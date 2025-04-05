package rview

import (
	"reflect"
	"sync"
)

type Tracker struct {
	refs   []Trackable
	mutex  *sync.Mutex
	effect func()
}

type ActiveTrackerMgr struct {
	trackers []*Tracker
}

type Ref[T any] struct {
	value    T
	trackers []*Tracker
}

type Trackable interface {
	AddTracker(tracker *Tracker)
	RemoveTracker(tracker *Tracker)
	Trigger()
}

var activeTrackerMgr *ActiveTrackerMgr = &ActiveTrackerMgr{}

func (atm *ActiveTrackerMgr) Push(effect *Tracker) {
	atm.trackers = append(atm.trackers, effect)
}

func (atm *ActiveTrackerMgr) Pop() *Tracker {
	count := len(atm.trackers)
	if count == 0 {
		return nil
	}

	tracker := atm.trackers[count-1]
	atm.trackers = atm.trackers[:count-1]

	return tracker
}

func (atm *ActiveTrackerMgr) ActiveTracker() *Tracker {
	count := len(atm.trackers)
	if count == 0 {
		return nil
	}

	return atm.trackers[count-1]
}

func (atm *ActiveTrackerMgr) Depth() int {
	return len(atm.trackers)
}

func NewTracker(effect func()) *Tracker {
	re := &Tracker{
		refs:   make([]Trackable, 0),
		effect: effect,
		mutex:  &sync.Mutex{},
	}
	re.RunAndWatch()
	return re
}

func (t *Tracker) RunAndWatch() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	activeTrackerMgr.Push(t)
	defer activeTrackerMgr.Pop()

	t.clean()
	t.effect()
}

func (t *Tracker) AddRef(ref Trackable) {
	t.refs = append(t.refs, ref)
}

func (t *Tracker) clean() {
	if len(t.refs) == 0 {
		return
	}

	for _, ref := range t.refs {
		ref.RemoveTracker(t)
	}
	t.refs = t.refs[:0]
}

func NewRef[T any](t T) *Ref[T] {
	return &Ref[T]{
		value: t,
	}
}

func (r *Ref[T]) Get() T {
	activeTracker := activeTrackerMgr.ActiveTracker()
	if activeTracker != nil {
		activeTracker.AddRef(r)
		r.AddTracker(activeTracker)
	}

	return r.value
}

func (r *Ref[T]) Set(val T) {
	r.value = val
	r.Trigger()
}

func (r *Ref[T]) isRef() bool {
	return true
}

func (r *Ref[T]) Type() reflect.Type {
	return reflect.TypeOf(r.value)
}

func (r *Ref[T]) AddTracker(tracker *Tracker) {
	for _, otracker := range r.trackers {
		if otracker == tracker {
			return
		}
	}

	r.trackers = append(r.trackers, tracker)
}

func (r *Ref[T]) RemoveTracker(tracker *Tracker) {
	for idx, otracker := range r.trackers {
		if otracker == tracker {
			r.trackers = append(r.trackers[:idx], r.trackers[idx+1:]...)
		}
	}
}

func (r *Ref[T]) Trigger() {
	for _, tracker := range r.trackers {
		tracker.RunAndWatch()
	}
}

func isRef(v interface{}) bool {
	_, ok := v.(interface{ isRef() bool })
	return ok
}
