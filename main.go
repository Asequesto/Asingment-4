package usermanagement

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	users []*User
}

func NewServer() *server {
	return &server{}
}

func (s *server) AddUser(ctx context.Context, user *User) (*UserID, error) {

	return &UserID{Id: 123}, nil
}

func (s *server) GetUser(ctx context.Context, userID *UserID) (*User, error) {

	return &User{
		Id:    userID.Id,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil
}

func (s *server) ListUsers(req *Empty, stream UserService_ListUsersServer) error {
	users := []*User{
		{Id: 1, Name: "User 1", Email: "user1@example.com"},
		{Id: 2, Name: "User 2", Email: "user2@example.com"},
	}
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func StartServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	RegisterUserServiceServer(srv, NewServer())
	log.Printf("Server listening on port %s", port)
	return srv.Serve(lis)
}

func ExampleServer_AddUser() {
	s := NewServer()

	user := &User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	userID, err := s.AddUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Error adding user: %v", err)
	}

	log.Printf("User added successfully. UserID: %d", userID.Id)
}

func ExampleServer_ListUsers() {
	s := NewServer()

	mockStream := &MockUserListStream{}

	if err := s.ListUsers(&Empty{}, mockStream); err != nil {
		log.Fatalf("Error listing users: %v", err)
	}

	listedUsers := mockStream.users
	log.Printf("Listed %d users:", len(listedUsers))
	for _, user := range listedUsers {
		log.Printf("- %s", user.Name)
	}
}
