package serviceservers

import (
	"context"
	"errors"
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
