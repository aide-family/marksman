package datasource

type Engine string

const (
	// Metric 引擎
	EnginePrometheus      Engine = "prometheus"
	EngineVictoriaMetrics Engine = "victoriametrics"
	EngineInfluxDB        Engine = "influxdb"

	// Logs 引擎
	EngineLoki          Engine = "loki"
	EngineElasticsearch Engine = "elasticsearch"
	EngineVMLog         Engine = "vmlog"

	// Trace 引擎
	EngineTrace Engine = "trace"

	// Event 引擎
	EngineMQTT     Engine = "mqtt"
	EngineKafka    Engine = "kafka"
	EngineRocketMQ Engine = "rocketmq"
	EngineRabbitMQ Engine = "rabbitmq"
	EngineK8sEvent Engine = "k8s_event"
)

func (e Engine) String() string {
	return string(e)
}
