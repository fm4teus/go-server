# Go Server - Servidor de Arquivos Simples

Este é um servidor de arquivos simples escrito em Go, que permite listar e fazer download de arquivos de um diretório especificado. O servidor foi implementado sem o uso da biblioteca http nativa do Go.

## Funcionalidades

O servidor possui as seguintes funcionalidades:

- Listagem de diretórios: Quando uma requisição GET é feita para o servidor com um caminho de diretório, o servidor retorna uma página HTML com a lista de arquivos e subdiretórios encontrados nesse diretório.
- Download de arquivos: O servidor permite que o usuário faça o download de arquivos específicos. Quando uma requisição GET é feita para o servidor com o caminho de um arquivo válido, o servidor retorna o arquivo como uma resposta HTTP.

## Conceitos de Redes de Computadores

Neste código, são utilizados alguns conceitos de redes de computadores. Abaixo, estão brevemente explicados:

1. Endereço IP: O endereço IP é uma identificação numérica única atribuída a cada dispositivo conectado a uma rede de computadores. No código, o endereço IP do servidor é definido como "localhost", o que significa que o servidor irá escutar apenas as conexões locais na máquina em que está sendo executado.

2. Porta: Uma porta é um número que identifica um serviço específico em um dispositivo. No código, a porta do servidor é definida como 8000, indicando que o servidor irá ouvir as conexões nessa porta.

3. Socket: Um socket é um ponto de extremidade em uma conexão de rede. No código, é utilizado um socket TCP para estabelecer a comunicação com os clientes. O servidor cria um socket TCP e aguarda por conexões de clientes.

4. Protocolo HTTP: O protocolo HTTP (Hypertext Transfer Protocol) é um protocolo de aplicação utilizado para a comunicação entre clientes e servidores na Web. O servidor implementa um parser HTTP básico para interpretar as requisições recebidas dos clientes e gerar as respostas apropriadas.

5. TCP/IP: O TCP/IP é uma família de protocolos de comunicação que são amplamente utilizados na Internet. O servidor utiliza o protocolo TCP/IP para estabelecer a conexão com os clientes e transmitir os dados entre eles.

## Execução do Servidor

Para executar o servidor:

```
./main
```

1. O servidor será iniciado e começará a escutar as conexões na porta especificada.
2. Você pode acessar o servidor abrindo um navegador da web e digitando o seguinte endereço na barra de endereços:

```
http://localhost:8000
```

Para testar alterações, compilar e executar:

```
go run main.go
```

Lembre-se de substituir a porta se você alterou a constante `PORT` no código.

## Contribuição

Sinta-se à vontade para contribuir com melhorias para este servidor.
