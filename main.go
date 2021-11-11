package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/wisdommatt/ecommerce-microservice-product-service/grpc/proto"
	servers "github.com/wisdommatt/ecommerce-microservice-product-service/grpc/service-servers"
	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
	"github.com/wisdommatt/ecommerce-microservice-product-service/services"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	mustLoadDotenv(log)

	db, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_CONNECTION")), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("failed to connect database", os.Getenv("MYSQL_CONNECTION"))
	}
	db.AutoMigrate(&products.Product{})

	natsConn, err := nats.Connect(os.Getenv("NATS_URI"))
	if err != nil {
		log.WithField("nats_uri", os.Getenv("NATS_URI")).WithError(err).
			Error("an error occured while connecting to nats")
	}
	defer natsConn.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "2424"
	}
	tracer := initTracer("product-service")
	opentracing.SetGlobalTracer(tracer)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.WithError(err).WithField("port", port).Fatal("an error occured while listening to tcp conn")
	}

	userServiceConn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDR"), grpc.WithInsecure())
	if err != nil {
		log.WithField("userServiceAddr", os.Getenv("USER_SERVICE_ADDR")).WithError(err).
			Fatal("an error occured while connecting to user service")
	}
	userServiceClient := proto.NewUserServiceClient(userServiceConn)
	productRepo := products.NewRepository(db, initTracer("mysql"))
	productService := services.NewProductService(productRepo, userServiceClient, natsConn)

	grpcServer := grpc.NewServer()
	proto.RegisterProductServiceServer(grpcServer, servers.NewProductServer(productService))
	log.WithField("port", port).Info("app running")
	grpcServer.Serve(lis)
}

func mustLoadDotenv(log *logrus.Logger) {
	err := godotenv.Load(".env", ".env-defaults")
	if err != nil {
		log.WithError(err).Fatal("Unable to load env files")
	}
}

func initTracer(serviceName string) opentracing.Tracer {
	return initJaegerTracer(serviceName)
}

func initJaegerTracer(serviceName string) opentracing.Tracer {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}
	tracer, _, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatal("ERROR: cannot init Jaeger", err)
	}
	return tracer
}
