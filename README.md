# Storage

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
mockgen --source ./repo/database/idatabase_repository.go --destination ./repo/database/database_mock.go
```

- Para gerar os mocks a partir da interface do file storage:

```
mockgen --source ./repo/file_storage/istorage_repository.go --destination ./repo/file_storage/file_storage_mock.go
```

Com esses comandos, os mocks antigos são sobrescritos por novos e atualizados.

> ## Subir o banco de teste com o Docker

Para conseguir testar as funcionalidades que acessam diretamente o banco de dados, em /repo/database/postgres_test.go, é necessário ter o banco de dados rodando. Os passos são:

- Renomeie o arquivo .env.example para .env
- Execute o seguinte comando:

```
docker compose up -d --build
```

Em caso de erro, você pode verificar os logs com o seguinte comando:

```
docker logs postgres_test
```

### Iniciando banco de dados local a partir de seu container

- Após levantar o banco de dados uma única vez, você poderá dar start nele, todas as vezes que ligar o computador, executando o seguinte comando:

```sh
docker start postgres_test
```

### Removendo banco de dados local (a partir do docker-compose deste repositório)

```sh
docker-compose  rm -sf
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
