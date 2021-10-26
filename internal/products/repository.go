package products

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
)

// Repository is the interface that describes a product repository
// object.
type Repository interface{}

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
