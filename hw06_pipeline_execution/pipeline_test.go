package hw06pipelineexecution

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}
func TestAnotherStagesPipeline(t *testing.T) {
	// helper function
	chanToSlice := func(c In) []interface{} {
		var slice []interface{}
		for v := range c {
			slice = append(slice, v)
		}
		return slice
	}

	in := make(chan interface{})
	done := make(chan interface{})
	defer close(done)

	go func() {
		in <- 1
		in <- 2
		in <- 3
		close(in)
	}()

	stages := []Stage{
		func(in In) (out Out) {
			outCh := make(Bi)
			go func() {
				defer close(outCh)
				for v := range in {
					outCh <- v.(int) * 2
				}
			}()
			return outCh
		},
		func(in In) (out Out) {
			outCh := make(Bi)
			go func() {
				defer close(outCh)
				for v := range in {
					outCh <- v.(int) + 1
				}
			}()
			return outCh
		},
	}

	out := ExecutePipeline(in, done, stages...)
	result := chanToSlice(out)

	expected := []interface{}{3, 5, 7} // (1*2)+1, (2*2)+1, (3*2)+1

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Got %v, want %v", result, expected)
	}
}
