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


func (a *Article) FindOne() (err error) {
	conn, err := mgo.Dial(HOST)
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	return articles.Find(bson.M{"Title":a.Title, "PostTime":bson.M{"$gte": a.PostTime, "$lt": a.UpdateTime}}).One(a)
}

func (a *Article) FindAll() (all []Article, err error){
	conn, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	db := conn.DB(DATABASE)
	articles := db.C("articles")
	iter := articles.Find(nil).Sort(bson.M{"PostTime":-1}).Iter()
    err = iter.All(&all)
	return all, err
} 

func (a *Article) FindAllByTag(tag string) (all []Article, err error) {
	conn, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	db := conn.DB(DATABASE)
	articles := db.C("articles")
	iter := articles.Find(bson.M{"Tags":tag}).Sort(bson.M{"PostTime":-1}).Iter()
	err = iter.All(&all)
	return all, err
}

func (a *Article) FindAllByArchive(archive string) (all []Article, err error) {
	year := archive[0:4]
	month := archive[5:]
	if len(month) == 1 {
		month = "0"+month
	}
	beginTime,err := time.Parse("2006-01-02", year+"-"+month+"-01")
	endTime := beginTime.AddDate(0,1,0)

	conn, err := mgo.Dial(HOST)
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	iter := articles.Find(bson.M{"PostTime":bson.M{"$gte": beginTime, "$lt": endTime}}).Sort(bson.M{"PostTime":-1}).Iter()
	
	err = iter.All(&all)
	return all, err
}