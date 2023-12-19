package functions

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type User_Status struct {
	UserID   string
	Logedin  bool
	UserName string
	Errosmsg string
}

type User_Post struct {
	Id           string
	UserPost     string
	Creator      string
	PostLikes    string
	PostDislikes string
	CommentCount int
	Picture      string
}

var AllUserPosts []User_Post
var UserStatus User_Status
var MyName string

func Homepage(w http.ResponseWriter, r *http.Request) {
	// Log in user if CheckSession(r) returns true.
	// If CheckSession(r) returns true, then user is logged in and Status.LoggedIn is true.
	CheckSession(r)

	// ALEX: fetching current user name out of DB
	if Status.LoggedIn {
		selectName := db.QueryRowContext(context.Background(), "select USERNAME from PEOPLE where ID = ?", Status.UserID)
		err := selectName.Scan(&MyName)
		if err != nil {
			Logout(w, r) // kill the cookie if its exists for some reason
			fmt.Println("Error scanning USERS 1:", err)
		}
	}

	// Update user status on homepage
	UserStatus = User_Status{UserID: Status.UserID, Logedin: Status.LoggedIn, UserName: MyName}

	if r.Method == "POST" {

		// If user is logged in, then user can post.
		if Status.LoggedIn {
			post := r.FormValue("post")
			likedPostId := r.FormValue("like")
			dislikePostID := r.FormValue("dislike")

			cat := r.FormValue("Categorie1")
			cat2 := r.FormValue("Categorie2")
			cat3 := r.FormValue("Categorie3")
			cat4 := r.FormValue("Categorie4")

			if cat == "" && cat2 == "" && cat3 == "" && cat4 == "" {
				UserStatus.Errosmsg = "You must choose atleast one categorie"
				post = ""
			}

			// Insert post into database if post is not empty and uploaded image file isnt bigger than 20MB.
			if post != "" && ImageUpload(w, r) != "Too big file" {
				statement, err := db.Prepare("INSERT INTO POSTS(CATEGORIE, CATEGORIE2, CATEGORIE3, CATEGORIE4, POST, USERID, IMG) VALUES (?,?,?,?,?,?,?)")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// Insert post into database. Post is inserted into database with the username of the user who posted it.
				_, err = statement.Exec(cat, cat2, cat3, cat4, post, Status.UserID, ImageUpload(w, r))
				
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// If user has liked or disliked a post, then PostLike function is called.
			// "like" and "dislike" are needed for PostLike function, so it knows which sqltables to use.
			if likedPostId != "" {
				LikeDislike(likedPostId, "like", &UserStatus.UserID)
				UserStatus.Errosmsg = ""
			}

			if dislikePostID != "" {
				LikeDislike(dislikePostID, "dislike", &UserStatus.UserID)
				UserStatus.Errosmsg = ""
			}
		}
	}

	// Send all posts to homepage
	SendPosts(r)

	// Create template for homepage, which displays all posts and user status.
	type Template struct {
		AllUserPosts []User_Post
		UserStatus   User_Status
	}

	CombinedTemplate := Template{
		UserStatus:   UserStatus,
		AllUserPosts: AllUserPosts,
	}

	custom_tmpl, err := template.ParseFiles("template/home.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	err = custom_tmpl.Execute(w, CombinedTemplate)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
	
}

func SendPosts(r *http.Request) {
	// Reset AllUserPosts to nil, so that posts are not duplicated.
	AllUserPosts = nil

	categorie_filter := r.FormValue("FilterDropdown")

	// Choose, which SQL command to pick, depneding on user input and insert it into sqlCommand.
	sqlCommand := ""
	allPosts := "select ID, POST, USERID, LIKECOUNT, DISLIKECOUNT, COMMENTCOUNT, IMG from POSTS"
	myposts := "select ID, POST, USERID, LIKECOUNT, DISLIKECOUNT, COMMENTCOUNT, IMG from POSTS where USERID = ?"
	myLikes := "select ID, POST, USERID, LIKECOUNT, DISLIKECOUNT, COMMENTCOUNT, IMG from POSTS where ID in (select POSTID from POSTLIKES where USERID = ?)"	
	filter := "SELECT ID, POST, USERID, LIKECOUNT, DISLIKECOUNT, COMMENTCOUNT, IMG FROM POSTS WHERE CATEGORIE = ? OR CATEGORIE2 = ? OR CATEGORIE3 = ? OR CATEGORIE4 = ?"

	// Displays posts created by logged in user.
	if categorie_filter == "MyPost" {
		sqlCommand = myposts
		categorie_filter = UserStatus.UserID // Change categorie_filter to username of logged in user.

		// Displays all posts if user hasnt selected anything yet or has selected "All posts"
	} else if categorie_filter == "AllCat" || categorie_filter == "" {
		sqlCommand = allPosts

		// Displays liked posts from logged in user.
	} else if categorie_filter == "myLikes" {
		sqlCommand = myLikes
		categorie_filter = UserStatus.UserID

		// Displays posts from selected categorie.
	} else {
		sqlCommand = filter
	}

	// Select data from database, depending what sqlCommand is.
	sql_posts, err := db.Query(sqlCommand, categorie_filter, categorie_filter, categorie_filter, categorie_filter)
	if err != nil {
		fmt.Println("Error filtering posts from database:", err)
	}

	// Loop over SQL query and append posts to AllUserPosts to display all posts on homepage.
	for sql_posts.Next() {
		var id string
		var allposts string
		var created string
		var createdName string
		var postLikes string
		var postDisLikes string
		var commentcount int
		var image string

		if err := sql_posts.Scan(&id, &allposts, &created, &postLikes, &postDisLikes, &commentcount, &image); err != nil {
			log.Fatal(err)
		}

		// ALEX: getting comment author name out of DB
		selectName := db.QueryRowContext(context.Background(), "select USERNAME from PEOPLE where ID = ?", created)
		err := selectName.Scan(&createdName)
		if err != nil {
			fmt.Println("Error scanning USERS 2:", err)
		}

		err = db.QueryRow("select count(*) from commentstable WHERE POSTID = ?", id).Scan(&commentcount)
		if err != nil {
			fmt.Println("Error counting comments.", err)
			return
		}
		AllUserPosts = append(AllUserPosts, User_Post{Id: id, UserPost: allposts, Creator: createdName, PostLikes: postLikes, PostDislikes: postDisLikes, CommentCount: commentcount, Picture: image})
	}
}

func Handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 not found")
}
