package main

// Adapted from: https://github.com/jakecoffman/go-angular-tutorial
// Adapted from: https://mschoebel.info/2014/03/09/snippet-golang-webapp-login-logout/
// Comments (Apart from the last one) are my own, info gathered from w3Schools and https://golang.org
import (
	//"flag"
	"fmt"
	"github.com/gorilla/mux" //Using Gorilla instead of macaroon
	"github.com/gorilla/securecookie" //For using session cookies
	//"log" // It defines a type, Logger, with methods for formatting output. 
	"net/http"
)

// cookies are handled here

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
//Reading cookies for username
func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}
//Saves username in map then encodes with value map and stores that in a cookie
func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}
//Returns to indexPage and clears cookies
func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler handles the login for stored users

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
		setSession(name, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler logs the current user out

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// index page contains code for login to application

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// internal page contains code for the page user will see once successfully logged in

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`
//Handles moving to internal page and bringing along userName
func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// server main method which handles all the port hosting and http methods

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	//port := flag.Int("port", 8080, "port to serve on")
	//dir := flag.String("directory", "web/", "directory of web files")
	//flag.Parse()

	//fs := http.Dir(*dir) // we have specified the root directory as the default: web
	//fileHandler := http.FileServer(fs) //FileServer is a built-in to GO to allow serving of the html file(s)
	//http.Handle("/", fileHandler) //Url routing, Gets file from path specified,
	

	//log.Printf("Running on port %d\n", *port) //Logger accessible through helper function Print[f|ln]
										    //That logger writes to standard error and prints the date and time of each logged message.

	//addr := fmt.Sprintf("127.0.0.1:%d", *port)
	
	http.Handle("/", router)
	// this call blocks -- the progam runs here forever
	//err := http.ListenAndServe(addr, nil)
	http.ListenAndServe(":8000", nil)
	//fmt.Println(err.Error())
}
