package main

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var loginTPL *template.Template
var mainTPL *template.Template
var adminTPL *template.Template
var editTPL *template.Template
var rssTPL *template.Template

func init() {
	var err error
	loginTPL, err = template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println("init loginTPL err ")
		log.Println(err)
	}

	mainTPL, err = template.ParseFiles(
		"templates/common/header.html",
		"templates/main.html",
		"templates/blog/article.html",
		"templates/blog/articles.html",
		"templates/common/sidebar.html",
		"templates/common/footer.html")
	if err != nil {
		log.Println("init mainTPL err ")
		log.Println(err)
	}

	adminTPL, err = template.ParseFiles(
		"templates/common/header.html",
		"templates/main.html",
		"templates/blog/article.html",
		"templates/admin/articles.html",
		"templates/common/sidebar.html",
		"templates/common/footer.html")
	if err != nil {
		log.Println("init adminTPL err ")
		log.Println(err)
	}

	editTPL, err = template.ParseFiles("templates/blog/edit_article.html")
	if err != nil {
		log.Println("init editTPL err ")
		log.Println(err)
	}

	rssTPL, err = template.ParseFiles("templates/rss.xml")
	if err != nil {
		log.Println("init rssTPL err ")
		log.Println(err)
	}
	if err == nil {
		log.Println("init templates ok ")
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		loginTPL.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm()
		fmt.Println(r.Form["user"])
		fmt.Println(r.Form["password"])
		http.Redirect(w, r, "/", http.StatusFound)
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}

var dateTime = regexp.MustCompile("^[0-9]{4}/[0-9]{2}/[0-9]{2}/+")

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/blog/"):]
	if !dateTime.MatchString(url) {
		http.NotFound(w, r)
		return
	}
	year := url[0:4]
	month := url[5:7]
	day := url[8:10]
	title := url[11:]
	postTime, _ := time.Parse("2006-01-02", year+"-"+month+"-"+day)
	updateTime := postTime.AddDate(0, 0, 1)
	article := &Article{Title: title, PostTime: postTime, UpdateTime: updateTime}
	err := article.FindOne()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	vars := make(map[string]interface{})
	vars["article"] = article
	vars["sidebar"] = GetSideBar()
	data := Data{}
	data.Flags.Single = true
	data.Flags.Sidebar = true
	data.Vars = vars
	mainTPL.ExecuteTemplate(w, "main", &data)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id != "" {
		article := &Article{Id: bson.ObjectIdHex(id)}
		article.FindById()
		editTPL.Execute(w, article)
	} else {
		editTPL.Execute(w, nil)
	}
}

func SaveArticleHandler(w http.ResponseWriter, r *http.Request) {
	author := "wendyeq"
	title := r.FormValue("title")
	tags := r.FormValue("tags")
	content := r.FormValue("content")
	tags = strings.TrimSpace(tags)
	tags = strings.Replace(tags, "，", ",", -1)
	tags = strings.Replace(tags, " ", "", -1)
	tag := strings.Split(tags, ",")
	time := time.Now()
	article := &Article{Id: bson.NewObjectId(), Author: author, Title: title, Content: content, Tags: tag, PostTime: time, UpdateTime: time}
	err := article.Save()
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	title := r.FormValue("title")
	tags := r.FormValue("tags")
	content := r.FormValue("content")
	tags = strings.TrimSpace(tags)
	tags = strings.Replace(tags, "，", ",", -1)
	tags = strings.Replace(tags, " ", ",", -1)
	tag := strings.Split(tags, ",")
	article := &Article{Id: bson.ObjectIdHex(id), Title: title, Content: content, Tags: tag}
	article.Update()
	http.Redirect(w, r, "/", http.StatusFound)
}

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/admin/blog/article/delete/"):]
	article := &Article{Id: bson.ObjectIdHex(id)}
	err := article.Delete()
	if err == nil {
		http.Redirect(w, r, "/admin/", http.StatusFound)
	}
}

func TagHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Path[len("/blog/tag/"):]
	article := &Article{}
	articles, err := article.FindAllByTag(tag)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	vars := make(map[string]interface{})
	vars["articles"] = articles
	vars["sidebar"] = GetSideBar()
	data := Data{}
	data.Flags.Home = true
	data.Flags.Sidebar = true
	data.Vars = vars
	mainTPL.ExecuteTemplate(w, "main", &data)
}

func ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	archive := r.URL.Path[len("/blog/archive/"):]
	article := Article{}
	articles, err := article.FindAllByArchive(archive)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	vars := make(map[string]interface{})
	vars["articles"] = articles
	vars["sidebar"] = GetSideBar()
	data := Data{}
	data.Flags.Home = true
	data.Flags.Sidebar = true
	data.Vars = vars
	mainTPL.ExecuteTemplate(w, "main", &data)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	article := Article{}
	articles, err := article.FindAll()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	vars := make(map[string]interface{})
	vars["articles"] = articles
	vars["sidebar"] = GetSideBar()
	data := Data{}
	data.Flags.Home = true
	data.Flags.Sidebar = true
	data.Vars = vars
	adminTPL.ExecuteTemplate(w, "main", &data)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	article := Article{}
	articles, err := article.FindAll()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	vars := make(map[string]interface{})
	vars["articles"] = articles
	vars["sidebar"] = GetSideBar()
	data := Data{}
	data.Flags.Home = true
	data.Flags.Sidebar = true
	data.Vars = vars

	//b, _ := json.Marshal(group)
	//w.Header().Add("Content-Type", "text/javascript")
	//w.Write(b)
	mainTPL.ExecuteTemplate(w, "main", &data)
}

func Rss(w http.ResponseWriter, r *http.Request) {
	article := Article{}
	articles, err := article.FindAll()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	rssTPL.Execute(w, articles)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	http.ServeFile(w, r, filepath.Join(".", path))
}
