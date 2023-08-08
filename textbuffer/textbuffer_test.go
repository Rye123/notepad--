package textbuffer

import (
	"fmt"
	"testing"
)

func TestTextBufferInsertion(t *testing.T) {
	buf := newBuffer()

	// Test length
	if buf.Length() != 0 {
		t.Fatalf("Test failed, expected length 0, buf.Length(): " + fmt.Sprintf("%d", buf.Length()))
	}

	// Test error with insertion on empty buffer
	err := buf.Insert(1, 'X')
	if err == nil {
		t.Fatalf("Expected error, error was nil. buf.String(): " + buf.String())
	}

	// Test insertion with empty buffer
	buf.Insert(0, 'H')
	if buf.String() != "H" {
		t.Fatalf("Expected \"H\", instead buf.String(): " + buf.String())
	}

	// Test insertion at the start
	buf.Insert(0, 'A')
	if buf.String() != "AH" {
		t.Fatalf("Expected \"AH\", instead buf.String(): " + buf.String())
	}

	// Test insertion in the middle
	buf.Insert(1, 'E')
	if buf.String() != "AEH" {
		t.Fatalf("Expected \"AEH\", instead buf.String(): " + buf.String())
	}

	// Test error with insertion on non-existent index
	err = buf.Insert(-1, 'X')
	if err == nil {
		t.Fatalf("Expected error, error was nil. buf.String(): " + buf.String())
	}
	err = buf.Insert(4, 'X')
	if err == nil {
		t.Fatalf("Expected error, error was nil. buf.String(): " + buf.String())
	}

	// Test insertion at the end
	buf.Insert(3, 'I')
	if buf.String() != "AEHI" {
		t.Fatalf("Expected \"AEHI\", instead buf.String(): " + buf.String())
	}

	// Test length
	if buf.Length() != 4 {
		t.Fatalf("Expected length 4, instead buf.Length(): " + fmt.Sprintf("%d", buf.Length()))
	}
}

func TestTextBufferDelete(t *testing.T) {
	buf := newBuffer()

	buf.Insert(buf.Length(), 'H')
	buf.Insert(buf.Length(), 'e')
	buf.Insert(buf.Length(), 'l')
	buf.Insert(buf.Length(), 'l')
	buf.Insert(buf.Length(), 'o')
	buf.Insert(buf.Length(), ' ')
	buf.Insert(buf.Length(), 'T')
	buf.Insert(buf.Length(), 'h')
	buf.Insert(buf.Length(), 'e')
	buf.Insert(buf.Length(), 'r')
	buf.Insert(buf.Length(), 'e')

	// Test: Delete first
	ch := buf.Delete(0)
	if ch != 'H' {
		t.Fatalf("Expected 'H', instead ch=" + string(ch))
	}
	if buf.String() != "ello There" {
		t.Fatalf("Expected \"ello There\", instead buf.String(): " + buf.String())
	}

	// Test: Delete Middle Values
	ch = buf.Delete(4)
	if ch != ' ' {
		t.Fatalf("Expected ' ', instead ch=" + string(ch))
	}
	if buf.String() != "elloThere" {
		t.Fatalf("Expected \"elloThere\", instead buf.String(): " + buf.String())
	}

	ch = buf.Delete(5)
	if ch != 'h' {
		t.Fatalf("Expected 'h', instead ch=" + string(ch))
	}
	if buf.String() != "elloTere" {
		t.Fatalf("Expected \"elloTere\", instead buf.String(): " + buf.String())
	}

	// Test: Delete Last
	ch = buf.Delete(buf.Length()-1)
	if ch != 'e' {
		t.Fatalf("Expected 'e', instead ch=" + string(ch))
	}
	if buf.String() != "elloTer" {
		t.Fatalf("Expected \"elloTer\", instead buf.String(): " + buf.String())
	}

	// Test: Length
	if buf.Length() != 7 {
		t.Fatalf("Expected length 7, instead buf.Length(): " + fmt.Sprintf("%d", buf.Length()))
	}

	// Test full deletion
	for buf.Length() > 0 {
		buf.Delete(0)
	}

	if buf.Length() != 0 {
		t.Fatalf("Expected length 0, instead buf.Length(): " + fmt.Sprintf("%d", buf.Length()))
	}
}
