# Installation

```
go get -u github.com/cclin81922/osbapi-baas/cmd/osbapibaas
export PATH=$PATH:~/go/bin
```

# Command Line Usage

```
osbapibaas -port=8443
```

# Deploy osbapibaas Using Helm

```
make deploy-baas
```

# For Developer

Run all tests

```
echo "127.0.0.1   localhost.localdomain" >> /etc/hosts
go test github.com/cclin81922/osbapi-baas/cmd/osbapibaas
```
