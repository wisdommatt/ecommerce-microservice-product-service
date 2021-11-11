package serviceservers

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
	"github.com/wisdommatt/ecommerce-microservice-product-service/mocks"
	productmock "github.com/wisdommatt/ecommerce-microservice-product-service/test/mocks/products"
	"google.golang.org/grpc/metadata"
)

func TestProductServer_AddProduct(t *testing.T) {
	productService := &mocks.ProductService{}
	productService.On("AddProduct", mock.Anything, "jwtToken", ProtoNewProductToInternal(&proto.NewProduct{
		Name:  "Product 1",
		Price: 10000,
	})).Return(nil, errors.New("an error occured"))

	productService.On("AddProduct", mock.Anything, "jwtToken", ProtoNewProductToInternal(&proto.NewProduct{
		Name:  "Product 2",
		Price: 20000,
	})).Return(&products.Product{
		Sku:   "sku.123",
		Name:  "Product 2",
		Price: 20000,
	}, nil)

	authMetaData := metadata.New(map[string]string{
		"Authorization": "jwtToken",
	})
	ctxWithMetadata := metadata.NewIncomingContext(context.TODO(), authMetaData)

	type args struct {
		ctx context.Context
		req *proto.NewProduct
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.Product
		wantErr bool
	}{
		{
			name:    "request without metadata",
			args:    args{ctx: context.Background(), req: nil},
			wantErr: true,
		},
		{
			name: "request metadata without authorization",
			args: args{
				ctx: metadata.NewIncomingContext(context.TODO(), metadata.New(map[string]string{})),
			},
			wantErr: true,
		},
		{
			name: "AddProduct service implementation with error",
			args: args{ctx: ctxWithMetadata, req: &proto.NewProduct{
				Name:  "Product 1",
				Price: 10000,
			}},
			wantErr: true,
		},
		{
			name: "AddProduct service implementation with error",
			args: args{ctx: ctxWithMetadata, req: &proto.NewProduct{
				Name:  "Product 2",
				Price: 20000,
			}},
			want: &proto.Product{
				Sku:   "sku.123",
				Name:  "Product 2",
				Price: 20000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductServer(productService)
			got, err := s.AddProduct(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductServer.AddProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductServer.AddProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductServer_GetProduct(t *testing.T) {
	type args struct {
		input *proto.GetProductInput
	}
	tests := []struct {
		name                  string
		args                  args
		serviceGetProductFunc func(ctx context.Context, sku string) (*products.Product, error)
		want                  *proto.Product
		wantErr               bool
	}{
		{
			name: "GetProduct service implementation with error",
			args: args{input: &proto.GetProductInput{Sku: "sku.eee"}},
			serviceGetProductFunc: func(ctx context.Context, sku string) (*products.Product, error) {
				return nil, errors.New("an error occured")
			},
			wantErr: true,
		},
		{
			name: "GetProduct service implementation without error",
			args: args{input: &proto.GetProductInput{Sku: "sku.eee"}},
			serviceGetProductFunc: func(ctx context.Context, sku string) (*products.Product, error) {
				return &products.Product{Name: "HP 2223", Description: "Slim HP laptop"}, nil
			},
			want: &proto.Product{Name: "HP 2223", Description: "Slim HP laptop"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductServer(&productmock.ServiceMock{
				GetProductFunc: tt.serviceGetProductFunc,
			})
			got, err := s.GetProduct(context.TODO(), tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductServer.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductServer.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
