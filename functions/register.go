package functions
import (
	"html/template"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)
type Authentication struct {
	Failed_msg             string
	Authentication_success bool
	LoggedIn               bool
}
func Register(w http.ResponseWriter, r *http.Request) {
	var Status Authentication
	// Insert user input into the database only if POST method is used. Otherwise, empty fields will be inserted into the database, everytime when the page is loaded.
	if r.Method == "POST" {
		// Get user input
		e_mail := r.FormValue("email")
		userName := r.FormValue("username")
		pw := r.FormValue("password")
		Status.Authentication_success = false
		Status.LoggedIn = false
		var email_userName_Count int
		// Check if email or username already exists in database. Count the number of rows that match the email or username.
		// Add the count to variable email_userName_Count.
		err := db.QueryRow("SELECT COUNT(*) FROM people WHERE EMAIL = ? OR USERNAME = ?", e_mail, userName).Scan(&email_userName_Count)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// If there are more than 0 rows that match the email or username, set registration_error to true and display error Status.
		if email_userName_Count > 0 && e_mail != "" && userName != "" && pw != "" {
			Status.Failed_msg = "Email or username is already taken!"
			// Check if any fields are empty
		} else if e_mail == "" || userName == "" || pw == "" {
			Status.Failed_msg = "Please fill in empty fields!"
			// If there were no problems with the registration, set Authentication_success to true.
		} else {
			Status.Authentication_success = true
		}
		// If all fields are filled and email/username duplicate is not found then instert user input into database.
		if Status.Authentication_success {
			// Hash the password before inserting it into the database.
			hashedPW, _ := HashPW(pw)
			statement, err := db.Prepare("INSERT INTO people(EMAIL, USERNAME, PASSWORD) VALUES (?, ?, ?)")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = statement.Exec(e_mail, userName, hashedPW)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	custom_tmpl, err := template.ParseFiles("template/register.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v\n", err)
		return
	}
	err = custom_tmpl.Execute(w, Status)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v\n", err)
	}
}
// HashPW hashes the password
func HashPW(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
