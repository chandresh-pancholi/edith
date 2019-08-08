package api

import (
	"context"
	"fmt"
	"github.com/chandresh-pancholi/edith/api/client"
	"github.com/chandresh-pancholi/edith/api/event"
	"github.com/chandresh-pancholi/edith/pkg/elasticsearch"
	es "github.com/chandresh-pancholi/edith/pkg/elasticsearch/config"
	"github.com/chandresh-pancholi/edith/pkg/kafka/admin"
	"github.com/chandresh-pancholi/edith/pkg/kafka/consumer"
	"github.com/chandresh-pancholi/edith/pkg/kafka/producer"
	"github.com/chandresh-pancholi/edith/pkg/mysql"
	"log"
	"net/http"
	"net/http/pprof"
	"os"


	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"

	"github.com/spf13/viper"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// Server info required to run MSG as backend service
type Server struct {
	logger      *zap.Logger
	port        string
	router      *http.ServeMux
	producer    *producer.Producer
	adminClient *admin.AdminClient
	esClient    *es.ESClient
	db          *mysql.DB
}

// NewServer . create new MSG server struct
func NewServer() *Server {
	s := new(Server)

	s.router = newRouter()
	if viper.GetBool("zap.log.enabled") {
		s.logger = s.enableLogger()
	} else {
		s.logger, _ = zap.NewProduction()
	}
	s.port = viper.GetString("port")

	p, err := producer.NewProducer(s.logger)
	if err != nil {
		s.logger.Error("exception while creating kafka producer ", zap.Error(err))
	}

	s.producer = p

	s.adminClient = admin.NewAdminClient(s.logger)
	db, err := mysql.NewDB()
	if err != nil {
		s.logger.Error("mysql connection failed ", zap.Error(err))
		os.Exit(1)
	}
	s.db = db
	defer s.db.Close()

	s.esClient, err = es.NewESClient()
	if err != nil {
		s.logger.Error("ES client connection failed", zap.Error(err))
	}
	return s
}

//TODO Adding close methods for Producer, Consumer, DB

// Close method close kafka admin client
func (s *Server) Close() {
	s.adminClient.Close()
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

// Run method is to run MSG application
func (s *Server) Run() {

	defer s.logger.Sync()

	s.observability()

	//s.router.Handle("/metrics", promhttp.Handler())
	s.routerHandler()

	h := &ochttp.Handler{Handler: s.router}
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		log.Fatal("Failed to register ochttp.DefaultServerViews")
	}
	s.initialiseConsumer()
	s.logger.Info("MSG is listening at port", zap.String("port", s.port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), h)
	if err != nil {
		s.logger.Error("MSG starting failed ", zap.String("port", s.port), zap.Error(err))
		return
	}

}
func (s *Server) routerHandler() {
	s.eventHandler()
	s.clientHandler()
	s.profiler()
}

func (s *Server) eventHandler() {
	eventHandler := event.NewEventHandler(s.producer, s.logger)

	s.router.HandleFunc("/event", eventHandler.Emit)
}

func (s *Server) clientHandler() {
	clientHandler := client.NewClientHandler(s.db, s.logger, s.adminClient, s.esClient)

	s.router.HandleFunc("/client", clientHandler.CreateClient)
	s.router.HandleFunc("/client/get", clientHandler.GetClient)
	s.router.HandleFunc("/client/stop", clientHandler.StopClient)
}

func (s *Server) profiler() {
	s.router.HandleFunc("/debug/pprof", pprof.Index)
	s.router.HandleFunc("/debug/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	s.router.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	s.router.Handle("/debug/heap", pprof.Handler("heap"))
	s.router.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	s.router.Handle("/debug//block", pprof.Handler("block"))
	s.router.Handle("/debug/allocs", pprof.Handler("allocs"))
	s.router.Handle("/debug/mutex", pprof.Handler("mutex"))
}

func (s *Server) observability() {
	s.logger.Info("enabling observability in postpaid jobs")

	jaegerServiceName := viper.GetString("jaeger.serviceName")
	collectorEndPoint := fmt.Sprintf("%s://%s:%d%s", "http", viper.GetString("jaeger.collector.host"), viper.GetInt("jaeger.collector.port"), viper.GetString("jaeger.collector.path"))
	agentEndPoint := fmt.Sprintf("%s:%d", viper.GetString("jaeger.agent.host"), viper.GetInt("jaeger.agent.port"))

	s.logger.Info("jaeger exporter", zap.String("service_name", jaegerServiceName), zap.String("collector_end_point", collectorEndPoint), zap.String("agent_end_point", agentEndPoint))
	je, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: collectorEndPoint,
		AgentEndpoint:     agentEndPoint,
		Process: jaeger.Process{
			ServiceName: jaegerServiceName,
		},
	})
	if err != nil {
		s.logger.Error("jaeger exporter initialisation failed", zap.Error(err))
	}

	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: viper.GetString("prometheus.namespace"),
	})
	if err != nil {
		s.logger.Error("failed to create Prometheus exporter", zap.Error(err))
	}

	s.logger.Info("registering prometheus exporter")
	view.RegisterExporter(pe)

	s.router.Handle("/metrics", pe)
}

func (s *Server) initialiseConsumer() {
	esResponse, err := elasticsearch.Search(s.esClient, viper.GetString("elasticsearch.client.index"), "topic", "")
	if err != nil {
		s.logger.Error("es search query failed in all consumer initialization", zap.Error(err))
	}

	consumerList, err := elasticsearch.ConsumerSearchResult(esResponse)
	if err != nil {
		s.logger.Error("elasticsearch response to consumer desrialisation failed", zap.Any("es_response", esResponse), zap.Error(err))
	}

	for _, c := range consumerList {
		s.logger.Info("initialising kafka consumer", zap.String("topic", c.Topic), zap.String("groupId", c.GroupID))

		consumerGroup, handler, err := consumer.NewKafkaConsumer(c.GroupID, s.esClient, s.logger)
		if err != nil {
			s.logger.Error("kafka consumer group creation failed", zap.String("groupId", c.GroupID), zap.Error(err))
		}
		go func() {
			if err := consumerGroup.Consume(context.Background(), []string{c.Topic}, handler); err != nil {
				s.logger.Error("fail to consume kafka message", zap.String("topic", c.Topic), zap.Error(err))
			}

		}()
	}
	s.logger.Info("kafka consumers up and running")

}

//enableLogger enables logs to a file. it only gets enable on production
func (s *Server) enableLogger() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  viper.GetString("zap.log.logfile"),
		MaxSize:   viper.GetInt("zap.log.maxsize"), // megabytes
		MaxAge:    viper.GetInt("zap.log.maxage"),  // days
		LocalTime: viper.GetBool("zap.log.localtime"),
		Compress:  viper.GetBool("zap.log.compress"),
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	return logger
}
