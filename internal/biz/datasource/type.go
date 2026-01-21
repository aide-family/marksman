package datasource

type Type string

const (
	TypeMetric Type = "metric"
	TypeLogs   Type = "logs"
	TypeTrace  Type = "trace"
	TypeEvent  Type = "event"
)

func (t Type) String() string {
	return string(t)
}
