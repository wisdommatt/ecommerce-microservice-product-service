package serviceservers

import (
	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

func ProtoNewProductToInternal(newProduct *proto.NewProduct) *products.Product {
	return &products.Product{
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Category:    newProduct.Category,
		Brand:       newProduct.Brand,
		Price:       newProduct.Price,
		ImageURL:    newProduct.ImageUrl,
	}
}

func InternalProductToProto(product *products.Product) *proto.Product {
	return &proto.Product{
		Sku:         product.Sku,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Brand:       product.Brand,
		Price:       product.Price,
		ImageUrl:    product.ImageURL,
		MerchantId:  product.MerchantID,
	}
}
