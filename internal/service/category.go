package service

import (
	"context"
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