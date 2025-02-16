package postgres

type Config struct {
	URL      string `env:"DB_URL"    envDefault:"postgres://su:su@postgres:5432/image?sslmode=disable"`
	MinConns int32  `env:"MIN_CONNS" envDefault:"1"`
	MaxConns int32  `env:"MAX_CONNS" envDefault:"3"`
}
