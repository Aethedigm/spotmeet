// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package middleware

import (
	"myapp/data"

	"github.com/aethedigm/celeritas"
)

type Middleware struct {
	App    *celeritas.Celeritas
	Models data.Models
}
