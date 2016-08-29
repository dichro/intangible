package async

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/net/context"
)

var canceled, cancel = context.WithCancel(context.Background())

func init() { cancel() }

func TestRingPeek(t *testing.T) {
	r := NewRingBuffer(5)
	for i := 0; i < 7; i++ {
		r.Push(i)
	}
	for i, test := range []struct {
		fromSeq, toSeq, dataMin, dataLen int
	}{
		{0, 7, 0, 0},
		{2, 7, 2, 5},
		{6, 7, 6, 1},
	} {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		toSeq, data := r.peek(ctx, test.fromSeq)
		if toSeq != test.toSeq || len(data) != test.dataLen {
			t.Errorf(`[%d] r.Peek(ctx, %d) = %d, len(%d); want %d, len(%d)`, i, test.fromSeq, toSeq, len(data), test.toSeq, test.dataLen)
		}
		if test.dataLen == 0 && data != nil {
			t.Errorf("[%d] r.Peek(ctx, %d) = %d, %s; want %d, <nil>", i, test.fromSeq, toSeq, data, test.toSeq)
		}
	}
}

func ExampleRing() {
	r := NewRingBuffer(10)
	r.Push(1) // ignored, because no watchers
	ch := make(chan []interface{})
	go r.Watch(context.Background(), ch)
	<-ch      // read ready signal
	r.Push(2) // this will be visible to the watcher
	fmt.Println(<-ch)
	r.Push(3) // rapid updates will be combined
	r.Push(4)
	fmt.Println(<-ch)
	// Output:
	// [2]
	// [3 4]
}
