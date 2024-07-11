/*
Copyright 2018 github.com/ucirello

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pglock

import (
	"context"
	"sync"
	"time"
)

// Lock is the mutex entry in the database.
type Lock struct {
	client           *Client
	name             string
	owner            string
	heartbeatWG      sync.WaitGroup
	heartbeatContext context.Context
	heartbeatCancel  context.CancelFunc
	leaseDuration    time.Duration

	replaceData   bool
	data          []byte
	failIfLocked  bool
	keepOnRelease bool

	mu                  sync.Mutex
	isReleased          bool
	recordVersionNumber int64
}

// Data returns the content of the lock, if any is available.
func (l *Lock) Data() []byte {
	return l.data
}

// Close releases the lock and interrupts the locks heartbeat, if configured.
func (l *Lock) Close() error {
	err := l.client.Release(l)
	return err
}

// IsReleased indicates whether the lock is either released or lost after
// heartbeat.
func (l *Lock) IsReleased() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.isReleased
}

// Owner returns who currently owns the lock.
func (l *Lock) Owner() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.owner
}

// RecordVersionNumber is the expectation that this lock entry has about its
// consistency in the database. If the RecordVersionNumber from the database
// mismatches the one in the lock, it means that some clock drift has taken
// place and this lock is no longer valid.
func (l *Lock) RecordVersionNumber() int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.recordVersionNumber
}

// LockOption reconfigures how the lock behaves on acquire and release.
type LockOption func(*Lock)

// WithCustomHeartbeatContext will override the context used for the heartbeats.
// It means the cancelation now is responsibility of the caller of the lock.
func WithCustomHeartbeatContext(ctx context.Context) LockOption {
	return func(l *Lock) {
		l.heartbeatContext, l.heartbeatCancel = context.WithCancel(ctx)
	}
}

// FailIfLocked will not retry to acquire the lock, instead returning.
func FailIfLocked() LockOption {
	return func(l *Lock) {
		l.failIfLocked = true
	}
}

// WithData creates lock with data.
func WithData(data []byte) LockOption {
	return func(l *Lock) {
		l.data = data
	}
}

// KeepOnRelease preserves the lock entry when Close() is called on the lock.
func KeepOnRelease() LockOption {
	return func(l *Lock) {
		l.keepOnRelease = true
	}
}

// ReplaceData will force the new content to be stored in the lock entry.
func ReplaceData() LockOption {
	return func(l *Lock) {
		l.replaceData = true
	}
}

// ReadOnlyLock holds a copy of the information of a lock in the database.
type ReadOnlyLock Lock

// Name returns the lock's name.
func (l *ReadOnlyLock) Name() string {
	return l.name
}

// Data returns the content of the lock, if any is available.
func (l *ReadOnlyLock) Data() []byte {
	return l.data
}

// Owner returns who currently owns the lock.
func (l *ReadOnlyLock) Owner() string {
	return l.owner
}
