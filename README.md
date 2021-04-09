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

Next we define a new "kind", aka custom-resource. 

```
$ operator-sdk create api --group cache --version v1alpha1 --kind Mysql --resource --controller       # kinds need to start with an uppercase
Writing scaffold for you to edit...
api/v1alpha1/mysql_types.go
controllers/mysql_controller.go
Running make:
$ make
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
```

Note, an operator can manage 1 or more CRDs, these CRDs are referred to as "APIs" as indicated in the above example.

see:  https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/

```
git status
On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   PROJECT
        modified:   go.mod
        modified:   main.go

Untracked files:
  (use "git add <file>..." to include in what will be committed)
        api/
        config/crd/
        config/rbac/mysql_editor_role.yaml
        config/rbac/mysql_viewer_role.yaml
        config/samples/
        controllers/

no changes added to commit (use "git add" and/or "git commit -a")
```

(git commit no2)

Now we can create our operator image (see https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/):

```
$ docker login quay.io -u sher.chowdhury@ibm.com -p xxxxxxxxx 
$ # export USERNAME=sher_chowdhury0
$ export OPERATOR_IMG="quay.io/sher_chowdhury0/mysql-operator2:v0.0.1"
$ make docker-build docker-push IMG=$OPERATOR_IMG
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen "crd:trivialVersions=true,preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
mkdir -p /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/testbin
test -f /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/testbin/setup-envtest.sh || curl -sSLo /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/testbin/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.7.0/hack/setup-envtest.sh
source /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/testbin/setup-envtest.sh; fetch_envtest_tools /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/testbin; setup_envtest_env /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/testbin; go test ./... -coverprofile cover.out
fetching envtest tools@1.19.2 (into '/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/testbin')
x bin/
x bin/etcd
x bin/kubectl
x bin/kube-apiserver
setting up env vars
?       github.com/Sher-Chowdhury/mysql-operator2       [no test files]
?       github.com/Sher-Chowdhury/mysql-operator2/api/v1alpha1  [no test files]
ok      github.com/Sher-Chowdhury/mysql-operator2/controllers   9.719s  coverage: 0.0% of statements
docker build -t quay.io/sher_chowdhury0/mysql-operator2:v0.0.1 .
[+] Building 90.5s (18/18) FINISHED                                                                                                                                                       
 => [internal] load build definition from Dockerfile                                                                                   0.1s
 => => transferring dockerfile: 835B                                                                                                   0.0s
 => [internal] load .dockerignore                                                                                                      0.0s
 => => transferring context: 193B                                                                                                      0.0s
 => [internal] load metadata for gcr.io/distroless/static:nonroot                                                                      1.4s
 => [internal] load metadata for docker.io/library/golang:1.15                                                                         1.2s
 => [auth] library/golang:pull token for registry-1.docker.io                                                                          0.0s
 => [stage-1 1/3] FROM gcr.io/distroless/static:nonroot@sha256:cd784033c94dd30546456f35de8e128390ae15c48cbee5eb7e3306857ec17631        1.7s
 => => resolve gcr.io/distroless/static:nonroot@sha256:cd784033c94dd30546456f35de8e128390ae15c48cbee5eb7e3306857ec17631                0.0s
 => => sha256:cd784033c94dd30546456f35de8e128390ae15c48cbee5eb7e3306857ec17631 1.67kB / 1.67kB                                         0.0s
 => => sha256:8f3b47c7984464f417f9d5f5e232ac3fae6453e84f053724fef457c4ba67ceaf 426B / 426B                                             0.0s
 => => sha256:fb7b4da47366a77c2cd1973d031835127eeb6efb5d255dd2ebf7ba573e601825 478B / 478B                                             0.0s
 => => sha256:5dea5ec2316d4a067b946b15c3c4f140b4f2ad607e73e9bc41b673ee5ebb99a3 657.65kB / 657.65kB                                     1.3s
 => => extracting sha256:5dea5ec2316d4a067b946b15c3c4f140b4f2ad607e73e9bc41b673ee5ebb99a3                                              0.2s
 => [builder 1/9] FROM docker.io/library/golang:1.15@sha256:a4dbaabb67af3cc2a41168a32cbd7035738692b38bbb5392498ec34dbee9216b          24.4s
 => => resolve docker.io/library/golang:1.15@sha256:a4dbaabb67af3cc2a41168a32cbd7035738692b38bbb5392498ec34dbee9216b                   0.0s
 => => sha256:a4dbaabb67af3cc2a41168a32cbd7035738692b38bbb5392498ec34dbee9216b 2.36kB / 2.36kB                                         0.0s
 => => sha256:2b7f8090d63788525c18f580b4ab30a7ca3fd381ca16227e5cfbe4d0443ee71e 1.79kB / 1.79kB                                         0.0s
 => => sha256:874f8671ee4ec24d7e8102cdcc3dd215027ded05c6789852bea4fc4d1135668e 7.10kB / 7.10kB                                         0.0s
 => => sha256:004f1eed87df3f75f5e2a1a649fa7edd7f713d1300532fd0909bb39cd48437d7 50.43MB / 50.43MB                                       4.0s
 => => sha256:5d6f1e8117dbb1c6a57603cb4f321a861a08105a81bcc6b01b0ec2b78c8523a5 7.83MB / 7.83MB                                         1.5s
 => => sha256:48c2faf66abec3dce9f54d6722ff592fce6dd4fb58a0d0b72282936c6598a3b3 10.00MB / 10.00MB                                       0.6s
 => => sha256:234b70d0479d7f16d7ee8d04e4ffdacc57d7d14313faf59d332f18b2e9418743 51.84MB / 51.84MB                                       7.4s
 => => sha256:f5e9f83ff9bcf98a081ea281823c299a293e0870c4fd132c22cf425b01bd310e 68.74MB / 68.74MB                                      14.9s
 => => sha256:44b45011f179e4d84c65eee524de7bd86bc9c3abca653bc7e87767c812a0d421 121.08MB / 121.08MB                                    15.5s
 => => extracting sha256:004f1eed87df3f75f5e2a1a649fa7edd7f713d1300532fd0909bb39cd48437d7                                              3.0s
 => => sha256:06aa4da2eeb6b7ecb94028f6e8eaf5ecc9ff4a8b60c1d168235e9e04031fe995 156B / 156B                                             7.5s
 => => extracting sha256:5d6f1e8117dbb1c6a57603cb4f321a861a08105a81bcc6b01b0ec2b78c8523a5                                              0.3s
 => => extracting sha256:48c2faf66abec3dce9f54d6722ff592fce6dd4fb58a0d0b72282936c6598a3b3                                              0.3s
 => => extracting sha256:234b70d0479d7f16d7ee8d04e4ffdacc57d7d14313faf59d332f18b2e9418743                                              3.2s
 => => extracting sha256:f5e9f83ff9bcf98a081ea281823c299a293e0870c4fd132c22cf425b01bd310e                                              1.9s
 => => extracting sha256:44b45011f179e4d84c65eee524de7bd86bc9c3abca653bc7e87767c812a0d421                                              5.8s
 => => extracting sha256:06aa4da2eeb6b7ecb94028f6e8eaf5ecc9ff4a8b60c1d168235e9e04031fe995                                              0.0s
 => [internal] load build context                                                                                                      0.0s
 => => transferring context: 77.40kB                                                                                                   0.0s
 => [builder 2/9] WORKDIR /workspace                                                                                                   0.3s
 => [builder 3/9] COPY go.mod go.mod                                                                                                   0.0s
 => [builder 4/9] COPY go.sum go.sum                                                                                                   0.0s
 => [builder 5/9] RUN go mod download                                                                                                 15.5s
 => [builder 6/9] COPY main.go main.go                                                                                                 0.0s
 => [builder 7/9] COPY api/ api/                                                                                                       0.0s
 => [builder 8/9] COPY controllers/ controllers/                                                                                       0.0s
 => [builder 9/9] RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go                             47.3s
 => [stage-1 2/3] COPY --from=builder /workspace/manager .                                                                             0.1s
 => exporting to image                                                                                                                 0.2s
 => => exporting layers                                                                                                                0.2s
 => => writing image sha256:1246fe9f0585f9f5b55a4f4f1c96ee557a2e55f92f24e28b8b279c81496f3556                                           0.0s
 => => naming to quay.io/sher_chowdhury0/mysql-operator2:v0.0.1                                                                        0.0s
docker push quay.io/sher_chowdhury0/mysql-operator2:v0.0.1
The push refers to repository [quay.io/sher_chowdhury0/mysql-operator2]
e79608ea9009: Pushed 
417cb9b79ade: Pushed 
v0.0.1: digest: sha256:8dd417a0b4241f245b6a44ab5e7204d30ddef2d29bab34b13671f6720614ade6 size: 739
```

