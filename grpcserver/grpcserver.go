package grpcserver

import (
	"context"
	"hezzl/protogrpc"
)

type GRPCServer struct{}

type UsersType map[string]string //Мок виде map для тестирования

var Users = make(UsersType, 10) //Мок виде map для тестирования

func (s *GRPCServer) AddUser(ctx context.Context, req *protogrpc.AddRequest) (*protogrpc.AddResponse, error) {
	_, ok := Users[req.User]

	if ok {
		return &protogrpc.AddResponse{AddUserResponse: "User already exists: " + req.GetUser()}, nil
	} else {
		Users[req.User] = ""
		return &protogrpc.AddResponse{AddUserResponse: "User added: " + req.GetUser()}, nil
	}
}

func (s *GRPCServer) DelUser(ctx context.Context, req *protogrpc.DelRequest) (*protogrpc.DelResponse, error) {
	_, ok := Users[req.User]

	if ok {
		delete(Users, req.User)
		return &protogrpc.DelResponse{DelUserResponse: "User deleted succesfully: " + req.GetUser()}, nil
	} else {
		return &protogrpc.DelResponse{DelUserResponse: "User doesn't exists: " + req.GetUser()}, nil
	}

}

func (s *GRPCServer) ListUsers(ctx context.Context, req *protogrpc.ListUsersRequest) (*protogrpc.ListUsersResponse, error) {
	/*keys := []string{}
	for key, _ := range Users {
		keys = append(keys, key)
	}*/

	return &protogrpc.ListUsersResponse{Listusers: Users}, nil
}

/*func (s *GRPCServer) mustEmbedUnimplementedUsersAdminServer() {
	//TODO implement me
	panic("implement me")
}
*/
