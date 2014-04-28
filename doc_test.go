package lunk

import "os"

func Example() {
	j := NewJSONEventLogger(os.Stdout)

	root := j.LogRoot(Message("doing something"))
	sub := j.Log(root, root, Message("doing some sub-action"))
	j.Log(root, sub, Message("the sub-action involved something else"))
}