This pushes up the image: https://quay.io/repository/sher_chowdhury0/mysql-operator2

Which you can pull down using: `docker pull quay.io/sher_chowdhury0/mysql-operator2:v0.0.1`

Also, very importantly this step created our CRD:

```
$ git status                                             
On branch main
Your branch is up to date with 'origin/main'.

Untracked files:
  (use "git add <file>..." to include in what will be committed)
        mysql-operator2/config/crd/bases/
        mysql-operator2/config/rbac/role.yaml

no changes added to commit (use "git add" and/or "git commit -a")


╭(☸️ |default/api-vocable-cp-fyre-ibm-com:6443/kube:admin:default) sher  ~/github/mysql-operator2   main ●  
╰➤ ls -l mysql-operator2/config/crd/bases/
total 8
-rw-r--r--  1 sherchowdhury  staff  1920  9 Apr 11:02 cache.codingbee.net_mysqls.yaml
╭(☸️ |default/api-vocable-cp-fyre-ibm-com:6443/kube:admin:default) sher  ~/github/mysql-operator2   main ●  
╰➤ cat mysql-operator2/config/crd/bases/cache.codingbee.net_mysqls.yaml

---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.4.1
  creationTimestamp: null
  name: mysqls.cache.codingbee.net
spec:
  group: cache.codingbee.net
  names:
    kind: Mysql
    listKind: MysqlList
    plural: mysqls
    singular: mysql
  scope: Namespaced
  versions:
  - name: v1alpha1
    schema:
      openAPIV3Schema:
        description: Mysql is the Schema for the mysqls API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: MysqlSpec defines the desired state of Mysql
            properties:
              foo:
                description: Foo is an example field of Mysql. Edit Mysql_types.go
                  to remove/update
                type: string
            type: object
          status:
            description: MysqlStatus defines the observed state of Mysql
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
```

