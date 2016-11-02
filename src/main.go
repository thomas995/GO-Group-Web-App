package main

 import (
         //to use mux
         "net/http"
         "fmt"
          "html/template"
 )

//to use Macaron
 import "gopkg.in/macaron.v1"

//Create a page struct with a public string property called name
type Page struct {
        Name string
}

 func main() {
         //Template.Must absorbs the error from ParseFiles and
         // stop the program from running if it cannot parse the template
         templates := template.Must(template.ParseFiles("../public/Index.html", "../Template/bootstrap-3.3.7-dist/dist/css/bootstrap.min.css"))
        
         //HandleFunc is from standard HTTP Library and takes two arguments
         http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
                 //create a new instane of page with the name Worker
        p:= Page{Name: "Worker"}
        //FormValue gives us access to query strings
        if name := r.FormValue("name"); name !="" {
                //update form with value contained in query parameter (?name="Your Name")
                p.Name = name
        }
        //Throws InternalServerError50 if template isn't found
        //pass page(p) as the third parameter to the ExecuteTemplate method
        if err := templates.ExecuteTemplate(w, "Index.html", p); err!= nil{
                http.Error(w, err.Error(), http.StatusInternalServerError)
                 }
        })


        m := macaron.Classic()
        m.Use(macaron.Renderer())
        //Wrapping in a fmt.println allows you to get info about errors such as port 8080 not being open
        fmt.Println(http.ListenAndServe(":8080", nil))
 }