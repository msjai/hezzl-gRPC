package grpcserver

import (
	"context"
	"github.com/go-redis/redis/v8"
	"hezzl/protogrpc"
	"log"
	"time"
)

type GRPCServer struct{}

type UsersType map[string]string //Мок виде map для тестирования

var (
	Users = make(UsersType, 10) //Мок виде map для тестирования
)

func InitRedisConnection(ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "158.160.10.60:6379",
		Password: "wiNNer4000", // no password set
		DB:       0,            // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("redis connected!")

	return rdb
}

func (s *GRPCServer) AddUser(ctx context.Context, req *protogrpc.AddRequest) (*protogrpc.AddResponse, error) {
	rdb := InitRedisConnection(ctx)
	defer rdb.Close()

	_, ok := Users[req.User]
	if ok { //если пользователь существует в базе
		return &protogrpc.AddResponse{AddUserResponse: "User already exists: " + req.GetUser()}, nil
	} else { //если пользователя в базе нет
		Users[req.User] = "" //записали пользователя в базу
		log.Printf("User %s added to db", req.User)

		if rdb.Exists(ctx, "listofusers").Val() > 0 { //если в кэш есть список пользователей, то чистим, чтобы обновить
			log.Print("There is list of users in cache!")
			result, err := rdb.Del(ctx, "listofusers").Result()
			if err != nil {
				log.Print(err)
			} else if result == 1 {
				log.Print("list of users in cache cleared!")
			}
		} else {
			log.Print("Cache is empty!")
		}
		keys := makeKeys()

		rdb.LPush(ctx, "listofusers", *keys) //полностью обновляем весь кэш
		rdb.Expire(ctx, "listofusers", 60*time.Second)
		log.Printf("after adding a user into db, cash has been refreshed! Added in to cash %s", *keys)

		listOfUsers := getListOfUsersFromCache(ctx, rdb)
		log.Printf("Cash is %s", *listOfUsers)

		return &protogrpc.AddResponse{AddUserResponse: "User added: " + req.GetUser()}, nil
	}
}

func (s *GRPCServer) DelUser(ctx context.Context, req *protogrpc.DelRequest) (*protogrpc.DelResponse, error) {
	rdb := InitRedisConnection(ctx)
	defer rdb.Close()

	_, ok := Users[req.User]
	if ok { //если пользователь существует в базе, удалем его из базы
		delete(Users, req.User)
		log.Printf("User %s deleted from db", req.User)

		if rdb.Exists(ctx, "listofusers").Val() > 0 { //если в кэш есть список пользователей, то чистим, чтобы обновить
			log.Print("There is list of users in cache!")
			result, err := rdb.Del(ctx, "listofusers").Result()
			if err != nil {
				log.Print(err)
			} else if result == 1 {
				log.Print("list of users in cache cleared!")
			}
		} else {
			log.Print("Cache is empty!")
		}
		keys := makeKeys()

		rdb.LPush(ctx, "listofusers", *keys) //полностью обновляем весь кэш
		rdb.Expire(ctx, "listofusers", 60*time.Second)
		log.Printf("after deleting a user from db, cash has been refreshed! Added in to cash %s", *keys)

		listOfUsers := getListOfUsersFromCache(ctx, rdb)
		log.Printf("Cash is %s", *listOfUsers)

		return &protogrpc.DelResponse{DelUserResponse: "User deleted succesfully: " + req.GetUser()}, nil
	} else { //если пользователя в базе не существовало
		return &protogrpc.DelResponse{DelUserResponse: "User doesn't exists: " + req.GetUser()}, nil
	}
}

func (s *GRPCServer) ListUsers(ctx context.Context, req *protogrpc.ListUsersRequest) (*protogrpc.ListUsersResponse, error) {
	rdb := InitRedisConnection(ctx)
	defer rdb.Close()

	if rdb.Exists(ctx, "listofusers").Val() > 0 {
		listOfUsers := getListOfUsersFromCache(ctx, rdb)
		log.Printf("List of users got from Cache!!! They are: %s", *listOfUsers)

		return &protogrpc.ListUsersResponse{Listusers: *listOfUsers}, nil
	} else {
		keys := makeKeys()

		rdb.LPush(ctx, "listofusers", *keys) //полностью обновляем весь кэш
		rdb.Expire(ctx, "listofusers", 60*time.Second)
		listOfUsers := getListOfUsersFromCache(ctx, rdb)
		log.Printf("List of users got from db to refresh cache! Cache refreshed. List of users got from cache! They are: %s", *listOfUsers)

		return &protogrpc.ListUsersResponse{Listusers: *listOfUsers}, nil
	}
}

func getListOfUsersFromCache(ctx context.Context, rdb *redis.Client) *[]string {
	val := rdb.LRange(ctx, "listofusers", 0, -1).Val()
	return &val
}

func makeKeys() *[]string {
	keys := []string{}
	for key, _ := range Users {
		keys = append(keys, key)
	}

	return &keys
}

/*func (s *GRPCServer) mustEmbedUnimplementedUsersAdminServer() {
	//TODO implement me
	panic("implement me")
}
*/