(git commit no 3)


Side note, to check if olm is installed on your cluster, first log into your cluster, then run:

```
$ operator-sdk olm status --olm-namespace=openshift-operator-lifecycle-manager
INFO[0004] Fetching CRDs for version "0.17.0"
INFO[0004] Using locally stored resource manifests
INFO[0008] Successfully got OLM status for version "0.17.0"

NAME                                            NAMESPACE    KIND                        STATUS
operators.operators.coreos.com                               CustomResourceDefinition    Installed
operatorgroups.operators.coreos.com                          CustomResourceDefinition    Installed
installplans.operators.coreos.com                            CustomResourceDefinition    Installed
clusterserviceversions.operators.coreos.com                  CustomResourceDefinition    Installed
subscriptions.operators.coreos.com                           CustomResourceDefinition    Installed
system:controller:operator-lifecycle-manager                 ClusterRole                 Installed
aggregate-olm-edit                                           ClusterRole                 Installed
aggregate-olm-view                                           ClusterRole                 Installed
catalogsources.operators.coreos.com                          CustomResourceDefinition    Installed
olm                                                          Namespace                   namespaces "olm" not found
olm-operator-binding-olm                                     ClusterRoleBinding          clusterrolebindings.rbac.authorization.k8s.io "olm-operator-binding-olm" not found
olm-operator                                    olm          Deployment                  deployments.apps "olm-operator" not found
catalog-operator                                olm          Deployment                  deployments.apps "catalog-operator" not found
olm-operator-serviceaccount                     olm          ServiceAccount              serviceaccounts "olm-operator-serviceaccount" not found
operators                                                    Namespace                   namespaces "operators" not found
global-operators                                operators    OperatorGroup               operatorgroups.operators.coreos.com "global-operators" not found
olm-operators                                   olm          OperatorGroup               operatorgroups.operators.coreos.com "olm-operators" not found
packageserver                                   olm          ClusterServiceVersion       clusterserviceversions.operators.coreos.com "packageserver" not found
```

olm is a bit like kubernete's version of a package manager, e.g. like yum, dnf, pip...etc, but for kubernetes. Olm is installed by default on openshift clusters. 

Next we create the bundle image for our operator (https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/):

