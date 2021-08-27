# TCTemplateBack

Este projeto consiste em uma POC que futuramente pode ser tratado como um plano de substituição do TCStartKit.


## 🚀 Começando

Essas instruções permitirão que você obtenha uma cópia do projeto em operação na sua máquina local para fins de desenvolvimento e teste.

### 📋 Pré-requisitos

Ferramentas: 

- [Docker](https://docs.docker.com/desktop/)
- [Golang](https://golang.org/doc/install)
- [Nodemon](https://nodemon.io/)

Projetos:

- [TCNats](https://github.com/tradersclub/TCNats)

### 🔧 Instalação

Nomear o projeto: basta renomear a pasta TCTemplateBack para o nome que desejar, as referências já serão substituidas.

```
$ mv TCTemplateBack ApiRenomeda

$ make rename
```

Configuração do projeto: no arquivo config.yml é necessário trocar o client-id para um id único do seu projeto, além das configurações adicionais.

```
nats:
  client-id: tc-client-poc-arquitetura-id
```


## ⚙️ Executando os testes

Para executar os testes unitários da api, utilizar o comando `make mock` e depois o `make test`.

- `make mock`: cria as implementações das interfaces, com o objetivo de realização da injeção de dependência para execução dos testes unitários.
- `make test`: executa os testes em si.


## 📦 Desenvolvimento

Existe dois comandos básicos para execução do projeto:

- `make run`: wrapper para o `go run main.go` injetando a varíavel de ambiente `VERSION` para ser listada a branch ou versão no endpoint de health.
- `make run-watch`: utiliza das mesmas funcionalidades dos `make run` porém adiciona o live-reload para o código com extenções .go.


## 🗂 Arquitetura

### Descrição dos diretórios e arquivos mais importantes:

- `./event`: O codígo relacionado com a inscrição dos eventos RECEBIDOS a partir NATS e NATS-STREAMING.
- `./event/event.go`: Nesse arquivo está toda parte de registros dos sub-modulos que existem nesse diretório.
- `./handler`: O codígo relacionado com as rotas, middlewares e versionamento da api.
- `./handler/handler.go`: Nesse arquivo está toda parte de registros dos sub-modulos que existem nesse diretório, incluindo versionamento de rotas e gerenciamento de middlewares.
- `./handler/v1`: Este diretório possui a configuração e registro de todos os sub-modulos.
- `./handler/v1/v1.go`: Nesse arquivo está toda parte de registros dos sub-modulos que existem nesse diretório com o path `/v1/**`.
- `./handler/middleware`: Aqui é aonde se encontra os middlewares em geral, como exemplo podemos citar os de injeção de sessão no contexto e o de autorização das rotas.
- `./model`: Este diretório possui todos os arquivos de modelos globais do projeto
- `./repository`: Aqui se encontra todo o código que é utilizado para consultas externas, geralmente usando banco de dados, consultas a outras apis e cache. Obs, nesta arquitetura proposta a localização e utilização do cache e consultas do através do NATS se encontra restrita ao domínio dos serviços.
- `./repository/repository.go`: Arquivo para o registro, configuração e injeção de depêndencias externas nos sub-modulos.
- `./service`: Aqui se encontra todo o código que é utilizado para orquestração e regras de negôcio do serviço. Obs, nesta arquitetura proposta a localização e utilização do cache e consultas do através do NATS se encontra restrita ao domínio dos serviços.
- `./service/service.go`: Arquivo para o registro, configuração e injeção de depêndencias como cache, conexão com o NATS e repositórios nos sub-modulos.
- `./scripts`: Arquivos de scritps em bash em geral.
- `./util`: Sub-modulos necessários para manutenção do projeto em geral.


## ☢️ Boas Práticas

1 - Centralize suas configurações no arquivo `main.go`, e injete o objeto aos modulos necessários.

2 - Somente utilize a pasta `./model` para modelos globais. Quando um modelo pertence a um escopo menor, como exemplo um modelo utilizado para retorno somente em uma única rota específica é aconselhado que seja criado um arquivo dentro desse modulo com a extensão `_model.go` para conter esse código.

`ERRADA`:
```go
// ./model/todo.go
package model

type ResponseTodoAdd struct {
    Add  bool        `json:"added"`
    Todo *model.Todo `json:"todo"`
}

```

`CORRETA`:
```go
// ./handler/v1/todo/todo_model.go
package todo

type responseTodoAdd struct {
    Add  bool        `json:"added"`
    Todo *model.Todo `json:"todo"`
}

```

3 - A boa prática número 2 pode ser extendida para qualquer funcionalidade do sistema, códigos que são utilizados em pacotes específicos devem ficar contidos nesses pacotes.

4 - NUNCA chamar um metódo irmão exportável. Com essa prática tentamos evitar que um código acabe dando voltas ao invés de seguir somente um fluxo além de previnir efeitos colateráis.

`ERRADA`:
```go
type serviceImpl struct {}

func (s *serviceImpl) Update(ctx context.Context, m *model.TODO) (*model.TODO, error) {
	td, err := s.ReadByID(ctx, m.ID) // JAMAIS FAÇA ISSO
    ...
}

func (s *serviceImpl) ReadByID(ctx context.Context, id string) (*model.TODO, error) {
	result := <-s.repository.TODO.ReadByID(ctx, id)
	...
}
```

`CORRETA`:
```go
type serviceImpl struct {}

func (s *serviceImpl) Update(ctx context.Context, m *model.TODO) (*model.TODO, error) {
	result := <-s.repository.TODO.ReadByID(ctx, m.ID)
    ...
}
```

`CORRETA, PORÉM NÃO É RECOMENDADO`:
```go
type serviceImpl struct {}

func (s *serviceImpl) readByID(ctx context.Context, id string) (*model.TODO, error) {
	result := <-s.repository.TODO.ReadByID(ctx, id)
	...
} 

func (s *serviceImpl) Update(ctx context.Context, m *model.TODO) (*model.TODO, error) {
	todo, err := s.readByID(ctx, m.ID)
    ...
}

func (s *serviceImpl) ReadByID(ctx context.Context, m *model.TODO) (*model.TODO, error) {
	todo, err := s.readByID(ctx, m.ID)
    ...
}

```

## 🛠️ Construído com

* [echo](https://echo.labstack.com/) - Framework Web
* [go mod](https://blog.golang.org/using-go-modules) - Dependência
* [viper](https://github.com/spf13/viper) - Configuração 
* [logrus](github.com/sirupsen/logrus) - Log
* [sqlx](https://github.com/jmoiron/sqlx) - Gereciamento de conexão de bancos relacionais
* [validator](github.com/go-playground/validator/v10) - Validador de structs
