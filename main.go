// package main

// import (
// 	"github.com/samkreter/dockdev/cmd"
// )

// func main() {
// 	cmd.Execute()
// }

// package main

// import (
// 	"github.com/samkreter/dockdev/cmd"
// )

// func main() {
// 	cmd.Execute()
// }

package main

import (
	"fmt"

	"k8s.io/client-go/pkg/api"
	_ "k8s.io/client-go/pkg/api/install"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const yaml = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: push
`

var deployment = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: my-nginx
spec:
  replicas: 2
  template:
    metadata:
      labels:
        run: my-nginx
    spec:
      containers:
      - name: my-nginx
        image: nginx
        ports:
        - containerPort: 80
`

func main() {
	// decode := api.Codecs.UniversalDecoder().Decode
	decode := api.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(deployment), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	t := obj.(*v1beta1.Deployment)
	fmt.Println(t.Spec.Replicas)
}
