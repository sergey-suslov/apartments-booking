package main

import (
	"booking/pkg/booking"
	"context"
	"flag"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"
)

func main() {
	fs := flag.NewFlagSet("bookingcli", flag.ExitOnError)
	var (
		port     = fs.String("port", "50052", "Port of Booking service")
		mongoURI = fs.String("mongo", "mongodb://user:password@localhost:27017/booking", "MongoDB connection string mongodb://...")
		help     = fs.Bool("h", false, "Show help")
		test     = fs.Bool("test", false, "Show help")
		logDebug = fs.Bool("debug", false, "Log debug info")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	_ = fs.Parse(os.Args[1:])
	if *help {
		fs.Usage()
		os.Exit(1)
	}

	mc, closeConn := connectMongo(*mongoURI)

	if *test {
		// do some test thing
	}

	logConfig := zap.NewProductionConfig()
	if *logDebug {
		logConfig.Level.SetLevel(zap.DebugLevel)
	}
	logger, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	repository := booking.NewRepository(mc.Database("booking"))
	service := booking.NewService(repository)
	service = booking.NewLoggingService(logger, service)

	fieldKeys := []string{"method"}
	service = booking.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "booking_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "booking_service",
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
	mux.Handle("/reservations", booking.MakeHttpHandler(service, httpLogger))

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

func connectMongo(uri string) (*mongo.Client, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client, func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
}
