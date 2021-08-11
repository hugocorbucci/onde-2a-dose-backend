# Onde 2a dose – backend

For documentation in English, look at [README.en-US.md](./README.en-US.md).

Esse programa é um proxy/cache para os dados em https://deolhonafila.prefeitura.sp.gov.br/processadores/dados.php.
Ele responde a dois endereços:
1. `POST /data.raw` que se comporta como a fonte tanto para pedidos quanto respostas
2. `GET /data` que incrementa os dados da fonte com latitude e longitude para uso com um aplicativo de mapeamento (como GoogleMaps)

## Desenvolvimento

Para qualquer desenvolvimento, é necessário ter [docker](https://www.docker.com) instalado além de uma versão de [Go](https://golang.org/) compatível (1.16+ em 2021).

Na raíz do projeto, para iniciar um servidor local:

```bash
make local-run
```

Para rodar os testes automatizados:

```bash
make test
```
