package server

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

// ? Implémentation du service grpc
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

// ? Méthode CreateNewuser
func (s *UserManagementServer) CreateNewuser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())

	// ? Création de l'id de l'user
	var user_id int32 = int32(rand.Int31n(1000))
	// ? Retourner un user avec la référence du service protobuf
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

func main() {
	// ? Initialiser l'écoute au port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// ? Créer un server
	server := grpc.NewServer()
	// ? Enregistrer le server
	pb.RegisterUserManagementServer(server, &UserManagementServer{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}