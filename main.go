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
	NULL = "\x00"
)

var (
	server = flag.String("server", "server.duinocoin.com:2817", "addr and port of server.")
	name = flag.String("name", os.Getenv("MINERNAME"), "wallet/miner name.")
	quiet = flag.Bool("quiet", false, "disable logging to console.")
	debug = flag.Bool("debug", false, "console log send/receive messages.")
	wait = flag.Int("wait", 10, "time to wait between task checks.")
	batch = flag.Int("batch", 10, "how many jobs to create.")
	version = "0.1"
)

func init() {
	//Following is needed or else numbers aren't random.
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

	logger("Starting ds-go-node version ", version, NEWLINE)

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
		cj := &CreateJob{
			User: *name,
		}

		err := cj.sync(conn)
		if err != nil {
			if err == io.EOF {
				conn, _ = connect()
			}
			continue
		}

		err = cj.createJobs()
		if err != nil {
			loggerDebug("createjob error ", err, NEWLINE)
			continue
		}

		err = cj.sendJobs(conn)
		if err != nil {
			loggerDebug("sendjob error ", err, NEWLINE)
			if err == io.EOF {
				conn, _ = connect()
			}
			continue
		}
	}
}

// Provides a marshal for unit testing.
func (j *CreateJob) marshal() (string, error) {
	res, err := json.Marshal(*j)
	return string(res), err
}

// sendJobs sends the result of the job over the connection.
func (j *CreateJob) sendJobs(conn net.Conn) (err error) {

	res, err := j.marshal()
	if err != nil {
		return
	}

	err = send(conn, res)
	if err != nil {
		return
	}

	resp, err := read(conn)
	if err != nil {
		return
	}

	logger("Submit Job Response: ", resp)

	return
}

func makeJob(job *Job, diff uint64) (err error) {
	job.Nonce = uint64(rand.Intn(int(diff * 100)))
	nonce := strconv.FormatUint(job.Nonce, 10)
	data := []byte(job.LastHash + nonce)

	h := sha1.New()
	h.Write(data)

	job.ExpectedHash = hex.EncodeToString(h.Sum(nil))
	loggerDebug("created job ", *job, NEWLINE)
	return
}

// createJobs loops to create (*batch) amount of jobs.
func (j *CreateJob) createJobs() (err error) {

	for i := 0; i < *batch; i++ {
		job := Job{
			LastHash: j.LastHash,
		}

		err = makeJob(&job, j.Difficulty)
		if err != nil {
			return
		}

		j.Jobs = append(j.Jobs, job)
	}

	return
}

func parseUint(str string) (ret uint64, err error) {
	str = strings.TrimRight(str, NULL)
	str = strings.TrimRight(str, NEWLINE)
	ret, err = strconv.ParseUint(str, 10, 64)
	return
}

// parseJobs parses the job request sent from the server.
func (j *CreateJob) parseJobs(buf *string) (err error) {

	str := strings.Split(string(*buf), SEPERATOR)
	if len(str) < 2 {
		loggerDebug("string split error ", *buf)
		return errors.New("str split error")
	}

	diff, err := parseUint(str[2])
	if err != nil {
		loggerDebug("unable to parse uint ", err, NEWLINE)
		return
	}

	switch str[0] {
	case "CREATE_JOBS":
		j.LastHash = str[1]
		j.Difficulty = diff
	case "NO_TASK":
		sleep := time.Duration(*wait) * time.Second
		logger("no_task sleep for ", sleep, NEWLINE)
		time.Sleep(sleep)
		err = errors.New("no_task error")
	default:
		loggerDebug("task command error ", str[0], NEWLINE)
		err = errors.New("task command error")
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

	resp, err := read(conn)
	if err != nil {
		return
	}

	logger("Connected to Server Version: ", resp, NEWLINE)

	return
}

// sync is used to request jobs.
func (j *CreateJob) sync(conn net.Conn) (err error) {

	err = send(conn, "NODE")
	if err != nil {
		return
	}

	resp, err := read(conn)
	if err != nil {
		return
	}

	logger("Get Job Response: ", resp)

	return j.parseJobs(&resp)
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

func loggerDebug(msg ...interface{}) {
	if (!*debug) {
		return
	}

	logger(msg...)
}

// read is a helper for reciving a string
func read(conn net.Conn) (ret string, err error) {
	buf := make([]byte, 128)
	_, err = conn.Read(buf)

	if err != nil {
		return
	}

	ret = string(buf)
	
	loggerDebug("read ", ret)
	return
}

// send is a helper for sending a string
func send(conn net.Conn, str string) (err error) {
	fmt.Fprintln(conn, str)
	loggerDebug("send ", str, NEWLINE)
	return
}
