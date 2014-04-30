package lunk

import "os"

func Example() {
	l := NewJSONEventLogger(os.Stdout)

	rootID := NewRootEventID()
	l.Log(rootID, Message("root action"))

	subID := NewEventID(rootID)
	l.Log(subID, Message("sub action"))

	leafID := NewEventID(rootID)
	l.Log(leafID, Message("leaf action"))

	// Produces something like this:
	// {
	//     "event": {
	//         "msg": "root action"
	//     },
	//     "pid": 44345,
	//     "host": "server1.example.com",
	//     "time": "2014-04-28T13:58:32.201883418-07:00",
	//     "id": "09c84ee90e7d9b74",
	//     "root": "09c84ee90e7d9b74",
	//     "schema": "message"
	// }
	// {
	//     "event": {
	//         "msg": "sub action"
	//     },
	//     "pid": 44345,
	//     "host": "server1.example.com",
	//     "time": "2014-04-28T13:58:32.202241745-07:00",
	//     "parent": "09c84ee90e7d9b74",
	//     "id": "794f8bde67a7f1a7",
	//     "root": "09c84ee90e7d9b74",
	//     "schema": "message"
	// }
	// {
	//     "event": {
	//         "msg": "leaf action"
	//     },
	//     "pid": 44345,
	//     "host": "server1.example.com",
	//     "time": "2014-04-28T13:58:32.202257354-07:00",
	//     "parent": "794f8bde67a7f1a7",
	//     "id": "33cff19e8bfb7cef",
	//     "root": "09c84ee90e7d9b74",
	//     "schema": "message"
	// }
}
