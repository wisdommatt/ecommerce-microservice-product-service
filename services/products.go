package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
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
	natsConn          *nats.Conn
}

// NewProductService returns a new product service object.
func NewProductService(
	productRepo products.Repository, userServiceClient proto.UserServiceClient, natsConn *nats.Conn,
) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo:       productRepo,
		userServiceClient: userServiceClient,
		natsConn:          natsConn,
	}
}

func (s *ProductServiceImpl) AddProduct(ctx context.Context, jwtToken string, newProduct *products.Product) (*products.Product, error) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span = opentracing.StartSpan("service.GetUsers")
	}
	userResponse, err := s.userServiceClient.GetUserFromJWT(ctx, &proto.GetUserFromJWTInput{JwtToken: jwtToken})
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("retrieving merchant details from jwt"))
		return nil, errors.New("you are not authenticated")
	}
	span.SetTag("merchant", userResponse.User)
	newProduct.MerchantID = userResponse.User.Id
	err = s.productRepo.SaveProduct(ctx, newProduct)
	if err != nil {
		return nil, errors.New("an error occured while adding product, please try again later")
	}
	s.publishProductAddedEmailEvent(span, userResponse.User.Email, newProduct)
	return newProduct, nil
}

func (s *ProductServiceImpl) publishProductAddedEmailEvent(span opentracing.Span, userEmail string, product *products.Product) {
	span = opentracing.StartSpan("publish-product-added-email-event", opentracing.ChildOf(span.Context()))
	defer span.Finish()
	natsMessage := map[string]interface{}{
		"to":      userEmail,
		"subject": "Product added successfully",
		"parameters": map[string]string{
			"productName":        product.Name,
			"productImageUrl":    product.ImageURL,
			"productCategory":    product.Category,
			"productPrice":       fmt.Sprintf("%.2f", product.Price),
			"productDescription": product.Description,
		},
	}
	natsMessageJSON, err := json.Marshal(natsMessage)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("converting object to json"), log.Object("object", natsMessage))
		return
	}
	span.LogFields(log.String("nats.message", string(natsMessageJSON)))
	err = s.natsConn.Publish("notification.SendProductAddedEmail", natsMessageJSON)
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("nats.notification.SendProductAddedEmail"))
	}
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
