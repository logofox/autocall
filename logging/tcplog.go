package logging

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"runtime"
	"sync"
)

// TcplogBackend is a simple logger to tcp backend. It automatically maps
// the internal log levels to appropriate syslog log levels.
type TcplogBackend struct {
	Color bool
	flag  int
	conns map[string]net.Conn
	l     sync.Mutex
}

// NewTcplogBackend Listening on a port, waiting for a connection, the connection to the broadcast log
// given prefix. If prefix is not given, the prefix will be derived from the
// launched command.
func NewTcplogBackend(addr string, prefix string, flag int) (b *TcplogBackend, err error) {
	var backend TcplogBackend
	backend.Color = false
	backend.conns = make(map[string]net.Conn)
	backend.flag = flag

	_listen := func() {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go func(c net.Conn) {
				for {
					r := bufio.NewReader(conn)
					msg, err := r.ReadString('\n')
					if err != nil {
						b.l.Lock()
						delete(b.conns, conn.RemoteAddr().String())
						b.l.Unlock()
						return
					}

					if msg == "conncentd\n" {
						b.l.Lock()
						b.conns[conn.RemoteAddr().String()] = conn
						b.l.Unlock()
					}
				}
			}(conn)
		}
	}

	go _listen()

	return &backend, err
}

func (b *TcplogBackend) broadcast(msg string) {
	defer b.l.Unlock()
	b.l.Lock()
	for addr, conn := range b.conns {
		if conn != nil {
			w := bufio.NewWriter(conn)
			if _, err := w.WriteString(msg); err == nil {
				w.Flush()
			} else {
				delete(b.conns, addr)
			}

		}
	}
}

func (b *TcplogBackend) Log(level Level, calldepth int, rec *Record) error {
	msg := "Unhandled log level"
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if ok {
		re, _ := regexp.Compile(`[^/\\\\]*$`)
		file = re.FindString(file)

		time := rec.Time.Format("2006-01-02 15:04:05.0000")
		formatted := rec.Formatted(calldepth + 1)

		if b.Color {
			msg = fmt.Sprintf("%s%s [%s] %s:%d %s\033[0m\n", colors[level], time, level, file, line, formatted)
		} else {
			msg = fmt.Sprintf("%s [%s] %s:%d %s\n", time, level, file, line, formatted)
		}
	}

	b.broadcast(msg)
	return nil
}
