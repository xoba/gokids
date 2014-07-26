package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const max = 100

type Messages interface {
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

func main() {

	maker := func(messages Messages, title, greeting, image string) func(w http.ResponseWriter, r *http.Request) {

		return func(w http.ResponseWriter, r *http.Request) {

			msg := r.URL.Query().Get("message")
			if msg != "" {
				log.Printf("we got a message: %q\n", msg)
				messages.Add(msg)
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
<br/>
<img src='%s' height='200px'/>
<ol>`, title, greeting, image)

			list := messages.Get()
			n := len(list)
			if n > max {
				n = max
			}
			for i := 0; i < n; i++ {
				m := list[len(list)-i-1]
				fmt.Fprintf(w, "<li> %s </li>", m)
			}
			fmt.Fprintf(w, `
</ol></body>
</html>
`)
		}
	}

	david := maker(&file{"david.json"}, "david's website", "howdy, this is david's website! if you have any thing to tell me, just tipe the following bellow. have a nice day!!", "http://www.shawnkinley.com/wp-content/uploads/have-a-nice-day.jpg")

	sylvia := maker(&file{"sylvia.json"}, "sylvia's website", "Hi!This is my website!I like to play piano!", "http://d2jngs55a0uns9.cloudfront.net/ad853ef37ed136a7ee5686c5b0c00cd0bf370f14_NONE_388153.jpg")

	kidFunction := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		log.Printf("%s %s\n", r.RemoteAddr, r.RequestURI)

		switch r.Host {
		case "gosylvia.ch":
			sylvia(w, r)
		case "godavid.ch":
			david(w, r)
		}
	}

	s := &http.Server{
		Handler: http.HandlerFunc(kidFunction),
	}

	log.Fatal(s.ListenAndServe())
}
