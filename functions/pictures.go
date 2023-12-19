package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ImageUpload(w http.ResponseWriter, r *http.Request) string {
	r.ParseMultipartForm(20 << 20)

	// Get file from user input. Send empty string if user decides not to upload file
	file, handler, err := r.FormFile("picture")
	if file == nil {
		return ""
	}

	if err != nil {
		fmt.Println("Error retrieving image from form.", err)
	}
	defer file.Close()

	// Check the file size. If its bigger than 20MB, then return error msg.
	if handler.Size > 20<<20 {
		UserStatus.Errosmsg = "The uploaded file is too large. Please upload a file smaller than 20MB."
		return "Too big file"
	}

	// Creates temporary file in folder and give it unique name.
	tempFile, err := ioutil.TempFile("./static/userpictures/", "upload-*.jpg")
	if err != nil {
		fmt.Println("Error downloading image.", err)
	}
	defer tempFile.Close()

	// Reads the file into a byte slice.
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading bytes from file.", err)
	}

	//Writes the byte slice to the temporary file.
	tempFile.Write(fileBytes)

	//Get image location and send it to POSTS SQLite database
	imagePath := tempFile.Name()
	return imagePath
}
