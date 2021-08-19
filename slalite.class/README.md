# SLALite #

[![License: Apache v2](https://img.shields.io/badge/License-Apache%20v2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![version](https://img.shields.io/badge/version-0.9.1-blue.svg)]()

This project is licensed under the Apache License 2.0

## Rotterdam & CLASS Project ##

This version of the SLALite is being used in the CLASS project: https://class-project.eu/

Rotterdam relies on the [SLALite asset](https://gitlab.atosresearch.eu/ari/SLALite) to manage the SLA Agreements generated and used in the [CLASS project](https://class-project.eu/) by other components and applicaitons. This (_temp_) version of the SLALite includes an adapter for [Prometheus](https://prometheus.io/) which should be merged with the original project / repository.

This version of the SLALite can be downloaded from the following URL as a docker image: https://cloud.docker.com/u/atosclass/repository/docker/atosclass/slalite

### Environment variables ###

To connect this SLALite with Prometheus and Rotterdam, define the following environment variables (e.g.):

- **MetricsPrometheus** "go_memstats_frees_total, metric_name_2, ..."
- **UrlPrometheus** "http://X.X.X.X:XXX"
- **UrlRotterdam** "http://X.X.X.X:XXX"

----------------------------

## Description ##

The SLALite is a lightweight implementation of an SLA system, inspired by the WS-Agreement standard. Its features are:

* REST interface to manage creation and update of agreements
* Agreements evaluation on background; any breach in the agreement terms generates an SLA violation.
* Configurable monitoring: a monitoring has to be provided externally.
* Configurable repository: a memory repository (for developing purposes) and a mongodb repository are provided, but more can be added.

An agreement is represented by a simple JSON structure (see examples in resources/samples):

```
{
    "id": "2018-000234",
    "name": "an-agreement-name",
    "details":{
        "id": "2018-000234",
        "type": "agreement",
        "name": "an-agreement-name",
        "provider": { "id": "a-provider", "name": "A provider" },
        "client": { "id": "a-client", "name": "A client" },
        "creation": "2018-01-16T17:09:45Z",
        "expiration": "2019-01-17T17:09:45Z",
        "guarantees": [
            {
                "name": "TestGuarantee",
                "constraint": "[execution_time] < 100"
            }
        ]
    }
}
```

## Quick usage guide ##

### Installation ###

Download repository.

Build the Docker image:

    make docker

Run the container:

    docker run -ti -p 8090:8090 slalite:<version>

Stop execution pressing CTRL-C

To run the service under HTTPs, you must change supply a different configuration file and the certificate files. You will find these files in docker/https for debugging purposes. DO NOT USE THE CERT.PEM and KEY.PEM in production!!

    docker run -ti -p 8090:8090 -v $PWD/docker/https:/etc/slalite slalite

### Configuration ###

The SLALite can be configured with a configuration file and with environment
variables. The configuration file is read by default from /etc/slalite and the current
working directory. The `-f` parameter can be used to set the config file location.

```
$ ./SLALite -h
Usage of SLALite:
  -b string
        Filename (w/o extension) of config file (default "slalite")
  -d string
        Directories where to search config files (default "/etc/slalite:.")
  -f string
        Path of configuration file. Overrides -b and -d
```

#### File settings ####

*General settings*

* `singlefile` (default: `false`). Sets if all file settings are read
  from a single file or from several files. For example, when `singlefile=false`,
  the MongoDB settings are read from the file `mongodb.yml`.
* `repository` (default: `memory`). Sets the repository type to use. Set this
  value to `mongodb` to use a MongoDB database.
* `externalIDs` (default: `false`). Set this to true if the repository auto assign
  the IDs of the saved entities.
* `checkPeriod` (default: `60`). Sets the period in seconds of assessments
  executions.
* `CAPath`. Sets the value of a file path containing certificates of trusted
  CAs; to be used to connect as client to SSL servers whose certificate is
  not trusted by default (e.g. self-signed certificates)

*REST interface settings*

* `port` (default: `8090`). Port of REST interface.
* `enableSsl` (default: `false`). Enables the use of SSL on the REST
  interface. The two following variables should be set.
* `sslCertPath` (default: `cert.pem`). Sets the certificate path.
* `sslKeyPath` (default: `key.pem`). Sets the private key path to access the
  certificate.

*MongoDB settings (default file: /etc/slalite/mongodb.yml)*

* `connection` (default: `localhost`). Sets the MongoDB host.
* `database` (default: `slalite`). Sets the MongoDB database name to use.
* `clear_on_boot` (default: `false`). Sets if the database is cleared on
  startup (useful for tests).

#### Env vars  ####

Every file setting can be overriden with the use of environment variables.
The name of the var is the uppercase setting name prefixed with `SLA_`. For
example, to override the check period, set the env var `SLA_CHECKPERIOD`.

### Usage ###

SLALite offers a usual REST API, with an endpoint on /agreements

Add an agreement:

    curl -k -X POST -d @resources/samples/agreement.json https://localhost:8090/agreements

Get agreements:

    curl -k https://localhost:8090/agreements
    curl -k https://localhost:8090/agreements/a02

Add a template:

    curl -k -X POST -d @resources/samples/template.json https://localhost:8090/templates

Get templates:

    curl -k https://localhost:8090/templates
    curl -k https://localhost:8090/templates/t01

Create agreement from template:

    curl -k -X POST -d @resources/samples/create-agreement.json https://localhost:8090/create-agreement

    {"template_id":"t01","agreement_id":"9be511e8-347f-4a40-b784-e80789e4c65b","parameters":{"M":1,"N":100,"agreementname":"An agreement name","client":{"id":"client01","name":"A name of a client"},"provider":{"id":"provider01","name":"A name of a provider"}}}
