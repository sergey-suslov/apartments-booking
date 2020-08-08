package main

import (
	"apartments/pkg/apartments"
	"context"
	"flag"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"
)

func main() {
	fs := flag.NewFlagSet("addcli", flag.ExitOnError)
	var (
		port = fs.String("port", "50051", "Port of Apartments service")
		help = fs.Bool("h", false, "Show help")
		test = fs.Bool("test", false, "Show help")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	_ = fs.Parse(os.Args[1:])
	if *help {
		fs.Usage()
		os.Exit(1)
	}

	mc, closeConn := connectMongo()

	if *test {
		createTestApartments(mc.Database("apartments"))
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	repository := apartments.NewRepository(mc.Database("apartments"))
	service := apartments.NewService(repository)
	service = apartments.NewLoggingService(logger, service)

	fieldKeys := []string{"method"}
	service = apartments.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "apartments_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "apartments_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		service)
	defer func() {
		_ = logger.Sync()
		closeConn()
	}()

	mux := http.NewServeMux()

	httpLogger := kitlog.With(kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr)), "component", "http")
	mux.Handle("/apartments", apartments.MakeHttpHandler(service, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Info("listening", zap.String("port", *port))
		errs <- http.ListenAndServe(":"+*port, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Info("terminated", zap.Error(<-errs))
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		_ = w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func connectMongo() (*mongo.Client, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://user:password@localhost:27017/apartments"))
	if err != nil {
		panic(err)
	}

	return client, func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
}

func createTestApartments(mc *mongo.Database) {
	cities := []string{"Dublin", "Munich", "London"}
	for i := 1; i < 5; i++ {
		city := cities[rand.Intn(len(cities))]
		_, _ = mc.Collection("apartments").InsertOne(context.Background(), bson.M{
			"title":   fmt.Sprintf("Test apartment from %s %d", city, i+1),
			"address": "Dublin, somewhere st. 25",
			"owner":   "Mike",
			"city":    city,
		})
	}
}
