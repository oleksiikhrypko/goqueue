package klist

import (
	"fmt"
	"testing"
	"time"

	"goqueue/pkg/storage/inmem"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func Test_KList_general(t *testing.T) {
	var err error
	ctx := context.Background()

	db := inmem.NewDB(ctx)

	name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
	list, err := New(ctx, name, db)
	require.NoError(t, err)

	in := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for _, v := range in {
		err = list.Add([]byte(v))
		require.NoError(t, err)
	}

	require.Equal(t, "1", string(list.GetFirst()))
	require.Equal(t, "9", string(list.GetLast()))
	v, err := list.GetNext([]byte("9"))
	require.NoError(t, err)
	require.Nil(t, v)
	v, err = list.GetPrev([]byte("1"))
	require.NoError(t, err)
	require.Nil(t, v)
	v, err = list.GetNext([]byte("2"))
	require.NoError(t, err)
	require.Equal(t, "3", string(v))
	v, err = list.GetPrev([]byte("4"))
	require.NoError(t, err)
	require.Equal(t, "3", string(v))
	b, err := list.IsItemExists([]byte("4"))
	require.NoError(t, err)
	require.True(t, b)
	b, err = list.IsItemExists([]byte("444"))
	require.NoError(t, err)
	require.False(t, b)
	b = list.IsItemFirst([]byte("1"))
	require.True(t, b)
	b = list.IsItemFirst([]byte("111"))
	require.False(t, b)
	b = list.IsItemLast([]byte("1"))
	require.False(t, b)
	b = list.IsItemLast([]byte("999"))
	require.False(t, b)
	b = list.IsItemLast([]byte("9"))
	require.True(t, b)
}

func TestKList_Add(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		in        []string
		expect    []string
		expectRev []string
	}{
		{
			name:      "add 1",
			in:        []string{"1"},
			expect:    []string{"1"},
			expectRev: []string{"1"},
		},
		{
			name:      "add",
			in:        []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "4", "1", "9"},
			expect:    []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
			expectRev: []string{"9", "8", "7", "6", "5", "4", "3", "2", "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmem.NewDB(ctx)

			name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
			l, err := New(ctx, name, db)
			require.NoError(t, err)

			for _, v := range tt.in {
				err = l.Add([]byte(v))
				require.NoError(t, err)
			}

			res := make([]string, 0, len(tt.expect))
			item := l.GetFirst()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetNext(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expect, res)
			require.Equal(t, len(tt.expect), int(l.GetCount()))

			res = make([]string, 0, len(tt.expect))
			item = l.GetLast()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetPrev(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectRev, res)
			require.Equal(t, len(tt.expectRev), int(l.GetCount()))
		})
	}
}

func TestKList_Pop(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		in     []string
		expect []string
	}{
		{
			name:   "pop",
			in:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "4", "1", "9"},
			expect: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmem.NewDB(ctx)
			name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
			l, err := New(ctx, name, db)
			require.NoError(t, err)

			for _, v := range tt.in {
				err = l.Add([]byte(v))
				require.NoError(t, err)
			}

			res := make([]string, 0, len(tt.expect))
			item, err := l.Pop()
			require.NoError(t, err)
			for item != nil {
				res = append(res, string(item))
				item, err = l.Pop()
				require.NoError(t, err)
			}
			require.Equal(t, tt.expect, res)
			require.Equal(t, 0, int(l.GetCount()))
		})
	}
}

