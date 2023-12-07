package staticlist

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/minio/kes-go"
	"testing"
)

func NewFakeStore() *Store {
	return &Store{
		config: &Config{
			Entries: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
}

func TestStore_Create(t *testing.T) {
	store := NewFakeStore()

	t.Run("create element that exists", func(t *testing.T) {
		err := store.Create(context.Background(), "key1", []byte("x"))
		assertErrorIs(t, err, kes.ErrKeyExists)
	})
	t.Run("create element that doesn't exist is not allowed", func(t *testing.T) {
		err := store.Create(context.Background(), "new-key", []byte("x"))
		assertErrorIs(t, err, kes.ErrNotAllowed)
	})
}

func TestStore_Set(t *testing.T) {
	store := NewFakeStore()

	t.Run("set value of element is not allowed", func(t *testing.T) {
		err := store.Set(context.Background(), "some-key", []byte("x"))
		assertErrorIs(t, err, kes.ErrNotAllowed)
	})
}

func TestStore_Get(t *testing.T) {
	store := NewFakeStore()

	t.Run("GET string value that exists", func(t *testing.T) {
		b, err := store.Get(context.Background(), "key1")
		assertNoError(t, err)
		assertEqualBytes(t, []byte("value1"), b)
	})
	t.Run("GET element that doesn't exist", func(t *testing.T) {
		_, err := store.Get(context.Background(), "new-key")
		assertErrorIs(t, err, kes.ErrKeyNotFound)
	})

}

func TestStore_Delete(t *testing.T) {
	store := NewFakeStore()

	t.Run("DELETE element that doesn't exist", func(t *testing.T) {
		err := store.Delete(context.Background(), "new-key")
		assertErrorIs(t, err, kes.ErrKeyNotFound)
	})
	t.Run("DELETE element that exists is not allowed", func(t *testing.T) {
		err := store.Delete(context.Background(), "key1")
		assertErrorIs(t, err, kes.ErrNotAllowed)
	})
}

func TestStore_List(t *testing.T) {
	store := NewFakeStore()
	t.Run("returns list of two elements", func(t *testing.T) {
		iter, err := store.List(context.Background())
		assertNoError(t, err)
		k, more := iter.Next()
		assertEqualComparable(t, k, "key1")
		assertEqualComparable(t, more, true)
		k, more = iter.Next()
		assertEqualComparable(t, k, "key2")
		assertEqualComparable(t, more, false)
		k, more = iter.Next()
		assertEqualComparable(t, k, "")
		assertEqualComparable(t, more, false)
	})
}

//=== tools:

func assertErrorIs(t *testing.T, err, target error) {
	if err == nil || target == nil {
		t.Fatal("error can't be null")
	}
	if !errors.Is(err, target) {
		t.Fatal(fmt.Sprintf("error '%v' isn't '%v'", err, target))
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func assertEqualComparable(t *testing.T, expected, got any) {
	if expected != got {
		t.Fatalf("expected '%v' got '%v'", expected, got)
	}
}
func assertEqualBytes(t *testing.T, expected, got []byte) {
	if !bytes.Equal(expected, got) {
		t.Fatalf("expected '%v' got '%v'", expected, got)
	}
}
