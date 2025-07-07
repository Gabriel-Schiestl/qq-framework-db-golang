package db

type IDbProvider interface {
	GetDriver() string
	GetHost() string
	GetPort() int
	GetName() string
	GetUser() string
	GetPass() string
	GetConnectionRetries() int
	GetMaxIdleConnections() int
	GetMaxOpenConnections() int
	GetConnectionMaxLifetime() int
	GetLimitePaginacao() int
	GetLimiteRotinas() int
}