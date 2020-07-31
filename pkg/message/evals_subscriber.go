package message

import (
	"fmt"
	"toggle/server/pkg/create"
	"toggle/server/pkg/evaluate"
)

// StartEvalEventsReceiever handles messages on the evaluations topic
func StartEvalEventsReceiever(c create.Service, m Service) {
	go func() {
		m.Subscribe("evaluations", func(e *evaluate.Evaluation) {
			fmt.Printf("---------------------Received a eval: %+v\n", e)
			e.Count++
			c.CreateEvaluation(e)
		})
	}()
}
