// Copyright 2023 - MinIO, Inc. All rights reserved.
// Use of this source code is governed by the AGPLv3
// license that can be found in the LICENSE file.

package staticlist

import (
	"context"
	"fmt"
	"github.com/minio/kes-go"
	"github.com/minio/kes/kv"
)

type Config struct {
	Entries map[string]string
}

func (c *Config) Validate() error {
	// TODO: validate keys
	return nil
}

type Store struct {
	config *Config
}

func NewStore(_ context.Context, config *Config) (*Store, error) {
	return &Store{config: config}, nil
}

// Status returns the current state of the Store or an error explaining why fetching status information failed.
func (s *Store) Status(_ context.Context) (kv.State, error) {
	return kv.State{Latency: 0}, nil
}

// Create creates a new entry at the storage if and only if no entry for the give key exists.
// If such an entry already exists, Create returns ErrExists.
func (s *Store) Create(_ context.Context, name string, _ []byte) error {
	_, exists := s.config.Entries[name]
	if exists {
		return fmt.Errorf("key '%s' already exists: %w", name, kes.ErrKeyExists)
	}
	return fmt.Errorf("key '%s' doesn't exist, create operation is not allowed on static list: %w", name, kes.ErrNotAllowed)
}

// Set writes the key-value pair to the storage.
func (s *Store) Set(_ context.Context, name string, _ []byte) error {
	return fmt.Errorf("set key '%s' value operation is not allowed on static list: %w", name, kes.ErrNotAllowed)
}

// Get returns the value associated with the given key.
// It returns ErrNotExists if no such entry exists.
func (s *Store) Get(_ context.Context, name string) ([]byte, error) {
	value, exists := s.config.Entries[name]
	if !exists {
		return nil, kes.ErrKeyNotFound
	}
	return []byte(value), nil
}

// Delete deletes the key and the associated value from the storage.
// It returns ErrNotExists if no such entry exists.
func (s *Store) Delete(_ context.Context, name string) error {
	_, exists := s.config.Entries[name]
	if !exists {
		return fmt.Errorf("key '%s' doesn't exist: %w", name, kes.ErrKeyNotFound)
	}
	return fmt.Errorf("key '%s' exists, but delete operation is not allowed on static list: %w", name, kes.ErrNotAllowed)
}

// List returns an Iter enumerating the stored entries.
func (s *Store) List(ctx context.Context) (kv.Iter[string], error) {
	keys := make([]string, len(s.config.Entries))
	i := 0
	for k := range s.config.Entries {
		keys[i] = k
		i = i + 1
	}
	return &keysIter{keys: keys, index: 0, ctx: ctx}, nil
}

type keysIter struct {
	keys  []string
	index int
	ctx   context.Context
}

func (s *keysIter) Next() (string, bool) {
	key := ""
	if s.index < len(s.keys) {
		key = s.keys[s.index]
		s.index++
	}
	return key, s.index < len(s.keys)
}

func (s *keysIter) Close() error {
	s.keys = nil
	return s.ctx.Err()
}

// Close  terminate or release resources that were opened or acquired.
func (s *Store) Close() error { return nil }

// This line ensures that the Store type implements the kv.Store interface.
// It will fail at compile time if the contract is not satisfied.
var _ kv.Store[string, []byte] = (*Store)(nil)
