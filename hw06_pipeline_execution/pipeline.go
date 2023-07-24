package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	input := in
	for _, stage := range stages {
		stageOutput := stage(input)
		newInput := make(Bi)
		go func() {
			defer close(newInput)
			for {
				select {
				case val, ok := <-stageOutput:
					if !ok {
						return
					}
					newInput <- val
				case <-done:
					return
				}
			}
		}()
		input = newInput
	}
	return input
}
