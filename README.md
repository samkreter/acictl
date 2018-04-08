# acictl
A simple way to interact with Azure Container Instance in a Kubernetes style. 

## Installation

With a working Go environment run: `go get github.com/samkreter/acictl`

## Usage

#### Create
acictl create allows for creating different ACI from a Kubernetes deployment spec.

If we have a deployment spec named test.yaml with the following,

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

then running `acictl create -g ResourceGroup -f test.yaml` will create 3 ACIs with names nginx-deploymet-<randomstring>

#### Delete 

To delete a deployment, simply run `acictl delete -g ResourceGroup -f test.yaml` and all instances will be deleted.

#### Convert

Convert allows you to generating an Azure ARM template from a Kubernetes deployment spec. 

Running `acictl convert -f test.yaml > template.json` will generate the ARM template. This can be used with the azure cli to create a container instaces using 

`az group deployment create -g <resourec-group> -n <container-group-name> --template-file template.json` 

Here is an example output to test.yaml
```json
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": [
    {
      "name": "nginx-deployment",
      "type": "Microsoft.ContainerInstance/containerGroups",
      "location": "westus",
      "properties": {
        "containers": [
          {
            "name": "nginx",
            "properties": {
              "image": "nginx:1.7.9",
              "ports": [
                {
                  "protocol": "TCP",
                  "port": 80
                }
              ],
              "instanceView": {
                "currentState": {
                  "startTime": "0001-01-01T00:00:00Z",
                  "finishTime": "0001-01-01T00:00:00Z"
                },
                "previousState": {
                  "startTime": "0001-01-01T00:00:00Z",
                  "finishTime": "0001-01-01T00:00:00Z"
                }
              },
              "resources": {
                "requests": {
                  "memoryInGB": 1,
                  "cpu": 1
                },
                "limits": {
                  "memoryInGB": 1,
                  "cpu": 1
                }
              }
            }
          }
        ],
        "ipAddress": {
          "ports": [
            {
              "protocol": "TCP",
              "port": 80
            }
          ],
          "type": "Public"
        },
        "osType": "Linux",
        "instanceView": {}
      },
      "apiVersion": "2018-04-01"
    }
  ]
}
```

