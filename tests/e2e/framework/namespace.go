package framework

import (
	"context"

	core "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_util "kmodules.xyz/client-go/meta"
)

func (f *Framework) Namespace() string {
	return f.namespace
}

func (f *Framework) CreateTestNamespace() error {
	obj := &core.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: f.namespace,
		},
	}
	_, err := f.kubeClient.CoreV1().Namespaces().Create(context.TODO(), obj, meta_v1.CreateOptions{})
	return err
}

func (f *Framework) DeleteTestNamespace() error {
	return f.kubeClient.CoreV1().Namespaces().Delete(context.TODO(), f.namespace, meta_util.DeleteInForeground())
}
