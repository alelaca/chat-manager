package entities

import (
	"testing"
	"time"

	"github.com/alelaca/chat-manager/src/entities/dtos"
)

func TestCreatePost_OK(t *testing.T) {
	msg := "Hello everyone!"
	sender := "Jane"
	room := "work-room"
	expectedTimeAfter := time.Now()

	postDTO := dtos.PostDTO{
		Message: msg,
		Sender:  sender,
		Room:    room,
	}

	post, err := CreatePost(postDTO)

	if err != nil {
		t.Errorf("Test failed with error '%s'", err.Error())
	}

	if post.ID == "" {
		t.Errorf("Test failed. Post id expected not to be empty")
	}

	if post.Message != msg {
		t.Errorf("Test failed. Post message expected '%s', got '%s'", msg, post.Message)
	}

	if post.Sender != sender {
		t.Errorf("Test failed. Post sender expected '%s', got '%s'", sender, post.Sender)
	}

	if post.Room != room {
		t.Errorf("Test failed. Post room expected '%s', got '%s'", room, post.Room)
	}

	if post.CreatedDate.Before(expectedTimeAfter) {
		t.Errorf("Test failed. Post created date not created properly")
	}
}

func TestCreatePost_ValidationFail(t *testing.T) {
	msg := ""
	sender := "Jane"
	room := "work-room"

	postDTO := dtos.PostDTO{
		Message: msg,
		Sender:  sender,
		Room:    room,
	}

	_, err := CreatePost(postDTO)

	if err == nil {
		t.Errorf("Test failed. Expected an error from function")
	}
}

func TestValidatePost_InvalidMessage(t *testing.T) {
	msg := ""
	sender := "Jane"
	room := "work-room"

	postDTO := dtos.PostDTO{
		Message: msg,
		Sender:  sender,
		Room:    room,
	}

	err := ValidatePost(postDTO)

	if err == nil {
		t.Errorf("Test failed. Expected an error from function")
	}
}

func TestValidatePost_InvalidSender(t *testing.T) {
	msg := "Hello everyone!"
	sender := ""
	room := "work-room"

	postDTO := dtos.PostDTO{
		Message: msg,
		Sender:  sender,
		Room:    room,
	}

	err := ValidatePost(postDTO)

	if err == nil {
		t.Errorf("Test failed. Expected an error from function")
	}
}

func TestValidatePost_InvalidRoom(t *testing.T) {
	msg := "Hello everyone!"
	sender := "Jane"
	room := ""

	postDTO := dtos.PostDTO{
		Message: msg,
		Sender:  sender,
		Room:    room,
	}

	err := ValidatePost(postDTO)

	if err == nil {
		t.Errorf("Test failed. Expected an error from function")
	}
}
