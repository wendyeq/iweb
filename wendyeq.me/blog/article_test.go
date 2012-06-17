package blog

import (
	"testing"
    "time"
	"launchpad.net/mgo/bson"
)

func TestSave(t *testing.T) {
	article := &Article{Id: bson.NewObjectId(), Author: "wendyeq",Title: "Test1", Content: "good content", Tags: nil, PostTime: time.Now(), UpdateTime: time.Now()}
	err := article.Save()
	if err != nil {
		t.Error("insert error: %s" + err)
	}

}