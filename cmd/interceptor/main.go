package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

const (
	KB = 1024
	MB = 1024 * KB
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:9001")
	failOnErr(err)

	tcpListener, err := net.ListenTCP("tcp", addr)
	failOnErr(err, "main(): Listening to port.")
	defer tcpListener.Close()

	fmt.Println("started the tcp listener", tcpListener.Addr())
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("handling new connection")
	fmt.Println("local addr", conn.LocalAddr())

	buf := make([]byte, 512*KB)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("error while reading:", err)
	}
	fmt.Printf("---body start---\n%s\n---body end---\n\n", buf[:n])

	res, err := fwdHTTPTraffic(&buf, "http://localhost:9002")
	if err != nil {
		fmt.Println("forwarded server error:", err)
		return
	}
	if _, err = conn.Write(*res); err != nil {
		fmt.Println("error while writing response:", err)
	}
}

func fwdHTTPTraffic(body *[]byte, fwdSrv string) (*[]byte, error) {
	bufReader := bufio.NewReader(bytes.NewReader(*body))
	parsedReq, err := http.ReadRequest(bufReader)
	if err != nil {
		return nil, err
	}
	newReq, err := http.NewRequest(
		parsedReq.Method,
		fmt.Sprintf("%s%s", fwdSrv, parsedReq.RequestURI),
		parsedReq.Body,
	)
	// add all the initial headers
	for k, values := range parsedReq.Header {
		for _, v := range values {
			newReq.Header.Add(k, v)
		}
	}
	// add all the initial cookies
	for _, c := range parsedReq.Cookies() {
		newReq.AddCookie(c)
	}

	// send the request back to the actual server
	res, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return nil, err
	}

	// convert the response struct into a text http response
	var buf bytes.Buffer
	// the first line for defining the response
	if _, err = buf.WriteString(fmt.Sprintf("HTTP/%d.%d %d %s\r\n",
		res.ProtoMajor, res.ProtoMinor, res.StatusCode, res.Status)); err != nil {
		return nil, err
	}
	// add all the headers
	for k, values := range res.Header {
		valuesStr := ""
		for _, v := range values {
			valuesStr += v
		}
		if _, err = buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, valuesStr)); err != nil {
			return nil, err
		}
	}
	// mandatory line break to separate the body
	if _, err := buf.WriteString("\r\n"); err != nil {
		return nil, err
	}
	// write the response body if available
	if res.Body == nil {
		b := buf.Bytes()
		return &b, nil
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(resBody); err != nil {
		return nil, err
	}

	b := buf.Bytes()
	fmt.Println(string(b))
	return &b, nil
}

func failOnErr(err error, msg ...string) {
	if err != nil {
		log.Fatalln(nil, msg)
	}
}
