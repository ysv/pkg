package cfip

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sethvargo/go-limiter/httplimit"
	rlredis "github.com/sethvargo/go-redisstore"
)

func New(lmtConf *LimiterConf) (*httplimit.Middleware, error) {
	store, err := rlredis.New(&rlredis.Config{
		Tokens:   lmtConf.Tokens,
		Interval: lmtConf.Interval,
		Dial:     lmtConf.Dialer,
	})

	if err != nil {
		return nil, err
	}

	return httplimit.NewMiddleware(store, keyFunc(lmtConf.Namespace, "rl", "ip", fmt.Sprintf("%s-%d", lmtConf.Interval.String(), lmtConf.Tokens)))
}

func keyFunc(prefixes ...string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		ip := r.Header.Get("CF-CONNECTING-IP")
		return strings.Join(append(prefixes, ip), ":"), nil
	}
}
