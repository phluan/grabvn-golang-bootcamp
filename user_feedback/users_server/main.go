package main

import (
	"context"
	"database/sql"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	status "google.golang.org/grpc/status"
	database "grab.com/luanpham/users_feedback/db"
	models "grab.com/luanpham/users_feedback/models"
	pb "grab.com/luanpham/users_feedback/pb"
	"log"
	"net"
	"regexp"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Create(ctx context.Context, in *pb.UserSignUpRequest) (*pb.UserResponse, error) {
	user := models.User{Email: in.Email, Password: in.Password, FirstName: in.FirstName, LastName: in.LastName}
	sqlStatement := `
		INSERT INTO users (email, password, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	err := database.DBCon.QueryRow(sqlStatement, user.Email, user.Password, user.FirstName, user.LastName).Scan(&id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot create user")
	}

	return &pb.UserResponse{Id: id, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}, nil
}

func (s *server) List(ctx context.Context, in *pb.PaginationRequest) (*pb.UsersResponse, error) {
	sqlStatement := `
		SELECT id, email, first_name, last_name FROM users
	`
	rows, err := database.DBCon.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var userResponses []*pb.UserResponse
	for rows.Next() {
		userResponse := pb.UserResponse{}
		err := rows.Scan(&userResponse.Id, &userResponse.Email, &userResponse.FirstName, &userResponse.LastName)
		if err != nil {
			panic(err)
		}
		userResponses = append(userResponses, &userResponse)
	}

	return &pb.UsersResponse{Users: userResponses}, nil
}

func (s *server) Get(ctx context.Context, in *pb.UserIdentifier) (*pb.UserResponse, error) {
	sqlStatement := `
		SELECT id, email, first_name, last_name
		FROM users
		WHERE id = $1
	`

	var user models.User
	row := database.DBCon.QueryRow(sqlStatement, in.Id)
	err := row.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName)
	log.Printf("Here")
	switch err {
	case sql.ErrNoRows:
		return nil, status.Errorf(codes.NotFound, "User not found")
	case nil:
		return &pb.UserResponse{Id: user.Id, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}, nil
	default:
		match, _ := regexp.MatchString("invalid input syntax for uuid", err.Error())
		if match {
			return nil, status.Errorf(codes.InvalidArgument, "Id should have uuid format")
		} else {
			panic(err)
		}
	}
}

func (s *server) Update(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *server) Delete(ctx context.Context, in *pb.UserIdentifier) (*pb.SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func main() {
	database.InitDB()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUsersServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
