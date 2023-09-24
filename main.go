package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// PageData struct holds the data to be passed to the template.
type PageData struct {
	AsciiArt string
}

func handleDefault(w http.ResponseWriter, r *http.Request) { // Check if the request URL path is exactly "/" and Send a 404 Not Found error page
	if r.URL.Path != "/" { // r.URL.Path
		renderErrorPage(w, http.StatusNotFound)
		return
	}
	renderTemplate(w, "templates/index.html", nil)
}

// handleAsciiArt is the handler for the "/ascii-art" route.// Only allow POST requests for this route
func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { //or if r.Mthod!="POST"
		renderErrorPage(w, http.StatusMethodNotAllowed) // Send a 405 Method Not Allowed error page
		return
	}
	// Get the form values for "text" and "banner"
	text := r.FormValue("text")
	if text <= " " || text >= "~" {
		http.Error(w, "400 Bad Request ", http.StatusBadRequest)
		return
	}
	banner := r.FormValue("banner")
	result := generateAsciiArt(text, banner) // Generate ASCII art based on the provided text and banner

	data := PageData{ // Create a PageData struct with the generated ASCII art
		AsciiArt: result,
	}
	renderTemplate(w, "templates/index.html", data) // Render the main.html template with the generated ASCII art
}

// renderTemplate renders the specified template with the provided data.
func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile) // Parse the template file
	if err != nil {
		renderErrorPage(w, http.StatusNotFound) // Send a 404 Not Found error page
		return
	}
	err = tmpl.Execute(w, data) // Execute the template with the provided data
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError) // Send a 500 Internal Server Error page
	}
}

func generateAsciiArt(text, banner string) string {
	words := strings.Split(text, "\n")
	rawBytes, err := os.ReadFile(fmt.Sprintf("templates/%s.txt", banner))
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.ReplaceAll(string(rawBytes), "\r\n", "\n"), "\n")
	// Create a strings.Builder to store the resulting ASCII art
	var result strings.Builder
	for i, word := range words {
		if word == "" {
			if i < len(words)-1 {
				result.WriteString("\n")
			}
			continue
		}
		for h := 1; h < 9; h++ {
			for _, l := range word {
				for lineIndex, line := range lines {
					if lineIndex == (int(l)-32)*9+h {
						result.WriteString(line)
					}
				}
			}
			result.WriteString("\n")
		}
	}

	return result.String()
}

// renderErrorPage renders an error page based on the status code.
func renderErrorPage(w http.ResponseWriter, status int) {

	w.WriteHeader(status) // Set the status code for the response
	// Parse and execute the error page template with the status code
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/%d.html", status))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, status)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
func main() {

	fs := http.FileServer(http.Dir("static")) // Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleDefault) // Define the handlers for different routes
	http.HandleFunc("/ascii-art", handleAsciiArt)
	fmt.Println("Server started. Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
