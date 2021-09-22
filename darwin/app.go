package darwin

// doc: https://blakewilliams.me/posts/handling-macos-url-schemes-with-go

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "register.h"
*/
import "C"

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chyroc/scheme_register/helper"
)

//go:embed Info.plist
var info string

func Register(name, scheme string, f func(url string)) error {
	scheme = strings.ToLower(scheme)

	if isRunInApp(name) {
		handle(scheme, f)
		return nil
	}

	if err := reInstallApp(name, scheme); err != nil {
		return err
	}

	return nil
}

var urlListener chan string = make(chan string)

func handle(scheme string, realHandle func(url string)) {
	log.Printf("[scheme_register] %s", scheme)

	go func() {
		timeout := time.After(4 * time.Second)
		select {
		case url := <-urlListener:
			log.Printf("[scheme_register] handle %s", url)
			realHandle(url)
			os.Exit(0)
		case <-timeout:
			log.Printf("[scheme_register] timeout")
			os.Exit(1)
		}
	}()

	C.RunApp()
}

func ShowError(title string, details string) {
	C.ShowAlert(
		C.CString(title),
		C.CString(details),
	)
}

//export HandleURL
func HandleURL(u *C.char) {
	urlListener <- C.GoString(u)
}

func reInstallApp(name, scheme string) error {
	infoData := info
	infoData = strings.ReplaceAll(infoData, "{{Name}}", name)
	infoData = strings.ReplaceAll(infoData, "{{Scheme}}", scheme)

	tempDir, err := ioutil.TempDir("", "scheme-register-*")
	if err != nil {
		return err
	}

	tempAppDir := fmt.Sprintf("%s/%s.app", tempDir, name)
	appDir := fmt.Sprintf("/Applications/%s.app", name)

	if err = os.MkdirAll(fmt.Sprintf("%s/Contents", tempAppDir), 0o777); err != nil {
		return err
	}

	if err = helper.Copy(os.Args[0], fmt.Sprintf("%s/Contents/%s", tempAppDir, name)); err != nil {
		return err
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s/Contents/Info.plist", tempAppDir), []byte(infoData), 0o666); err != nil {
		return err
	}

	if err = os.RemoveAll(appDir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if err = os.Rename(tempAppDir, fmt.Sprintf("/Applications/%s.app", name)); err != nil {
		return err
	}

	return nil
}

func isRunInApp(name string) bool {
	return strings.HasPrefix(os.Args[0], fmt.Sprintf("/Applications/%s.app", name))
}
