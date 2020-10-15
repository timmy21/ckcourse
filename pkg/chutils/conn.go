package chutils

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	defaultConnOptions = ConnOptions{
		Port: 9000,
	}
)

type ConnOptions struct {
	Port  int
	Debug bool
}

func (o ConnOptions) DataSourceName(addr string) string {
	s := fmt.Sprintf("tcp://%s:%d", addr, o.Port)
	if o.Debug {
		s = fmt.Sprintf("%s?debug=true", s)
	}
	return s
}

type ConnOption func(o *ConnOptions)

func WithPort(v int) ConnOption {
	return func(o *ConnOptions) {
		o.Port = v
	}
}

func WithDebug(v bool) ConnOption {
	return func(o *ConnOptions) {
		o.Debug = v
	}
}

func CreateConnect(addr string, options ...ConnOption) (*sqlx.DB, error) {
	opts := defaultConnOptions
	for _, v := range options {
		v(&opts)
	}

	connect, err := sqlx.Open("clickhouse", opts.DataSourceName(addr))
	if err != nil {
		return nil, err
	}
	if err := connect.Ping(); err != nil {
		return nil, err
	}
	return connect, nil
}
