package services

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

// ProductService is the interface that describes a product service.
type ProductService interface {
	AddProduct(ctx context.Context, jwtToken string, newProduct *products.Product) (*products.Product, error)
	GetProduct(ctx context.Context, sku string) (*products.Product, error)
}

// ProductServiceImpl is the default implementation for ProductService
// interface.
type ProductServiceImpl struct {
	productRepo       products.Repository
	userServiceClient proto.UserServiceClient
}

// NewProductService returns a new product service object.
func NewProductService(productRepo products.Repository, userServiceClient proto.UserServiceClient) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo:       productRepo,
		userServiceClient: userServiceClient,
	}
}

func (s *ProductServiceImpl) AddProduct(ctx context.Context, jwtToken string, newProduct *products.Product) (*products.Product, error) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span = opentracing.StartSpan("service.GetUsers")
	}
	response, err := s.userServiceClient.GetUserFromJWT(ctx, &proto.GetUserFromJWTInput{JwtToken: jwtToken})
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("retrieving merchant details from jwt"))
		return nil, errors.New("you are not authenticated")
	}
	span.SetTag("merchant", response.User)
	newProduct.MerchantID = response.User.Id
	err = s.productRepo.SaveProduct(ctx, newProduct)
	if err != nil {
		return nil, errors.New("an error occured while adding product, please try again later")
	}
	return newProduct, nil
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context, sku string) (*products.Product, error) {
	if sku == "" {
		return nil, errors.New("sku must be provided")
	}
	product, err := s.productRepo.GetProductBySKU(ctx, sku)
	if err != nil {
		return nil, errors.New("product does not exist")
	}
	return product, nil
}
