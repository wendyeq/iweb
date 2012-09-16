package main

import (
	"log"
	"net/http"
)

func InitRouter() {
	http.HandleFunc("/static/", StaticHandler)
	http.HandleFunc("/ueditor/", StaticHandler)
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	http.HandleFunc("/blog/", ViewHandler)
	http.HandleFunc("/blog/tag/", TagHandler)
	http.HandleFunc("/blog/archive/", ArchiveHandler)

	http.HandleFunc("/admin/", AdminHandler)
	http.HandleFunc("/admin/blog/article/post", SaveArticleHandler)
	http.HandleFunc("/admin/blog/article/edit", EditHandler)
	http.HandleFunc("/admin/blog/article/update", UpdateArticleHandler)
	http.HandleFunc("/admin/blog/article/delete/", DeleteArticleHandler)

	http.HandleFunc("/feed", Rss)
	http.HandleFunc("/feed/atom", Rss)
	http.HandleFunc("/rss.xml", Rss)
	log.Println("init router ok")
}
