package controllers

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "k8s.io/api/core/v1"
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
	authConfig, registryHost, err := getAuthConfigForRegistry(ctx, kubeClient)
	if err != nil {
		return nil, err
	}
	for index, container := range containers {
		imgName := container.Image

		if isBackupRegistry(imgName, authConfig.Username) {
			continue
		}

		backupRegistryImageName := fmt.Sprintf("%s/%s", authConfig.Username, strings.ReplaceAll(imgName, "/", "_"))
		if registryHost != "" {
			backupRegistryImageName = fmt.Sprintf("%s/%s", registryHost, backupRegistryImageName)
		}

		err := crane.Copy(container.Image, backupRegistryImageName, crane.WithAuth(authn.FromConfig(authConfig)))
		if err != nil {
			return nil, err
		}

		containers[index].Image = backupRegistryImageName
	}
	return containers, nil
}

func getAuthConfigForRegistry(ctx context.Context, kubeClient client.Client) (authn.AuthConfig, string, error) {
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
