package textbuffer

import (
	"errors"
)

type Textbuffer interface {
	String() string // Returns contents of the textbuffer as a string
	Insert(index int, ch rune) error // Inserts `ch` into the string at `index`. Error raised if index is not in the range [0, Textbuffer.Length()] (inclusive)
	Delete(index int) rune // Deletes and returns the character at `index`.
	Length() int // Returns the size of the buffer
}

type buffer struct {
	arr []rune
}

func newBuffer() *buffer {
	return &buffer{make([]rune, 0)}
}

func (buf *buffer) String() string {
	return string(buf.arr)
}

func (buf *buffer) Insert(index int, ch rune) error {
	if index < 0 || index > len(buf.arr) {
		return errors.New("Index error")
	}
	
	// Case: Insertion at end
	if index == len(buf.arr) {
		buf.arr = append(buf.arr, ch)
		return nil
	}

	// Case: Insertion at beginning or middle
	buf.arr = append(buf.arr, buf.arr[len(buf.arr)-1])
	copy(buf.arr[index+1:], buf.arr[index:])
	buf.arr[index] = ch
	return nil
}

func (buf *buffer) Delete(index int) rune {
	ch := buf.arr[index]
	for i := index; i < len(buf.arr) - 1; i++ {
		buf.arr[i] = buf.arr[i+1]
	}
	buf.arr = buf.arr[:len(buf.arr) - 1]

	return ch
}

func (buf *buffer) Length() int {
	return len(buf.arr)
}

// A dynamic array with efficient insertion at a particular index
// [0 1 2 3 ...][ GAP ][... 3 2 1 0]
type gapbuffer struct {
	left []rune
	right []rune
	cursorIndex int
}

