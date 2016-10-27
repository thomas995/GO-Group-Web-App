package main

 import (
         "net/http"
 )

 import "gopkg.in/macaron.v1"

/* Dropdown list with Tradesperson's job titles
 adapted from https://www.socketloop.com/tutorials/golang-populate-dropdown-with-html-template-example */
 func SimpleSelectTag(w http.ResponseWriter, r *http.Request) {
         html := `<!DOCTYPE html>
 <html>
 <body>
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
 </body>
 </html>`

         w.Write([]byte(html))
 }

 func main() {
         mux := http.NewServeMux()
         mux.HandleFunc("/", SimpleSelectTag)

        m := macaron.Classic()
        m.Use(macaron.Renderer())
        http.ListenAndServe(":8080", mux)
 }