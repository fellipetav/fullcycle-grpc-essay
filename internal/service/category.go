package service

import (
	"context"
	"io"
	"github.com/devfullcycle/14-gRPC/internal/database"
	"github.com/devfullcycle/14-gRPC/internal/pb"
)

type CategoryService struct{
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

// método de criação do database.Category
func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

// para criar a Categoria eu preciso de um método com uma assinatura exata para eu implementar minha interface
func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category {
		Id: 			category.ID,
		Name:			category.Name,
		Description: 	category.Description,
	}

	return &pb.CategoryResponse { // eu poderia retornar uma Category se eu quisesse
		Category: categoryResponse,
	}, nil
}

// método para Listar as categorias usando o método [FindAll()] do /database/category.go
func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var categoriesResponse []*pb.Category

	for _, category := range categories {
		categoryResponse := &pb.Category {
			Id:				category.ID,
			Name:			category.Name,
			Description:  	category.Description,
		}

		categoriesResponse = append(categoriesResponse, categoryResponse)
	}

	return &pb.CategoryList{Categories: categoriesResponse}, nil
}

// GetCategory returno uma categoria específica inteira em função do id fornecido na request.
//
// ctx - O objeto de contexto.
// in - O objeto GetCategoryRequest que contém os parâmetros de entrada (o Id).
// *pb.CategoryResponse - O objeto de resposta que contém as informações da categoria.
// erro - Um objeto de erro que será nulo caso tenhamos sucesso.
func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Find(in.Id)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category {
		Id:				category.ID,
		Name:			category.Name,
		Description:  	category.Description,
	}
	
	return 	categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	// este loop se encarregará de mandar o stream a todo o tempo
	for {
		// recebendo o stream de dados
		category, err := stream.Recv()
		// se chegou no final e não tem mais nada para mandar, envia-nos todos os dados e fecha
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		// se encontrar erro, retorna o erro esai
		if err != nil {
			return err
		}

		// se não saiu, retorne a categoria e continue
		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		// se encontrar erro, ele sai
		if err != nil {
			return err
		}

		// darei um append na categoria recém-criada para encher a nossa lista
		categories.Categories = append(categories.Categories, &pb.Category {
			Id: 			categoryResult.ID,
			Name: 			categoryResult.Name,
			Description: 	categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		// se esse dado chegar no final, ele sai
		if err == io.EOF {
			return nil
		}
		// se tiver algum erro, ele retorna o erro
		if err != nil {
			return err
		}

		// crio minha categoria
		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		// envio o dado de resposta via stream conforme vou recebendo a chamada (o dado) via stream
		err = stream.Send(&pb.Category {
			Id: 			categoryResult.ID,
			Name: 			categoryResult.Name,
			Description: 	categoryResult.Description,
		})
		if err != nil {
			return err
		}
	}
}