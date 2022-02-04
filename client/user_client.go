package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Titouan-Esc/golang-grpc/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// ? Connection au server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	// ? Fermer la connection
	defer conn.Close()

	// ? Définir un nouveau client
	client := pb.NewUserManagementClient(conn)

	// ? Définir une nouveux context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// ? Créer un map de users
	var new_users = make(map[string]int32)
	new_users["Gege"] = 62
	new_users["Titouan"] = 21

	// ? Créer un boucle pour créer les user avec le CreateNewuser
	for name, age := range new_users {
		res, err := client.CreateNewuser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("Could not create user: %v", err)
		}
		log.Printf(`User details:
		Name: %s
		Age: %d
		Id: %d`, res.GetName(), res.GetAge(), res.GetId())
	}
}