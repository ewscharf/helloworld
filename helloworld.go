package main

import (
	"html/template"
	_ "image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type dados struct {
	Name   string
	Foto   multipart.File
	Handle *multipart.FileHeader
	Src    string
}

var tpl *template.Template

func main() {
	tpl, _ = tpl.ParseGlob("tpl/*.html")
	http.HandleFunc("/hello", hello)
	http.Handle("/fotos/", http.StripPrefix("/fotos/", http.FileServer(http.Dir("./fotos"))))
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "form.html", nil)
	case "POST":

		r.ParseMultipartForm(10 << 20)
		var info dados
		info.Name = r.FormValue("name")
		info.Foto, info.Handle, _ = r.FormFile("foto")
		defer info.Foto.Close()
		ioutil.TempFile("images", "*.jpg")
		if info.Handle.Header["Content-Type"][0] == "image/jpeg" {
			newFoto, _ := ioutil.TempFile("fotos", "img*.jpg")
			defer newFoto.Close()
			defer os.Remove(newFoto.Name())
			fileBytes, _ := ioutil.ReadAll(info.Foto)
			newFoto.Write(fileBytes)
			info.Src = filepath.ToSlash(newFoto.Name())

		}

		tpl.ExecuteTemplate(w, "done.html", info)
	}

}
