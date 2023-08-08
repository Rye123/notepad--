package textbuffer

import (
	"errors"
)

type TextBuffer interface {
	String() string // Returns contents of the textbuffer as a string
	Insert(index int, ch rune) error // Inserts `ch` into the string at `index`. Error raised if index is not in the range [0, Textbuffer.Length()] (inclusive)
	Append(s string) error // Appends a string `s` into the end of the buffer.
	Delete(index int) rune // Deletes and returns the character at `index`.
	Clear() // Clears the buffer.
	Length() int // Returns the size of the buffer
	GetIndex() int // Returns the current index
	MoveIndex(newIndex int) // Moves index to a new index
}

// Simple buffer for testing Textbuffer functions
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

func (buf *buffer) Append(s string) error {
	for _, c := range(s) {
		err := buf.Insert(buf.Length(), c)
		if err != nil {
			return err
		}
	}
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

func (buf *buffer) Clear() {
	buf.arr = make([]rune, 0)
}

func (buf *buffer) Length() int {
	return len(buf.arr)
}

func (buf *buffer) GetIndex() int {
	return 0
}

func (buf *buffer) MoveIndex(newIndex int) {
	return
}

// A dynamic array with efficient insertion at a particular index
// [0 1 2 3 ...][ GAP ][... 3 2 1 0]
// Invariant: cursorIndex is always at len(left)
type GapBuffer struct {
	left []rune
	right []rune
	cursorIndex int
}

func NewGapBuffer() *GapBuffer {
	return &GapBuffer{
		make([]rune, 0),
		make([]rune, 0),
		0,
	}
}

func (buf *GapBuffer) GetIndex() int {
	return buf.cursorIndex
}

// Moves index to `newIndex`. If `newIndex` is beyond the bounds of the buffer, it stays at the closest bound.
func (buf *GapBuffer) MoveIndex(newIndex int) {
	if newIndex < 0 {
		newIndex = 0
	} else if newIndex > buf.Length() {
		newIndex = buf.Length()
	}
	buf.cursorIndex = newIndex

	for buf.cursorIndex > len(buf.left) {
		// Move characters from buf.right to buf.left
		ch := buf.right[len(buf.right) - 1]
		buf.right = buf.right[:len(buf.right) - 1]
		buf.left = append(buf.left, ch)
	}

	for buf.cursorIndex < len(buf.left) {
		// Move characters from buf.left to buf.right
		ch := buf.left[len(buf.left) - 1]
		buf.left = buf.left[:len(buf.left) - 1]
		buf.right = append(buf.right, ch)
	}
}

func (buf *GapBuffer) String() string{
	fullBuffer := make([]rune, buf.Length())
	copy(fullBuffer, buf.left)
	for i, ch := range(buf.right) {
		fullBufferIndex := len(fullBuffer) -1 - i
		fullBuffer[fullBufferIndex] = ch
	}
	return string(fullBuffer)
}

func (buf *GapBuffer) Insert(index int, ch rune) error {
	if index < 0 || index > buf.Length() {
		return errors.New("Index error")
	}
	
	if index != buf.cursorIndex {
		buf.MoveIndex(index)
	}

	// invariant after MoveIndex: insertion always appends to left stack
	buf.left = append(buf.left, ch)
	buf.cursorIndex++

	return nil
}

func (buf *GapBuffer) Append(s string) error {
	for _, c := range(s) {
		err := buf.Insert(buf.Length(), c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (buf *GapBuffer) Delete(index int) rune {
	if index != buf.cursorIndex {
		buf.MoveIndex(index)
	}

	// invariant after MoveIndex: deletion always deletes from right stack
	if len(buf.right) == 0 {
		return rune(0) // \0 character
	}

	ch := buf.right[len(buf.right)-1]
	buf.right = buf.right[:len(buf.right)-1]
	return ch
}

func (buf *GapBuffer) Clear() {
	buf.left = make([]rune, 0)
	buf.right = make([]rune, 0)
	buf.cursorIndex = 0
}

func (buf *GapBuffer) Length() int {
	return len(buf.left) + len(buf.right)
}
