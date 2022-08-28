package grpcserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"hezzl/broker"
	"hezzl/protogrpc"
	"log"
	"time"
)

type GRPCServer struct{}

type UsersType map[string]string //Мок виде map для тестирования

var (
	Users       = make(UsersType, 10) //Мок виде map для тестирования
	ErrNoRecord = errors.New("record not found")
)

const (
	host     = "158.160.10.60"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "hezzlusers"
)

func InitPostgresConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("postgres connected!")

	return db
}

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

func UserExistsPostgres(db *sql.DB, userName string) bool {
	var userFromDB string

	stmt := `SELECT "user_name" FROM "users" WHERE user_name = $1`
	row := db.QueryRow(stmt, userName)

	err := row.Scan(&userFromDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("there is no such user %s in data base", userName)
			return false
		} else {
			log.Print(err)
			return false
		}
	}

	return true
}

func AddUserPostgres(db *sql.DB, userName string) error {
	stmt := `INSERT INTO users (user_name) VALUES($1)`

	_, err := db.Exec(stmt, userName)
	if err != nil {
		return err
	}

	return nil
}

func DelUserPostgres(db *sql.DB, userName string) error {
	stmt := `DELETE FROM users WHERE user_name=$1`

	_, err := db.Exec(stmt, userName)
	if err != nil {
		return err
	}

	return nil
}

func (s *GRPCServer) AddUser(ctx context.Context, req *protogrpc.AddRequest) (*protogrpc.AddResponse, error) {
	rdb := InitRedisConnection(ctx)
	defer rdb.Close()

	db := InitPostgresConnection()
	defer db.Close()

	ok := UserExistsPostgres(db, req.User)
	/*_, ok := Users[req.User]*/
	if ok { //если пользователь существует в базе
		return &protogrpc.AddResponse{AddUserResponse: "User already exists: " + req.GetUser()}, nil
	} else { //если пользователя в базе нет
		/*Users[req.User] = "" //записали пользователя в базу*/
		err := AddUserPostgres(db, req.User) //записали пользователя в базу
		if err != nil {
			return nil, err
		}
		log.Printf("User %s added to db", req.User)
		broker.Produce(ctx, fmt.Sprintf("User %s added to db", req.User))

		if rdb.Exists(ctx, "listofusers").Val() > 0 { //если в кэш есть список пользователей, то чистим, чтобы обновить
			log.Print("There is list of users in cache!")
			result, err := rdb.Del(ctx, "listofusers").Result()
			if err != nil {
				log.Print(err)
				return nil, err
			} else if result == 1 {
				log.Print("list of users in cache cleared!")
			}
		} else {
			log.Print("Cache is empty!")
		}
		keys := makeKeys(db)

		rdb.LPush(ctx, "listofusers", *keys) //полностью обновляем весь кэш
		rdb.Expire(ctx, "listofusers", 60*time.Second)
		log.Printf("after adding a user %s into db, cash has been refreshed! Added in to cash %s", req.User, *keys)

		listOfUsers := getListOfUsersFromCache(ctx, rdb)
		log.Printf("Cash is %s", *listOfUsers)

		return &protogrpc.AddResponse{AddUserResponse: "User added: " + req.GetUser()}, nil
	}
}

func (s *GRPCServer) DelUser(ctx context.Context, req *protogrpc.DelRequest) (*protogrpc.DelResponse, error) {
	rdb := InitRedisConnection(ctx)
	defer rdb.Close()

	db := InitPostgresConnection()
	defer db.Close()

	ok := UserExistsPostgres(db, req.User)
	/*_, ok := Users[req.User]*/
	if ok { //если пользователь существует в базе, удалем его из базы
		err := DelUserPostgres(db, req.User) //удаляем пользователя из базы
		if err != nil {
			return nil, err
		}
		log.Printf("User %s deleted from db", req.User)

		if rdb.Exists(ctx, "listofusers").Val() > 0 { //если в кэш есть список пользователей, то чистим, чтобы обновить
			log.Print("There is list of users in cache!")
			result, err := rdb.Del(ctx, "listofusers").Result()
			if err != nil {
				log.Print(err)
				return nil, err
			} else if result == 1 {
				log.Print("list of users in cache cleared!")
			}
		} else {
			log.Print("Cache is empty!")
		}
		keys := makeKeys(db)

		rdb.LPush(ctx, "listofusers", *keys) //полностью обновляем весь кэш
		rdb.Expire(ctx, "listofusers", 60*time.Second)
		log.Printf("after deleting a user %s from db, cash has been refreshed! Added in to cash %s", req.User, *keys)

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

	db := InitPostgresConnection()
	defer db.Close()

	if rdb.Exists(ctx, "listofusers").Val() > 0 {
		listOfUsers := getListOfUsersFromCache(ctx, rdb)
		log.Printf("List of users got from Cache!!! They are: %s", *listOfUsers)

		return &protogrpc.ListUsersResponse{Listusers: *listOfUsers}, nil
	} else {
		keys := makeKeys(db)

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

func makeKeys(db *sql.DB) *[]string {
	var (
		userFromDB string
		Users      []string
	)

	stmt := `SELECT user_name FROM users`
	rows, err := db.Query(stmt)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userFromDB)
		if err != nil {
			log.Print(err)
			return nil
		}
		Users = append(Users, userFromDB)
	}

	if err = rows.Err(); err != nil {
		log.Print(err)
		return nil
	}

	return &Users
}

/*func (s *GRPCServer) mustEmbedUnimplementedUsersAdminServer() {
	//TODO implement me
	panic("implement me")
}
*/
