package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var (
	server_pass = "SYSPRO1018262\r\n"
	wavfile     = "quote2020enc.wav"
	offset      = 78
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr,
			"Usage: %s port\n", os.Args[0])
		os.Exit(1)
	}
	portno := os.Args[1]
	server_loop(portno)
}

func server_loop(portno string) {
	listener, err := tcp_listen_port(portno)
	if err != nil {
		panic(err)
	}

	wavdata, err := read_file(wavfile)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := syscall.Munmap(wavdata); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Accepting incoming connections...")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func() {
			fmt.Println("connected from", conn.RemoteAddr().String())
			defer func() {
				fmt.Println("disconnected from", conn.RemoteAddr().String())
				conn.Close()
			}()
			in, out := fdopen_sock(conn)

			//step 1:receive pass(FUN2020)
			client_pass, err := recv(in)
			if err != nil {
				fmt.Println(err)
				return
			}
			if client_pass != "FUN2020" {
				fmt.Println("invaild pass phrase")
				return
			}

			//step 2:send pass(SYSPRO10*****)
			send(out, []byte(server_pass))

			//step 3:receive key(0x**)
			s, err := recv(in)
			if err != nil {
				fmt.Println(err)
				return
			}
			key, err := getkey(s)
			if err != nil {
				fmt.Println(err)
				return
			}

			//step 4:send wav data
			buf := make([]byte, len(wavdata)-offset)
			for i, v := range wavdata[offset:] {
				buf[i] = v ^ key
			}
			send(out, buf)
		}()
	}
}

func getkey(s string) (byte, error) {
	r := regexp.MustCompile(`0x[0-9A-Fa-f]{2}`)
	if !r.MatchString(s) {
		return 0x00, errors.New("Invailed key")
	}
	s = strings.TrimLeft(s, "0x")
	b, err := strconv.ParseUint(s, 16, 8)
	return byte(b), err
}

func read_file(filename string) ([]byte, error) {
	//Open input file
	fi, err := syscall.Open(filename, syscall.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := syscall.Close(fi); err != nil {
			panic(err)
		}
	}()

	//Get status of input file
	var sti syscall.Stat_t
	if err := syscall.Fstat(fi, &sti); err != nil {
		return nil, err
	}

	//Read input file by syscall.Mmap
	return syscall.Mmap(fi, 0, int(sti.Size), syscall.PROT_READ, syscall.MAP_SHARED)
}

func recv(in *bufio.Reader) (string, error) {
	line, _, err := in.ReadLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}

func send(out *bufio.Writer, data []byte) {
	out.Write(data)
	out.Flush()
}

func fdopen_sock(conn *net.TCPConn) (*bufio.Reader, *bufio.Writer) {
	in := bufio.NewReader(conn)
	out := bufio.NewWriter(conn)
	return in, out
}

func tcp_listen_port(portno string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+portno)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	return listener, err
}
