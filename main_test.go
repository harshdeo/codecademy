package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

const filename = "1.jpg"

func TestUpload(t *testing.T) {
	// Prepare request body with an image file for testing
	body := &bytes.Buffer{}
	formWriter := multipart.NewWriter(body)
	fileWriter, err := formWriter.CreateFormFile("image", "./test_images/"+filename)
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	// Open the image file and copy its content to the form file
	file, err := os.Open("./test_images/" + filename)
	if err != nil {
		t.Fatalf("Error opening image file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Fatalf("Error copying file content: %v", err)
	}

	formWriter.Close()

	// Create a request with the prepared body
	req, err := http.NewRequest("POST", "/upload", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", formWriter.FormDataContentType())
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/upload", Upload).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := fmt.Sprintf("File - %v  was uploaded successfully", filename)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Clean up by deleting the uploaded file
	err = os.Remove("./images/" + filename)
	if err != nil {
		t.Fatalf("Error deleting test file: %v", err)
	}

	fmt.Println("TestUpload has passed successfully")
}

func TestUploadFailure(t *testing.T) {
	// Prepare a request body with a valid image file
	body := &bytes.Buffer{}
	formWriter := multipart.NewWriter(body)

	// Create a temporary file with content (use a real file or create a temporary one)
	fileContent := []byte("file content")
	fileName := "test.jpg"
	fileWriter, err := formWriter.CreateFormFile("image", fileName)
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	fileWriter.Write(fileContent)

	formWriter.Close()

	// Create a request with the prepared body
	req, err := http.NewRequest("POST", "/upload", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", formWriter.FormDataContentType())
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// Simulate a failure by returning a BadRequest
		http.Error(w, "Simulated failure", http.StatusBadRequest)
	}).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	fmt.Println("TestUploadFailure has passed successfully")		
}

func TestDelete(t *testing.T) {
	// Prepare a test file to delete
	Filename := "testfile.jpg"
	content := []byte("This is a test image content")
	err := os.WriteFile("./images/"+Filename, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request with a route variable for the filename
	req, err := http.NewRequest("DELETE", "/delete/"+Filename, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/delete/{name}", Delete).Methods("DELETE")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := fmt.Sprintf("File - %v  was deleted successfully", Filename)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	fmt.Println("TestDelete has passed successfully")
}

func TestDeleteNonExistentFile(t *testing.T) {
	// Prepare a non-existent test file
	Filename := "nonexistent.jpg"

	// Create a request with a route variable for the non-existent filename
	req, err := http.NewRequest("DELETE", "/delete/"+Filename, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/delete/{name}", Delete).Methods("DELETE")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	fmt.Println("TestDeleteNonExistentFile has passed successfully")
}

func TestUpdate(t *testing.T) {
	// Prepare a test file to delete
	content := []byte("This is a test image content")
	err := os.WriteFile("./images/"+filename, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare request body with an image file for testing
	body := &bytes.Buffer{}
	formWriter := multipart.NewWriter(body)
	fileWriter, err := formWriter.CreateFormFile("image", "./test_images/"+filename)
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	// Open the image file and copy its content to the form file
	file, err := os.Open("./test_images/" + filename)
	if err != nil {
		t.Fatalf("Error opening image file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Fatalf("Error copying file content: %v", err)
	}

	formWriter.Close()

	// Create a request with the prepared body
	req, err := http.NewRequest("POST", "/update", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", formWriter.FormDataContentType())
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/update", Update).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := fmt.Sprintf("File - %v  was updated successfully", filename)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Clean up by deleting the uploaded file
	err = os.Remove("./images/" + filename)
	if err != nil {
		t.Fatalf("Error deleting test file: %v", err)
	}

	fmt.Println("TestUpdate has passed successfully")

}

func TestUpdateNonExistentFile(t *testing.T) {
	// Prepare request body with an image file for testing
	body := &bytes.Buffer{}
	formWriter := multipart.NewWriter(body)
	fileWriter, err := formWriter.CreateFormFile("image", "./test_images/"+filename)
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	// Open the image file and copy its content to the form file
	file, err := os.Open("./test_images/" + filename)
	if err != nil {
		t.Fatalf("Error opening image file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Fatalf("Error copying file content: %v", err)
	}

	formWriter.Close()

	// Create a request with the prepared body
	req, err := http.NewRequest("POST", "/update", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", formWriter.FormDataContentType())
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/update", Update).Methods("POST")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := fmt.Sprintf("File - %v - didn't exist, so the file was uploaded successfully", filename)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Clean up by deleting the uploaded file
	err = os.Remove("./images/" + filename)
	if err != nil {
		t.Fatalf("Error deleting test file: %v", err)
	}

	fmt.Println("TestUpdateNonExistentFile has passed successfully")
}

func TestFetch(t *testing.T) {
	// Prepare a test file to fetch
	Filename := "testfile.jpg"
	content := []byte("This is a test file content")
	err := os.WriteFile("./images/"+Filename, content, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove("./images/" + Filename)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Create a request with a route variable for the filename
	req, err := http.NewRequest("GET", "/fetch/"+Filename, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/fetch/{name}", Fetch).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response headers
	if contentType := rr.Header().Get("Content-Type"); contentType != "image/jpg" {
		t.Errorf("Handler returned wrong Content-Type: got %v want %v", contentType, "image/jpg")
	}

	if contentDisposition := rr.Header().Get("Content-Disposition"); contentDisposition != `attachment;filename="`+Filename+`"` {
		t.Errorf("Handler returned wrong Content-Disposition: got %v want %v", contentDisposition, `attachment;filename="`+Filename+`"`)
	}

	fmt.Println("TestFetch has passed successfully")
}

func TestFetchNonExistentFile(t *testing.T) {
	// Prepare a non-existent test file
	Filename := "nonexistent.jpg"

	// Create a request with a route variable for the non-existent filename
	req, err := http.NewRequest("GET", "/fetch/"+Filename, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/fetch/{name}", Fetch).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	fmt.Println("TestFetchNonExistentFile has passed successfully")
}

func TestFetchlist(t *testing.T) {

	// Create some test files in the directory
	dirPath := "./images/"
	fileNames := []string{"file1.jpg", "file2.jpg", "file3.jpg"}
	for _, fileName := range fileNames {
		filePath := dirPath + fileName
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatal(err)
		}

		//Remove the files after testing
		defer os.Remove(filePath)
		defer file.Close()
	}

	// Create a request
	req, err := http.NewRequest("GET", "/fetchlist", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a test router
	router := mux.NewRouter()
	router.HandleFunc("/fetchlist", Fetchlist).Methods("GET")

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response body
	var response []string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing JSON response: %v", err)
	}

	// Check if the response contains the expected file names
	for _, fileName := range fileNames {
		found := false
		for _, respFileName := range response {
			if fileName == respFileName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in the response", fileName)
		}
	}
	
	fmt.Println("TestFetchlist has passed successfully")
}