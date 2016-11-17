package main

// Adapted from: https://github.com/jakecoffman/go-angular-tutorial
// Adapted from: https://mschoebel.info/2014/03/09/snippet-golang-webapp-login-logout/
// Comments (Apart from the last one) are my own, info gathered from w3Schools and https://golang.org
import (
	//"flag"
	"fmt"
	//"log" // It defines a type, Logger, with methods for formatting output. 
	"net/http"
	"github.com/gorilla/mux" //need to type "go get github.com/gorilla/mux" into cmder to use (without the quotations obviously)
	"github.com/gorilla/securecookie" //need to type "github.com/gorilla/securecookie" into cmder to use (without the quotations obviously)
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
			userName = cookieValue["email"]
		}
	}
	return userName
}

//Saves username in map then encodes with value map and stores that in a cookie
func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"email": userName,
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
	email := request.FormValue("email")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if email != "" && pass != "" {
		// .. check credentials ..
		setSession(email, response)
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
<!-- Incorporating some HTML -->
<head>
<!-- Nav bar -->
<form class="navbar-form navbar-left">
	<title>WorkTracker</title>
	</form>
	<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
	<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<div class="container">
		
	</div>

	<nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="#">Work Tracker</a>
        </div>
        <div id="navbar" class="navbar-collapse collapse">

        <form method="post" action="/login">
    <input type="email" placeholder="Enter your email" id="email" name="email">
    <input type="password" placeholder="Password" id="password" name="password">
    <button type="submit" class="btn btn-success">Login</button>
</form>    
</div><!--/.navbar-collapse -->
</div>
</nav>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// internal page contains code for the page user will see once successfully logged in

const internalPage = `
<h1>-</h1>
<hr>
<small>User: %s</small>
<html ng-app> <!-- 'ng-app'' placed within a tag (in this case, the HTML tag)
				allows HTML to become the route element for AngularJS.
				All AngularJS applications must have a root element. Only one instance allowed. --> 
<head>
	<title>WorkTracker</title>
	<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
	<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>

<body>
	<!-- Main jumbotron for a primary marketing message or call to action -->
    <div class="jumbotron">
      <!--Changed the top header here to include a name being entered and shown back to the user with a friendly hello-->
      <div class="container">
        <h2><input type="text" placeholder="Your Name here" ng-model="yes"></h2>
		      <h3><p>Hello <span ng-bind="yes"></span></p></h3>
        <p>This is a place where you can record all of your daily duties in one place, ready to show the boss. </p>
      </div>
    </div>

	<div class="container">
		
	</div>

	<nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
                     <form method="post" action="/logout">
    <button type="submit" class="btn btn-success">Logout</button>
          </button>
          <a class="navbar-brand" href="#">Work Tracker</a>
        </div>
        <div id="navbar" class="navbar-collapse collapse">
          <form class="navbar-form navbar-right">
  
</form>
          </form>
        </div><!--/.navbar-collapse -->
      </div>
    </nav>

	<div class="container">
      <!-- Example row of columns -->
      <div class="row">
        <div class="col-md-4">
          <h2>Job Title</h2>
          <!-- Dropdown list with Tradesperson's job titles
  adapted from https://www.socketloop.com/tutorials/golang-populate-dropdown-with-html-template-example -->
          <select>
    <option value="Beautician">Beautician</option>
    <option value="Builder">Builder</option>
    <option value="Carpenter">Carpenter</option>
    <option value="Cleaner">Cleaner</option>
    <option value="Delivery">Delivery</option>
    <option value="Electrician">Electrician</option>
    <option value="Farmer">Farmer</option>
    <option value="Gardener">Gardener</option>
    <option value="Hair and Beauty">Hair and Beauty</option>
    <option value="Hairdresser">Hairdresser</option>
    <option value="Mechanic">Mechanic</option>
    <option value="Painter">Painter</option>
    <option value="Plumber">Plumber</option>
    <option value="Technician">Technician</option>
    <option value="Tiler">Tiler</option>
    <option value="Transport">Transport</option>
    <option value="Other">Other</option>
 
  </select>
          <p><a class="btn btn-default" href="#" role="button">View details &raquo;</a></p>
        </div>
<!-- Code adapted from http://www.w3schools.com/html/html_form_elements.asp -->
        <div class="col-md-4">
          <h2>Hours Worked This Week</h2>
          <select name="HoursWorked">
  <option value=">10 hours">>10 Hours</option>
  <option value="10-15">10-15 Hours</option>
  <option value="15-20">15-20 Hours</option>
  <option value="25-30">25-30 Hours</option>
  <option value="30-40">30-40 Hours</option>
  <option value="40+">40+ Hours</option>

</select>
          <p><a class="btn btn-default" href="#" role="button">View details &raquo;</a></p>
       </div>
       
        <div class="col-md-4">
          <h2>Additonal Information</h2>
          <!--Adapted from http://www.w3schools.com/tags/tag_input.asp-->
          <form action="demo_form.asp">
              <input type="text" name="AddInfo" value=""><br>
              <input type="submit" value="Submit">
</form><p><a class="btn btn-default" href="#" role="button">View details &raquo;</a></p>
        </div>
      </div>

	  <hr>

      <footer>
        <p>&copy; 2016 WorkTracker, Inc.</p>
      </footer>
    </div> <!-- /container -->

<script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.2.3/angular.min.js"></script>
</body>
</html>

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
//	dir := flag.String("directory", "web/", "directory of web files")
	//flag.Parse()

//	fs := http.Dir(*dir) // we have specified the root directory as the default: web
//	fileHandler := http.FileServer(fs) //FileServer is a built-in to GO to allow serving of the html file(s)
//	http.Handle("/", fileHandler) //Url routing, Gets file from path specified.

	//log.Printf("Running on port %d\n", *port) //Logger accessible through helper function Print[f|ln]
										    //That logger writes to standard error and prints the date and time of each logged message.
	//addr := fmt.Sprintf("127.0.0.1:%d", *port)

	http.Handle("/", router)
	// this call blocks -- the progam runs here forever
	//err := http.ListenAndServe(addr, nil)
	http.ListenAndServe(":8000", nil)
	//fmt.Println(err.Error())
}
