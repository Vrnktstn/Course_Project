package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type PigeonList struct {
	Id    int
	Name  string
	Breed string
	Cute  int
}

var database *sql.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from pigeon_list_db.pigeons")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	pigeonlist := []PigeonList{}

	for rows.Next() {
		pl := PigeonList{}
		err := rows.Scan(&pl.Id, &pl.Name, &pl.Breed, &pl.Cute)
		if err != nil {
			fmt.Println(err)
			continue
		}
		pigeonlist = append(pigeonlist, pl)
	}

	tmpl, _ := template.ParseFiles("templates/home.html")
	tmpl.Execute(w, pigeonlist)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		name := r.FormValue("name")
		breed := r.FormValue("breed")
		cute := r.FormValue("cute")

		_, err = database.Exec("insert into pigeon_list_db.pigeons (name, breed, cute) values (?, ?, ?)",
			name, breed, cute)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/add.html")
	}
}

func editPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.QueryRow("select * from pigeon_list_db.pigeons where id = ?", id)
	pl := PigeonList{}
	err := row.Scan(&pl.Id, &pl.Name, &pl.Breed, &pl.Cute)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/edit.html")
		tmpl.Execute(w, pl)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	name := r.FormValue("name")
	breed := r.FormValue("breed")
	cute := r.FormValue("cute")

	_, err = database.Exec("update pigeon_list_db.pigeons set name=?, breed=?, cute = ? where id = ?",
		name, breed, cute, id)

	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

func viewPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.QueryRow("select * from pigeon_list_db.pigeons where id = ?", id)
	pl := PigeonList{}
	err := row.Scan(&pl.Id, &pl.Name, &pl.Breed, &pl.Cute)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/view.html")
		tmpl.Execute(w, pl)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec("delete from pigeon_list_db.pigeons where id = ?", id)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

// func upload(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Println("method:", r.Method)
// 	if r.Method == "GET" {
// 		crutime := time.Now().Unix()
// 		h := md5.New()
// 		io.WriteString(h, strconv.FormatInt(crutime, 10))
// 		token := fmt.Sprintf("%x", h.Sum(nil))

// 		t, _ := template.ParseFiles("templates/upload.html")
// 		t.Execute(w, token)
// 	} else {
// 		r.ParseMultipartForm(32 << 20)
// 		file, handler, err := r.FormFile("uploadfile")
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		defer file.Close()
// 		fmt.Fprintf(w, "%v", handler.Header)
// 		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		defer f.Close()
// 		io.Copy(f, file)
// 	}
// }

func main() {

	db, err := sql.Open("mysql", "root:123456@/pigeon_list_db")

	if err != nil {
		log.Println(err)
	}

	database = db
	defer db.Close()

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/add", addHandler)
	r.HandleFunc("/edit/{id:[0-9]+}", editPage).Methods("GET")
	r.HandleFunc("/edit/{id:[0-9]+}", editHandler).Methods("POST")
	r.HandleFunc("/view/{id:[0-9]+}", viewPage).Methods("GET")
	r.HandleFunc("/delete/{id:[0-9]+}", deleteHandler)
	// r.HandleFunc("/upload", upload)

	http.Handle("/", r)

	http.ListenAndServe(":9000", nil)
}
