package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/naoina/genmai"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type Failure struct {
	ID        int64     `db:"pk"`
	Title     string    `db:"unique" json:"title"`
	Body      string    `json:"body"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var indexTmpl *template.Template = template.Must(template.ParseFiles("temp/index.tmpl"))
var dbObj *genmai.DB

func main() {
	db, err := genmai.New(&genmai.SQLite3Dialect{}, "db/Failures.db")
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.CreateTableIfNotExists(&Failure{}); err != nil {
		log.Fatalln(err)
	}
	dbObj = db

	goji.Get("/", indexPage)
	goji.Post("/", indexPage)

	staticPattern := regexp.MustCompile("^/(css|js)")
	goji.Handle(staticPattern, http.FileServer(http.Dir("./asset")))

	goji.Serve()
}

type ViewData struct {
	Add      bool
	Error    string
	Failures []Failure
	Title    string
	Body     string
}

func indexPage(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "text/html: charset=utf-8")

	// viewデータの作成
	var view ViewData
	if val := r.FormValue("add"); val != "" {
		view.Add = true
	} else if val := r.FormValue("insert"); val != "" {
		title := r.FormValue("title")
		body := r.FormValue("body")
		if view.Error = insertDBIfNeed(title, body); view.Error != "" {
			view.Title = title
			view.Body = body
			view.Add = true
		}
	} else if val := r.FormValue("deleteIdx"); val != "" {
		id, err := strconv.Atoi(val)
		if err != nil {
			log.Printf("delete err id=>%s", val)
		} else {
			obj := Failure{ID: int64(id)}
			if _, err := dbObj.Delete(&obj); err != nil {
				panic(err)
			}
		}
	}

	// データを取得
	if err := dbObj.Select(&view.Failures); err != nil {
		panic(err)
	}
	log.Printf(":%v", view)
	if err := indexTmpl.Execute(w, view); err != nil {
		panic(err)
	}
}

func insertDBIfNeed(title string, body string) string {
	if title == "" || body == "" {
		return "input not enough"
	}
	failture := &Failure{Title: title, Body: body}
	n, err := dbObj.Insert(failture)
	if err != nil {
		return err.Error()
	}
	log.Printf("add %d=>%v ", n, failture)
	return ""
}
