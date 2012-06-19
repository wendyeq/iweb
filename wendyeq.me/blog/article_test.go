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
		t.Errorf("insert error: %s", err)
	}

}

func TestFindAll(t *testing.T) {
	article := new(Article)
	all, err := article.FindAll()
	if err != nil {
		t.Errorf("len %d, find all error: %s", len(all), err)
	}
}

func TestFindAllByTag(t *testing.T) {
	article := new(Article)
	tag := "读书笔记"
	all, err := article.FindAllByTag(tag)
	if err != nil {
		t.Errorf("len %d, find all by tag error: %s", len(all), err)
	}
}