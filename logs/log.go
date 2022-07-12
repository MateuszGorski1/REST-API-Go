package logs

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// C acts as a list of http.Handler constructors
var C alice.Chain

// ServeLogs enables user to see application logs
func ServeLogs() {
	log := zerolog.New(os.Stdout).With().
		Logger()

	C = alice.New()
	C = C.Append(hlog.NewHandler(log))
	C = C.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		params := strings.Split(r.URL.Path, "/")
		var param1, param2 string
		if len(params) == 3 {
			param1 = params[2]
		} else if len(params) == 4 {
			param1 = params[2]
			param2 = params[3]
		}
		e := hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Str("host", r.Host)

		if len(param1) != 0 {
			e = e.Str("param1", param1)
		}
		if len(param2) != 0 {
			e = e.Str("param2", param2)
		}

		e.Send()
	}))

}
