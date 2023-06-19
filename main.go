package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	HOST      = "localhost"                // Endereço IP do servidor
	PORT      = 8000                       // Porta que o servidor irá ouvir
	DIRECTORY = "/home/fmateus/Documents/" // Diretório a ser listado
)

func generateFileLinks() string {
	// Gera os links para os arquivos no diretório especificado.
	files, err := os.ReadDir(DIRECTORY)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	links := []string{}

	for _, file := range files {
		//filePath := fmt.Sprintf("%s/%s", DIRECTORY, file.Name())
		link := fmt.Sprintf("<a href=\"/download?file=%s\">%s</a>", file.Name(), file.Name())
		links = append(links, link)
	}

	return strings.Join(links, "\n")
}

func handleRequest(request string) []byte {
	// Lida com a requisição recebida e retorna uma resposta.
	if strings.HasPrefix(request, "GET / HTTP/1.1") {
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n%s", generateFileLinks())
		return []byte(response)
	} else if strings.HasPrefix(request, "GET /download?file=") {
		fileName := strings.Split(request, "=")[1]
		filePath := fmt.Sprintf("%s/%s", DIRECTORY, fileName)

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
	} else if strings.HasPrefix(request, "GET /HEADER") {
		return []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n%s", request))
	} else {
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

	fmt.Println(response)
	if _, err := conn.Write(response); err != nil {
		fmt.Println(err)
	}
}

func main() {
	runServer()
}
