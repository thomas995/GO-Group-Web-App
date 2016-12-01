//https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/
// Adapted from: https://github.com/jakecoffman/go-angular-tutorial
// Adapted from: https://mschoebel.info/2014/03/09/snippet-golang-webapp-login-logout/
// Info gathered from w3Schools and https://golang.org
package main

import (
   //need to type 'go get database/sql' into cmder to use
	"database/sql" // This allows us to use the SQL Database
	"log"
	"github.com/gorilla/securecookie" //Using gorilla to handle cookies
	"net/http" // http allows HTTP client & server use/implementations.
	"net/smtp" //- For sending emails. Adapted from https://dinosaurscode.xyz/go/2016/06/21/sending-email-using-golang/
	"github.com/zenazn/goji/web" //Allows you to parse information from the html page into the email
//need to type ' go get github.com/go-sql-driver/mysql' into cmder to use
	_ "github.com/go-sql-driver/mysql"  // This installs a driver in order to be able to use mySQL
	//need to type 'go get golang.org/x/crypto/bcrypt' into cmder to use
	"golang.org/x/crypto/bcrypt" // This import statement allows encryption and decryption of passwords
)

//A variable which initialises the use of sql
var db *sql.DB

// A variable to handle errors
var err error

// cookies are handled here

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

//Reading cookies for username
func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["username"]
		}
	}
	return userName
}

//Saves username in map then encodes with value map and stores that in a cookie
func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
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

//We start off with a Signup function so the user can Signup
// We initialise the use of the response (res) using the ResponseWriter
// and the request (req) by using http.Request
func signupPage(res http.ResponseWriter, req *http.Request) {
	//If the request is not a POST method (the POST request method requests that a web server accept and store
	// the data enclosed in the body of the request message. It is used when submitting a completed web Form)
	if req.Method != "POST" {
		//then serve the SignUp page
		http.ServeFile(res, req, "signup.html")
		return
	}

	//The username and password are set up as a formValue
	username := req.FormValue("username")
	password := req.FormValue("password")

	//the user is initialised as a string
	var user string

	// Query the database to see if the user signing in is clashing credentials
	//with any existing users within the database
	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	//a switch statement to handle the signin
	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		//If the error list is not equal to null, in other words if there is (are) error(s)
		if err != nil {
			//There will be a server error 500 and your account will not be created
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		// The "_" (blank identifier) avoids having to declare all the variables for the return values.
		// Try inserting a username and password into the users table and encrypt the password
		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		//If there are errors
		if err != nil {
			// The response (res) will say Server error 500 and you can't create your account
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		// Otherwise A user is created
		res.Write([]byte("User created!"))
		return
		//checking for additional server errors(server not available, table deleated etc)
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
		// by default the root page (/index.html) is shown to the user
		// error 301 checks if the page has been moved permanantly
	default:
		http.Redirect(res, req, "/", 301)
	}
}

// This is a function to allow a user to log into the website
func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	//This time we have a Username variable and a Password variable
	var databaseUsername string
	var databasePassword string

	//Scan the database to find the returning user as they login with their username and password
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
	//If there are errors with the user logging in
	if err != nil {
		//Error 301 (page moved permanantly)
		http.Redirect(res, req, "/login", 301)
		return
	}
	//comparing the password the user is entering with the encrypted one in the database
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	//If there are no errors then it redirects to the login page, otherwise if the page has been moved - 301 error
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}

	//If all is successful, serve up the internal.html file which only a logged-in user should be able to see
	http.ServeFile(res, req, "Internal.html")
	//res.Write([]byte("Hello" + databaseUsername))
}

//A function to serve up the homepage called index.html which includes the login and signup buttons
func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

//- A function that sends an email from workforcegroupproject@gmail.com to an employee's employer.
//- The body of the message will be decided by what hours the employee logs and if they put anything in the information box.
//func sendEmail(body string, to string)
func sendEmail(c web.C, w http.ResponseWriter, r *http.Request) {
	//This is the email address we created ourselves
	from := "workforcegroupproject@gmail.com"
	//This is the password associated with that email we created
	password := "GoGroup2016!"
	//- For now emails / form submission from the website will be sent to the email workforcegroupproject@gmail.com.
	to := "workforcegroupproject@gmail.com"

//Parses the data from html into the form
	err := r.ParseForm()
	if err != nil {
		// Handle error here via logging and then return
		log.Printf("Error: %s", err)
		return
	}

//The calculation for hoursWorked * wage was output to the label lblRes in a previous function
//So we are going to take this label and the "additional information textbox content and send it
//to the 'employer' who hypothetically has the email "workforcegroupproject@gmail.com"
	body := r.PostFormValue("lblRes" + "AddInfo")
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
//This is the setup for the email, so you've got "From" , "To" and "body" (which are declared above)
// and we are going to use these to send the 'msg' (email)
	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: WorkForce hours worked" + "\r\n\r\n" +
		body + "\r\n"

//This handles the errors in sending the email and uses go's 'smtp.SendMail' function to actually send the email
//"smtp.gmail.com:587" accesses the gmail server
	sendError := smtp.SendMail("smtp.gmail.com:587" , auth, from, []string{to}, []byte(msg))
	if sendError != nil {
		log.Printf("Error: %s", sendError)
		return
	}

	log.Print("message sent")

}

