package main

import (
	"os"

	"bitbucket.org/codahale/lunk"
)

func main() {
	l := lunk.NewJSONEventLogger(os.Stdout)
	root := l.LogRoot(lunk.Message("root action"))
	sub := l.Log(root, root, lunk.Message("sub action"))
	l.Log(root, sub, lunk.Message("leaf action"))
}
