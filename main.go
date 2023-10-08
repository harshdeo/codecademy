package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	maxFileSize = 10 << 20 // 10*(2^20) Bytes = 10MB
	apiKey = "123456789"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	// Limit of 10MB on the size of the file
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the cat image from the request
	catimage, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error fetching the file from the http request - %v", err.Error()))
		http.Error(w, "Image can't be fetched. Check if the parameter key is 'image'.", http.StatusBadRequest)
		return
	}
	defer catimage.Close()

	// Create a file to store the image on disk
	// Assumption - no two files have different names
	catfile, err := os.Create("./images/" + handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer catfile.Close()

	// Copy the uploaded image to the created file
	_, err = io.Copy(catfile, catimage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("File - %v  was uploaded successfully", handler.Filename))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	// Handle preflight request
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")	

	// Get the filename from the route variables
	filename := mux.Vars(r)["name"]

	// Delete the file
	err := os.Remove("./images/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("File - %v  was deleted successfully", filename))
}

func Update(w http.ResponseWriter, r *http.Request) {
	// Limit of 10MB on the size of the file
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the cat image from the request
	catimage, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error fetching the file from the http request - %v", err.Error()))
		http.Error(w, "Image can't be fetched. Check if the parameter key is 'image'.", http.StatusBadRequest)
		return
	}
	defer catimage.Close()

	_, err = os.Stat("./images/" + handler.Filename)

	if errors.Is(err, os.ErrNotExist) {
		// If file doesn't exist, create a file to store the image on disk
		// Assumption - no two files have different names
		catfile, err := os.Create("./images/" + handler.Filename)
		if err != nil {
			fmt.Println(err)
		}
		defer catfile.Close()

		// Copy the uploaded image to the created file
		_, err = io.Copy(catfile, catimage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, fmt.Sprintf("File - %v - didn't exist, so the file was uploaded successfully", handler.Filename))
	} else {
		// If file exists, then open the existing file and update it
		catfile, err := os.OpenFile("./images/"+handler.Filename, os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer catfile.Close()

		// Copy the uploaded image to the existing file
		_, err = io.Copy(catfile, catimage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, fmt.Sprintf("File - %v  was updated successfully", handler.Filename))
	}
}

func Fetch(w http.ResponseWriter, r *http.Request) {
	// Get the filename from the route variables
	filename := mux.Vars(r)["name"]
	_, err := os.Stat("./images/" + filename)

	if errors.Is(err, os.ErrNotExist) {
		// If file doesn't exist, bad request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		// If file exists, then read the file and send it
		iBytes, err := os.ReadFile("./images/" + filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/jpg")
		w.Header().Set("Content-Disposition", `attachment;filename="`+filename+`"`)
		w.Write(iBytes)
	}

	fmt.Fprintf(w, fmt.Sprintf("File - %v  was sent successfully", filename))
}

func Fetchlist(w http.ResponseWriter, r *http.Request) {
	// Open the images directory
	list, err := os.Open("./images/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer list.Close()

	// Reads the list of files
	files, err := list.Readdir(-1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parses the base file name
	filenames := make([]string, len(files))
	for i, f := range files {
		filenames[i] = f.Name()
	}
	err = json.NewEncoder(w).Encode(filenames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the API key in the request header
		key := r.Header.Get("API-Key")

		if key != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Continue with the next handler if authentication is successful
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	// Enable CORS with gorilla/handlers
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.Use(handlers.CORS(headers, methods, origins))
    r.Use(authenticate)

	r.HandleFunc("/upload", Upload).Methods("POST")
	r.HandleFunc("/delete/{name}", Delete).Methods("DELETE")
	r.HandleFunc("/update", Update).Methods("POST")
	r.HandleFunc("/fetch/{name}", Fetch).Methods("GET")
	r.HandleFunc("/fetchlist", Fetchlist).Methods("GET")

	log.Print("Running Cat server http://localhost:8080")
	log.Println(http.ListenAndServe(":8080", r))
}
