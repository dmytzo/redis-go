package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	CommandPING = "PING"
	CommandECHO = "ECHO"
	CommandGET  = "GET"
	CommandSET  = "SET"

	SignRespArray      = '*'
	SignRespBulkString = '$'

	SignString = "+"
)

func main() {
	fmt.Println("redis")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Fatalf("bind to port: %s", err.Error())
	}

	defer l.Close()

	storage := map[string]string{}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("accept connection: %s", err.Error())
		}

		go func() {
			for {
				reader := bufio.NewReader(c)

				commands, err := parse(reader)

				if err != nil {
					if errors.Is(err, io.EOF) {
						continue
					}

					log.Fatalf("read input: %s", err.Error())
				}

				if len(commands) == 0 {
					continue
				}

				var resp string

				switch strings.ToUpper(commands[0]) {
				case CommandPING:
					resp = "PONG"
				case CommandECHO:
					resp = strings.Join(commands[1:], " ")
				case CommandSET:
					if len(commands) < 3 {
						continue
					}

					storage[commands[1]] = commands[2]
					resp = "OK"
				case CommandGET:
					if len(commands) < 2 {
						continue
					}

					resp = storage[commands[1]]
				}

				if _, err := c.Write(strb(resp)); err != nil {
					log.Fatalf("write: %s", err.Error())
				}
			}
		}()
	}
}

func strb(s string) []byte {
	return []byte(fmt.Sprintf("%s%s\r\n", SignString, s))
}

type cmd []string

func parse(b *bufio.Reader) (cmd, error) {
	respType, err := b.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read byte: %w", err)
	}

	itemNumRaw, err := b.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read bytes: %w", err)
	}

	itemNum, err := strconv.Atoi(string(itemNumRaw[:len(itemNumRaw)-2]))
	if err != nil {
		return nil, fmt.Errorf("strconv atoi: %w", err)
	}

	var command cmd

	switch respType {
	case SignRespArray:
		for i := 0; i < itemNum; i++ {
			cmds, err := parse(b)
			if err != nil {
				return nil, fmt.Errorf("parse: %w", err)
			}

			command = append(command, cmds...)
		}
		return command, nil

	case SignRespBulkString:
		itemRaw, err := b.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("peak: %w", err)
		}

		item := itemRaw[:len(itemRaw)-2]
		if itemLen := len(item); itemLen != itemNum {
			return nil, fmt.Errorf("wrond len: %d != %d", itemLen, itemNum)
		}

		return cmd{string(item)}, nil
	}

	return command, nil
}
