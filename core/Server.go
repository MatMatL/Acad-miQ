package academiq

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

// port definition
const port = ":8080"

// templates definition
var index = template.Must(template.ParseFiles("template/index.html"))
var create = template.Must(template.ParseFiles("template/create.html"))

//var load = template.Must(template.ParseFiles("template/load.html"))

func LaunchServer() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(port, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, nil)
}

type CreateData struct {
	Text   string
	IsText bool
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var data CreateData = CreateData{"", false}

	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20) // 10 Mo limit
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		datasetName := r.FormValue("dataset-name")
		file, handler, _ := r.FormFile("dataset-image")

		if file != nil {
			defer file.Close()
		}

		if datasetName == "" {
			data.Text = "Nom requis"
			data.IsText = true
		} else if NameAlreadyUsed(datasetName) {
			data.Text = "Nom déjà utilisé"
			data.IsText = true
		} else {
			if file != nil {
				filePath := "./uploads/" + handler.Filename
				out, err := os.Create(filePath)
				if err != nil {
					http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
					return
				}
				defer out.Close()

				_, err = io.Copy(out, file)
				if err != nil {
					http.Error(w, "Error occurred while saving the file", http.StatusInternalServerError)
					return
				}

				request := `INSERT INTO Sets (SetName, SetImagePath) VALUES (?, ?);`
				_, err = db.Exec(request, datasetName, filePath)
				if err != nil {
					fmt.Println(err)
				} else {
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
			} else {
				request := `INSERT INTO Sets (SetName) VALUES (?);`
				_, err = db.Exec(request, datasetName)
				if err != nil {
					fmt.Println(err)
				} else {
					http.Redirect(w, r, "/newPost", http.StatusSeeOther)
					return
				}
			}
		}
	}

	create.Execute(w, data)
}
