package postgetter

import (
	"encoding/json"
	configurer "github.com/phluan/GrabGoTrainingWeek5Assignment/configurer"
	"io/ioutil"
)

const (
	getPostsEndpoint = "https://my-json-server.typicode.com/typicode/demo/posts"
)

type PostGetter interface {
	GetPosts() ([]Post, error)
}

type Post struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type PostGetterImpl struct {
	configurer configurer.Configurer
}

func (impl *PostGetterImpl) GetPosts() ([]Post, error) {
	resp, err := impl.configurer.HTTPClient().Get(getPostsEndpoint)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	var posts []Post
	if err = json.Unmarshal(body, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func New(configurer configurer.Configurer) (PostGetter, error) {
	return &PostGetterImpl{
		configurer: configurer,
	}, nil
}
