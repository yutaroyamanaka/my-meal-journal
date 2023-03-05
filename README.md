# My-Meal-Journal

![coverage](docs/coverage.svg) ![coverage](docs/ratio.svg)

My-Meal-Journal is a simple API for recording meal information. This project uses Terraform to set up AWS resources and deploys the application through GitHub Actions.

## Deployment Overview

### Prerequisites

- IAM user with AdministratorAccess policy
- Repository Secrets and Variables

Here is AWS architecture overview of this project.

<img width="960" alt="Screen Shot 2023-03-05 at 14 48 12" src="https://user-images.githubusercontent.com/34850155/222943883-9a0d6eaa-85dc-4ae2-9996-6b4aa03b2439.png">

The primary future task is to migrate to **GitOps** from CIOps for more security.
See also [backlogs](https://github.com/yutaroyamanaka/my-meal-journal/issues/39) for other future work.

> "Note: Although GitOps is recognized as a superior approach, this project currently uses CIOps due to its simplicity and resource availability at the time.  I acknowledge the potential benefits of GitOps and may explore this approach in the future, but for the current project objectives, CIOps was deemed the best fit."

## Application Usage

### Production

```sh
$ curl -X POST <LOAD_BALANCER_ENDPOINT>/add -H 'Content-Type: application/json' -d '{"name": "pizza", "category": 1}' 
```

### Local

```sh
# you need docker to debug locally
$ make up
$ curl -X POST localhost:8080/add -H 'Content-Type: application/json' -d '{"name": "pizza", "category": 1}'
```

category means the type of meal as follows.

| Category | Description |
| --- | ----------- |
| 0 | Breakfast |
| 1 | Lunch |
| 2 | Dinner |
| 3 | Others |

## Run tests

```sh
# you need docker to run tests
$ make test
```
