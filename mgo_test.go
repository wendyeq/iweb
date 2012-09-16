package main

import (
	"labix.org/v2/mgo"
	"testing"
)

func TestSession(t *testing.T) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		t.Error(err)
	}
	defer session.Close()
}
