# ROTTERDAM CAAS

&copy; Atos Spain S.A. 2018

[![License: Apache v2](https://img.shields.io/badge/License-Apache%20v2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![version](https://img.shields.io/badge/version-1.1.2-blue.svg)]()


-----------------------

[Description](#description)

[Installation Guide](#installation-guide)

[Docker image](#docker-image)

[Documentation](#documentation)

[Usage Guide](#usage-guide)

[Relation to other components](#porject-class:-relation-to-other-components)

[LICENSE](#license)

-----------------------

### Description

**Rotterdam** is a native-cloud Infrastructure-as-a-Service (IaaS) facade which facilitates the deployment and life cycle management of containerized tasks on container orchestration platforms. Its main purpose is to upload, organize, run, manage and stop sets of containers (named tasks) through API calls, and abstract all the resource infrastructure details, even the concept of cluster of machines/instances, to micro-service developers (in the case of CLASS, data analytics application/service developers).

The Rotterdam's CaaS API application is a Golang component, responsible for the deployment of these tasks in a Kubernetes cluster.

A docker image can be downloaded from https://cloud.docker.com/u/atosclass/repository/docker/atosclass/rotterdam-caas.

-----------------------

### Installation Guide

#### Requirements

- Docker (https://docs.docker.com/install/)
- Kubernetes (https://kubernetes.io/docs/setup/) or Openshift-OKD (https://docs.okd.io/latest/install/index.html)

#### Installation

- To run `Rotterdam` in `Docker`:

```bash
docker pull atosclass/rotterdam-caas:tagname
docker run [OPTIONS] atosclass/rotterdam-caas:tagname [COMMAND] [ARG...]
```

- To run `Rotterdam` in Openshift, deploy the [image](https://cloud.docker.com/u/atosclass/repository/docker/atosclass/rotterdam-caas) using the OKD UI.


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

#### Golang

To start application run:

```bash
cd path_to_app
go run main.go
```

To start API doc web in port `8001`:

```bash
swagger serve ./swagger.yaml /p 8001
```


#### Docker

#### Kubernetes

#### Openshift


-----------------------

### LICENSE

Libraries used in this project:

| library                         | license | url                                   | description |
|---------------------------------|---------|---------------------------------------|-------------|
| github.com/tidwall/buntdb | MIT | https://github.com/tidwall/buntdb | BuntDB is a low-level, in-memory, key/value store in pure Go |
| github.com/nikunjy/rules/parser | MIT | https://github.com/nikunjy/rules | Rules engine written in golang with the help of antlr |


`Rotterdam` is licensed under [Apache License, version 2](LICENSE.TXT).
