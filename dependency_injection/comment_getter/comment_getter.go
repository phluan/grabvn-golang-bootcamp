package commentgetter

import (
	"encoding/json"
	configurer "github.com/phluan/GrabGoTrainingWeek5Assignment/configurer"
	"io/ioutil"
)

const (
	getCommentsEndpoint = "https://my-json-server.typicode.com/typicode/demo/comments"
)

type Comment struct {
	ID     int64  `json:"id"`
	Body   string `json:"body"`
	PostID int64  `json:"postId"`
}

type CommentGetter interface {
	GetComments() ([]Comment, error)
}

type CommentGetterImpl struct {
	configurer configurer.Configurer
}

func (impl *CommentGetterImpl) GetComments() ([]Comment, error) {
	resp, err := impl.configurer.HTTPClient().Get(getCommentsEndpoint)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	var comments []Comment
	if err = json.Unmarshal(body, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func New(configurer configurer.Configurer) (CommentGetter, error) {
	return &CommentGetterImpl{
		configurer: configurer,
	}, nil
}
