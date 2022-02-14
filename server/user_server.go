package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	pb "github.com/Titouan-Esc/golang-grpc/proto"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/grpc"
)

// ? Définir le port du server
const (
	port = ":50051"
	dsn = "postgres://titouanescorneboueu@localhost:5432/test?sslmode=disable"
)

var pgdb *sql.DB
var db *bun.DB

// ! Func pour la structure du UserManagementServer type
func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
	}
}

// ! Implémentation du service grpc
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer // Ce connecte au server grpc
}

// ! Méthode CreateNewuser
func (s *UserManagementServer) CreateNewuser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}

	// ? INSERT
	db.NewInsert().
		Model(created_user).
		Exec(ctx)

	return created_user, nil
}

// ! Méthode Getusers
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var users_list *pb.UserList = &pb.UserList{}

	// ? SELECT
	db.NewSelect().
		Model(users_list).
		Exists(ctx)

	return users_list, nil
}

// ! Fonction qui vas run le code, ce qui nous permet de simplifier le code
func (s *UserManagementServer) Run() error {
	// ? Initialiser l'écoute au port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// ? Créer un server
	server := grpc.NewServer()
	// ? Enregistrer le server
	pb.RegisterUserManagementServer(server, s)
	log.Printf("Server listening at %v", lis.Addr())

	return server.Serve(lis)
}

func main() {
	ctx := context.Background()
	pgdb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(pgdb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	_, err := db.NewCreateTable().Model((*pb.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// ? Instancier un nouveau UserManagementServer
	var user_server *UserManagementServer = NewUserManagementServer()
	// ? Appeler la fonction Run avec la variable au dessus 
	if err := user_server.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}