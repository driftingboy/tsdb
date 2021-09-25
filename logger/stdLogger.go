package logger

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

// 可变化的抽一个 option config，由option func 配置
func NewStdLogger(w io.Writer, bufferSize int) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				buff := new(strings.Builder)
				buff.Grow(bufferSize)
				return buff
			},
		},
	}
}

func (s stdLogger) Log(Level Level, kvs ...interface{}) {
	if len(kvs) == 0 {
		return
	}
	if len(kvs)&1 == 1 {
		kvs = append(kvs, "")
	}
	buff := s.pool.Get().(*strings.Builder)
	defer s.pool.Put(buff)
	buff.WriteString(Level.String())
	for i := 0; i < len(kvs); i += 2 {
		buff.WriteString(fmt.Sprintf(" %s=%v", kvs[i], kvs[i+1]))
	}
	_ = s.log.Output(4, buff.String())
	buff.Reset()
}
