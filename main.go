package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

const (
	HOST      = "localhost"                                    // Endereço IP do servidor
	PORT      = 8000                                           // Porta que o servidor irá ouvir
	DIRECTORY = "/home/fmateus/Pictures/server-example-files/" // Diretório a ser listado
)

var (
	directoryRegex = regexp.MustCompile(`^GET /([A-Za-z0-9\-/]+)* HTTP/1\.1`)
	downloadRegex  = regexp.MustCompile(`^GET /([A-Za-z0-9\-/]+)*download\?file\=`)
)

func renderHTMLContent(content string) string {
	template := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Go Server</title>
		<style>
			body {
				background-color: #f1f1f1;
				font-family: Arial, sans-serif;
			}
			
			h1 {
				color: #333;
			}
			
			p {
				color: #666;
			}


			ul {
				list-style-type: none;
				padding: 0;
			}
			
			li {
				margin-bottom: 5px;
			}
			
			a {
				color: #000;
				text-decoration: none;
				transition: color 0.3s ease;
			}
			
			a:hover {
				color: #ff0000;
			}

			.container {
				max-width: 800px;
				margin: 0 auto;
				padding: 20px;
				background-color: #fff;
			}

		</style>
	</head>
	<body>
		<div class="container">
			<div style="display: flex; align-items: center;">
				<img src="https://raw.githubusercontent.com/betandr/gophers/master/Gopher.png" alt="Gopher" style="width: 50px; height: 66px;">
				<h1 style="margin-left: 10px;">Go Server</h1>
			</div>
			<p>Este é um servidor de arquivos simples.</p>
			<p>Escrito em Go sem o uso da biblioteca http nativa.</p>
			<ul>
			%s
			</ul>
		</div>
	</body>
	</html>
	`

	return fmt.Sprintf(template, content)
}

func generateFileLinks(subdir string) string {
	// Gera os links para os arquivos no diretório especificado.
	files, err := os.ReadDir(DIRECTORY + subdir)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	links := []string{}

	for _, file := range files {
		if file.IsDir() {
			link := fmt.Sprintf("<li><a href=\"%s/%s\">%s</a></li>", subdir, file.Name(), file.Name())
			links = append(links, link)
			continue
		}
		//filePath := fmt.Sprintf("%s/%s", DIRECTORY, file.Name())
		link := fmt.Sprintf("<li><a href=\"/%s/download?file=%s\">⬇️ %s</a></li>", subdir, file.Name(), file.Name())
		links = append(links, link)
	}

	return renderHTMLContent(strings.Join(links, "\n"))
}

func handleRequest(request string) []byte {
	switch {
	case regexp.MustCompile(`^GET /HEADER`).MatchString(request):
		return []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n%s", request))
	case directoryRegex.MatchString(request):

		subdir := "/"
		sm := directoryRegex.FindStringSubmatch(request)
		if len(sm) > 1 {
			subdir = sm[1]
		}

		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n%s", generateFileLinks(subdir))
		return []byte(response)
	case downloadRegex.MatchString(request):
		subdir := "/"
		sm := downloadRegex.FindStringSubmatch(request)
		fmt.Println("sm: ", sm)
		if len(sm) > 1 {
			subdir = sm[1]
		}
		fileName := strings.Split(strings.Split(request, "=")[1], " HTTP/1.1")[0]
		filePath := fmt.Sprintf("%s%s%s", DIRECTORY, subdir, fileName)

		fmt.Println("path: {", filePath, "}")

		if fileInfo, err := os.Stat(filePath); err == nil && !fileInfo.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println(err)
				return []byte("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/html\r\n\r\nInternal Server Error.")
			}
			defer file.Close()

			fileData := make([]byte, fileInfo.Size())
			if _, err := file.Read(fileData); err != nil {
				fmt.Println(err)
				return []byte("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/html\r\n\r\nInternal Server Error.")
			}

			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Disposition: attachment; filename=%s\r\n\r\n%s", fileName, fileData)
			return []byte(response)
		} else {
			return []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\nFile not found.")
		}
	default:
		return []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\nNot Found.")
	}
}

func runServer() {
	// Inicia o servidor e fica ouvindo por conexões.
	listenAddr := fmt.Sprintf("%s:%d", HOST, PORT)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	fmt.Printf("Servidor HTTP em execução em %s\n", listenAddr)
	fmt.Printf("Diretório: %s\n", DIRECTORY)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Conexão recebida de %s\n", clientAddr)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	request := string(buffer[:n])
	response := handleRequest(request)

	if _, err := conn.Write(response); err != nil {
		fmt.Println(err)
	}
}

func main() {
	runServer()
}
