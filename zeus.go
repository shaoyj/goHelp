package goHelp

import (
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	DefaultKafkaZipkinTopic = "metrics.zipkin"
	ClientServerSameSpan    = false
	TraceID128Bit           = false
	globalCollector         zipkin.Collector
)

func NewZipkinTracerFromKafka(kafkaEndpoints []string, debug bool, hostPort, serviceName string) {
	// create kafka zipkin collector
	collector, err := zipkin.NewKafkaCollector(kafkaEndpoints, zipkin.KafkaTopic(DefaultKafkaZipkinTopic))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("create kafka collector fail")
		return
	}
	globalCollector = collector
	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

	t, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(ClientServerSameSpan),
		zipkin.TraceID128Bit(TraceID128Bit),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("unable to create Zipkin tracer")
		return
	}

	// Explicitly set our tracer to be the default tracer.
	opentracing.SetGlobalTracer(t)
}

func GetGlobalTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}

func GetZipkinCollector() zipkin.Collector {
	return globalCollector
}

func CloseCollector() {
	if globalCollector != nil {
		globalCollector.Close()
	}
}

//Init zeus
func InitZeus(devPoints, proPoints, hostPort, serverName string) {
	if strings.EqualFold(GetEnv(), "FAT") {
		NewZipkinTracerFromKafka(strings.Split(devPoints, ","), false,
			hostPort, serverName)
	} else {
		NewZipkinTracerFromKafka(strings.Split(proPoints, ","), false,
			hostPort, serverName)
	}
}

func FinishZeus(span *opentracing.Span) {
	(*span).SetTag("hostname", GetHostname()).Finish()
}

func CommitZeus(text string) {
	opentracing.StartSpan(text).SetTag("hostname", GetHostname()).Finish()
}
