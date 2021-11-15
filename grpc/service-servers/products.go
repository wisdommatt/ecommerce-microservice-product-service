package serviceservers

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	"github.com/wisdommatt/ecommerce-microservice-product-service/services"
	"google.golang.org/grpc/metadata"
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
	span, _ := opentracing.StartSpanFromContext(ctx, "AddProduct")
	defer span.Finish()
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("request.body", req)

	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(errors.New("no meta data in grpc context")))
		return nil, errors.New("no metadata sent, please try again later")
	}
	jwtToken := extractAuthorizationFromMetaData(metaData)
	if jwtToken == "" {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(errors.New("no authorization token in metadata")))
		return nil, errors.New("no authorization token found in metadata")
	}
	ctx = opentracing.ContextWithSpan(ctx, span)
	newProduct, err := s.productService.AddProduct(ctx, jwtToken, ProtoNewProductToInternal(req))
	if err != nil {
		return nil, err
	}
	return InternalProductToProto(newProduct), nil
}

func (s *ProductServer) GetProduct(ctx context.Context, input *proto.GetProductInput) (*proto.Product, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetProduct")
	defer span.Finish()
	ext.SpanKindRPCServer.Set(span)
	span.SetTag("param.input", input)

	ctx = opentracing.ContextWithSpan(ctx, span)
	product, err := s.productService.GetProduct(ctx, input.Sku)
	if err != nil {
		return nil, err
	}
	return InternalProductToProto(product), nil
}

func extractAuthorizationFromMetaData(md metadata.MD) string {
	values := md.Get("Authorization")
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
