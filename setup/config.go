package setup

import (
	"context"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
)

// API
type API struct {
	RootToken string `env:"ROOT_TOKEN"`
}

// DB
type DB struct {
	Username    string `env:"USERNAME, default=dino"`
	Password    string `env:"PASSWORD, default=dino"`
	Host        string `env:"HOST, default=postgres"`
	Port        string `env:"PORT, default=5432"`
	Name        string `env:"NAME, default=dino"`
	ExtraParams string `env:"EXTRA_PARAMS, default=sslmode=disable"`
}

// Dial
func (dbc *DB) Dial() string {
	var b strings.Builder
	b.WriteString("postgresql://")
	b.WriteString(dbc.Username + ":")
	b.WriteString(dbc.Password + "@")
	b.WriteString(dbc.Host + ":")
	b.WriteString(dbc.Port + "/")
	b.WriteString(dbc.Name + "?")
	b.WriteString(dbc.ExtraParams)
	return b.String()
}

// Logger
type Logger struct {
	Level string `env:"LEVEL,default=DEBUG"`
}

// Proxy
type Proxy struct {
	Host string `env:"HOST, default=127.0.0.1"`
	Port string `env:"PORT, default=8080"`

	ReadTimeout       time.Duration `env:"READ_TIMEOUT, default=15"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT, default=15"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT, default=15"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT, default=30"`
}

// Server
type Server struct {
	Host string `env:"HOST, default=127.0.0.1"`
	Port string `env:"PORT, default=50051"`

	QuicHost string `env:"QUIC_HOST, default=127.0.0.1"`
	QuicPort string `env:"QUIC_PORT, default=4242"`

	CertPath string `env:"SSL_CERT_PATH"`
	KeyPath  string `env:"SSL_KEY_PATH"`
}

// JWT
type JWT struct {
	Duration int64    `env:"DURATION, default=-1"` // jwt token duration in seconds
	Issuer   string   `env:"ISSUER, default=dino.local"`
	Audience []string `env:"AUD,default=dino"`
}

// Authenticator
type Authenticator struct {
	JWT *JWT `env:",prefix=JWT_"`
}

// Tunnel
type Tunnel struct {
	ID       string `env:"ID"`
	Token    string `env:"TOKEN"`
	Endpoint string `env:"ENDPOINT, default=tunnel.dino.local:4242"`
}

// Configs
type Config struct {
	fx.Out

	API  *API           `env:",prefix=API_"`
	Auth *Authenticator `env:",prefix=AUTH_"`

	DB  *DB     `env:",prefix=DB_"`
	Log *Logger `env:",prefix=LOG_"`

	Proxy *Proxy `env:",prefix=PROXY_"`

	Server *Server `env:",prefix=SERVER_"`
	Tunnel *Tunnel `env:",prefix=TUNNEL_"`
}

// Module
var Module = fx.Module("setup_config", fx.Provide(newModule))

func newModule() (Config, error) {
	var cfg Config
	err := envconfig.Process(context.TODO(), &cfg)
	return cfg, err
}
