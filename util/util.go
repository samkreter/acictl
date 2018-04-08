package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"

	kirix "github.com/samkreter/Kirix/providers/aci"
	client "github.com/virtual-kubelet/virtual-kubelet/providers/azure/client/aci"
)

type ArmTemplate struct {
	Schema         string        `json:"$schema"`
	ContentVersion string        `json:"contentVersion"`
	Resources      []interface{} `json:"resources"`
}

var (
	RandStringLength          = 5
	ArmTemplateSchema         = "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#"
	ArmTemplateContentVersion = "1.0.0.0"
)

func GenerateArmTemplate(cg *client.ContainerGroup) *ArmTemplate {
	cgWithAPIVersion := struct {
		client.ContainerGroup
		APIVersion string `json:"apiVersion"`
	}{
		*cg,
		"2018-04-01",
	}

	return &ArmTemplate{
		Schema:         ArmTemplateSchema,
		ContentVersion: ArmTemplateContentVersion,
		Resources:      []interface{}{cgWithAPIVersion},
	}
}

func Delete(deploymentFile string, resourceGroup string) error {
	deployment, err := GetDeploymentFromFile(deploymentFile)
	if err != nil {
		return fmt.Errorf("Parse deployment error: %s", err)
	}

	aciClient, err := kirix.CreateACIClient()
	if err != nil {
		return err
	}

	cgList, err := aciClient.ListContainerGroups(resourceGroup)
	if err != nil {
		return fmt.Errorf("Container group list error: %s", err)
	}

	for _, cg := range cgList.Value {
		if strings.Contains(cg.Name, deployment.Name) {
			fmt.Printf("Deleting container group %s\n", cg.Name)
			err := aciClient.DeleteContainerGroup(resourceGroup, cg.Name)
			if err != nil {
				return fmt.Errorf("Delete container group error: %s", err)
			}
		}
	}

	return nil
}

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

	aciClient, err := kirix.CreateACIClient()
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

	cg, err := kirix.GetACIFromK8sPod(pod, region, "Linux")
	if err != nil {
		return err
	}

	template := GenerateArmTemplate(cg)

	jsonData, err := json.MarshalIndent(template, "", "  ")
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
