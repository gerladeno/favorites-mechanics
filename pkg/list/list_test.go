//nolint:paralleltest,funlen
package list_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gerladeno/favorites-mechanics/pkg/list"
)

func TestDeLinkedList(t *testing.T) {
	l := list.DeLinkedList[int]{} //nolint:varnamelen

	t.Run("delete nil from empty list", func(t *testing.T) {
		l.DeleteElement(nil)
		require.Equal(t, 0, l.Len())
	})

	t.Run("delete node from empty list", func(t *testing.T) {
		l.DeleteElement(&list.Node[int]{})
		require.Equal(t, 0, l.Len())
	})

	t.Run("add first node", func(t *testing.T) {
		node := l.AddElement(2, nil, nil)
		require.Equal(t, 1, l.Len())
		require.Equal(t, 2, node.Value)
	})

	t.Run("add another node before", func(t *testing.T) {
		node := l.AddElement(1, nil, l.Head)
		require.Equal(t, 2, l.Len())
		require.Equal(t, 1, l.Head.Value)
		require.Equal(t, 2, l.Tail.Value)
		require.Equal(t, 1, node.Value)
	})

	t.Run("add another after", func(t *testing.T) {
		node := l.AddElement(5, nil, nil)
		require.Equal(t, []int{1, 2, 5}, l.List())
		require.Equal(t, 3, l.Len())
		require.Equal(t, 5, node.Value)
	})

	t.Run("add one in the middle after", func(t *testing.T) {
		node := l.AddElement(3, l.Head.Next, nil)
		require.Equal(t, []int{1, 2, 3, 5}, l.List())
		require.Equal(t, 4, l.Len())
		require.Equal(t, 3, node.Value)
	})

	t.Run("add one in the middle before", func(t *testing.T) {
		node := l.AddElement(4, nil, l.Tail)
		require.Equal(t, []int{1, 2, 3, 4, 5}, l.List())
		require.Equal(t, 5, l.Len())
		require.Equal(t, 4, node.Value)
	})

	t.Run("swap inner", func(t *testing.T) {
		l.SwapItems(l.Head.Next, l.Tail.Prev)
		require.Equal(t, []int{1, 4, 3, 2, 5}, l.List())
	})

	t.Run("swap first with inner", func(t *testing.T) {
		l.SwapItems(l.Head, l.Tail.Prev)
		require.Equal(t, []int{2, 4, 3, 1, 5}, l.List())
	})

	t.Run("swap last with inner", func(t *testing.T) {
		l.SwapItems(l.Head.Next, l.Tail)
		require.Equal(t, []int{2, 5, 3, 1, 4}, l.List())
	})

	t.Run("swap first with last", func(t *testing.T) {
		l.SwapItems(l.Head, l.Tail)
		require.Equal(t, []int{4, 5, 3, 1, 2}, l.List())
	})

	t.Run("move first to last", func(t *testing.T) {
		l.MoveItem(l.Head, l.Tail, nil)
		require.Equal(t, []int{5, 3, 1, 2, 4}, l.List())
	})

	t.Run("move inner to first", func(t *testing.T) {
		l.MoveItem(l.Tail.Prev, nil, l.Head)
		require.Equal(t, []int{2, 5, 3, 1, 4}, l.List())
	})

	t.Run("move last to inner", func(t *testing.T) {
		l.MoveItem(l.Tail, l.Head, nil)
		require.Equal(t, []int{2, 4, 5, 3, 1}, l.List())
	})

	t.Run("move inner to last", func(t *testing.T) {
		l.MoveItem(l.Head.Next.Next, l.Tail, nil)
		require.Equal(t, []int{2, 4, 3, 1, 5}, l.List())
	})

	t.Run("move last to first", func(t *testing.T) {
		l.MoveItem(l.Tail, nil, l.Head)
		require.Equal(t, []int{5, 2, 4, 3, 1}, l.List())
	})

	t.Run("move inner to inner", func(t *testing.T) {
		l.MoveItem(l.Head.Next, nil, l.Tail.Prev)
		require.Equal(t, []int{5, 4, 2, 3, 1}, l.List())
	})

	t.Run("delete first", func(t *testing.T) {
		l.DeleteElement(l.Head)
		require.Equal(t, []int{4, 2, 3, 1}, l.List())
		require.Equal(t, 4, l.Len())
	})

	t.Run("delete last", func(t *testing.T) {
		l.DeleteElement(l.Tail)
		require.Equal(t, []int{4, 2, 3}, l.List())
		require.Equal(t, 3, l.Len())
	})

	t.Run("delete the middle one", func(t *testing.T) {
		l.DeleteElement(l.Head.Next)
		require.Equal(t, []int{4, 3}, l.List())
		require.Equal(t, 2, l.Len())
	})

	t.Run("delete one of the rest two", func(t *testing.T) {
		l.DeleteElement(l.Head)
		require.Equal(t, l.Head, l.Tail)
		require.Equal(t, []int{3}, l.List())
		require.Equal(t, 1, l.Len())
	})

	t.Run("delete the last one", func(t *testing.T) {
		l.DeleteElement(l.Tail)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Head)
		require.Nil(t, l.Tail)
	})
}
