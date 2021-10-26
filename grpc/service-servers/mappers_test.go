package serviceservers

import (
	"reflect"
	"testing"

	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

func TestProtoNewProductToInternal(t *testing.T) {
	type args struct {
		newProduct *proto.NewProduct
	}
	tests := []struct {
		name string
		args args
		want *products.Product
	}{
		{
			name: "incomplete fields",
			args: args{newProduct: &proto.NewProduct{
				Name:        "Leather Shoe",
				Description: "Leather strong shoe",
			}},
			want: &products.Product{
				Name:        "Leather Shoe",
				Description: "Leather strong shoe",
			},
		},
		{
			name: "complete fields",
			args: args{newProduct: &proto.NewProduct{
				Name:        "Pink Slippers",
				Description: "Cute pink slippers",
				Category:    "slippers",
				Brand:       "Nike",
				Price:       100000,
			}},
			want: &products.Product{
				Name:        "Pink Slippers",
				Description: "Cute pink slippers",
				Category:    "slippers",
				Brand:       "Nike",
				Price:       100000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoNewProductToInternal(tt.args.newProduct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoNewProductToInternal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalProductToProto(t *testing.T) {
	type args struct {
		product *products.Product
	}
	tests := []struct {
		name string
		args args
		want *proto.Product
	}{
		{
			name: "incomplete fields",
			args: args{product: &products.Product{
				Name: "White Canvas",
				Sku:  "whiteCanvas.11234",
			}},
			want: &proto.Product{
				Name: "White Canvas",
				Sku:  "whiteCanvas.11234",
			},
		},
		{
			name: "complete fields",
			args: args{product: &products.Product{
				Name:        "Blue Canvas",
				Sku:         "blue.canvas.123",
				Description: "This is a nice blue canvas",
				Category:    "Fashion",
				Brand:       "Addidas",
				Price:       9009999,
			}},
			want: &proto.Product{
				Name:        "Blue Canvas",
				Sku:         "blue.canvas.123",
				Description: "This is a nice blue canvas",
				Category:    "Fashion",
				Brand:       "Addidas",
				Price:       9009999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InternalProductToProto(tt.args.product); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InternalProductToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}
