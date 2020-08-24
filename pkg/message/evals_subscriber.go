package message

import (
	"time"
	"toggle/server/pkg/create"
	"toggle/server/pkg/models"
)

// StartEvalEventsReceiever persists the evaluation into database and updates
// the flag's last evaluated property
func StartEvalEventsReceiever(c create.Service, m Service) {
	go func() {
		m.Subscribe("evaluations", func(e *models.Evaluation) {
			e.CreatedAt = time.Now()
			c.CreateEvaluation(e)
			flag := e.Flag
			flag.Evaluated = time.Now()
			c.UpdateFlag(&flag)
		})
	}()
}
