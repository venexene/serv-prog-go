package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"os/signal"
	"time"
)

type HTTPServer struct {
	root   string
	ln     net.Listener
	ctx    context.Context
	cancel context.CancelFunc

	wg  sync.WaitGroup
	log *os.File
	mu  sync.Mutex
}

type HTTPRequest struct {
	Method string
	Path   string
	Proto  string
}

func main() {
	root, _ := os.Getwd()

	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("log error:", err)
		return
	}
	defer logFile.Close()

	ctx, cancel := context.WithCancel(context.Background())

	s := &HTTPServer{
		root:   root,
		ctx:    ctx,
		cancel: cancel,
		log:    logFile,
	}

	go s.run(":8080")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("\nShutting down...")
	s.shutdown()
}

func (s *HTTPServer) run(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	s.ln = ln

	fmt.Println("Listening on", addr)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		if tcp, ok := ln.(*net.TCPListener); ok {
			tcp.SetDeadline(time.Now().Add(time.Second))
		}

		conn, err := ln.Accept()
		if err != nil {
			if op, ok := err.(*net.OpError); ok && op.Timeout() {
				continue
			}
			if s.ctx.Err() != nil {
				return
			}
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handle(conn)
		}()
	}
}

func (s *HTTPServer) shutdown() {
	s.cancel()
	if s.ln != nil {
		s.ln.Close()
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All requests done")
	case <-time.After(10 * time.Second):
		fmt.Println("Shutdown timeout")
	}
}

func (s *HTTPServer) handle(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second))

	ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

	req, err := readRequest(conn)
	if err != nil {
		s.writeError(conn, 400, "Bad Request")
		s.logReq(ip, "-", 400)
		return
	}

	if req.Method != "GET" {
		s.writeError(conn, 405, "Method Not Allowed")
		s.logReq(ip, req.Path, 405)
		return
	}

	status := s.serveFile(conn, req)
	s.logReq(ip, req.Path, status)
}

func readRequest(conn net.Conn) (*HTTPRequest, error) {
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) != 3 {
		return nil, fmt.Errorf("bad request")
	}

	// skip headers
	for {
		l, _ := reader.ReadString('\n')
		if l == "\r\n" || l == "\n" {
			break
		}
	}

	return &HTTPRequest{
		Method: parts[0],
		Path:   parts[1],
		Proto:  parts[2],
	}, nil
}

func (s *HTTPServer) serveFile(conn net.Conn, r *HTTPRequest) int {
	path, err := url.PathUnescape(r.Path)
	if err != nil {
		s.writeError(conn, 400, "Bad Request")
		return 400
	}

	full := filepath.Join(s.root, filepath.Clean(path))

	rootAbs, _ := filepath.Abs(s.root)
	fileAbs, _ := filepath.Abs(full)

	if !strings.HasPrefix(fileAbs, rootAbs) {
		s.writeError(conn, 403, "Forbidden")
		return 403
	}

	info, err := os.Stat(fileAbs)
	if err != nil {
		if os.IsNotExist(err) {
			s.writeError(conn, 404, "Not Found")
			return 404
		}
		s.writeError(conn, 500, "Internal Server Error")
		return 500
	}

	if info.IsDir() {
		index := filepath.Join(fileAbs, "index.html")
		if i, err := os.Stat(index); err == nil && !i.IsDir() {
			return s.sendFile(conn, index, i)
		}
		return s.dirList(conn, fileAbs, r.Path)
	}

	return s.sendFile(conn, fileAbs, info)
}

func (s *HTTPServer) sendFile(conn net.Conn, path string, info os.FileInfo) int {
	f, err := os.Open(path)
	if err != nil {
		s.writeError(conn, 500, "Internal Server Error")
		return 500
	}
	defer f.Close()

	ct := mime.TypeByExtension(filepath.Ext(path))
	if ct == "" {
		ct = "application/octet-stream"
	}

	fmt.Fprintf(conn,
		"HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\nConnection: close\r\n\r\n",
		ct, info.Size())

	io.Copy(conn, f)
	return 200
}

func (s *HTTPServer) dirList(conn net.Conn, dir, reqPath string) int {
	files, _ := os.ReadDir(dir)

	var b strings.Builder
	b.WriteString("<html><body><h1>Index of " + reqPath + "</h1><ul>")

	for _, f := range files {
		name := f.Name()
		if f.IsDir() {
			name += "/"
		}
		fmt.Fprintf(&b, `<li><a href="%s">%s</a></li>`, name, name)
	}

	b.WriteString("</ul></body></html>")

	body := b.String()

	fmt.Fprintf(conn,
		"HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s",
		len(body), body)

	return 200
}

func (s *HTTPServer) writeError(conn net.Conn, code int, msg string) {
	body := fmt.Sprintf("<h1>%d %s</h1>", code, msg)

	fmt.Fprintf(conn,
		"HTTP/1.1 %d %s\r\nContent-Type: text/html\r\nContent-Length: %d\r\n\r\n%s",
		code, msg, len(body), body)
}

func (s *HTTPServer) logReq(ip, path string, code int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s | %s | %s | %d\n", t, ip, path, code)

	fmt.Print(line)
	s.log.WriteString(line)
}