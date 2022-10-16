package middleware

import (
	"myapp/data"

	"github.com/aethedigm/celeritas"
)

type Middleware struct {
	App    *celeritas.Celeritas
	Models data.Models
}