```
$ echo $OPERATOR_IMG
quay.io/sher_chowdhury0/mysql-operator2:v0.0.1
make bundle IMG=$OPERATOR_IMG

/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen "crd:trivialVersions=true,preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go: creating new go.mod: module tmp
Downloading sigs.k8s.io/kustomize/kustomize/v3@v3.8.7
go: downloading sigs.k8s.io/kustomize/kustomize/v3 v3.8.7
go: downloading sigs.k8s.io/kustomize/api v0.6.5
go: downloading sigs.k8s.io/kustomize/cmd/config v0.8.5
go: downloading k8s.io/client-go v0.18.10
go: downloading sigs.k8s.io/kustomize/kyaml v0.9.4
go: downloading k8s.io/apimachinery v0.18.10
go: downloading github.com/hashicorp/go-multierror v1.1.0
go: downloading github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00
go: downloading k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
go: downloading github.com/yujunz/go-getter v1.4.1-lite
go: downloading go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5
go: downloading golang.org/x/net v0.0.0-20200625001655-4c5254603344
go: downloading github.com/qri-io/starlib v0.4.2-0.20200213133954-ff2e8cd5ef8d
go: downloading github.com/go-openapi/strfmt v0.19.5
go: downloading k8s.io/api v0.18.10
go: downloading github.com/go-openapi/validate v0.19.8
go: downloading go.mongodb.org/mongo-driver v1.1.2
go: downloading github.com/hashicorp/go-version v1.1.0
go: downloading github.com/hashicorp/go-safetemp v1.0.0
go: downloading github.com/mitchellh/go-testing-interface v1.0.0
go: downloading github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
go: downloading github.com/go-openapi/errors v0.19.2
go: downloading github.com/ulikunitz/xz v0.5.5
go: downloading github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d
go: downloading github.com/hashicorp/go-cleanhttp v0.5.0
go: downloading github.com/go-openapi/analysis v0.19.5
go: downloading github.com/go-openapi/runtime v0.19.4
go: downloading golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd
go: downloading github.com/go-stack/stack v1.8.0
go: downloading github.com/go-openapi/loads v0.19.4
operator-sdk generate kustomize manifests -q

Display name for the operator (required): 
> mysql

Description for the operator (required): 
> install mysql dbs

Provider's name for the operator (required): 
> mysqlprovider

Any relevant URL for the provider name (optional): 
>       

Comma-separated list of keywords for your operator (required): 
> db,mysql

Comma-separated list of maintainers and their emails (e.g. 'name1:email1, name2:email2') (required): 
> sher.chowdhury@ibm.com
cd config/manager && /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/kustomize edit set image controller=quay.io/sher_chowdhury0/mysql-operator2:v0.0.1
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/kustomize build config/manifests | operator-sdk generate bundle -q --overwrite --version 0.0.1  
INFO[0000] Building annotations.yaml                    
INFO[0000] Writing annotations.yaml in /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bundle/metadata 
INFO[0000] Building Dockerfile                          
INFO[0000] Writing bundle.Dockerfile in /Users/sherchowdhury/github/mysql-operator2/mysql-operator2 
operator-sdk bundle validate ./bundle
INFO[0000] Found annotations file                        bundle-dir=bundle container-tool=docker
INFO[0000] Could not find optional dependencies file     bundle-dir=bundle container-tool=docker
INFO[0000] All validation tests have completed successfully 
```

This created the following:

```
$ git status
On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   config/manager/kustomization.yaml

Untracked files:
  (use "git add <file>..." to include in what will be committed)
        bundle.Dockerfile
        bundle/
        config/manifests/

no changes added to commit (use "git add" and/or "git commit -a")
```

This is the point where the clusterserviceversion.yaml got created:

```
$ cat config/manifests/bases/mysql-operator2.clusterserviceversion.yaml 
apiVersion: operators.coreos.com/v1alpha1
kind: ClusterServiceVersion
metadata:
  annotations:
    alm-examples: '[]'
    capabilities: Basic Install
  name: mysql-operator2.v0.0.0
  namespace: placeholder
spec:
  apiservicedefinitions: {}
  customresourcedefinitions:
    owned:
    - description: Mysql is the Schema for the mysqls API
      displayName: Mysql
      kind: Mysql
      name: mysqls.cache.codingbee.net
      version: v1alpha1
  description: install mysql dbs
  displayName: mysql
  icon:
  - base64data: ""
    mediatype: ""
  install:
    spec:
      deployments: null
    strategy: ""
  installModes:
  - supported: false
    type: OwnNamespace
  - supported: false
    type: SingleNamespace
  - supported: false
    type: MultiNamespace
  - supported: true
    type: AllNamespaces
  keywords:
  - db
  - mysql
  links:
  - name: Mysql Operator2
    url: https://mysql-operator2.domain
  maturity: alpha
  provider:
    name: mysqlprovider
  version: 0.0.0
```

