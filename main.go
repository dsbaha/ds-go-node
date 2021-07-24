package main

import (
	"io"
	"os"
	"fmt"
	"net"
	"flag"
	"time"
	"errors"
	"strconv"
	"strings"
	"math/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

const (
	SEPERATOR = ","
	NEWLINE = "\n"
	ENCODING = "utf-8"
	NULL = "\x00"
)

var (
	server = flag.String("server", "server.duinocoin.com:2817", "addr and port of server.")
	name = flag.String("name", os.Getenv("MINERNAME"), "wallet/miner name.")
	quiet = flag.Bool("quiet", false, "disable logging to console.")
	wait = flag.Int("wait", 10, "time to wait between task checks.")
)

func init() {
	//Following is needed or else numbers aren't random.
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

	if *name == "" {
		logger("Name Not Set", NEWLINE)
		flag.PrintDefaults()
		os.Exit(1)
	}

	conn, err := connect()
	if err != nil {
		logger("Connection Error ", err)
		os.Exit(1)
	}

	for {
		cj := &CreateJob{}

		err := cj.sync(conn)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				conn, _ = connect()
			} else {
				time.Sleep(time.Duration(*wait) * time.Second)
			}
			continue
		}

		err = cj.createJobs()
		if err != nil {
			logger("createjob error ", err)
			continue
		}

		err = cj.sendJobs(conn)
		if err != nil {
			logger("sendjob error", err)
			if err == io.EOF {
				conn.Close()
				conn, _ = connect()
			}
			continue
		}
	}
}

// Provides a marshal for unit testing.
func (j *CreateJob) marshal() (res []byte, err error) {
	res, err = json.Marshal(*j)
	return
}

// sendJobs sends the result of the job over the connection.
func (j *CreateJob) sendJobs(conn net.Conn) (err error) {

	res, err := j.marshal()
	if err != nil {
		return
	}

	fmt.Fprintln(conn, string(res))

	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		return
	}

	logger("Submit Job Response: ", string(buf))

	return
}

// createJobs loops to create 10 jobs.
func (j *CreateJob) createJobs() (err error) {
	for i := 0; i < 10; i++ {
		random := uint64(rand.Intn(int(j.Difficulty * 100)))
		nonce := strconv.FormatUint(random, 10)
		data := []byte(j.LastHash + nonce)

		h := sha1.New()
		h.Write(data)

		job := Job{
			LastHash: j.LastHash,
			ExpectedHash: hex.EncodeToString(h.Sum(nil)),
			Nonce: random,
		}

		j.Jobs = append(j.Jobs, job)
	}

	return
}

// parseJobs parses the job request sent from the server.
func (j *CreateJob) parseJobs(buf *[]byte) (err error) {

	str := strings.Split(string(*buf), SEPERATOR)
	if len(str) < 2 {
		return errors.New("str split error")
	}

	str[2] = strings.TrimRight(str[2], NULL)
	str[2] = strings.TrimRight(str[2], NEWLINE)

	difficulty, err := strconv.ParseUint(str[2], 10, 64)
	if err != nil {
		return
	}

	switch str[0] {
	case "CREATE_JOBS":
		j.LastHash = str[1]
		j.Difficulty = difficulty
	case "NO_TASK":
		fallthrough
	default:
		err = errors.New("task error")
	}

	return
}

// connect is used to connect to the server.
func connect() (conn net.Conn, err error) {
	logger("Connecting to Server: ", *server, NEWLINE)

	conn, err = net.Dial("tcp", *server)
	if err != nil {
		return
	}

	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		return
	}

	logger("Connected to Server Version: ", string(buf), NEWLINE)

	return
}

// sync is used to request jobs.
func (j *CreateJob) sync(conn net.Conn) (err error) {

	fmt.Fprintln(conn, "NODE")

	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		return
	}

	logger("Get Job Response: ", string(buf))

	return j.parseJobs(&buf)
}

// logger is the general purpose logger
// which can be turned off w/ cmd line switch
func logger(msg ...interface{}) {
	if (*quiet) {
		return
	}

	tm := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] ", tm)

	for _, v := range msg {
		fmt.Print(v)
	}
}
