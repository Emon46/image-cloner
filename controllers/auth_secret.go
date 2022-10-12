package controllers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func getRegistryAuthCred(ctx context.Context, kubeClient client.Client) (string, string, string, error) {
	secretName := os.Getenv(RegistrySecretName)
	secretNamespace := os.Getenv(RegistrySecretNamespace)
	if secretName == "" || secretNamespace == "" {
		return "", "", "", fmt.Errorf("set Env key \"%s\" and \"%s\" inside controller container", RegistrySecretName, RegistrySecretNamespace)
	}

	secret := &v1.Secret{}
	objMeta := types.NamespacedName{
		Name:      secretName,
		Namespace: secretNamespace,
	}
	err := kubeClient.Get(ctx, objMeta, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Infoln(fmt.Sprintf("secret %s not exist", objMeta.String()))
			return "", "", "", err
		}
		return "", "", "", err
	}

	return validateAndGetRegistryAuthCreds(secret)
}

func validateAndGetRegistryAuthCreds(secret *v1.Secret) (string, string, string, error) {
	authByte, ok := secret.Data["auth"]
	if !ok {
		return "", "", "", fmt.Errorf("\"auth\" key field not exist inside secret")
	}

	authSlice := strings.Split(string(authByte), ":")
	if len(authSlice) != 2 {
		return "", "", "", fmt.Errorf(fmt.Sprintf("assigned invalid value to auth key inside %s/%s secret", secret.Namespace, secret.Name))
	}
	registryHost := ""
	if _, ok := secret.Data[RegistryHostName]; ok {
		registryHost = string(secret.Data[RegistryHostName])
	}
	username := authSlice[0]
	password := authSlice[1]
	return username, password, registryHost, nil
}
