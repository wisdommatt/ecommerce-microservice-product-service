package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
	productmock "github.com/wisdommatt/ecommerce-microservice-product-service/test/mocks/products"
)

// test will be resolved when mocks are replaced with the one generated
// by mockery.

// func TestProductServiceImpl_AddProduct(t *testing.T) {
// 	type args struct {
// 		newProduct *products.Product
// 	}
// 	tests := []struct {
// 		name                string
// 		args                args
// 		repoSaveProductFunc func(ctx context.Context, product *products.Product) error
// 		wantErr             bool
// 	}{
// 		{
// 			name: "SaveProduct repository implementation with error",
// 			args: args{newProduct: &products.Product{
// 				Name:        "Yellow Speaker",
// 				Description: "Speaker with yellow sounds",
// 			}},
// 			repoSaveProductFunc: func(ctx context.Context, product *products.Product) error {
// 				return errors.New("an error occured while saving product")
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "SaveProduct repository implementation without error",
// 			args: args{newProduct: &products.Product{
// 				Name:        "Green Speaker",
// 				Description: "Speaker with green charger",
// 				Price:       10009999,
// 			}},
// 			repoSaveProductFunc: func(ctx context.Context, product *products.Product) error {
// 				return nil
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewProductService(&productmock.RepositoryMock{
// 				SaveProductFunc: tt.repoSaveProductFunc,
// 			})
// 			_, err := s.AddProduct(context.TODO(), tt.args.newProduct)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ProductServiceImpl.AddProduct() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

func TestProductServiceImpl_GetProduct(t *testing.T) {
	type args struct {
		ctx context.Context
		sku string
	}
	tests := []struct {
		name                    string
		args                    args
		repoGetProductBySKUFunc func(ctx context.Context, sku string) (*products.Product, error)
		want                    *products.Product
		wantErr                 bool
	}{
		{
			name:    "empty sku",
			args:    args{sku: ""},
			wantErr: true,
		},
		{
			name: "GetProductBySKU repository implementation with error",
			args: args{sku: "sku.111222"},
			repoGetProductBySKUFunc: func(ctx context.Context, sku string) (*products.Product, error) {
				return nil, errors.New("an error occured while finding product")
			},
			wantErr: true,
		},
		{
			name: "GetProductBySKU repository implementation without error",
			args: args{sku: "sku.222333"},
			repoGetProductBySKUFunc: func(ctx context.Context, sku string) (*products.Product, error) {
				return &products.Product{Name: "Apple Watch", Price: 1999288}, nil
			},
			want: &products.Product{Name: "Apple Watch", Price: 1999288},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(&productmock.RepositoryMock{
				GetProductBySKUFunc: tt.repoGetProductBySKUFunc,
			}, nil)
			got, err := s.GetProduct(tt.args.ctx, tt.args.sku)
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
