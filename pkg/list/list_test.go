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
		l.AddElement(2, nil, nil)
		require.Equal(t, 1, l.Len())
	})

	t.Run("add another node before", func(t *testing.T) {
		l.AddElement(1, nil, l.Head)
		require.Equal(t, 2, l.Len())
		require.Equal(t, 1, l.Head.Value)
		require.Equal(t, 2, l.Tail.Value)
	})

	t.Run("add another after", func(t *testing.T) {
		l.AddElement(5, l.Tail, nil)
		require.Equal(t, []int{1, 2, 5}, l.List())
		require.Equal(t, 3, l.Len())
	})

	t.Run("add one in the middle after", func(t *testing.T) {
		l.AddElement(3, l.Head.Next, nil)
		require.Equal(t, []int{1, 2, 3, 5}, l.List())
		require.Equal(t, 4, l.Len())
	})

	t.Run("add one in the middle before", func(t *testing.T) {
		l.AddElement(4, nil, l.Tail)
		require.Equal(t, []int{1, 2, 3, 4, 5}, l.List())
		require.Equal(t, 5, l.Len())
	})

	t.Run("delete first", func(t *testing.T) {
		l.DeleteElement(l.Head)
		require.Equal(t, []int{2, 3, 4, 5}, l.List())
		require.Equal(t, 4, l.Len())
	})

	t.Run("delete last", func(t *testing.T) {
		l.DeleteElement(l.Tail)
		require.Equal(t, []int{2, 3, 4}, l.List())
		require.Equal(t, 3, l.Len())
	})

	t.Run("delete the middle one", func(t *testing.T) {
		l.DeleteElement(l.Head.Next)
		require.Equal(t, []int{2, 4}, l.List())
		require.Equal(t, 2, l.Len())
	})

	t.Run("delete one of the rest two", func(t *testing.T) {
		l.DeleteElement(l.Head)
		require.Equal(t, l.Head, l.Tail)
		require.Equal(t, []int{4}, l.List())
		require.Equal(t, 1, l.Len())
	})

	t.Run("delete the last one", func(t *testing.T) {
		l.DeleteElement(l.Tail)
		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Head)
		require.Nil(t, l.Tail)
	})
}