package grpcserver

import (
	"context"
	"github.com/go-redis/redis/v8"
	"hezzl/protogrpc"
	"log"
)

type GRPCServer struct{}

type UsersType map[string]string //Мок виде map для тестирования

var Users = make(UsersType, 10) //Мок виде map для тестирования
var rdb *redis.Client

func InitRedisConnection() {
	ctx := context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     "158.160.10.60:6379",
		Password: "wiNNer4000", // no password set
		DB:       0,            // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("redis connected!")

}

func (s *GRPCServer) AddUser(ctx context.Context, req *protogrpc.AddRequest) (*protogrpc.AddResponse, error) {
	_, ok := Users[req.User]
	if ok { //если пользователь существует в базе
		return &protogrpc.AddResponse{AddUserResponse: "User already exists: " + req.GetUser()}, nil
	} else { //если пользователя в базе нет
		Users[req.User] = "" //записали пользователя в базу
		log.Printf("User %s added", req.User)

		if rdb.Exists(ctx, "listofusers").Val() == 1 { //если в кэш есть список пользователей, то чистим, чтобы обновить
			log.Print("There is list of users in cache!")
			result, err := rdb.Del(ctx, "listofusers").Result()
			if err != nil {
				log.Print(err)
			} else if result == 1 {
				log.Print("list of users in cache cleared!")
			}
		}

		keys := []string{}
		for key, _ := range Users {
			keys = append(keys, key)
		}

		rdb.LPush(ctx, "listofusers", keys) //полностью обновляем весь кэш
		log.Printf("after adding a user into db, cash has been refreshed! Added in to cash %s", keys)

		listofusers := rdb.LRange(ctx, "listofusers", 0, -1).Val()
		log.Printf("Cash is %s", listofusers)

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
	keys := []string{}
	for key, _ := range Users {
		keys = append(keys, key)
	}

	return &protogrpc.ListUsersResponse{Listusers: keys}, nil
}

/*func (s *GRPCServer) mustEmbedUnimplementedUsersAdminServer() {
	//TODO implement me
	panic("implement me")
}
*/
