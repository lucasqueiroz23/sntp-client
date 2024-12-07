# sntp-client

Este repositório contém o código e a documentação do trabalho 1 da disciplina FGA0211 - Fundamentos de Redes de Computadores.
O trabalho foi realizado individualmente pelo aluno Lucas Henrique Lima de Queiroz, de matrícula 190091703.

## Informações sobre o desenvolvimento

- Sistema operacional utilizado: [Ubuntu 20.04.6 LTS](https://releases.ubuntu.com/20.04/) (amd64)
- Ambiente de desenvolvimento utilizado: [Neovim v0.10.0](https://github.com/neovim/neovim)
- Linguagem de programação utilizada: [Golang v1.23.3](https://go.dev/)

## Como compilar o programa?

Para compilar, existem duas opções: de maneira automatizada, utilizando o Makefile, ou manualmente, com o comando `go build`.
Para as duas opções, é necessário [instalar o go na máquina local](https://go.dev/dl/).

As duas alternativas irão gerar um arquivo executável, chamado `client`.

### Opção 1: utilizando o Makefile

1. Tenha o [GNU Make](https://www.gnu.org/software/make/) instalado na sua máquina.
2. Utilize o comando `make`.

### Opção 2: manual

1. Utilize o comando `go build -o client main.go`.

## Como executar o programa?

Para executá-lo, utilize o comando `./client <ip>`, onde `<ip>` deve ser substituído pelo endereço de IP ou hostname do servidor NTP de interesse. 

Se o endereço especificado for o de um servidor NTP válido, o processo mostrará, na stdout, a data/hora fornecida pelo servidor. 



**OBS: Pode-se utilizar hostnames para os servidores, não sendo necessário utilizar o IP.**
**OBS: Caso o IP seja utilizado, deve ser utilizado um endereço IPv4.**

### Testes de execução

Para testar o programa em diversos servidores, utilizei [a lista](https://tf.nist.gov/tf-cgi/servers.cgi) disponibilizada pelo [NIST](https://www.nist.gov/).
Para facilitar o teste, criei os seguintes scripts:

- `run_test.sh`: compilar o programa e executá-lo em um dos servidores disponibilizados;
- `run_a_lot_of_tests.sh`: compilar o programa e executá-lo em vários dos servidores disponibilizados.

## Detalhes de implementação e explicação dos arquivos

Basicamente, o programa recebe o ip do servidor pela linha de comando, envia um datagrama UDP para o servidor especificado, lê o pacote de retorno e realiza um parsing com base no timestamp retornado. Abaixo, segue a explicação do que cada arquivo faz:

- `main.go`: chama os packages definidos abaixo, que realizam a lógica do programa, e no final retorna a mensagem de data e hora.
- `command-line/cli.go`: lê os argumentos de linha de comando e retorna o endereço de IP para o `clientSocket`. 
- `client-socket/client-socket.go`: a partir do IP obtido pelo package `commandLine`, este realiza a comunicação com o servidor. Caso, após 20 segundos, não houver resposta do servidor, será realizada outra tentativa. Se, após 20 segundos, não houver resposta novamente, o programa será fechado com mensagem de erro. Caso a resposta venha, o pacote de retorno será tratado pelo `parser`.
- `parser/parser.go`: a partir do pacote obtido pelo `clientSocket`, realiza a leitura do timestamp em segundos a partir de 1/1/1900 e retorna uma string com a data que o servidor mandou.
- `error-handling/error-handling.go`: realiza um tratamento geral de erros. É chamado por todos os packages acima.

É importante observar que o `clientSocket` utiliza o pacote [`net`](https://pkg.go.dev/net), que faz parte da linguagem Go.

## Limitações conhecidas

- Não há conversão de [timezones](https://en.wikipedia.org/wiki/Coordinated_Universal_Time). Todos os servidores que testei, por exemplo, retornam um horário em UTC+0. Não será realizada a conversão para UTC-3 (horário de Brasília).
- Existe a possibilidade da data não estar 100% correta em caso de anos bissextos.
- Na linha de comando, é importante colocar **apenas o endereço de ip do servidor**. Ou seja: **não se deve especificar a porta**, pois ela já foi especificada no enunciado do trabalho (porta 123), e, portanto, foi hard-coded.
- Caso um endereço de IP (e não o hostname) seja fornecido, ele deve estar na forma IPv4. Endereços IPv6 não funcionarão (mas seus hostnames funcionarão, caso existam).
