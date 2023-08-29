package a

import (
	_ "io/ioutil" // want "\"io/ioutil\" is not allowed to be imported"

	oldContext "golang.org/x/net/context"
)

func a(_ oldContext.Context) {
}
