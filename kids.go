package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const max = 100

type HasMessages interface {
	Add(m string)
	Get() []string
}

type memory struct {
	list []string
}

func (x *memory) Add(m string) {
	x.list = append(x.list, m)
}
func (x *memory) Get() []string {
	return x.list
}

type storage struct {
	List []string
}

type file struct {
	path string
}

func (x *file) Add(m string) {
	list := x.Get()
	list = append(list, m)
	f, err := os.Create(x.path)
	check(err)
	defer f.Close()
	e := json.NewEncoder(f)
	check(e.Encode(storage{list}))
}

func (x *file) Get() []string {
	f, err := os.Open(x.path)
	if err != nil {
		return nil
	}
	defer f.Close()
	d := json.NewDecoder(f)
	var m storage
	check(d.Decode(&m))
	return m.List
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (h Website) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	msg := r.URL.Query().Get("message")
	if msg != "" {
		log.Printf("we got a message: %q\n", msg)
		h.Add(msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="initial-scale=1.0,width=device-width,user-scalable=no">
    <title>%s</title> 
  </head>
  <body>
<p>%s</p>
<form>
<input style='padding:0.5em;width:90%%;' autocorrect='off' autocapitalize='none' autocomplete='off' type='text' name='message' autofocus/>
</form>
`, h.Title, h.Greeting)

	if len(h.Links) > 0 {
		fmt.Fprintln(w, "<p>here's some of my favorite links:</p>")

		fmt.Fprintln(w, "<ol>")
		for _, link := range h.Links {
			fmt.Fprintf(w, "<li>see <a target='_blank' href='%s'>%s</a></li>", link, link)
		}
		fmt.Fprintln(w, "</ol>")
	}

	fmt.Fprintf(w, `<br/>
<img src='%s' height='200px'/>
`, h.Image)

	list := h.Get()
	n := len(list)
	if n > 0 {
		fmt.Fprintln(w, "<p> here's some messages left in the past:</p> <ol>")
		if n > max {
			n = max
		}
		for i := 0; i < n; i++ {
			m := list[len(list)-i-1]
			fmt.Fprintf(w, "<li> %s </li>", m)
		}
		fmt.Fprintln(w, "</ol>")
	}
	fmt.Fprintf(w, `
</body>
</html>
`)

}

type Website struct {
	Title    string
	Greeting string
	Image    string
	Links    []string
	HasMessages
}

func main() {

	david := Website{
		Title:       "david's website",
		Greeting:    "Howdy, this is David's website! If you have any thing to tell me, just tipe the following bellow. Have a nice day!!",
		Image:       "http://www.shawnkinley.com/wp-content/uploads/have-a-nice-day.jpg",
		HasMessages: &file{"david.json"},
		Links: []string{
			"http://www.totaljerkface.com/happy_wheels.tjf",
			"http://www.sheppardsoftware.com/",
		},
	}

	sylvia := Website{
		Title:       "sylvia's website",
		Greeting:    "Hi!This is my website!I like to play piano!I love to Program!!!",
		Image:       "http://d2jngs55a0uns9.cloudfront.net/ad853ef37ed136a7ee5686c5b0c00cd0bf370f14_NONE_388153.jpg",
		HasMessages: &file{"sylvia.json"},
		Links: []string{
			"http://www.poptropica.com/",
			"http://www.abcya.com/",
		},
	}

	kidFunction := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		log.Printf("%s %s\n", r.RemoteAddr, r.RequestURI)

		switch r.Host {
		case "gosylvia.ch":
			sylvia.ServeHTTP(w, r)
		case "godavid.ch":
			david.ServeHTTP(w, r)
		}
	}

	s := &http.Server{
		Handler: http.HandlerFunc(kidFunction),
	}

	log.Fatal(s.ListenAndServe())
}
