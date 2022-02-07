package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/Titouan-Esc/golang-grpc/proto"
	"google.golang.org/grpc"
)

// ? Définir le port du server
const (
	port = ":50051"
)

// ! Fonc pour la structure du UserManagementServer type
func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UserList{},
	}
}

// ! Implémentation du service grpc
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer // Ce connecte au server grpc

	user_list *pb.UserList // Variable pour la liste des users
}

// ! Méthode CreateNewuser
func (s *UserManagementServer) CreateNewuser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())

	// ? Création de l'id de l'user
	var user_id int32 = int32(rand.Int31n(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	s.user_list.Users = append(s.user_list.Users, created_user)
	// ? Retourner un user avec la référence du service protobuf
	return created_user, nil
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

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return s.user_list, nil
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