package main

import (
	"labix.org/v2/mgo/bson"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	postTime, _ := time.Parse("2006-01-02", "2012-06-01")
	updateTime, _ := time.Parse("2006-01-02", "2012-06-02")
	article := &Article{Id: bson.NewObjectId(), Author: "wendyeq", Title: "Test1", Content: "good content" + time.Now().Format(time.RFC3339), Tags: nil, PostTime: postTime, UpdateTime: updateTime}
	err := article.Save()
	if err != nil {
		t.Errorf("insert error: %s", err)
	}

}

func TestFindOne(t *testing.T) {
	postTime, _ := time.Parse("2006-01-02", "2012-06-01")
	updateTime, _ := time.Parse("2006-01-02", "2012-06-02")
	article := &Article{Id: bson.NewObjectId(), Author: "wendyeq", Title: "Test1", Content: "good", Tags: nil, PostTime: postTime, UpdateTime: updateTime}
	err := article.FindOne()
	if err != nil || article.Content == "good" {
		t.Errorf("find one article content is: %s, error: %s", article.Content, err)
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

func TestFindAllByArchive(t *testing.T) {
	article := new(Article)
	archive := "201206"
	all, err := article.FindAllByArchive(archive)
	if err != nil {
		t.Errorf("len %d, find all by archive error: %s", len(all), err)
	}
}
