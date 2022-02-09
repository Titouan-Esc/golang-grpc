package main

import (
	"context"
	"database/sql"
	"log"
	"net"

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
)

// ! Fonc pour la structure du UserManagementServer type
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

	return created_user, nil
}

// ! Méthode Getusers
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var users_list *pb.UserList = &pb.UserList{}

	return users_list, nil
}

// ! Fonction qui vas run le code, ce qui nous permet de simplifier le code
func (s *UserManagementServer) Run() error {
	ctx := context.Background()

	dsn := "postgres://titouanescorneboueu@localhost:5432/test?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(pgdb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	_, err := db.NewCreateTable().Model((*pb.User)(nil)).Exec(ctx)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}


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
	/*
		C'est la seule chose que nous avons besoin d'avoir dans la func main
	*/
	// ? Instancier un nouveau UserManagementServer
	var user_server *UserManagementServer = NewUserManagementServer()
	// ? Appeler la fonction Run avec la variable au dessus 
	if err := user_server.Run(); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}