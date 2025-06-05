# Serviço de Encurtamento de URLs

Este é um serviço simples de encurtamento de URLs escrito em Go. Ele permite que os usuários encurtem URLs longas e sejam redirecionados para as URLs originais usando os links encurtados

## Funcionalidades
- Encurta URLs fornecidas via requisições HTTP GET.
- Redireciona de URLs curtas para as URLs originais.
- Utiliza criptografia AES (modo CTR) para armazenar URLs com segurança.
- Gera IDs curtos de 6 caracteres usando caracteres alfanuméricos.
- Armazenamento de URLs seguro para threads usando mutex.

## Pré-requisitos
- Go 1.16 ou superior
- Nenhuma dependência externa além da biblioteca padrão do Go

## Uso
1. Inicie o servidor:
   ```bash
   go run main.go
   ```
   O servidor será iniciado em `http://localhost:8080`.

2. Encurtar uma URL:
   - Envie uma requisição GET para `http://localhost:8080/shorten?url=<sua-url>`.
   - Exemplo:
     ```bash
     curl "http://localhost:8080/shorten?url=https://example.com"
     ```
   - A resposta fornecerá uma URL encurtada, por exemplo, `http://localhost:8080/abc123`.

3. Acessar a URL encurtada:
   - Navegue até a URL encurtada (por exemplo, `http://localhost:8080/abc123`) em um navegador ou usando uma ferramenta como `curl`.
   - O servidor redirecionará para a URL original.

## Endpoints
- **GET /shorten?url=<url-original>**  
  Cria uma URL encurtada para a `url-original` fornecida. A URL deve começar com `http://` ou `https://`. Retorna a URL encurtada.

- **GET /<id-curto>**  
  Redireciona para a URL original associada ao `id-curto`. Retorna um erro 404 se o ID curto não for encontrado.

## Segurança
- As URLs são criptografadas usando AES-256 no modo CTR com uma chave de 32 bytes.
- Um vetor de inicialização (IV) aleatório é gerado para cada criptografia.
- As URLs criptografadas são armazenadas em memória com seus respectivos IDs curtos.
- A chave secreta está codificada como `secretaeskey12345678901234567890` (32 bytes). Em um ambiente de produção, substitua por uma chave gerada e armazenada de forma segura.

## Exemplo
```bash
# Encurtar uma URL
curl "http://localhost:8080/shorten?url=https://www.example.com/very/long/url"

# Resposta
URL encurtada: http://localhost:8080/xyz789

# Acessar a URL encurtada
curl -L http://localhost:8080/xyz789
# Redireciona para https://www.example.com/very/long/url
```

## Tratamento de Erros
- Se a URL fornecida estiver vazia ou não começar com `http://` ou `https://`, um erro 400 Bad Request é retornado.
- Se o ID curto não for encontrado, um erro 404 Not Found é retornado.
