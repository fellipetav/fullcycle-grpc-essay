package main

import (
	"database/sql"
	"net"

	"github.com/devfullcycle/14-gRPC/internal/database"
	"github.com/devfullcycle/14-gRPC/internal/pb"
	"github.com/devfullcycle/14-gRPC/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

// minha função de entrypoint (função principal)
func main() {
	// abrindo a conexão com o banco de dados
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	// depois de estar no banco de dados, agora preciso criar o meu service (nosso servidor)
	categoryService := service.NewCategoryService(*categoryDb)

	grpcServer := grpc.NewServer()
	// atachando nosso serviço no nosso gRPC server
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	
	// vamos usar um client para reflection
	reflection.Register(grpcServer)

	// preciso abrir uma conexão TCP para eu conseguir falar com o gRPC
	lis, err := net.Listen("tcp", ":50052") // a porta padrão do gRPC é 50051
	if err != nil {
		panic(err)
	}

	// caso eu não consiga inicializar o servidor, traga um erro
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
