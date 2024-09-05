package trace

type Config struct {
	Enabled  bool   `json:"enabled"`
	Endpoint string `json:"endpoint"` // otel collector endpoint
	Name     string `json:"name"`     // service name
}
