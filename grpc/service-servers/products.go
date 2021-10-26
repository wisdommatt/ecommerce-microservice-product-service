package serviceservers

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/services"
)

type ProductServer struct {
	proto.UnimplementedProductServiceServer
	productService services.ProductService
}

// NewProductServer returns a new product server object.
func NewProductServer(productService services.ProductService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

func (s *ProductServer) AddProduct(ctx context.Context, req *proto.NewProduct) (*proto.Product, error) {
	span := opentracing.StartSpan("AddProduct")
	defer span.Finish()
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("request.body", req)

	ctx = opentracing.ContextWithSpan(ctx, span)
	newProduct, err := s.productService.AddProduct(ctx, ProtoNewProductToInternal(req))
	if err != nil {
		return nil, err
	}
	return InternalProductToProto(newProduct), nil
}

func (s *ProductServer) GetProduct(ctx context.Context, input *proto.GetProductInput) (*proto.Product, error) {
	span := opentracing.StartSpan("GetProduct")
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	product, err := s.productService.GetProduct(ctx, input.Sku)
	if err != nil {
		return nil, err
	}
	return InternalProductToProto(product), nil
}
