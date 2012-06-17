package blog

import (
	"time"
	"launchpad.net/mgo/bson"
	"launchpad.net/mgo"
)
const HOST = "localhost"
const DATABASE = "iwebdb"

type Article struct {
	Id bson.ObjectId "_id"
	Author string "Author"
	Title string "Title"
	PostTime time.Time "PostTime"
	UpdateTime time.Time "UpdateTime"
	Tags []string "Tags"
	Content string "Content"
}

func (a *Article) Save() (err error ) {
	conn, err := mgo.Dial(HOST)
    if err != nil {
		panic(err)
	}
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	return articles.Insert(a)
}

