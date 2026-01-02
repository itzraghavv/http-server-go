package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itzraghavv/httpWebServer/internal/request"
	"github.com/itzraghavv/httpWebServer/internal/response"
	"github.com/itzraghavv/httpWebServer/internal/server"
)

const port = 42069

func handler(res *response.Writer, req *request.Request) {

	switch {
	case req.RequestLine.RequestTarget == "/yourproblem":
		res.WriteStatusLine(response.BadRequest)
		responseBody := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody(responseBody)
	case req.RequestLine.RequestTarget == "/myproblem":
		res.WriteStatusLine(response.InternalServerError)
		responseBody := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders((headers))
		res.WriteBody(responseBody)
	default:
		res.WriteStatusLine(response.OK)
		responseBody := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody(responseBody)
	}
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
