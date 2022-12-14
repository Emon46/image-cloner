package controllers

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func getRegistryNameFromImage(imageName string) string {
	imageParts := strings.Split(imageName, "/")
	if len(imageParts) >= 2 {
		return imageParts[len(imageParts)-2]
	}
	return ""
}

func isBackupRegistry(imageName, registry string) bool {
	if getRegistryNameFromImage(imageName) == registry {
		return true
	}
	return false
}

func pushContainersToBackupRegistry(ctx context.Context, kubeClient client.Client, containers []v1.Container) ([]v1.Container, error) {
	authConfig, registryHost, err := newAuthConfigForRegistry(ctx, kubeClient)
	if err != nil {
		return nil, err
	}
	for index, container := range containers {
		imgName := container.Image

		if isBackupRegistry(imgName, authConfig.Username) {
			// this image is already the backup registry image
			// ignore this image
			continue
		}

		backupRegistryImageName := fmt.Sprintf("%s/%s", authConfig.Username, strings.ReplaceAll(imgName, "/", "_"))
		if registryHost != "" {
			backupRegistryImageName = fmt.Sprintf("%s/%s", registryHost, backupRegistryImageName)
		}

		// copy the modified image to own repository using crane.Pull and crane.Push
		image, err := crane.Pull(container.Image)
		if err != nil {
			return nil, err
		}
		err = crane.Push(image, backupRegistryImageName, crane.WithAuth(authn.FromConfig(authConfig)))
		if err != nil {
			return nil, err
		}
		klog.Infoln(fmt.Sprintf("successfully pushed the docker image: %s to backup registry as %s", container.Image, backupRegistryImageName))

		// ISSUE: crane.Copy() will not work if the original image is from one hub like: (gcr, docker) and backup image is pushing in other hub like: (gcr, docker)
		// we are passing only one auth but they are trying to use the same auth in both pull and push

		//err = crane.Copy(container.Image, backupRegistryImageName, crane.WithAuth(authn.FromConfig(authConfig)),
		//	crane.WithContext(ctx))
		//if err != nil {
		//	return nil, err
		//}

		// update the docker image with the backup registry image
		containers[index].Image = backupRegistryImageName
	}
	return containers, nil
}

// get the auth config from registry cred
func newAuthConfigForRegistry(ctx context.Context, kubeClient client.Client) (authn.AuthConfig, string, error) {
	username, password, registryHost, err := getRegistryAuthCred(ctx, kubeClient)
	if err != nil {
		return authn.AuthConfig{}, registryHost, err
	}
	return authn.AuthConfig{
		Username: username,
		Password: password,
		Auth:     fmt.Sprintf("%s:%s", username, password),
	}, registryHost, nil
}
