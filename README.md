# Storage

[![Coverage Status](https://coveralls.io/repos/github/dadosjusbr/storage/badge.svg?branch=master)](https://coveralls.io/github/dadosjusbr/storage?branch=master)

APIs de armazenamento do dadosjusbr

# Como contribuir com os testes e executa-los.

> ## Mocks

Os mocks são importantes peças para conseguirmos simular comportamentos de objetos. Nesse repositório, os mocks são utilizados no arquivo de teste "client_test.go", onde simulamos o comportamento dos diversos métodos que temos nas interfaces dos subdiretórios do diretório "repo".

Sempre que as interfaces forem modificadas, é necessário gerar os mocks novamente, sobrescrevendo os arquivos antigos pelos novos. Eles são gerados automaticamente, utilizando o passo a passo descrito logo a seguir.

> ### Gerando mocks

Estamos utilizando a biblioteca [gomock](https://github.com/golang/mock) para gerar nossos mocks. Siga a documentação do gomock para conseguir instalar o mockgen no seu computador.

Com o mockgen instalado, basta executar os seguintes comandos:

- Para gerar os mocks a partir da interface do database:

```
mockgen --source ./repo/database/interface.go --destination ./repo/database/database_mock.go
```

- Para gerar os mocks a partir da interface do file storage:

```
mockgen --source ./repo/file_storage/interface.go --destination ./repo/file_storage/file_storage_mock.go
```

Com esses comandos, os mocks antigos são sobrescritos por novos e atualizados.

> ## Subir o banco de teste com o Docker

Para conseguir testar as funcionalidades que acessam diretamente o banco de dados, em /repo/database/postgres_test.go, é necessário ter o banco de dados rodando. Execute os seguintes comandos:

Para buildar a imagem do banco de teste:

```
docker build -t dadosjusbr_test repo/database
```

Para subir o banco de dados:

```
docker run -d --name dadosjusbr_test -p 5432:5432 dadosjusbr_test
```

Em caso de erro, você pode verificar os logs com o seguinte comando:

```
docker logs dadosjusbr_test
```

Para parar o container com o banco de dados, utilize o seguinte comando:

```
docker stop dadosjusbr_test
```

Para remover o container, utilize o seguinte comando:

```
docker rm dadosjusbr_test
```

### Iniciando banco de dados local a partir de seu container

- Após levantar o banco de dados uma única vez, você poderá dar start nele, todas as vezes que ligar o computador, executando o seguinte comando:

```sh
docker start dadosjusbr_test
```

> ## Rodando os testes

Com todas as configurações feitas, basta executar os seguintes comandos para executar os testes:

- Para executar todos:

```
$ go test -v ./...
```

- Para executar um:

```
$ go test -v ${caminho para o arquivo de teste}
```

Executando o comando, você poderá ver as estatisticas relacionadas aos testes, como tempo que demorou a ser concluido, status, diretório, etc...
