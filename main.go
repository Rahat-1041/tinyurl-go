package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

var urls = map[string]string{}
var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<h2>Tiny URL Shortener</h2>
		<form action="/shorten" method="POST">
			<input type="text" name="url" placeholder="Enter long URL" size="40">
			<button type="submit">Shorten</button>
		</form>`)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	longURL := r.FormValue("url")
	if !strings.HasPrefix(longURL, "http") {
		longURL = "https://" + longURL
	}

	code := genCode()
	urls[code] = longURL

	shortURL := "http://" + r.Host + "/" + code
	fmt.Fprintf(w, `<p>Short URL: <a href="%s">%s</a></p><a href="/">Back</a>`, shortURL, shortURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	longURL, ok := urls[code]
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			homeHandler(w, r)
		} else {
			redirectHandler(w, r)
		}
	})
	http.HandleFunc("/shorten", shortenHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
