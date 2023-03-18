package server

import (
	"bufio"
	"context"

	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"tcp-server/internal/pkg/config"
	"tcp-server/internal/pkg/pow"
	"tcp-server/internal/pkg/protocol"
	"time"
)

var quote = []string{
	"Do something that interests you and do it to the , " +
		"absolute best of your ability whatever your limitations, " +
		"you will almost certainly still do better than anyone else",

	"And shall find wisdom and great treasures of knowledge, even hidden treasures",

	"The best way to predict your future is to create it",

	"When one door of happiness closes, another opens " +
		"but often we look so long at the closed door " +
		"that we do not see the one that has opened for us",
}
var ErrQuit = errors.New("client request to close connection")

type Clock interface {
	Now() time.Time
}

type Cache interface {
	Add(int, int64) error
	Get(int) (bool, error)
	Delete(int)
}

//run main function

func Run(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println("listening", listener.Addr())

	//client send new request every 5
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection:%w", err)
		}
		go handleConnection(ctx, conn)

	}
}

//handle connection

func handleConnection(ctx context.Context, conn net.Conn) {
	fmt.Println("new client", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error read connection", err)
			return
		}
		msg, err := ProcessRequest(ctx, req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("error process request", err)
			return
		}

		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				fmt.Println("error send message", err)
			}
		}
	}

}

func ProcessRequest(ctx context.Context, msgStr string, clientInfo string) (*protocol.Message, error) {

	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}

	switch msg.Header {
	case protocol.Quit:
		return nil, ErrQuit
	case protocol.RequestChallenge:
		fmt.Printf("client %s request challenge\n", clientInfo)

		conf := ctx.Value("config").(*config.Config)
		clock := ctx.Value("clock").(Clock)
		cache := ctx.Value("cache").(Cache)
		date := clock.Now()

		randValue :=
			rand.Intn(100000)
		err := cache.Add(randValue, conf.HashcashDuration)
		if err != nil {
			return nil, fmt.Errorf("err add rand to cache")
		}

		hashcash := pow.HashcashData{
			Version:    1,
			ZerosCount: conf.HashcashZerosCount,
			Date:       date.Unix(),
			Resource:   clientInfo,
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
			Counter:    0,
		}

		hashcashMarshaled, err := json.Marshal(hashcash)

		if err != nil {
			return nil, fmt.Errorf("error marshal:%v", err)
		}

		msg := protocol.Message{
			Header:  protocol.ResponseChallenge,
			Payload: string(hashcashMarshaled),
		}
		return &msg, nil

	case protocol.RequestResource:

		fmt.Printf("client %s requests resource payload %s\n", clientInfo, msg.Payload)

		var hashcash pow.HashcashData
		err := json.Unmarshal([]byte(msg.Payload), &hashcash)

		if err != nil {
			return nil, fmt.Errorf("invalid unmarshal hashcash:%w", err)
		}

		if hashcash.Resource != clientInfo {
			return nil, fmt.Errorf("invalid hashcash")
		}

		conf := ctx.Value("config").(*config.Config)
		clock := ctx.Value("clock").(Clock)
		cache := ctx.Value("cache").(Cache)

		randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
		if err != nil {
			return nil, fmt.Errorf("error decode rand:%w", err)
		}

		randValue, err := strconv.Atoi(string(randValueBytes))
		if err != nil {
			return nil, fmt.Errorf("error decode rand:%w", err)
		}

		exists, err := cache.Get(randValue)
		if err != nil {
			return nil, fmt.Errorf("error  rand value from cache:%w", err)
		}

		if !exists {
			return nil, fmt.Errorf("challenge expired or not sent")
		}
		//sent solution should not be outdated

		if clock.Now().Unix()-hashcash.Date > conf.HashcashDuration {
			return nil, fmt.Errorf("challenge expired")
		}

		maxIter := hashcash.Counter
		if maxIter == 0 {
			maxIter = 1
		}

		_, err = hashcash.ComputeHashCash(maxIter)
		if err != nil {
			return nil, fmt.Errorf("invalid hashcash")
		}

		//get random quote'

		fmt.Printf("client %s successfuly computed %s\n", clientInfo, msg.Payload)

		msg := protocol.Message{
			Header:  protocol.ResponseResource,
			Payload: quote[rand.Intn(4)],
		}
		cache.Delete(randValue)
		return &msg, nil

	default:
		return nil, fmt.Errorf("unknown header")

	}

}

func sendMsg(msg protocol.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}
