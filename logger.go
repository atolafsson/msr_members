package main

import (
	"log"
	"net/http"
	"time"
)

// Logger -- Log text with current timestamp
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		cUser := "n/a"
		//uCookie, _ := r.Cookie("username")
		//aCookie, _ := r.Cookie("uaccess")
		ok := true
		/*if uCookie != nil && aCookie != nil {
			cUser = uCookie.Value
		} else if r.RequestURI != "/login" {
			ok = false
		}*/
		if ok {
			inner.ServeHTTP(w, r)
			log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, cUser, time.Since(start))
		} /*else {
			log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, cUser, "***OOPS needs to log in")
			http.Redirect(w, r, "/static/Login.html", http.StatusFound)
		}*/
	})
}
