/*
package scheme_register Register a custom scheme to the operating system and complete the custom processing

import (
	"log"
	"os"

	"github.com/chyroc/scheme_register"
)

func main() {
	// setup log
	f, err := os.OpenFile("/tmp/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	log.SetPrefix("[example] ")
	log.SetFlags(log.LstdFlags)

	// run
	err = scheme_register.Register(&scheme_register.RegisterReq{
		Name:   "MyApp",
		Scheme: "myapp",
		Handler: func(url string) {
			log.Printf("get: %s", url)
		},
	})
	if err != nil {
		panic(err)
	}
}

and then, your app can handle: `myapp://path?a=b` uri
*/

package scheme_register

