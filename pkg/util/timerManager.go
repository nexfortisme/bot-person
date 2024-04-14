package util

import (
	"sync"
	"time"
)

type TimerManager struct {
    timers map[string]*time.Timer
    lock   sync.Mutex
}

// Singleton instance variables
var (
    instance *TimerManager
    once     sync.Once
)

// GetInstance returns the singleton instance of TimerManager.
func GetInstance() *TimerManager {
    once.Do(func() {
        instance = &TimerManager{
            timers: make(map[string]*time.Timer),
        }
    })
    return instance
}

// NewTimerManager creates a new TimerManager.
func NewTimerManager() *TimerManager {
    return &TimerManager{
        timers: make(map[string]*time.Timer),
    }
}

// SetTimer sets a timer that calls f after d. It associates the timer with a key.
func (m *TimerManager) SetTimer(key string, d time.Duration, f func()) {
    m.lock.Lock()
    defer m.lock.Unlock()
    if timer, exists := m.timers[key]; exists {
        timer.Stop() // Stop existing timer if any
    }
    m.timers[key] = time.AfterFunc(d, func() {
        f()
        m.ClearTimer(key) // Optionally clear timer after firing
    })
}

// StopTimer stops the timer associated with the key.
func (m *TimerManager) StopTimer(key string) {
    m.lock.Lock()
    defer m.lock.Unlock()
    if timer, exists := m.timers[key]; exists {
        timer.Stop()
        delete(m.timers, key)
    }
}

// ClearTimer clears a timer entry without stopping it, if it has already stopped.
func (m *TimerManager) ClearTimer(key string) {
    m.lock.Lock()
    defer m.lock.Unlock()
    delete(m.timers, key)
}

func(m *TimerManager) CheckTimer(key string) bool {
    m.lock.Lock()
    defer m.lock.Unlock()
    _, exists := m.timers[key]
    return exists
}

func(m *TimerManager) ExecTimerFunction(key string) {
    // Get the function from timers map and execute it
    m.lock.Lock()
    defer m.lock.Unlock()
    if timer, exists := m.timers[key]; exists {
        timer.Reset(0)
    }
}
