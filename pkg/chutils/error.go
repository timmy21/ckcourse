package chutils

import (
	"fmt"

	"github.com/ClickHouse/clickhouse-go"
)

func PrintError(err error) {
	if err == nil {
		return
	}
	if exception, ok := err.(*clickhouse.Exception); ok {
		fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
	} else {
		fmt.Println(err)
	}
}
