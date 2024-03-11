package error

import (
	"log"
	"net/http"
	"text/template"
)

type errorData struct {
	Template   string
	StatusText string
	StatusCode int
}

func ErrorHandler(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	data := errorData{
		Template:   "error",
		StatusText: http.StatusText(status),
		StatusCode: status,
	}
	tmpl, err := template.ParseGlob("./templates/*.html")
	if err != nil {
		log.Println("internal server error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
