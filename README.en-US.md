# Onde 2a dose – backend

Para documentação em Português, veja [README.md](./README.md).

This application serves as a proxy/cache for the data in https://deolhonafila.prefeitura.sp.gov.br/processadores/dados.php.
It provides two endpoints:
1. `POST /data.raw` which mimics the source's behavior for requests and responses
2. `GET /data` which augments the data from the source with latitude and longitude information to be used with a map application (like GoogleMaps)

## Development/Desenvolvimento

For any development, [docker](https://www.docker.com) and a compatible version of [Go](https://golang.org/) (1.16+ in 2021) are required.

In the root of the project, to start a local server:

```bash
make local-run
```

To run automated tests:

```bash
make test
```
