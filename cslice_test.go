package cslice

import (
	"slices"
	"testing"
)

func TestFront(t *testing.T) {
	c := New[int]()
	c.PushFront(2, 1)

	if item, ok := c.Front(); !ok || item != 2 {
		t.Error("expected 2, got ", item)
	}
	if item, ok := c.PopFront(); !ok || item != 2 {
		t.Error("expected 2, got ", item)
	}
	if item, ok := c.PopFront(); !ok || item != 1 {
		t.Error("expected 1, got ", item)
	}
	if item, ok := c.PopFront(); ok {
		t.Error("expected empty, got ", item)
	}
	if item, ok := c.Front(); ok {
		t.Error("expected empty, got ", item)
	}
}

func TestBack(t *testing.T) {
	c := New[int]()
	c.PushBack(1, 2)

	if item, ok := c.Back(); !ok || item != 2 {
		t.Error("expected 2, got ", item)
	}
	if item, ok := c.PopBack(); !ok || item != 2 {
		t.Error("expected 2, got ", item)
	}
	if item, ok := c.PopBack(); !ok || item != 1 {
		t.Error("expected 1, got ", item)
	}
	if item, ok := c.PopBack(); ok {
		t.Error("expected empty, got ", item)
	}
	if item, ok := c.Back(); ok {
		t.Error("expected empty, got ", item)
	}
}

func TestGetSet(t *testing.T) {
	c := New[int]()
	c.Set(1, 2)
	c.Set(0, 1)
	if item, ok := c.Get(0); !ok || item != 1 {
		t.Error("expected 1, got ", item)
	}
	if item, ok := c.Get(1); !ok || item != 2 {
		t.Error("expected 2, got ", item)
	}
	if item, ok := c.Get(2); ok {
		t.Error("expected empty, got ", item)
	}
}

func TestRemove(t *testing.T) {
	c := From[int]([]int{1, 2})
	if item, ok := c.Remove(1); !ok || item != 2 {
		t.Error("expected 2, got ", item)
	}
	if item, ok := c.Remove(0); !ok || item != 1 {
		t.Error("expected 1, got ", item)
	}
	if item, ok := c.Remove(0); ok {
		t.Error("expected empty, got ", item)
	}
}

func TestCopy(t *testing.T) {
	c := New[int]()
	copied := c.Copy([]int{1, 2, 3, 4, 5})
	if copied != 5 {
		t.Error("expected 5, got ", copied)
	}
}

func TestRange(t *testing.T) {
	c := From[int]([]int{1, 2, 3, 4, 5})
	count := 0
	c.Range(func(index int, item int) bool {
		if item != index+1 {
			t.Error("expected ", index+1, ", got ", item)
		}
		count += 1
		return true
	})
	if count != 5 {
		t.Error("expected 5, got ", count)
	}
}

func TestIter(t *testing.T) {
	c := From([]int{1, 2, 3, 4, 5})
	count := 0
	for item := range c.Iter() {
		if item != count+1 {
			t.Error("expected ", count+1, ", got ", item)
		}
		count += 1
	}
	if count != 5 {
		t.Error("expected 5, got ", count)
	}
}

func TestClear(t *testing.T) {
	c := From([]int{1, 2, 3, 4, 5})
	c.Clear()
	if c.Empty() != true {
		t.Error("expected empty, got not empty")
	}
}

func TestSlice(t *testing.T) {
	c := From([]int{1, 2, 3, 4, 5})
	fullCopy := c.Slice()
	if slices.Equal(fullCopy, []int{1, 2, 3, 4, 5}) != true {
		t.Error("expected [1 2 3 4 5], got ", fullCopy)
	}
	middleToEnd := c.Slice(2)
	if slices.Equal(middleToEnd, []int{3, 4, 5}) != true {
		t.Error("expected [3 4 5], got ", middleToEnd)
	}
	middleToStart := c.Slice(0, 2)
	if slices.Equal(middleToStart, []int{1, 2}) != true {
		t.Error("expected [1 2], got ", middleToStart)
	}
	negativeIndex := c.Slice(0, -1)
	if slices.Equal(negativeIndex, []int{1, 2, 3, 4, 5}) != true {
		t.Error("expected [1 2 3 4 5], got ", negativeIndex)
	}
}
