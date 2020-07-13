# Rotterdam-CaaS

[![version](https://img.shields.io/badge/version-0.0.9--SNAPSHOT-blue.svg)]()

&copy; Atos Spain S.A. 2016

Rotterdam-CaaS (_version 0.0.9_) is a component of the European Project class (https://class-project.eu/).

-----------------------

[Description](#description)

[Installation Guide](#installation-guide)

[Usage Guide](#usage-guide)

[LICENSE](#license)

-----------------------

### Description


#### Requirements

#### 1. Installation in Kubernetes

1. Kubernetes

#### 2. Server installation (Linux / Windows)

1. Java version 8 or higher
2. leiningen
3. Connection to a Kubernetes REST API


-----------------------

### Installation Guide

This application can be installed as a service / deployment in kubernetes, or as a server in a Linux / Windows environment.

#### 1. Installation in Kubernetes (_Deployment_)

Deployment in K8s using image from docker: `rsucasas/class-k8-app:0.4`

- Deployment (change the environment variables values according to your configuration):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: class-k8-app
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: class-k8-app
  template:
    metadata:
      labels:
        app: class-k8-app
    spec:
      containers:
        - image: rsucasas/class-k8-app:0.5
          name: class-k8-app
          imagePullPolicy: Always
          ports:
            - containerPort: 8083
          env:
            - name: K8s_API_URL
              value: http://192.168.7.24:8001
            - name: K8s_EXT_IP
              value: 192.168.7.24

```

**K8s_API_URL** Kubernetes REST API URL
**K8s_EXT_IP** External IP used to access the application

- Service

```yaml
kind: Service
apiVersion: v1
metadata:
  name: service-class-k8-app
spec:
  ports:
    - name: http
      port: 8082
      protocol: TCP
      targetPort: 8083
  selector:
    app: class-k8-app
  externalIPs:
    - 192.168.7.24
```

#### 2. Installation in Kubernetes (_StatefulSet_)

- Create a persistent volumen

```yaml
kind: PersistentVolume
apiVersion: v1
metadata:
  name: pv-2
  labels:
    type: local
spec:
  storageClassName: manual
  capacity:
    storage: 2Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/home/rsucasas/k8s-data/v2"
```

- Create a volumen claim. K8s automatically binds this claim to an existing persistent volumen

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-1
spec:
  storageClassName: manual
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

- StatefulSet and Service

```yaml
kind: Service
apiVersion: v1
metadata:
  name: service-class-k8-app
spec:
  ports:
    - name: http
      port: 8082
      protocol: TCP
      targetPort: 8083
  selector:
    app: class-k8-app
  externalIPs:
    - 192.168.7.24
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: class-k8-app
spec:
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: class-k8-app
  template:
    metadata:
      labels:
        app: class-k8-app
    spec:
      containers:
        - image: rsucasas/class-k8-app:0.5
          name: class-k8-app
          imagePullPolicy: Always
          ports:
            - containerPort: 8083
          env:
            - name: K8s_API_URL
              value: http://192.168.7.24:8001
            - name: K8s_EXT_IP
              value: 192.168.7.24
          volumeMounts:
            - mountPath: /tmp/store
              name: class-k8-vol
      volumes:
         - name: class-k8-vol
           persistentVolumeClaim:
              claimName: pvc-1
```

#### 3. Server installation (Linux / Windows)

1. Clone repository (_rotterdam_path_)

-----------------------

### Usage Guide

##### 1. Installation in Kubernetes

If the application was successfully deployed in K8s, it can be accessed in port `8082`. This can be changed in the _Service_ yaml (parameters _ports-port_ and _externalIPs_).

##### 2. Server installation (Linux / Windows)

Launch application:

```bash
cd rotterdam_path
lein ring server
```

By default the API can be accessed in port `8083`. This can be modified in `project.clj` file, before launching the application.


---------------------------------

### LICENSE

`Rotterdam` is licensed under [Apache License, version 2](../../LICENSE).

