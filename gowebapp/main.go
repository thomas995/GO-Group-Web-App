//https://dinosaurscode.xyz/go/2016/06/19/golang-mysql-authentication/
// Adapted from: https://github.com/jakecoffman/go-angular-tutorial
// Adapted from: https://mschoebel.info/2014/03/09/snippet-golang-webapp-login-logout/
// Info gathered from w3Schools and https://golang.org
package main

//- For sending emails. Adapted from https://github.com/golang/go/wiki/SendingMail

//- For sending emails. Adapted from https://dinosaurscode.xyz/go/2016/06/21/sending-email-using-golang/

//Using gorilla to handle cookies
import (
	"database/sql"
	"log"

	"github.com/gorilla/securecookie"

	// This allows us to use the SQL Database
	//need to type 'database/sql' into cmder to use
	// This installs a driver in order to be able to use mySQL
	"net/http"
	"net/smtp"

	"github.com/zenazn/goji/web"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//need to type 'github.com/go-sql-driver/mysql' into cmder to use
// This import statement allows encryption and decryption of passwords
//need to type 'golang.org/x/crypto/bcrypt' into cmder to use

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

//- A function that sends an email from workforcegroupproject@gmail.com to an employees employer.
//- The body of the message will be decided by what hours the employee logs and if they put anything in the information box.
//func sendEmail(body string, to string)
func sendEmail(c web.C, w http.ResponseWriter, r *http.Request) {
	from := "workforcegroupproject@gmail.com"
	password := "GoGroup2016!"
	//- For now emails will be sent to the project email.
	to := "workforcegroupproject@gmail.com"

	err := r.ParseForm()
	if err != nil {
		// Handle error here via logging and then return
		log.Printf("Error: %s", err)
		return
	}

	body := r.PostFormValue("lblRes" + "AddInfo")

	//Set body to equal the salary due and additional info.
	//body := req.FormValue("AddInfo")

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: WorkForce hours worked" + "\r\n\r\n" +
		body + "\r\n"

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Error: %s", err)
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
