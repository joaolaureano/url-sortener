package hash

import "testing"

func TestEncodeIsDeterministic(t *testing.T) {
	t.Parallel()

	got := Encode("hello")
	if again := Encode("hello"); got != again {
		t.Fatalf("Encode() = %q then %q, want a stable digest", got, again)
	}

	if want := Size * 2; len(got) != want {
		t.Fatalf("len(Encode()) = %d, want %d hex chars", len(got), want)
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	hash := Encode("hello")

	if !Compare("hello", hash) {
		t.Fatal("Compare() = false for matching content, want true")
	}

	if Compare("world", hash) {
		t.Fatal("Compare() = true for different content, want false")
	}
}

func TestDigestLength(t *testing.T) {
	t.Parallel()

	if got := len(Digest("hello")); got != Size {
		t.Fatalf("len(Digest()) = %d, want %d", got, Size)
	}
}

func TestMD5HasherSatisfiesOperator(t *testing.T) {
	t.Parallel()

	var _ Operator = MD5Hasher{}
}
