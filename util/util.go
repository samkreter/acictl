package util

import (
	"fmt"

	"k8s.io/client-go/pkg/api"
	_ "k8s.io/client-go/pkg/api/install"
	_ "k8s.io/client-go/pkg/apis/extensions/install"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func Convert(deploymentFile string) error {
	// decode := api.Codecs.UniversalDecoder().Decode
	decode := api.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(deployment), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	t := obj.(*v1beta1.Deployment)
	fmt.Println(t.Spec.Replicas)
	return fmt.Errorf("Not implemented.")
}
