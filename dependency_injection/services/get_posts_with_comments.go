package services

import (
	"errors"
	"fmt"
	commentgetter "github.com/phluan/GrabGoTrainingWeek5Assignment/comment_getter"
	configurer "github.com/phluan/GrabGoTrainingWeek5Assignment/configurer"
	postgetter "github.com/phluan/GrabGoTrainingWeek5Assignment/post_getter"
)

type PostWithCommentsResponse struct {
	Posts []PostWithComments `json:"posts"`
}
type PostWithComments struct {
	ID       int64                   `json:"id"`
	Title    string                  `json:"string"`
	Comments []commentgetter.Comment `json:"comments,omitempty"`
}

func GetPostWithComments(configurer configurer.Configurer) ([]byte, error) {
	postGetter, postGetterErr := postgetter.New(configurer)
	if postGetterErr != nil {
		return nil, errors.New(fmt.Sprintf("Cannot setup post getter: %s", postGetterErr))
	}

	commentGetter, commentGetterErr := commentgetter.New(configurer)
	if commentGetterErr != nil {
		return nil, errors.New(fmt.Sprintf("Cannot setup comment getter: %s", commentGetterErr))
	}

	posts, err := postGetter.GetPosts()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get posts failed with error: ", err))
	}

	// Get comments from api
	comments, err := commentGetter.GetComments()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get comments failed with error: ", err))
	}

	// Combine and return response
	postWithComments := combinepostwithcomments(posts, comments)
	resp := PostWithCommentsResponse{Posts: postWithComments}
	buf, err := configurer.Serializer().Render(resp)
	// buf, err := json.Marshal(resp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse response: ", err))
	}

	return buf, nil
}

func combinepostwithcomments(posts []postgetter.Post, comments []commentgetter.Comment) []PostWithComments {
	commentsByPostID := map[int64][]commentgetter.Comment{}
	for _, comment := range comments {
		commentsByPostID[comment.PostID] = append(commentsByPostID[comment.PostID], comment)
	}

	result := make([]PostWithComments, 0, len(posts))
	for _, post := range posts {
		result = append(result, PostWithComments{
			ID:       post.ID,
			Title:    post.Title,
			Comments: commentsByPostID[post.ID],
		})
	}

	return result
}
