package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"

	kirix "github.com/samkreter/Kirix/providers/aci"
	client "github.com/virtual-kubelet/virtual-kubelet/providers/azure/client/aci"
)

var (
	RandStringLength = 5
)

func Create(deploymentFile string, resourceGroup string, region string) error {

	deployment, err := GetDeploymentFromFile(deploymentFile)
	if err != nil {
		return err
	}

	pod := &v1.Pod{
		Spec: deployment.Spec.Template.Spec,
	}

	pod.Name = deployment.GetName()

	containerGroup, err := kirix.GetACIFromK8sPod(pod, region, "Linux")
	if err != nil {
		return err
	}

	aciClient, err := client.NewClient()
	if err != nil {
		return err
	}

	replicas := *deployment.Spec.Replicas
	for i := int32(0); i < replicas; i++ {
		containerGroup.Name = pod.Name + "-" + randSeq(RandStringLength)

		fmt.Printf("Creating Container Group %s.\n", containerGroup.Name)

		_, err = aciClient.CreateContainerGroup(
			resourceGroup,
			containerGroup.Name,
			*containerGroup,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func randSeq(n int) string {
	randChars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = randChars[rand.Intn(len(randChars))]
	}
	return string(b)
}

func Convert(deploymentFile string, resourceGroup string, region string) error {

	deployment, err := GetDeploymentFromFile(deploymentFile)
	if err != nil {
		return err
	}

	pod := &v1.Pod{
		Spec: deployment.Spec.Template.Spec,
	}

	pod.Name = deployment.GetName()

	aci, err := kirix.GetACIFromK8sPod(pod, region, "Linux")
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(aci)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	return nil
}

func GetDeploymentFromFile(deploymentFile string) (*v1beta1.Deployment, error) {
	deploymentData, err := ioutil.ReadFile(deploymentFile)
	if err != nil {
		return &v1beta1.Deployment{}, fmt.Errorf("Could not find deployment file: %s", deploymentFile)
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode(deploymentData, nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	deployment := obj.(*v1beta1.Deployment)

	return deployment, nil
}
