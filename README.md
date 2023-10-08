# Cat Server

Cat Server is a simple HTTP server designed to manage cat images. It supports uploading, deleting, updating, fetching, and listing cat images.
It does authentication using an API-Key. Further, dockerization is provided at the end.

## Installation

To install the Cat Server, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/harshdeo/codecademy.git
   cd codecademy
   ```

2. Install dependencies:

   ```bash
   go get -d ./...
   ```

## API Endpoints

### 1. Upload

   **Endpoint:** `/upload`  
   **Method:** POST  
   **Description:** Upload a cat image.  
   **Request:**
   - Method: `POST`
   - Headers: 
     - `Content-Type: multipart/form-data`
     - `API-Key: YOUR_API_KEY`
   - Body: Form data with a field named `image` containing the cat image file.
   
   **Response:**
   - Status Code: 200 OK
   - Body: JSON response indicating the success of the upload.

### 2. Delete

   **Endpoint:** `/delete/{name}`  
   **Method:** DELETE  
   **Description:** Delete a cat image by filename.  
   **Request:**
   - Method: `DELETE`
   - Headers: 
     - `API-Key: YOUR_API_KEY`
   
   **Response:**
   - Status Code: 200 OK
   - Body: JSON response indicating the success of the deletion.

### 3. Update

   **Endpoint:** `/update`  
   **Method:** POST  
   **Description:** Update a cat image. If the image doesn't exist, it will be created.  
   **Request:**
   - Method: `POST`
   - Headers: 
     - `Content-Type: multipart/form-data`
     - `API-Key: YOUR_API_KEY`
   - Body: Form data with a field named `image` containing the cat image file.
   
   **Response:**
   - Status Code: 200 OK
   - Body: JSON response indicating the success of the update.

### 4. Fetch

   **Endpoint:** `/fetch/{name}`  
   **Method:** GET  
   **Description:** Fetch a cat image by filename.  
   **Request:**
   - Method: `GET`
   - Headers: 
     - `API-Key: YOUR_API_KEY`
   
   **Response:**
   - Status Code: 200 OK
   - Body: The cat image file.

   The fetched image can be downloaded to local storage using `save response to file` option in Postman.

### 5. Fetchlist

   **Endpoint:** `/fetchlist`  
   **Method:** GET  
   **Description:** Fetch a list of all cat image filenames.  
   **Request:**
   - Method: `GET`
   - Headers: 
     - `API-Key: YOUR_API_KEY`
   
   **Response:**
   - Status Code: 200 OK
   - Body: JSON array containing filenames of cat images.

## Authentication

All requests to the Cat Server API must include the `API-Key` header with a valid API key for authentication.
The value for API-Key has been set to a constant value `123456789` for testing purpose.

## Configuration

- The maximum file size for uploads is set to 10MB (`maxFileSize` constant).

## Testing

To test the functions, execute the following command:

```bash
go test
```


## Running the Server

To run the Cat Server, execute the following command:

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

## Dockerization

First ensure that the docker daemon is running.
To containerize the Cat Server using Docker, use the provided `Dockerfile`:

```bash
docker build -t cat-server .
docker run -p 8080:8080 cat-server
```

The server will be accessible on `http://localhost:8080` within the Docker container.