func TestKList_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		in        []string
		expect    []string
		expectRev []string
		run       func(*KList)
	}{
		{
			name:      "delete first",
			in:        []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "4", "1", "9"},
			expect:    []string{"2", "3", "4", "5", "6", "7", "8", "9"},
			expectRev: []string{"9", "8", "7", "6", "5", "4", "3", "2"},
			run: func(l *KList) {
				err := l.Delete([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "delete last",
			in:        []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "4", "1", "9"},
			expect:    []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			expectRev: []string{"8", "7", "6", "5", "4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.Delete([]byte("9"))
				require.NoError(t, err)
			},
		},
		{
			name:      "delete mid",
			in:        []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "4", "1", "9"},
			expect:    []string{"1", "2", "3", "4", "6", "7", "8", "9"},
			expectRev: []string{"9", "8", "7", "6", "4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.Delete([]byte("5"))
				require.NoError(t, err)
			},
		},
		{
			name:      "delete mid len3",
			in:        []string{"1", "2", "3"},
			expect:    []string{"1", "3"},
			expectRev: []string{"3", "1"},
			run: func(l *KList) {
				err := l.Delete([]byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "delete last len2",
			in:        []string{"1", "2"},
			expect:    []string{"1"},
			expectRev: []string{"1"},
			run: func(l *KList) {
				err := l.Delete([]byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "delete first len2",
			in:        []string{"1", "2"},
			expect:    []string{"2"},
			expectRev: []string{"2"},
			run: func(l *KList) {
				err := l.Delete([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "delete len1",
			in:        []string{"1"},
			expect:    []string{},
			expectRev: []string{},
			run: func(l *KList) {
				err := l.Delete([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "delete len0",
			in:        []string{},
			expect:    []string{},
			expectRev: []string{},
			run: func(l *KList) {
				err := l.Delete([]byte("1"))
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmem.NewDB(ctx)
			name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
			l, err := New(ctx, name, db)
			require.NoError(t, err)

			for _, v := range tt.in {
				err = l.Add([]byte(v))
				require.NoError(t, err)
			}

			tt.run(l)

			res := make([]string, 0, len(tt.expect))
			item := l.GetFirst()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetNext(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expect, res)
			require.Equal(t, len(tt.expect), int(l.GetCount()))

			res = make([]string, 0, len(tt.expect))
			item = l.GetLast()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetPrev(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectRev, res)
			require.Equal(t, len(tt.expectRev), int(l.GetCount()))
		})
	}
}

func TestKList_SetToBegin(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		in        []string
		expect    []string
		expectRev []string
		run       func(*KList)
	}{
		{
			name:      "len0",
			in:        []string{},
			expect:    []string{"1"},
			expectRev: []string{"1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "len1",
			in:        []string{"2"},
			expect:    []string{"1", "2"},
			expectRev: []string{"2", "1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "len2",
			in:        []string{"2", "3"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "len3",
			in:        []string{"2", "3", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from first",
			in:        []string{"1", "2", "3"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from last",
			in:        []string{"2", "3", "1"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from mid len3",
			in:        []string{"2", "1", "3"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from mid len4",
			in:        []string{"2", "3", "1", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToBegin([]byte("1"))
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmem.NewDB(ctx)
			name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
			l, err := New(ctx, name, db)
			require.NoError(t, err)

			for _, v := range tt.in {
				err = l.Add([]byte(v))
				require.NoError(t, err)
			}

			tt.run(l)

			res := make([]string, 0, len(tt.expect))
			item := l.GetFirst()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetNext(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expect, res)
			require.Equal(t, len(tt.expect), int(l.GetCount()))

			res = make([]string, 0, len(tt.expect))
			item = l.GetLast()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetPrev(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectRev, res)
			require.Equal(t, len(tt.expectRev), int(l.GetCount()))
		})
	}
}

func TestKList_SetToEnd(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		in        []string
		expect    []string
		expectRev []string
		run       func(*KList)
	}{
		{
			name:      "len0",
			in:        []string{},
			expect:    []string{"1"},
			expectRev: []string{"1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "len1",
			in:        []string{"1"},
			expect:    []string{"1", "2"},
			expectRev: []string{"2", "1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "len2",
			in:        []string{"1", "2"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "len3",
			in:        []string{"1", "2", "3"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("4"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from last",
			in:        []string{"1", "2", "3"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from first",
			in:        []string{"3", "1", "2"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from mid len3",
			in:        []string{"1", "3", "2"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move from mid len4",
			in:        []string{"1", "4", "2", "3"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetToEnd([]byte("4"))
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmem.NewDB(ctx)
			name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
			l, err := New(ctx, name, db)
			require.NoError(t, err)

			for _, v := range tt.in {
				err = l.Add([]byte(v))
				require.NoError(t, err)
			}

			tt.run(l)

			res := make([]string, 0, len(tt.expect))
			item := l.GetFirst()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetNext(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expect, res)
			require.Equal(t, len(tt.expect), int(l.GetCount()))

			res = make([]string, 0, len(tt.expect))
			item = l.GetLast()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetPrev(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectRev, res)
			require.Equal(t, len(tt.expectRev), int(l.GetCount()))
		})
	}
}

func TestKList_SetAfter(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		in        []string
		expect    []string
		expectRev []string
		run       func(*KList)
	}{
		{
			name:      "len0",
			in:        []string{},
			expect:    []string{},
			expectRev: []string{},
			run: func(l *KList) {
				err := l.SetAfter([]byte("1"), []byte("1"))
				require.Error(t, err)
				err = l.SetAfter([]byte("1"), []byte("2"))
				require.Error(t, err)
			},
		},
		{
			name:      "len1",
			in:        []string{"1"},
			expect:    []string{"1", "2"},
			expectRev: []string{"2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("2"), []byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "after last len2",
			in:        []string{"1", "2"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("3"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "after first len3",
			in:        []string{"1", "3", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("2"), []byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "before last len3",
			in:        []string{"1", "2", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("3"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move first to last len2",
			in:        []string{"2", "1"},
			expect:    []string{"1", "2"},
			expectRev: []string{"2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("2"), []byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move first to last len3",
			in:        []string{"3", "1", "2"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("3"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 1->2 len3",
			in:        []string{"2", "1", "3"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("2"), []byte("1"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 2->3 len3",
			in:        []string{"1", "3", "2"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("3"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 2->3 len4",
			in:        []string{"1", "3", "2", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("3"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 2->4 len4",
			in:        []string{"1", "4", "2", "3"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("4"), []byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 1->3 len4",
			in:        []string{"3", "1", "2", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetAfter([]byte("3"), []byte("2"))
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmem.NewDB(ctx)
			name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
			l, err := New(ctx, name, db)
			require.NoError(t, err)

			for _, v := range tt.in {
				err = l.Add([]byte(v))
				require.NoError(t, err)
			}

			tt.run(l)

			res := make([]string, 0, len(tt.expect))
			item := l.GetFirst()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetNext(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expect, res)
			require.Equal(t, len(tt.expect), int(l.GetCount()))

			res = make([]string, 0, len(tt.expect))
			item = l.GetLast()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetPrev(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectRev, res)
			require.Equal(t, len(tt.expectRev), int(l.GetCount()))
		})
	}
}

func TestKList_SetBefore(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		in        []string
		expect    []string
		expectRev []string
		run       func(*KList)
	}{
		{
			name:      "len0",
			in:        []string{},
			expect:    []string{},
			expectRev: []string{},
			run: func(l *KList) {
				err := l.SetBefore([]byte("1"), []byte("1"))
				require.Error(t, err)
				err = l.SetBefore([]byte("1"), []byte("2"))
				require.Error(t, err)
			},
		},
		{
			name:      "len1",
			in:        []string{"2"},
			expect:    []string{"1", "2"},
			expectRev: []string{"2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("1"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "before last len2",
			in:        []string{"1", "3"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("2"), []byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "after first len3",
			in:        []string{"1", "3", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("2"), []byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "before last len3",
			in:        []string{"1", "2", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("3"), []byte("4"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move last to first len2",
			in:        []string{"2", "1"},
			expect:    []string{"1", "2"},
			expectRev: []string{"2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("1"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move last to first len3",
			in:        []string{"2", "3", "1"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("1"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 1<-2 len3",
			in:        []string{"2", "1", "3"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("1"), []byte("2"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 2<-3 len3",
			in:        []string{"1", "3", "2"},
			expect:    []string{"1", "2", "3"},
			expectRev: []string{"3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("2"), []byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 2<-3 len4",
			in:        []string{"1", "3", "2", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("2"), []byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 2<-4 len4",
			in:        []string{"1", "3", "4", "2"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("2"), []byte("3"))
				require.NoError(t, err)
			},
		},
		{
			name:      "move 1<-3 len4",
			in:        []string{"2", "3", "1", "4"},
			expect:    []string{"1", "2", "3", "4"},
			expectRev: []string{"4", "3", "2", "1"},
			run: func(l *KList) {
				err := l.SetBefore([]byte("1"), []byte("2"))
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := inmem.NewDB(ctx)
			name := fmt.Sprintf("q:%d:", time.Now().UnixMicro())
			l, err := New(ctx, name, db)
			require.NoError(t, err)

			for _, v := range tt.in {
				err = l.Add([]byte(v))
				require.NoError(t, err)
			}

			tt.run(l)

			res := make([]string, 0, len(tt.expect))
			item := l.GetFirst()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetNext(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expect, res)
			require.Equal(t, len(tt.expect), int(l.GetCount()))

			res = make([]string, 0, len(tt.expect))
			item = l.GetLast()
			for item != nil {
				res = append(res, string(item))
				item, err = l.GetPrev(item)
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectRev, res)
			require.Equal(t, len(tt.expectRev), int(l.GetCount()))
		})
	}
}
