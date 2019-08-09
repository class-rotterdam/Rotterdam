# ROTTERDAM REST API

&copy; Atos Spain S.A. 2018


-----------------------

[Description](#description)

[Installation Guide](#installation-guide)

[Docker image](#docker-image)

[Documentation](#documentation)

[Usage Guide](#usage-guide)

[PROJECT CALLS: Relation to other components](#porject-class:-relation-to-other-components)

[LICENSE](#license)

-----------------------

### Description

**Rotterdam** is a native-cloud Infrastructure-as-a-Service (IaaS) facade which facilitates the deployment and life cycle management of containerized tasks on container orchestration platforms. Its main purpose is to upload, organize, run, manage and stop sets of containers (named tasks) through API calls, and abstract all the resource infrastructure details, even the concept of cluster of machines/instances, to micro-service developers (in the case of CLASS, data analytics application/service developers).

The Rotterdam's CaaS API application is a Golang component, responsible for the deployment of these tasks in a Kubernetes cluster.

-----------------------

### Installation Guide

#### Requirements

1. Install `Golang`: https://golang.org/doc/install

#### Installation


-----------------------

### Docker image

-----------------------

### Documentation

##### SWAGGER DOC

Edit file `swagger.yaml`

```bash
swagger serve ./swaggerui/swagger.yaml /p 8001
```

##### SWAGGER UI

Edit file `swaggerui/swagger.json`



-----------------------

### Usage Guide

Edit 'swagger.json' file:

```json
"host": "localhost:8333",
```

```json
 "host": "#H_IP#",
```

To start application run:

```bash
cd path_to_app
go run main.go
```

To start API doc web in port `8001`:

```bash
swagger serve ./swagger.yaml /p 8001
```

-----------------------

### PROJECT CALLS: Relation to other components


-----------------------

### LICENSE

Libraries used in this project:

| library                         | license | url                                   | description |
|---------------------------------|---------|---------------------------------------|-------------|
| github.com/peterbourgon/diskv | MIT     | https://github.com/peterbourgon/diskv | database: persistent key-value store
| github.com/nanobox-io/golang-scribble | Mozilla Public License, version 2.0 | https://github.com/nanobox-io/golang-scribble | A tiny JSON database in Golang |
| - | - | - | - |









export KUBERNETES_SERVICE_HOST=kube4
export KUBERNETES_SERVICE_PORT=8443
