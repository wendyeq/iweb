package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

var HOST = GetConfig()["host"]
var DATABASE = GetConfig()["database"]

type Article struct {
	Id         bson.ObjectId "_id"
	Author     string        "Author"
	Title      string        "Title"
	PostTime   time.Time     "PostTime"
	UpdateTime time.Time     "UpdateTime"
	Tags       []string      "Tags"
	Content    string        "Content"
}

type CRUD interface {
	Save() error
	FindById() error
	Update() error
	Delete() error
}

func (a *Article) Save() (err error) {
	conn, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	return articles.Insert(a)
}

func (a *Article) FindById() (err error) {
	conn, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	err = articles.Find(bson.M{"_id": a.Id}).One(&a)
	return err
}

func (a *Article) Update() (err error) {
	conn, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")

	return articles.Update(bson.M{"_id": a.Id}, bson.M{"$set": bson.M{"Title": a.Title, "Tags": a.Tags, "Content": a.Content, "UpdateTime": time.Now()}})
}

func (a *Article) Delete() (err error) {
	conn, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	return articles.Remove(bson.M{"_id": a.Id})
}

func (a *Article) FindOne() (err error) {
	conn, err := mgo.Dial(HOST)
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	return articles.Find(bson.M{"Title": a.Title, "PostTime": bson.M{"$gte": a.PostTime, "$lt": a.UpdateTime}}).One(a)
}

func (a *Article) FindAll() (all []Article, err error) {
	conn, err := mgo.Dial(HOST)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	db := conn.DB(DATABASE)
	articles := db.C("articles")
	iter := articles.Find(nil).Sort("-PostTime").Iter()
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
	iter := articles.Find(bson.M{"Tags": tag}).Sort("-PostTime").Iter()
	err = iter.All(&all)
	return all, err
}

func (a *Article) FindAllByArchive(archive string) (all []Article, err error) {
	year := archive[0:4]
	month := archive[5:]
	if len(month) == 1 {
		month = "0" + month
	}
	beginTime, err := time.Parse("2006-01-02", year+"-"+month+"-01")
	endTime := beginTime.AddDate(0, 1, 0)

	conn, err := mgo.Dial(HOST)
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	iter := articles.Find(bson.M{"PostTime": bson.M{"$gte": beginTime, "$lt": endTime}}).Sort("-PostTime").Iter()

	err = iter.All(&all)
	return all, err
}

type Flags struct {
	Single  bool
	Home    bool
	Admin   bool
	Sidebar bool
}

type Data struct {
	Flags Flags
	Vars  interface{}
}

type Result struct {
	Key   string "_id"
	Value int
}

func GetSideBar() interface{} {
	sidebar := make(map[string]interface{})
	sidebar["archives"] = GetArchives()
	sidebar["tags"] = GetTags()
	return sidebar
}

func GetTags() (result []Result) {
	job := &mgo.MapReduce{
		Map: "function() { " +
			"    this.Tags.forEach( " +
			"        function(z){emit(z,1);})}",
		Reduce: "function(key, values) { " +
			"    var total=0; " +
			"    for(var i=0;i<values.length;i++){ " +
			"        total += values[i];} " +
			"    return total;}",
	}
	conn, err := mgo.Dial(HOST)
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	_, err = articles.Find(nil).MapReduce(job, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func GetArchives() (result []Result) {
	job := &mgo.MapReduce{
		Map: "function() { " +
			"    emit(this.PostTime.getFullYear()+'-'+eval(this.PostTime.getMonth()+1),1);}",
		Reduce: "function(key, values) { " +
			"    var total=0; " +
			"    for(var i=0;i<values.length;i++){ " +
			"        total += values[i];} " +
			"    return total;}",
	}
	conn, err := mgo.Dial(HOST)
	defer conn.Close()
	db := conn.DB(DATABASE)
	articles := db.C("articles")
	_, err = articles.Find(nil).MapReduce(job, &result)
	if err != nil {
		panic(err)
	}

	return result
}
