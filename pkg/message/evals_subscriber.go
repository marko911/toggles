package message

import (
	"fmt"
	"time"
	"toggle/server/pkg/create"
	"toggle/server/pkg/models"
)

// StartEvalEventsReceiever persists the evaluation into database and updates
// the flag's last evaluated property
func StartEvalEventsReceiever(c create.Service, m Service) {
	go func() {
		m.Subscribe("evaluations", func(e *models.Evaluation) {
			fmt.Printf("---------------------Received a eval: %+v\n", e)
			e.CreatedAt = time.Now()
			c.CreateEvaluation(e)
			flag := e.Flag
			flag.Evaluated = time.Now()
			c.UpdateFlag(&flag)
		})
	}()
}
