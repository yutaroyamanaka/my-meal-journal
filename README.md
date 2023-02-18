# My-Meal-Journal (in Progress)

![coverage](docs/coverage.svg) ![coverage](docs/ratio.svg)

## Usage

```sh
$ make up
$ curl -X POST localhost:8080/add -H 'Content-Type: application/json' -d '{"name": "pizza", "category": 1}'
```

| Category | Description |
| --- | ----------- |
| 0 | Breakfast |
| 1 | Lunch |
| 2 | Dinner |
| 3 | Others |

### Run tests

```sh
$ make test
```
