package blog

import (
    "testing"
	"launchpad.net/mgo"
)

func TestSession(t *testing.T) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		t.Error(err)
	}

	defer session.Close()
}
