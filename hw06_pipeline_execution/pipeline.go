package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = stage(wrapWithDone(out, done))
	}

	return out
}

// wrapWithDone прокси в канал out+проверка отмены через done.
func wrapWithDone(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case data, ok := <-in:
				if !ok {
					return
				}
				// прокси в канал out
				out <- data

			case <-done:
				return
			}
		}
	}()
	return out
}
