package scheme_register

import (
	"fmt"
	"runtime"

	"github.com/chyroc/scheme_register/darwin"
)

type RegisterReq struct {
	Name    string
	Scheme  string
	Handler func(url string)
}

func Register(req *RegisterReq) error {
	switch runtime.GOOS {
	case "darwin":
		return darwin.Register(req.Name, req.Scheme, req.Handler)
	default:
		return fmt.Errorf("unsupport " + runtime.GOOS)
	}
}
