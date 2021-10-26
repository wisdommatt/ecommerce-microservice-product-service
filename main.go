package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
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

	_, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_CONNECTION")), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("failed to connect database", os.Getenv("MYSQL_CONNECTION"))
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "2424"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.WithError(err).WithField("port", port).Fatal("an error occured while listening to tcp conn")
	}
	grpcServer := grpc.NewServer()
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
