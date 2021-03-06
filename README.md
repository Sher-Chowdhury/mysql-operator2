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
$ operator-sdk init --domain codingbee.net --repo github.com/Sher-Chowdhury/mysql-operator2
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

For more info, see - https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

This created:

```
$ cd ..
$ tree .
.
????????? LICENSE
????????? README.md
????????? mysql-operator2
    ????????? Dockerfile
    ????????? Makefile
    ????????? PROJECT
    ????????? bin
    ???   ????????? controller-gen
    ???   ????????? manager
    ????????? config
    ???   ????????? certmanager
    ???   ???   ????????? certificate.yaml
    ???   ???   ????????? kustomization.yaml
    ???   ???   ????????? kustomizeconfig.yaml
    ???   ????????? default
    ???   ???   ????????? kustomization.yaml
    ???   ???   ????????? manager_auth_proxy_patch.yaml
    ???   ???   ????????? manager_config_patch.yaml
    ???   ????????? manager
    ???   ???   ????????? controller_manager_config.yaml
    ???   ???   ????????? kustomization.yaml
    ???   ???   ????????? manager.yaml
    ???   ????????? prometheus
    ???   ???   ????????? kustomization.yaml
    ???   ???   ????????? monitor.yaml
    ???   ????????? rbac
    ???   ???   ????????? auth_proxy_client_clusterrole.yaml
    ???   ???   ????????? auth_proxy_role.yaml
    ???   ???   ????????? auth_proxy_role_binding.yaml
    ???   ???   ????????? auth_proxy_service.yaml
    ???   ???   ????????? kustomization.yaml
    ???   ???   ????????? leader_election_role.yaml
    ???   ???   ????????? leader_election_role_binding.yaml
    ???   ???   ????????? role_binding.yaml
    ???   ????????? scorecard
    ???       ????????? bases
    ???       ???   ????????? config.yaml
    ???       ????????? kustomization.yaml
    ???       ????????? patches
    ???           ????????? basic.config.yaml
    ???           ????????? olm.config.yaml
    ????????? go.mod
    ????????? go.sum
    ????????? hack
    ???   ????????? boilerplate.go.txt
    ????????? main.go

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

Note, an operator can manage 1 or more CRDs, these CRDs are referred to as "APIs" as indicated in the above example. In this example we created
the kind "mysql" but maybe we can generate other secondary-like crds, e.g. `mysqlconfig` and `mysqllogincreds`. 

each "kind" (aka operand), e.g. mysqlconfig comes with a pair of files:

- api/v1alpha1/mysqlconfig_types.go       # *_types.go file is where we define the crd structure 
- controllers/mysqlconfig_controller.go   # *_controller.go file is where we define what child resources should be created. 

So far we have created a kind called "mysql" so we have:

- api/v1alpha1/mysql_types.go       # *_types.go file is where we define the crd structure 
- controllers/mysql_controller.go   # *_controller.go file is where we define what child resources should be created. 



see:  https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/

Also see: https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/#create-a-new-api-and-controller


The above command created the following  file:

```
$ cat api/v1alpha1/mysql_types.go
/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MysqlSpec defines the desired state of Mysql
type MysqlSpec struct {
        // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
        // Important: Run "make" to regenerate code after modifying this file

        // Foo is an example field of Mysql. Edit Mysql_types.go to remove/update
        Foo string `json:"foo,omitempty"`
}