//Main function
func main() {

	//Open an sql connection, and handle errors,
	// The database is mysql, then we have the username of the server on Azure (b71da173aea4cf)
	// (05606ea1)
	//TCP stands for Transmission Control Protocol which is a set of networking protocols that allows
	// two or more computers to communicate
	// (eu-cdbr-azure-west-a.cloudapp.net)
	// Port Name on Azure (3306)
	// The name of the database (godatabase)
	db, err = sql.Open("mysql", "b71da173aea4cf:05606ea1@tcp(eu-cdbr-azure-west-a.cloudapp.net:3306)/godatabase")
	//If there are errors connecting to the server
	if err != nil {
		//output a panic error
		panic(err.Error())
	}
	//close the connection to the database
	defer db.Close()

	//Checking is the connection to the database still alive
	err = db.Ping()
	//otherwise painic
	if err != nil {
		panic(err.Error())
	}

	//Handle all of our functions
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)

	//serve on the port 8000 forever
	http.ListenAndServe(":8000", nil)
}
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
<<<<<<< HEAD
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">

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
  <script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.7.1/jquery.min.js"></script>
<!--Adapted from: https://www.sitepoint.com/ -->
  <script language="javascript">
		function Calculate()
		{
			var h = document.getElementById('hoursWorked').value;
			var t = document.getElementById('hourlyPay').value;
			var result = h * t;
			document.getElementById('lblRes').innerHTML = result;
		}
	</script>
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
 <!-- Code adapted from http://stackoverflow.com/questions/30717105/adding-view-details-button-with-js -->
  </select>
          <details>
    <summary>View Details</summary>
    <p>
        Please select your correct job title here from the dropdown list.
    </p>
</details>
</div>

<!-- Code adapted from http://www.w3schools.com/html/html_form_elements.asp -->
  <div class="col-md-4">
     <h2>Hours Worked This Week</h2>
  <select Id="hoursWorked" onChange="Calculate();">
			  <option value=10>10</option>
			  <option value=11>11</option>
			  <option value=12>12</option>
			  <option value=13>13</option>
			  <option value=14>14</option>
			  <option value=15>15</option>
			  <option value=16>16</option>
			  <option value=17>17</option>
			  <option value=18>18</option>
			  <option value=19>19</option>
			  <option value=20>20</option>
			  <option value=21>21</option>
			  <option value=22>22</option>
			  <option value=23>23</option>
			  <option value=24>24</option>
			  <option value=25>25</option>
			  <option value=26>26</option>
			  <option value=27>27</option>
			  <option value=28>28</option>
			  <option value=29>29</option>
			  <option value=30>30</option>
			  <option value=31>31</option>
			  <option value=32>32</option>
			  <option value=33>33</option>
			  <option value=34>34</option>
			  <option value=35>35</option>
			  <option value=36>36</option>
			  <option value=37>37</option>
			  <option value=38>38</option>
			  <option value=39>39</option>
			  <option value=40>40</option>

			</select>
			  &nbsp;&nbsp;
			* &nbsp;
			Hourly Pay:
			<select Id="hourlyPay" onChange="Calculate();">
			  <option value=6.24>6.24</option>
			  <option value=7.25>7.25</option>
			  <option value=9.15>9.15</option>
			  <option value=11.25>11.25</option>
			  <option value=20.83>20.83</option>
			</select>
			&nbsp;&nbsp;
			
			= Total:&nbsp;&nbsp;
			<label id="lblRes"> <!-- result of Calculating hoursWorked * hourlyPay -->
			 100
			</label>
          <details>
    <summary>View Details</summary>
    <p>
        Please select the correct amount of hours you worked this week. 
        <br> Please note: In order to use this site you have to be working a minimum of 10 hours a week and a maximum of 40 hours. </br>
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
          <details>
    <summary>View Details</summary>
    <p>
        Please select the correct amount of hours you worked this week.
    </p>
</details>
       </div>

        <div class="col-md-4">
          <h2>Additonal Information</h2>
          <!--Adapted from http://www.w3schools.com/tags/tag_input.asp-->
          <form action="demo_form.asp">
              <input type="text" name="AddInfo" value=""><br>
              <input type="submit" value="Submit">
</form>
<details>
    <summary>View Details</summary>
    <p>
        Please add any additional information that you feel is relevent for your boss to know about regarding this week's work.
    </p>
</details>
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
