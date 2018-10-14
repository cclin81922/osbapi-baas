# Installation

```
go get -u github.com/cclin81922/osbapi-baas/cmd/osbapibaas
export PATH=$PATH:~/go/bin
```

# Command Line Usage

```
osbapibaas -port=8443
```

# Deploy baas using Helm

```console
$ TAG=latest PULL=Never make deploy-baas
```
