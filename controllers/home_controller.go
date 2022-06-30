package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

var MobileDevices = [...]string{"Mobile Explorer", "Palm", "Motorola", "Nokia", "Palm", "Apple iPhone", "iPad", "Apple iPod Touch", "Sony Ericsson", "Sony Ericsson", "BlackBerry", "O2 Cocoon", "Treo", "LG", "Amoi", "XDA", "MDA", "Vario", "HTC", "Samsung", "Sharp", "Siemens", "Alcatel", "BenQ", "HP iPaq", "Motorola", "PlayStation Portable", "PlayStation 3", "PlayStation Vita", "Danger Hiptop", "NEC", "Panasonic", "Philips", "Sagem", "Sanyo", "SPV", "ZTE", "Sendo", "Nintendo DSi", "Nintendo DS", "Nintendo 3DS", "Nintendo Wii", "Open Web", "OpenWeb", "Android", "Symbian", "SymbianOS", "Palm", "Symbian S60", "Windows CE", "Obigo", "Netfront Browser", "Openwave Browser", "Mobile Explorer", "Opera Mini", "Opera Mobile", "Firefox Mobile", "Digital Paths", "AvantGo", "Xiino", "Novarra Transcoder", "Vodafone", "NTT DoCoMo", "O2", "mobile", "wireless", "j2me", "midp", "cldc", "up.link", "up.browser", "smartphone", "cellphone", "Generic Mobile"}

// RenderWithLayout renders a layout with an additional page layout
func RenderWithLayout(pageTemplate string, w http.ResponseWriter, data interface{}) error {
	var err error
	mainlayout, err := os.Open("views/main_layout.html")
	if err != nil {
		return errors.Wrap(err, "opening the main layout file")
	}
	page, err := os.Open("views/" + pageTemplate + ".html")
	if err != nil {
		return errors.Wrap(err, "opening the page layout file")
	}

	var tmpl *template.Template
	layoutContent, err := ioutil.ReadAll(mainlayout)
	if err != nil {
		return errors.Wrap(err, "reading the main layout file")
	}
	tmpl = template.New(pageTemplate)
	tmpl = tmpl.Delims("[[", "]]")
	tmpl.Parse(string(layoutContent))
	if err != nil {
		panic(err)
	}

	if pageTemplate != "" {
		pageContent, err := ioutil.ReadAll(page)
		if err != nil {
			return errors.Wrap(err, "reading the page layout file")
		}
		tmpl, err = tmpl.Parse(string(pageContent))
		if err != nil {
			return errors.Wrap(err, "parsing the page layout")
		}
	}

	//w.WriteHeader(http.StatusOK)
	tmpl.ExecuteTemplate(w, "main_layout", data)

	return nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	viewData := struct {
		IsMobileAgent bool
	}{}
	viewData.IsMobileAgent = detectMobile(r)

	err := RenderWithLayout("home", w, viewData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// detectMobile returns true if the user agent (aka the browser) is a mobile one.
// The list of devices was taken from:
// https://github.com/bcit-ci/CodeIgniter/blob/develop/system/libraries/User_agent.php
func detectMobile(r *http.Request) bool {
	agent := r.Header.Get("User-Agent")

	for _, device := range MobileDevices {
		if strings.Contains(agent, device) {
			return true
		}
	}
	return false
}
