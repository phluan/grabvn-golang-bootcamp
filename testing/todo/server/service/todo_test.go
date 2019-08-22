package service

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/xuanit/testing/todo/pb"
	"github.com/xuanit/testing/todo/server/repository/mocks"
	"testing"
)

func TestGetToDo(t *testing.T) {
	mockToDoRep := &mocks.ToDo{}
	toDo := &pb.Todo{}
	req := &pb.GetTodoRequest{Id: "123"}
	mockToDoRep.On("Get", req.Id).Return(toDo, nil)
	service := ToDo{ToDoRepo: mockToDoRep}

	res, err := service.GetTodo(nil, req)

	expectedRes := &pb.GetTodoResponse{Item: toDo}

	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
	mockToDoRep.AssertExpectations(t)
}

func TestCreateToDo_success(t *testing.T) {
	mockToDoRep := &mocks.ToDo{}
	toDo := &pb.Todo{}
	mockToDoRep.On("Insert", toDo).Return(nil)

	req := &pb.CreateTodoRequest{Item: toDo}
	service := ToDo{ToDoRepo: mockToDoRep}
	res, err := service.CreateTodo(nil, req)

	assert.Nil(t, err)

	_, uuid_test_err := uuid.FromString(res.Id)
	assert.Nil(t, uuid_test_err)

	mockToDoRep.AssertExpectations(t)
}

func TestCreateToDo_fail(t *testing.T) {
	mockToDoRep := &mocks.ToDo{}
	toDo := &pb.Todo{}
	mockError := errors.New("Error message")
	mockToDoRep.On("Insert", mock.Anything).Return(mockError)

	req := &pb.CreateTodoRequest{Item: toDo}
	service := ToDo{ToDoRepo: mockToDoRep}
	res, err := service.CreateTodo(nil, req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "Could not insert item", err.Error())
}
