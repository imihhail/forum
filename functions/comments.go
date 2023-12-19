package functions

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type SinglePost struct {
	PostID         string
	Post           string
	UserID         string
	Username       string
	LoggedUserID   string
	LoggedUserName string
	LoggedIn       bool
	CommentPostID  string
	Picture        string
}

type AllComments struct {
	CommentID          string
	AllComments        string
	CommentCreator     string
	CommentCreatorName string
	CommentLikes       string
	CommentDislikes    string
}

var rememberPostID string

func CommentsPage(w http.ResponseWriter, r *http.Request) {
	CheckSession(r)

	var SelectedPostComments []AllComments
	var Single_Post SinglePost

	//Update comment page with user status
	Single_Post.LoggedUserID = Status.UserID
	Single_Post.LoggedIn = Status.LoggedIn

	// ALEX: fetching current user name out of DB
	if Status.LoggedIn {
		selectName := db.QueryRowContext(context.Background(), "select USERNAME from PEOPLE where ID = ?", Single_Post.LoggedUserID)
		err := selectName.Scan(&Single_Post.LoggedUserName)
		if err != nil {
			fmt.Println("Error scanning USERS 1:", err)
		}
	}

	// Get Post ID from user input and give it to comment, so comments are connected with post
	postID := r.FormValue("PostID")
	user_comment := r.FormValue("user_comment")

	// Remember postID, so user can like/dislike post on comments page
	if postID != "" {
		rememberPostID = postID
	}
	postID = rememberPostID

	if r.Method == "POST" {
		if Status.LoggedIn {

			likedPostId := r.FormValue("like")
			dislikePostID := r.FormValue("dislike")

			if likedPostId != "" {
				LikeDislike(likedPostId, "commentLike", &Single_Post.LoggedUserID)
			}

			if dislikePostID != "" {
				LikeDislike(dislikePostID, "commentDislike", &Single_Post.LoggedUserID)
			}
		}
	}

	// Display post on comments page, which user selecte from homepage
	selectPost := db.QueryRowContext(context.Background(), "select ID, POST, USERID, IMG from POSTS where ID = ?", postID)
	err := selectPost.Scan(&Single_Post.PostID, &Single_Post.Post, &Single_Post.UserID, &Single_Post.Picture)
	if err != nil {
		fmt.Println("Error scanning post:", err)
	}

	// ALEX: fetching post author name out of DB
	selectName := db.QueryRowContext(context.Background(), "select USERNAME from PEOPLE where ID = ?", Single_Post.UserID)
	err = selectName.Scan(&Single_Post.Username)
	if err != nil {
		fmt.Println("Error scanning USERS 3:", err)
	}

	// Insert comment into database if comment is not empty
	if r.Method == "POST" && user_comment != "" {
		insertComment, err := db.Prepare("INSERT INTO COMMENTSTABLE(POSTID, COMMENT, USERID) VALUES(?,?,?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Comment is inserted into database with the username of the user who posted it.
		_, err = insertComment.Exec(postID, user_comment, Status.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Select all comments from database, which are connected with postID
	sqlComment, err := db.Query("Select ID, COMMENT, USERID, LIKECOUNT, DISLIKECOUNT from COMMENTSTABLE where POSTID = ?", postID)
	if err != nil {
		fmt.Println("Error scanning post comment:", err)
		log.Fatal(err)
	}

	// Loop over all comments which are connected with postID and append them to SelectedPostComments
	// ALEX: ***HERE WAS DIFFERENT UserName DECLARATION FOR SOME REASON***
	for sqlComment.Next() {
		var id string
		var commment string
		var userID string
		var tempCreatorName string
		var commentLikes string
		var commentDislikes string

		if err := sqlComment.Scan(&id, &commment, &userID, &commentLikes, &commentDislikes); err != nil {
			log.Fatal(err)
		}

		// ALEX: fetching comment author name out of DB
		selectName := db.QueryRowContext(context.Background(), "select USERNAME from PEOPLE where ID = ?", userID)
		err = selectName.Scan(&tempCreatorName)
		if err != nil {
			fmt.Println("Error scanning USERS 3:", err)
		}

		SelectedPostComments = append(SelectedPostComments, AllComments{CommentID: id, AllComments: commment, CommentCreator: userID, CommentCreatorName: tempCreatorName, CommentLikes: commentLikes, CommentDislikes: commentDislikes})
	}

	type Template struct {
		SelectedPostComments []AllComments
		Single_Post          SinglePost
	}

	CombinedTemplate := Template{
		Single_Post:          Single_Post,
		SelectedPostComments: SelectedPostComments,
	}

	custom_tmpl, err := template.ParseFiles("template/comment.html")
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
