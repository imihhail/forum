package functions

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginStatus struct {
	Login_msg   string
	UserName    string
	UserID      string
	Open_dialog bool
	LoggedIn    bool
}

var Status LoginStatus

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		// Get user input from login page
		e_mail := r.FormValue("email")
		input_password := r.FormValue("password")

		// Get password and username from database if users inserted email exists in database.
		// Alex ID change 1: "SELECT USERNAME, PASSWORD -> "SELECT ID, PASSWORD
		pw_userID_check := db.QueryRowContext(context.Background(), "SELECT ID, PASSWORD FROM people WHERE EMAIL = ?", e_mail)

		var db_password string
		email_found := true

		// Scan pw_userN_check for username and password and store them into db_password and Status.UserName.
		// If err != nil, then email is not found from database.
		if err := pw_userID_check.Scan(&Status.UserID, &db_password); err != nil && e_mail != "" && input_password != "" {
			Status.Login_msg = "Email: " + e_mail + " is not registered on this website"
			email_found = false

			// Check if any fields are empty
		} else if e_mail == "" || input_password == "" {
			Status.Login_msg = "Please fill in empty fields"

			// Compare the hashed password with the password entered by the user
			// If email is found from database and user inserted pw is does not match with database password, then display error message.
		} else if err := bcrypt.CompareHashAndPassword([]byte(db_password), []byte(input_password)); err != nil && e_mail != "" && email_found {
			Status.Login_msg = "Incorrect password!"

		} else {
			// If user input is valid and is not already logged in, then login is successful.
			if len(sessions) == 0 {
				Status.Open_dialog = true
				// Create cookie for the user
				CreateCookie(w)
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			}

			for _, id := range sessions {
				if id == Status.UserID {
					Status.Login_msg = "User is already logged in!"

				} else {
					Status.Open_dialog = true
					// Create cookie for the user
					CreateCookie(w)
					http.Redirect(w, r, "/home", http.StatusSeeOther)
				}
			}
		}
	}

	custom_tmpl, err := template.ParseFiles("template/login.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		fmt.Printf("Error parsing template: %v\n", err)
	}

	err = custom_tmpl.Execute(w, Status)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}

	// Setting open_dialog to false closes the dialog box, so it doesnt hover on the screen.
	Status.Open_dialog = false
	Status.Login_msg = ""
}
