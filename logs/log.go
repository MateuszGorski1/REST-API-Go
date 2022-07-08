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

var C alice.Chain

//Enables user to see application logs
func ServeLogs() {
	log := zerolog.New(os.Stdout).With().
		Logger()

	C = alice.New()
	C = C.Append(hlog.NewHandler(log))
	C = C.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		params := strings.Split(r.URL.Path, "/")
		var param1, param2 string
		if len(params) < 3 {
			param2 = "null"
			param1 = "null"
		} else if len(params) < 4 {
			param2 = "null"
			param1 = params[2]
		} else {
			param1 = params[2]
			param2 = params[3]
		}
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Str("host", r.Host).
			Str("param 1", param1).
			Str("param 2", param2).
			Msg("")
	}))

}
