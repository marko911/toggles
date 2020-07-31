package message

import (
	"fmt"
	"toggle/server/pkg/create"
	"toggle/server/pkg/models"
)

// StartEvalEventsReceiever handles messages on the evaluations topic
func StartEvalEventsReceiever(c create.Service, m Service) {
	go func() {
		m.Subscribe("evaluations", func(e *models.Evaluation) {
			fmt.Printf("---------------------Received a eval: %+v\n", e)
			c.CreateEvaluation(e)
		})
	}()
}
