package products

import (
	"time"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

// Repository is the interface that describes a product repository
// object.
type Repository interface {
	SaveProduct(ctx context.Context, product *Product) error
	GetProductBySKU(ctx context.Context, sku string) (*Product, error)
}

// ProductRepo is the default implementation for Repository inteface.
type ProductRepo struct {
	db     *gorm.DB
	tracer opentracing.Tracer
}

// NewRepository returns a new product repository object.
func NewRepository(db *gorm.DB, tracer opentracing.Tracer) *ProductRepo {
	return &ProductRepo{
		db:     db,
		tracer: tracer,
	}
}

func (r *ProductRepo) setMySqlComponentTags(span opentracing.Span, tableName string) {
	ext.DBInstance.Set(span, tableName)
	ext.DBType.Set(span, "mysql")
	ext.SpanKindRPCClient.Set(span)
}

// SaveProduct saves a new product to the database.
func (r *ProductRepo) SaveProduct(ctx context.Context, product *Product) error {
	product.Sku = uuid.NewString()
	product.TimeAdded = time.Now()

	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "SaveProduct")
	defer span.Finish()
	r.setMySqlComponentTags(span, "products")
	span.LogFields(
		log.Object("param.product", product),
	)

	err := r.db.Create(product).Error
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("gorm.db.Create"))
		return err
	}
	return nil
}

func (r *ProductRepo) GetProductBySKU(ctx context.Context, sku string) (*Product, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, r.tracer, "GetProductBySKU")
	defer span.Finish()
	r.setMySqlComponentTags(span, "products")
	span.SetTag("param.sku", sku)

	product := &Product{}
	err := r.db.Where("sku = ?", sku).First(product).Error
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err), log.Event("gorm.db.Where.First"))
		return nil, err
	}
	span.SetTag("response.product", product)
	return product, nil
}
