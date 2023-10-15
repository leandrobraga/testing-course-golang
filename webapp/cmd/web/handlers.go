package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/leandrobraga/testing-course-golang/webapp/pkg/data"
)

var pathToTemplates = "./templates/"

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]any)

	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.Get(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "Hit this page at "+time.Now().UTC().String())
	}
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
}

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {

	_ = app.render(w, r, "profile.page.gohtml", &TemplateData{})
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !app.authenticate(r, user, password) {
		app.Session.Put(r.Context(), "error", "Invalid login!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// prevent fixation atack
	_ = app.Session.RenewToken(r.Context())

	// redirect to some other page
	app.Session.Put(r.Context(), "flash", "Successfuly logged in!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *application) UploadProfilePic(w http.ResponseWriter, r *http.Request) {
	files, err := app.UploadFiles(r, "./static/img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := app.Session.Get(r.Context(), "user").(data.User)

	var i = data.UserImage{
		UserID:   user.ID,
		FileName: files[0].OriginalFileName,
	}

	_, err = app.DB.InsertUserImage(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateUser, err := app.DB.GetUser(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Session.Put(r.Context(), "user", updateUser)

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)

}

type TemplateData struct {
	IP    string
	Data  map[string]any
	Error string
	Flash string
	User  data.User
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	parsedFiles, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(r.Context())
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	err = parsedFiles.Execute(w, td)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) authenticate(r *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}
	app.Session.Put(r.Context(), "user", user)
	return true
}

type UploadedFile struct {
	OriginalFileName string
	FileSize         int64
}

func (app *application) UploadFiles(r *http.Request, uploadDir string) ([]*UploadedFile, error) {
	var uploadedFiles []*UploadedFile

	err := r.ParseMultipartForm(int64(1024 * 1024 * 5))
	if err != nil {
		return nil, fmt.Errorf("the uploaded file is too big, and must be less than %d bytes", 1024*1024*5)
	}

	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				uploadedFile.OriginalFileName = hdr.Filename

				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.OriginalFileName)); nil != err {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}

				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil
			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}
		}
	}

	return uploadedFiles, nil
}
