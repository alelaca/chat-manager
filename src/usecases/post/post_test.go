package post

import (
	"errors"
	"testing"
	"time"

	"github.com/alelaca/chat-manager/src/entities/dtos"
	"github.com/alelaca/chat-manager/src/repository"
	"github.com/alelaca/chat-manager/src/topics"
)

func TestCreatePost_OK(t *testing.T) {
	postDTO := dtos.PostDTO{
		Message: "hello world!",
		Sender:  "Jane",
		Room:    "work-room",
	}

	handler := Handler{
		TopicsHandler: &topics.TopicMock{},
		Repository:    &repository.RepositoryMock{},
	}

	timestamp := time.Now()

	post, err := handler.CreatePost(postDTO)

	if err != nil {
		t.Errorf("Test failed with error: '%s'", err.Error())
	}

	if post.ID == "" {
		t.Errorf("Test failed. Post id not set")
	}

	if post.Message != postDTO.Message {
		t.Errorf("Test failed. Post message expected '%s', got '%s'", postDTO.Message, post.Message)
	}

	if post.Sender != postDTO.Sender {
		t.Errorf("Test failed. Post sender expected '%s', got '%s'", postDTO.Sender, post.Sender)
	}

	if post.Room != postDTO.Room {
		t.Errorf("Test failed. Post room expected '%s', got '%s'", postDTO.Room, post.Room)
	}

	if post.CreatedDate.Before(timestamp) {
		t.Errorf("Test failed. Post CreatedDate not configured properly")
	}
}

func TestCreatePost_CommandOK(t *testing.T) {
	postDTO := dtos.PostDTO{
		Message: "/stock=aapl.us",
		Sender:  "Jane",
		Room:    "work-room",
	}

	handler := Handler{
		TopicsHandler: &topics.TopicMock{},
		Repository:    &repository.RepositoryMock{},
	}

	timestamp := time.Now()

	post, err := handler.CreatePost(postDTO)

	if err != nil {
		t.Errorf("Test failed with error: '%s'", err.Error())
	}

	if post.ID == "" {
		t.Errorf("Test failed. Post id not set")
	}

	if post.Message != postDTO.Message {
		t.Errorf("Test failed. Post message expected '%s', got '%s'", postDTO.Message, post.Message)
	}

	if post.Sender != postDTO.Sender {
		t.Errorf("Test failed. Post sender expected '%s', got '%s'", postDTO.Sender, post.Sender)
	}

	if post.Room != postDTO.Room {
		t.Errorf("Test failed. Post room expected '%s', got '%s'", postDTO.Room, post.Room)
	}

	if post.CreatedDate.Before(timestamp) {
		t.Errorf("Test failed. Post CreatedDate not configured properly")
	}
}

func TestCreatePost_EntitiesFail(t *testing.T) {
	postDTO := dtos.PostDTO{
		Message: "",
		Sender:  "Jane",
		Room:    "work-room",
	}

	handler := Handler{
		TopicsHandler: &topics.TopicMock{},
		Repository:    &repository.RepositoryMock{},
	}

	_, err := handler.CreatePost(postDTO)

	if err == nil {
		t.Errorf("Test failed, error expected")
	}
}

func TestCreatePost_RepositoryFail(t *testing.T) {
	postDTO := dtos.PostDTO{
		Message: "Hello world!",
		Sender:  "Jane",
		Room:    "work-room",
	}

	handler := Handler{
		TopicsHandler: &topics.TopicMock{},
		Repository: &repository.RepositoryMock{
			SavePostError: errors.New("repo error"),
		},
	}

	_, err := handler.CreatePost(postDTO)

	if err == nil {
		t.Errorf("Test failed, error expected")
	}
}

func TestCreatePost_CommandNotifyFail(t *testing.T) {
	postDTO := dtos.PostDTO{
		Message: "/stock=aapl.us",
		Sender:  "Jane",
		Room:    "work-room",
	}

	handler := Handler{
		TopicsHandler: &topics.TopicMock{
			NotifyMessageError: errors.New("error"),
		},
		Repository: &repository.RepositoryMock{},
	}

	_, err := handler.CreatePost(postDTO)

	if err == nil {
		t.Errorf("Test failed, error expected")
	}
}
