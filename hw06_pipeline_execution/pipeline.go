package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func Worker(in Bi, done In, out Out) {
	defer close(in)
	for {
		select {
		case <-done:
			return
		case tmp, ok := <-out:
			if !ok {
				return
			}
			in <- tmp
		}
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		input := make(Bi)
		go Worker(input, done, in)
		in = stage(input)
	}
	return in
}
