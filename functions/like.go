package functions

import (
	"fmt"
)

func LikeDislike(likedPostId string, likechoice string, userID *string) {
	var likeCount int
	var userLikeCount int

	// Variables for sql statements
	var table string
	var updateTable string
	var count string
	var deleteTable string
	var updateCount string

	// Change sql command depending if user likes or dislikes post or comment.
	if likechoice == "like" {
		table = "POSTLIKES"
		updateTable = "POSTS"
		count = "LIKECOUNT"
		deleteTable = "POSTDISLIKES"
		updateCount = "DISLIKECOUNT"

	} else if likechoice == "dislike" {
		table = "POSTDISLIKES"
		updateTable = "POSTS"
		count = "DISLIKECOUNT"
		deleteTable = "POSTLIKES"
		updateCount = "LIKECOUNT"

	} else if likechoice == "commentLike" {
		table = "COMMENTLIKES"
		updateTable = "COMMENTSTABLE"
		count = "LIKECOUNT"
		deleteTable = "COMMENTDISLIKES"
		updateCount = "DISLIKECOUNT"

	} else if likechoice == "commentDislike" {
		table = "COMMENTDISLIKES"
		updateTable = "COMMENTSTABLE"
		count = "DISLIKECOUNT"
		deleteTable = "COMMENTLIKES"
		updateCount = "LIKECOUNT"
	}

	// Count how many times user liked or disliked a post or comment.
	err := db.QueryRow("select count(*) from "+table+" WHERE POSTID = ? and USERID = ?", likedPostId, *userID).Scan(&userLikeCount)
	if err != nil {
		fmt.Println("Error selecting user like count.", err)
		return
	}

	// If usercountLikeCount is 0, then user has not liked or disliked a post or comment.
	if userLikeCount == 0 {
		recieveLike, err := db.Prepare("INSERT INTO " + table + " (POSTID, USERID) VALUES (?,?)")
		if err != nil {
			fmt.Println("Error preparing insert like statement.", err)
			return
		}

		_, err = recieveLike.Exec(likedPostId, *userID)
		if err != nil {
			fmt.Println("Error executing insert like statement.", err)
			return
		}

		// Count how many times a post or comment has been liked or disliked.
		err = db.QueryRow("select count(*) from "+table+" WHERE POSTID = ?", likedPostId).Scan(&likeCount)
		if err != nil {
			fmt.Println("Error selecting like count.", err)
			return
		}

		// Update post or comment with new like or dislike count.
		likeToPost, err := db.Prepare("UPDATE " + updateTable + " SET " + count + " = ? WHERE ID = ?")
		if err != nil {
			fmt.Println("Error preparing update like statement.", err)
			return
		}

		_, err = likeToPost.Exec(likeCount, likedPostId)
		if err != nil {
			fmt.Println("Error executing update like statement.", err)
			return
		}

		// If user has liked or disliked a post or comment, then delete the opposite.
		deleteLike, err := db.Prepare("DELETE FROM " + deleteTable + " WHERE POSTID = ? and USERID = ?")
		if err != nil {
			fmt.Println("Error preparing delete like statement.", err)
			return
		}

		_, err = deleteLike.Exec(likedPostId, *userID)
		if err != nil {
			fmt.Println("Error executing delete like statement.", err)
			return
		}

		// After deleting an opposite row, count how many times a post or comment has been liked or disliked.
		err = db.QueryRow("select count(*) from "+deleteTable+" WHERE POSTID = ?", likedPostId).Scan(&likeCount)
		if err != nil {
			fmt.Println("Error selecting like count after deleting a row.", err)
			return
		}

		// Update opposite post or comment with new like or dislike count.
		dlikeToPost, err := db.Prepare("UPDATE " + updateTable + " SET " + updateCount + " = ? WHERE ID = ?")
		if err != nil {
			fmt.Println("Error preparing update like statement after deleting a row.", err)
			return
		}

		_, err = dlikeToPost.Exec(likeCount, likedPostId)
		if err != nil {
			fmt.Println("Error executing update like statement after deleting a row.", err)
			return
		}
	}
}
