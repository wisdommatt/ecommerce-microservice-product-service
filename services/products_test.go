package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
	"github.com/wisdommatt/ecommerce-microservice-product-service/mocks"
)

func TestProductServiceImpl_AddProduct(t *testing.T) {
	productRepo := &mocks.Repository{}
	productRepo.On("SaveProduct", mock.Anything, &products.Product{
		Sku:        "123456",
		Name:       "Product 1",
		Price:      10000,
		MerchantID: "valid.user",
	}).Return(errors.New("an error occured"))

	productRepo.On("SaveProduct", mock.Anything, &products.Product{
		Sku:        "123456",
		Name:       "Product 2",
		Price:      15000,
		MerchantID: "valid.user",
	}).Return(nil)

	userServiceClient := &mocks.UserServiceClient{}
	userServiceClient.On("GetUserFromJWT", mock.Anything, &proto.GetUserFromJWTInput{JwtToken: "invalidJwt"}).
		Return(nil, errors.New("invalid jwt"))

	userServiceClient.On("GetUserFromJWT", mock.Anything, &proto.GetUserFromJWTInput{JwtToken: "validJwt"}).
		Return(&proto.GetUserFromJWTResponse{
			User: &proto.User{
				Id:       "valid.user",
				FullName: "Valid User",
			},
		}, nil)

	type args struct {
		jwtToken   string
		newProduct *products.Product
	}
	tests := []struct {
		name    string
		args    args
		want    *products.Product
		wantErr bool
	}{
		{
			name:    "invalid jwt token",
			args:    args{jwtToken: "invalidJwt", newProduct: nil},
			wantErr: true,
		},
		{
			name: "SaveProduct repo implementation with error",
			args: args{jwtToken: "validJwt", newProduct: &products.Product{
				Sku:   "123456",
				Name:  "Product 1",
				Price: 10000,
			}},
			wantErr: true,
		},
		{
			name: "SaveProduct repo implementation without error",
			args: args{jwtToken: "validJwt", newProduct: &products.Product{
				Sku:   "123456",
				Name:  "Product 2",
				Price: 15000,
			}},
			want: &products.Product{
				Sku:        "123456",
				Name:       "Product 2",
				Price:      15000,
				MerchantID: "valid.user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(productRepo, userServiceClient)
			got, err := s.AddProduct(context.Background(), tt.args.jwtToken, tt.args.newProduct)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductServiceImpl.AddProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductServiceImpl.AddProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductServiceImpl_GetProduct(t *testing.T) {
	productRepo := &mocks.Repository{}
	productRepo.On("GetProductBySKU", mock.Anything, "sku.111222").Return(nil, errors.New("an erorr occured"))
	productRepo.On("GetProductBySKU", mock.Anything, "sku.222333").Return(&products.Product{
		Name:  "Apple Watch",
		Price: 1999288,
	}, nil)

	type args struct {
		sku string
	}
	tests := []struct {
		name    string
		args    args
		want    *products.Product
		wantErr bool
	}{
		{
			name:    "empty sku",
			args:    args{sku: ""},
			wantErr: true,
		},
		{
			name:    "GetProductBySKU repository implementation with error",
			args:    args{sku: "sku.111222"},
			wantErr: true,
		},
		{
			name: "GetProductBySKU repository implementation without error",
			args: args{sku: "sku.222333"},
			want: &products.Product{Name: "Apple Watch", Price: 1999288},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(productRepo, nil)
			got, err := s.GetProduct(context.Background(), tt.args.sku)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductServiceImpl.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductServiceImpl.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
