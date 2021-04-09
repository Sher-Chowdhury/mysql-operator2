# mysql-operator2

## Overview

This guide is broken into 2 stages. 

1. build a mysql operator
2. build a wordpress operator

## The approach

We are going to start with a simple hello world example. And incrementally modify it into a fully functional mysql operator. 


We will start with implementing quick+dirty techniques, but then gradually improve and build on that to incorporate best practice and achieve level 1 maturity. 





mysql kubernetes operator built using the operator-sdk

```
brew install operator-sdk
```

Also see: https://sdk.operatorframework.io/docs/building-operators/golang/installation/


```
$ operator-sdk version
operator-sdk version: "v1.4.0", commit: "67f9c8b888887d18cd38bb6fd85cf3cf5b94fd99", kubernetes version: "1.19.4", go version: "go1.15.5", GOOS: "darwin", GOARCH: "amd64"
```

```
$ git remote -v
origin  git@github.com:Sher-Chowdhury/mysql-operator2.git (fetch)
origin  git@github.com:Sher-Chowdhury/mysql-operator2.git (push)
$ pwd
/Users/sherchowdhury/github/mysql-operator2
```

Then followed instructions in: https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/

The output is:

```
$ mkdir mysql-operator2
$ cd mysql-operator2 
╭(☸️ |default/api-vocable-cp-fyre-ibm-com:6443/kube:admin:default) sher  ~/github/mysql-operator2/mysql-operator2   main  
╰➤ operator-sdk init --domain codingbee.net --repo github.com/Sher-Chowdhury/mysql-operator2
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.7.0
go: downloading sigs.k8s.io/controller-runtime v0.7.0
go: downloading k8s.io/apimachinery v0.19.2
go: downloading k8s.io/utils v0.0.0-20200912215256-4140de9c8800
go: downloading k8s.io/client-go v0.19.2
go: downloading k8s.io/api v0.19.2
go: downloading k8s.io/apiextensions-apiserver v0.19.2
go: downloading k8s.io/klog/v2 v2.2.0
go: downloading golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
go: downloading github.com/imdario/mergo v0.3.10
go: downloading github.com/googleapis/gnostic v0.5.1
go: downloading golang.org/x/sys v0.0.0-20200622214017-ed371f2e16b4
go: downloading sigs.k8s.io/structured-merge-diff/v4 v4.0.1
go: downloading k8s.io/component-base v0.19.2
go: downloading google.golang.org/protobuf v1.24.0
go: downloading golang.org/x/net v0.0.0-20200707034311-ab3426394381
go: downloading golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
go: downloading github.com/golang/groupcache v0.0.0-20191227052852-215e87163ea7
go: downloading github.com/google/go-cmp v0.5.2
go: downloading k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
go: downloading gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
Update go.mod:
$ go mod tidy
go: downloading github.com/onsi/ginkgo v1.14.1
go: downloading github.com/Azure/go-autorest/autorest v0.9.6
go: downloading github.com/go-logr/zapr v0.2.0
go: downloading go.uber.org/goleak v1.1.10
go: downloading github.com/onsi/gomega v1.10.2
go: downloading golang.org/x/tools v0.0.0-20200616133436-c1934b75d054
go: downloading github.com/Azure/go-autorest/autorest/adal v0.8.2
go: downloading cloud.google.com/go v0.51.0
Running make:
$ make
go: creating new go.mod: module tmp
Downloading sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1
go: downloading sigs.k8s.io/controller-tools v0.4.1
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.4.1
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
Next: define a resource with:
$ operator-sdk create api
```

This created:

```
$ cd ..
$ tree .
.
├── LICENSE
├── README.md
└── mysql-operator2
    ├── Dockerfile
    ├── Makefile
    ├── PROJECT
    ├── bin
    │   ├── controller-gen
    │   └── manager
    ├── config
    │   ├── certmanager
    │   │   ├── certificate.yaml
    │   │   ├── kustomization.yaml
    │   │   └── kustomizeconfig.yaml
    │   ├── default
    │   │   ├── kustomization.yaml
    │   │   ├── manager_auth_proxy_patch.yaml
    │   │   └── manager_config_patch.yaml
    │   ├── manager
    │   │   ├── controller_manager_config.yaml
    │   │   ├── kustomization.yaml
    │   │   └── manager.yaml
    │   ├── prometheus
    │   │   ├── kustomization.yaml
    │   │   └── monitor.yaml
    │   ├── rbac
    │   │   ├── auth_proxy_client_clusterrole.yaml
    │   │   ├── auth_proxy_role.yaml
    │   │   ├── auth_proxy_role_binding.yaml
    │   │   ├── auth_proxy_service.yaml
    │   │   ├── kustomization.yaml
    │   │   ├── leader_election_role.yaml
    │   │   ├── leader_election_role_binding.yaml
    │   │   └── role_binding.yaml
    │   └── scorecard
    │       ├── bases
    │       │   └── config.yaml
    │       ├── kustomization.yaml
    │       └── patches
    │           ├── basic.config.yaml
    │           └── olm.config.yaml
    ├── go.mod
    ├── go.sum
    ├── hack
    │   └── boilerplate.go.txt
    └── main.go

12 directories, 34 files
```

(git commit no1)