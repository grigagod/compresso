package postgres

// Config could be used for configuring new postgres connections.
type Config struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	Driver   string
}
