import socket
import os

HOST = 'localhost'  # Endereço IP do servidor
PORT = 8000  # Porta que o servidor irá ouvir
DIRECTORY = '.'  # Diretório a ser listado

def generate_file_links():
    """Gera os links para os arquivos no diretório especificado."""
    files = os.listdir(DIRECTORY)
    links = []

    for file in files:
        file_path = os.path.join(DIRECTORY, file)
        link = f'<a href="/download?file={file}">{file}</a>'
        links.append(link)

    return '\n'.join(links)

def handle_request(request):
    """Lida com a requisição recebida e retorna uma resposta."""
    if request.startswith('GET / HTTP/1.1'):
        response = f'HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n{generate_file_links()}'.encode()
    elif request.startswith('GET /download?file='):
        file_name = request.split('=')[1]
        file_path = os.path.join(DIRECTORY, file_name)
        
        if os.path.isfile(file_path):
            with open(file_path, 'rb') as file:
                file_data = file.read()
            response = f'HTTP/1.1 200 OK\r\nContent-Disposition: attachment; filename={file_name}\r\n\r\n'
            response = response.encode() + file_data
        else:
            response = b'HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\nFile not found.'
    elif request.startswith('GET /HEADER'):
        response = b'HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n' + request.encode()
    else:
        response = b'HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\n\r\nNot Found.'

    return response

def run_server():
    """Inicia o servidor e fica ouvindo por conexões."""
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
        server_socket.bind((HOST, PORT))
        server_socket.listen(1)

        print(f'Servidor HTTP em execução em {HOST}:{PORT}')
        print(f'Diretório: {DIRECTORY}')

        while True:
            client_socket, client_address = server_socket.accept()
            print(f'Conexão recebida de {client_address[0]}:{client_address[1]}')

            request_data = client_socket.recv(1024).decode('utf-8')
            response_data = handle_request(request_data)

            print(response_data)
            client_socket.sendall(response_data)
            client_socket.close()

if __name__ == '__main__':
    run_server()
