package storage

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestCreateReferenceThenRecover(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()

	if created, _ := store.CreateReference("abc", "https://example.com"); !created {
		t.Fatal("CreateReference() created = false, want true")
	}

	got, err := store.RecoverReference("abc")
	if err != nil {
		t.Fatalf("RecoverReference() error = %v, want nil", err)
	}

	if want := "https://example.com"; got != want {
		t.Fatalf("RecoverReference() = %q, want %q", got, want)
	}
}

func TestRecoverReferenceMissingKey(t *testing.T) {
	t.Parallel()

	got, err := NewMemoryStore().RecoverReference("does-not-exist")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("RecoverReference() error = %v, want ErrNotFound", err)
	}

	if got != "" {
		t.Fatalf("RecoverReference() = %q, want empty string", got)
	}
}

func TestCreateReferenceRefusesCollision(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	store.CreateReference("dup", "https://first.example")

	created, current := store.CreateReference("dup", "https://second.example")
	if created {
		t.Fatal("CreateReference() created = true on collision, want false")
	}

	if want := "https://first.example"; current != want {
		t.Fatalf("CreateReference() current = %q, want %q", current, want)
	}

	got, _ := store.RecoverReference("dup")
	if want := "https://first.example"; got != want {
		t.Fatalf("stored value = %q, want the original %q preserved", got, want)
	}
}

func TestCreateReferenceIsIdempotent(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	store.CreateReference("same", "https://same.example")

	if created, _ := store.CreateReference("same", "https://same.example"); !created {
		t.Fatal("CreateReference() created = false on identical rewrite, want true")
	}
}

func TestMemoryStoreIsConcurrencySafe(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			key := fmt.Sprintf("k%d", i)
			store.CreateReference(key, "https://example.com")

			if _, err := store.RecoverReference(key); err != nil {
				t.Errorf("RecoverReference() error = %v, want nil", err)
			}
		}()
	}

	wg.Wait()
}

func TestMemoryStoreSatisfiesStore(t *testing.T) {
	t.Parallel()

	var _ Store = NewMemoryStore()
}
