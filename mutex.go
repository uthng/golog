package golog

import (
	"sync"
)

// Mutex is struct of mutex and locked to know
// if mutex is locked or not
type Mutex struct {
	mutex  sync.Mutex
	locked bool
}

// LockOnce locks mutex if it isnt yet
func (m *Mutex) LockOnce() {
	if !m.locked {
		m.mutex.Lock()
		m.locked = true
	}
}

// UnlockOnce unlocks mutex if it is
func (m *Mutex) UnlockOnce() {
	if m.locked {
		m.mutex.Unlock()
		m.locked = false
	}
}

// IsLocked returns if mutex is locked
func (m *Mutex) IsLocked() bool {
	return m.locked
}