// MysqlStatus defines the observed state of Mysql
type MysqlStatus struct {
        // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
        // Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Mysql is the Schema for the mysqls API
type Mysql struct {
        metav1.TypeMeta   `json:",inline"`
        metav1.ObjectMeta `json:"metadata,omitempty"`

        Spec   MysqlSpec   `json:"spec,omitempty"`
        Status MysqlStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MysqlList contains a list of Mysql
type MysqlList struct {
        metav1.TypeMeta `json:",inline"`
        metav1.ListMeta `json:"metadata,omitempty"`
        Items           []Mysql `json:"items"`
}

func init() {
        SchemeBuilder.Register(&Mysql{}, &MysqlList{})
}
```

This file is used to define your CRD's structure. 

The crd file itself will live in `config/crd/bases/cache.codingbee.net_mysqls.yaml`. But you don't edit crd files directly, instead you edit the corresponding *_types.go file (e.g. api/v1alpha1/mysql_types.go) and then run the `make generate` command to it generated for you. 

The `mysql-operator2/controllers/mysql_controller.go` contains the code that is run to create to create the necessary child resources of a given CR. 

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


$ ls -l mysql-operator2/config/crd/bases/
total 8
-rw-r--r--  1 sherchowdhury  staff  1920  9 Apr 11:02 cache.codingbee.net_mysqls.yaml

$ cat mysql-operator2/config/crd/bases/cache.codingbee.net_mysqls.yaml

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
$ make bundle IMG=$OPERATOR_IMG

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

(git commit no4)

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
????????? bases
???   ????????? mysql-operator2.clusterserviceversion.yaml
????????? kustomization.yaml

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
```

This updates the `bundle/manifests/cache.codingbee.net_mysqls.yaml` file to reflect any crd changes. 


Now push up the bundle image:

```


$ make docker-push IMG=$BUNDLE_IMG

docker push quay.io/sher_chowdhury0/mysql-operator2-bundle:v0.0.1
The push refers to repository [quay.io/sher_chowdhury0/mysql-operator2-bundle]
d02927543e04: Pushed 
9f25142cbc88: Pushed 
7fdabab67fdf: Pushed 
v0.0.1: digest: sha256:aebdf842a4c5da68535a3fb9271e711820fe8f90422a265936ee66da26637ade size: 939 
```

This bundle image is specific to the version of the operator in question. When you install the operator into your namespace, the operator's catalogsource pod works out which bundle image you need, pulls down that bundle image and run's a pod-job with that bundle image. This job is what installs the operator into your namespace. by installing I mean, creating the serviceaccount, role, rolebinding, and deployment resources in your namespace. 


Now in order to do the next part, go to the quay.io web console, and for each operator and bundle repo, go to repo's settings and make them public. 

Now do `oc login ...` to log into your test cluster, and create a new project `oc new-project xxx`

Note, I also had to do the following to get the next section working (otherwise an error shows when describing replicasets):


```
oc adm policy add-scc-to-user privileged -z default -n YOUR-NAMESPACE
```


Now you can test this bundle image by running:

```
$ echo $BUNDLE_IMG                                                      
quay.io/sher_chowdhury0/mysql-operator2-bundle:v0.0.1
$ operator-sdk run bundle $BUNDLE_IMG
INFO[0017] Successfully created registry pod: quay-io-sher-chowdhury0-mysql-operator2-bundle-v0-0-1 
INFO[0017] Created CatalogSource: mysql-operator2-catalog 
INFO[0018] OperatorGroup "operator-sdk-og" created      
INFO[0018] Created Subscription: mysql-operator2-v0-0-1-sub 
INFO[0022] Approved InstallPlan install-wgfx4 for the Subscription: mysql-operator2-v0-0-1-sub 
INFO[0023] Waiting for ClusterServiceVersion "ace-sher2/mysql-operator2.v0.0.1" to reach 'Succeeded' phase 
INFO[0024]   Waiting for ClusterServiceVersion "ace-sher2/mysql-operator2.v0.0.1" to appear 
INFO[0037]   Found ClusterServiceVersion "ace-sher2/mysql-operator2.v0.0.1" phase: Pending 
INFO[0039]   Found ClusterServiceVersion "ace-sher2/mysql-operator2.v0.0.1" phase: InstallReady 
INFO[0040]   Found ClusterServiceVersion "ace-sher2/mysql-operator2.v0.0.1" phase: Installing 
INFO[0048]   Found ClusterServiceVersion "ace-sher2/mysql-operator2.v0.0.1" phase: Succeeded 
INFO[0048] OLM has successfully installed "mysql-operator2.v0.0.1" 
```

This ended up creating:

```
$ oc get catalogsource,pods,deployments,role,rolebinding,serviceaccount | grep mysql
catalogsource.operators.coreos.com/mysql-operator2-catalog   mysql-operator2   grpc   operator-sdk   4m29s
pod/mysql-operator2-controller-manager-675764c59d-hhs8f               2/2     Running     0          3m54s
pod/quay-io-sher-chowdhury0-mysql-operator2-bundle-v0-0-1             1/1     Running     0          4m26s
deployment.apps/mysql-operator2-controller-manager   1/1     1            1           3m54s
role.rbac.authorization.k8s.io/mysql-operator2.v0.0.1                                            2021-04-09T13:31:29Z
role.rbac.authorization.k8s.io/mysql-operator2.v0.0.1-default-6864f55b5d                         2021-04-09T13:31:30Z
rolebinding.rbac.authorization.k8s.io/mysql-operator2.v0.0.1                                            Role/mysql-operator2.v0.0.1                                            3m57s
rolebinding.rbac.authorization.k8s.io/mysql-operator2.v0.0.1-default-6864f55b5d                         Role/mysql-operator2.v0.0.1-default-6864f55b5d                         3m55s
```



It also created:

```
$ oc get jobs
NAME                                                              COMPLETIONS   DURATION   AGE
ab6564ddccbd77495d1ac29273ebd7208d3725ff3d1e04ed662e41e0a2936fa   1/1           16s        4m55s
```

This job used the bundle image for one of it's init containers. 


The 2 pods that are created, one relates to the catalogsource, and the other is the controller pod:

```
$ oc get pods
NAME                                                              READY   STATUS      RESTARTS   AGE
ab6564ddccbd77495d1ac29273ebd7208d3725ff3d1e04ed662e41e0a2h9mmt   0/1     Completed   0          8m22s (job created this pod)
mysql-operator2-controller-manager-675764c59d-hhs8f               2/2     Running     0          8m4s   <== controller pod
quay-io-sher-chowdhury0-mysql-operator2-bundle-v0-0-1             1/1     Running     0          8m36s  <== catalogsource
```

Since this is just for testing, all these pods are using the `default` service account:

```
$ oc get pods mysql-operator2-controller-manager-675764c59d-hhs8f -o json | jq '.spec.serviceAccount'
"default"
$ oc get pods quay-io-sher-chowdhury0-mysql-operator2-bundle-v0-0-1 -o json | jq '.spec.serviceAccount'
"default"
$ oc get pods ab6564ddccbd77495d1ac29273ebd7208d3725ff3d1e04ed662e41e0a2h9mmt -o json | jq '.spec.serviceAccount'
"default"
```

Hence why a mysql specific serviceaccount wasn't created. the default serviceaacount is assigned with thmy mysql-related-cr privileges via the rolebindings:

```
$ oc get rolebindings mysql-operator2.v0.0.1-default-6864f55b5d -o wide
NAME                                        ROLE                                             AGE   USERS   GROUPS   SERVICEACCOUNTS
mysql-operator2.v0.0.1-default-6864f55b5d   Role/mysql-operator2.v0.0.1-default-6864f55b5d   20m                    ace-sher2/default

$ oc get rolebindings mysql-operator2.v0.0.1 -o wide
NAME                     ROLE                          AGE   USERS   GROUPS   SERVICEACCOUNTS
mysql-operator2.v0.0.1   Role/mysql-operator2.v0.0.1   22m                    /default, /default
```

Not sure why we have 2 rolebindings here though. 

It also created the crd:

```
$ oc get crds mysqls.cache.codingbee.net
NAME                         CREATED AT
mysqls.cache.codingbee.net   2021-04-09T13:16:43Z
```

Now I can create a mysql cr (https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/):

```
$ oc apply -f mysql-operator2/config/samples/cache_v1alpha1_mysql.yaml 
mysql.cache.codingbee.net/mysql-sample created

$ oc get mysql
NAME           AGE
mysql-sample   12m
```

To delete your test run, do (https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/):

```
$ operator-sdk cleanup mysql-operator2   

INFO[0002] subscription "mysql-operator2-v0-0-1-sub" deleted 
INFO[0002] customresourcedefinition "mysqls.cache.codingbee.net" deleted 
INFO[0003] clusterserviceversion "mysql-operator2.v0.0.1" deleted 
INFO[0003] catalogsource "mysql-operator2-catalog" deleted 
INFO[0004] operatorgroup "operator-sdk-og" deleted      
INFO[0004] Operator "mysql-operator2" uninstalled       
```


Here's how to do a direct deployment (https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/#direct-deployment)

Note, I hit a bug, and had to do the following first:

```
$ oc adm policy add-scc-to-user privileged -z default -n mysql-operator2-system    # <operator-name>-system
```

Then ran:

```
$ echo $OPERATOR_IMG
quay.io/sher_chowdhury0/mysql-operator2:v0.0.1
$ make deploy IMG=$OPERATOR_IMG
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen "crd:trivialVersions=true,preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && /Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/kustomize edit set image controller=quay.io/sher_chowdhury0/mysql-operator2:v0.0.1
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/kustomize build config/default | kubectl apply -f -
namespace/mysql-operator2-system created
customresourcedefinition.apiextensions.k8s.io/mysqls.cache.codingbee.net created
role.rbac.authorization.k8s.io/mysql-operator2-leader-election-role created
clusterrole.rbac.authorization.k8s.io/mysql-operator2-manager-role created
Warning: resource clusterroles/mysql-operator2-metrics-reader is missing the kubectl.kubernetes.io/last-applied-configuration annotation which is required by kubectl apply. kubectl apply should only be used on resources created declaratively by either kubectl create --save-config or kubectl apply. The missing annotation will be patched automatically.
clusterrole.rbac.authorization.k8s.io/mysql-operator2-metrics-reader configured
clusterrole.rbac.authorization.k8s.io/mysql-operator2-proxy-role created
rolebinding.rbac.authorization.k8s.io/mysql-operator2-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/mysql-operator2-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/mysql-operator2-proxy-rolebinding created
configmap/mysql-operator2-manager-config created
service/mysql-operator2-controller-manager-metrics-service created
deployment.apps/mysql-operator2-controller-manager created
```

This created the following namespace+resources:

```
$ oc get catalogsource,pods,deployments,role,rolebinding,serviceaccount | grep mysql
pod/mysql-operator2-controller-manager-55668f9bb7-2lg2t   2/2     Running   0          3m31s
deployment.apps/mysql-operator2-controller-manager   1/1     1            1           8m59s
role.rbac.authorization.k8s.io/mysql-operator2-leader-election-role   2021-04-09T14:01:17Z
rolebinding.rbac.authorization.k8s.io/mysql-operator2-leader-election-rolebinding   Role/mysql-operator2-leader-election-role     9m1s
```

It looks like this deploys the operator more directly, without any olm stuff, i.e. catalogsource and bundleimages aren't used here.  

To delete everything again:

```
$ make undeploy                
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/kustomize build config/default | kubectl delete -f -
namespace "mysql-operator2-system" deleted
customresourcedefinition.apiextensions.k8s.io "mysqls.cache.codingbee.net" deleted
role.rbac.authorization.k8s.io "mysql-operator2-leader-election-role" deleted
clusterrole.rbac.authorization.k8s.io "mysql-operator2-manager-role" deleted
clusterrole.rbac.authorization.k8s.io "mysql-operator2-metrics-reader" deleted
clusterrole.rbac.authorization.k8s.io "mysql-operator2-proxy-role" deleted
rolebinding.rbac.authorization.k8s.io "mysql-operator2-leader-election-rolebinding" deleted
clusterrolebinding.rbac.authorization.k8s.io "mysql-operator2-manager-rolebinding" deleted
clusterrolebinding.rbac.authorization.k8s.io "mysql-operator2-proxy-rolebinding" deleted
configmap "mysql-operator2-manager-config" deleted
service "mysql-operator2-controller-manager-metrics-service" deleted
deployment.apps "mysql-operator2-controller-manager" deleted
```


Another way to test code changes more quickly is to do:

```
$ oc apply -f config/crd/bases/cache.codingbee.net_mysqls.yaml  
customresourcedefinition.apiextensions.k8s.io/mysqls.cache.codingbee.net created


$ make run                                                    
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen "crd:trivialVersions=true,preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go run ./main.go
I0411 15:33:21.425894   29646 request.go:645] Throttling request took 1.000815463s, request: GET:https://api.locate.cp.fyre.ibm.com:6443/apis/admissionregistration.k8s.io/v1beta1?timeout=32s
2021-04-11T15:33:23.212+0100    INFO    controller-runtime.metrics      metrics server is starting to listen    {"addr": ":8080"}
2021-04-11T15:33:23.213+0100    INFO    setup   starting manager
2021-04-11T15:33:23.213+0100    INFO    controller-runtime.manager      starting metrics server {"path": "/metrics"}
2021-04-11T15:33:23.213+0100    INFO    controller-runtime.manager.controller.mysql     Starting EventSource    {"reconciler group": "cache.codingbee.net", "reconciler kind": "Mysql", "source": "kind source: /, Kind="}
2021-04-11T15:33:23.518+0100    INFO    controller-runtime.manager.controller.mysql     Starting Controller     {"reconciler group": "cache.codingbee.net", "reconciler kind": "Mysql"}
2021-04-11T15:33:23.518+0100    INFO    controller-runtime.manager.controller.mysql     Starting workers        {"reconciler group": "cache.codingbee.net", "reconciler kind": "Mysql", "worker count": 1}
.
.
.
.
.
.
...etc
```


## Indepth tutorial 
(https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)


At the moment, our CR doesn't do anything i.e. it doesn't create any child resources, such as pods, secrets,...etc. 

At the moment our operator only manages one type of crd, called:

```
$ oc get crds mysqls.cache.codingbee.net        
NAME                         CREATED AT
mysqls.cache.codingbee.net   2021-04-10T13:29:15Z
```

Operator-sdk provided this initial boiler-plate structure for this crd:



```
$ oc explain mysql --recursive 
KIND:     Mysql
VERSION:  cache.codingbee.net/v1alpha1

DESCRIPTION:
     Mysql is the Schema for the mysqls API

FIELDS:
   apiVersion   <string>
   kind <string>
   metadata     <Object>
      annotations       <map[string]string>
      clusterName       <string>
      creationTimestamp <string>
      deletionGracePeriodSeconds        <integer>
      deletionTimestamp <string>
      finalizers        <[]string>
      generateName      <string>
      generation        <integer>
      labels    <map[string]string>
      managedFields     <[]Object>
         apiVersion     <string>
         fieldsType     <string>
         fieldsV1       <map[string]>
         manager        <string>
         operation      <string>
         time   <string>
      name      <string>
      namespace <string>
      ownerReferences   <[]Object>
         apiVersion     <string>
         blockOwnerDeletion     <boolean>
         controller     <boolean>
         kind   <string>
         name   <string>
         uid    <string>
      resourceVersion   <string>
      selfLink  <string>
      uid       <string>
   spec <Object>
      foo       <string>
   status       <map[string]>
```

Under `mysqls.spec` we are provided with a single dummy setting `foo`. 

Now I want to add a new `mysqls.spec.environment` setting. This setting in turn will store these 4 key-values pairs:

- mysql_root_password
- mysql_database
- mysql_user
- mysql_password

(These relates to: https://hub.docker.com/_/mysql#Environment_Variables)

So that our cr can look something like:

```
apiVersion: cache.codingbee.net/v1alpha1
kind: Mysql
metadata:
  name: mysql-sample
spec:
  # Add fields here
  foo: bar
  environment:
    mysql_root_password: xxxxxx
    mysql_database: wordpressDB
    mysql_user: admin
    mysql_password: xxxx
```

To implement we updated the mysql_types.go file (git commit no5). 

Note: You can't write any comments in *_types.go because comments have special meaning in this file. 

Background info about Go struct tags - https://medium.com/@sher-chowdhury/struct-tags-in-go-ca05c09d4249


Mow apply our changes:

```
$ make generate
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
```

This makes some internal changes:

```
git status
On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   api/v1alpha1/zz_generated.deepcopy.go

no changes added to commit (use "git add" and/or "git commit -a")
```

(git commit no6)

Now we update our crd:

```
$ make manifests
/Users/sherchowdhury/github/mysql-operator2/mysql-operator2/bin/controller-gen "crd:trivialVersions=true,preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```

This updates the crd:

```
$ git status
On branch main
Your branch is up to date with 'origin/main'.

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
        modified:   config/crd/bases/cache.codingbee.net_mysqls.yaml

no changes added to commit (use "git add" and/or "git commit -a")
```

(git commit no7)


Now we can start using updating the `controllers/mysql_controller.go` to create resources ( e.g. pods, secrets, configmaps) from our CR's data. 

First we add code to capture the CR's yaml data into a variable, see:

(git commit no8)

In our example the cr's data is getting captured into a variable called "instance". 

Now I'm going use that data to create a pod definition.

To keep the code, organised, I'll create a nested package called 'pod', which will live in new folder called,
`resources/pods`. 


(git commit n10)

note, I had an error when trying to import the pod package in in the controller file. I fixed that by running:

```
go mod tidy
```

This in turn ended up updating the go.mod file. 


At this point we now have a pod definition, so the next step is write code that will do the commandline eqiuvalent
of `oc apply ...`. In Go we don't have an `oc apply ...` equivalent. Instead we only have the `oc create ...`

That's why we have to write some if-(else-if) logic to mimic the `oc apply` capability. I.e. first perform a check that a pod matching the pod-definition doesn't already exist. If it does then we should take no further action, if it doesn't exist, then do the Go equivalent of `oc create ...`.

Note, I hit more problem importing some packages. Fixed this by running the `Restart Language Server` vs-code command. 



