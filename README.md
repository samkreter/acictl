# acictl
A simple way to interact with Azure Container Instance in a Kubernetes style. 

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