package controllers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func ImageUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpFile, status, err := ReceiveFile(w, r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	fmt.Printf("tmpFile = %+v\n", tmpFile)
	//defer os.Remove(tmpFile)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusTemporaryRedirect)
}

func ReceiveFile(w http.ResponseWriter, r *http.Request) (string, int, error) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	tmpFile, err := ioutil.TempFile("", "oor-upload")
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(err, "creating a tmp file")
	}

	// in your case file would be fileupload
	file, _, err := r.FormFile("image-file")
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(err, "reading the uploaded image")
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(err, "reading the file type")
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" {
		return "", http.StatusBadRequest, errors.New("The provided file format is not allowed. Please upload a JPEG image")
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", http.StatusInternalServerError, errors.Wrap(err, "moving file cursor")
	}

	// Copy the file data to my buffer
	io.Copy(tmpFile, file)

	return tmpFile.Name(), http.StatusOK, nil
}
