# Cliente/Servidor em Go

## Executar Servidor

Acesse o diretório `server/`.

    $ go run main.go diffiehellman.go -name server_name -port 8000

## Executar o Cliente

Acesse o diretório `client/`.

    $ go run main.go diffiehellman.go -name client_name -ip localhost -port 8000 -n_messages 100


## Testar Erros

Para testar os erros de `HMAC`, basta descomentar a linha 124 em `client/main.go`.

Para testar os erros de `Nonce`, basta comentar a linha 134 em `client/main.go`.