Also the following looks interesting too:

```
$ cat config/manager/kustomization.yaml 
resources:
- manager.yaml

generatorOptions:
  disableNameSuffixHash: true

configMapGenerator:
- files:
  - controller_manager_config.yaml
  name: manager-config
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
images:
- name: controller
  newName: quay.io/sher_chowdhury0/mysql-operator2
  newTag: v0.0.1
```

There's also this file too:

```
$ tree config/manifests                
config/manifests
├── bases
│   └── mysql-operator2.clusterserviceversion.yaml
└── kustomization.yaml

1 directory, 2 files

$ cat config/manifests/kustomization.yaml 
resources:
- ../default
- ../samples
- ../scorecard
```

These paths seems to refer to:

```
$ ls -l config
total 0
drwx------   5 sherchowdhury  staff  160  9 Apr 10:31 certmanager
drwx------   6 sherchowdhury  staff  192  9 Apr 11:02 crd
drwx------   5 sherchowdhury  staff  160  9 Apr 10:31 default              <==
drwx------   5 sherchowdhury  staff  160  9 Apr 10:31 manager
drwxr-xr-x   4 sherchowdhury  staff  128  9 Apr 13:03 manifests
drwx------   4 sherchowdhury  staff  128  9 Apr 10:31 prometheus
drwx------  13 sherchowdhury  staff  416  9 Apr 11:02 rbac
drwx------   4 sherchowdhury  staff  128  9 Apr 10:44 samples              <==
drwxr-xr-x   5 sherchowdhury  staff  160  9 Apr 10:31 scorecard            <==
```

Next we create the bundle image and push it up (https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/):

```
$ # Note the "-bundle" component in the image name below.
$ export BUNDLE_IMG="quay.io/sher_chowdhury0/mysql-operator2-bundle:v0.0.1"

$ make bundle-build BUNDLE_IMG=$BUNDLE_IMG               

docker build -f bundle.Dockerfile -t quay.io/sher_chowdhury0/mysql-operator2-bundle:v0.0.1 .
[+] Building 0.2s (7/7) FINISHED                                                                                          
 => [internal] load build definition from bundle.Dockerfile                                                          0.0s
 => => transferring dockerfile: 44B                                                                                  0.0s
 => [internal] load .dockerignore                                                                                    0.0s
 => => transferring context: 35B                                                                                     0.0s
 => [internal] load build context                                                                                    0.0s
 => => transferring context: 727B                                                                                    0.0s
 => CACHED [1/3] COPY bundle/manifests /manifests/                                                                   0.0s
 => CACHED [2/3] COPY bundle/metadata /metadata/                                                                     0.0s
 => CACHED [3/3] COPY bundle/tests/scorecard /tests/scorecard/                                                       0.0s
 => exporting to image                                                                                               0.0s
 => => exporting layers                                                                                              0.0s
 => => writing image sha256:92526e73c1cca5140b8c6113be4d9dbb05cb797eef6ce0d7db7255e2395c60c8                         0.0s
 => => naming to quay.io/sher_chowdhury0/mysql-operator2-bundle:v0.0.1                                               0.0s

$ make docker-push IMG=$BUNDLE_IMG

docker push quay.io/sher_chowdhury0/mysql-operator2-bundle:v0.0.1
The push refers to repository [quay.io/sher_chowdhury0/mysql-operator2-bundle]
d02927543e04: Pushed 
9f25142cbc88: Pushed 
7fdabab67fdf: Pushed 
v0.0.1: digest: sha256:aebdf842a4c5da68535a3fb9271e711820fe8f90422a265936ee66da26637ade size: 939 
```

This bundle image is specific to the version of the operator in question. When you install the operator into your namespace, the operator's catalogsource pod works out which bundle image you need, pulls down that bundle image and run's a pod-job with that bundle image. This job is what installs the operator into your namespace. by installing I mean, creating the serviceaccount, role, rolebinding, and deployment resources in your namespace. 


Now you can test this bundle image by running:

```
operator-sdk run bundle $BUNDLE_IMG
```


