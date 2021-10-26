package serviceservers

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
	productmock "github.com/wisdommatt/ecommerce-microservice-product-service/test/mocks/products"
)

func TestProductServer_AddProduct(t *testing.T) {
	type args struct {
		req *proto.NewProduct
	}
	tests := []struct {
		name                  string
		args                  args
		serviceAddProductFunc func(ctx context.Context, newProduct *products.Product) (*products.Product, error)
		wantErr               bool
	}{
		{
			name: "AddProduct service implementation with error",
			args: args{req: &proto.NewProduct{
				Name:  "Bluetooth Headphone",
				Price: 19929444,
			}},
			serviceAddProductFunc: func(ctx context.Context, newProduct *products.Product) (*products.Product, error) {
				return nil, errors.New("invalid product")
			},
			wantErr: true,
		},
		{
			name: "AddProduct service implementation without error",
			args: args{req: &proto.NewProduct{
				Name:  "Bluetooth Headphone",
				Price: 19929444,
			}},
			serviceAddProductFunc: func(ctx context.Context, newProduct *products.Product) (*products.Product, error) {
				return &products.Product{Name: "Bluetooth Headphone"}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductServer(&productmock.ServiceMock{
				AddProductFunc: tt.serviceAddProductFunc,
			})
			_, err := s.AddProduct(context.TODO(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductServer.AddProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
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
