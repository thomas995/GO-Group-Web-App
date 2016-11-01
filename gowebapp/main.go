package main

// Adapted from: https://github.com/jakecoffman/go-angular-tutorial
// Comments (Apart from the last one) are my own, info gathered from w3Schools and https://golang.org
import (
	"flag"
	"fmt"
	"log" // It defines a type, Logger, with methods for formatting output. 
	"net/http"
)

func main() {
	port := flag.Int("port", 80, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	fs := http.Dir(*dir) // we have specified the root directory as the default: web
	fileHandler := http.FileServer(fs) //FileServer is a built-in to GO to allow serving of the html file(s)
	http.Handle("/", fileHandler) //Url routing, Gets file from path specified,
	

	log.Printf("Running on port %d\n", *port) //Logger accessible through helper function Print[f|ln]
										    //That logger writes to standard error and prints the date and time of each logged message.

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
