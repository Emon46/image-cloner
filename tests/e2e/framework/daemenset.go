package framework

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	"strings"
	"time"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_util "kmodules.xyz/client-go/meta"
)

func (i *Invocation) GetDaemonSet(name string) *appsv1.DaemonSet {
	image := "prom/node-exporter"

	return &appsv1.DaemonSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: i.Namespace(),
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "node-exporter",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app": "node-exporter",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "node-exporter",
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (f *Framework) CreateDaemonset(obj *appsv1.DaemonSet) error {
	_, err := f.kubeClient.AppsV1().DaemonSets(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func (f *Framework) DeleteDaemonSet(meta metav1.ObjectMeta) error {
	return f.kubeClient.AppsV1().DaemonSets(meta.Namespace).Delete(context.TODO(), meta.Name, meta_util.DeleteInBackground())
}

func (f *Framework) EventuallyImageClonedForDaemonSetToBackupRegistry(meta metav1.ObjectMeta, registry string) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			daemonset, err := f.kubeClient.AppsV1().DaemonSets(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			tmp := true

			for _, container := range daemonset.Spec.Template.Spec.Containers {
				img := container.Image
				if !strings.HasPrefix(img, registry) {
					tmp = false
					break
				}
			}
			return tmp
		},
		time.Minute*10,
		time.Second*10,
	)
}

func (f *Framework) EventuallyDaemonSetDeleted(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			_, err := f.kubeClient.AppsV1().DaemonSets(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
			return errors.IsNotFound(err)
		},
		time.Minute*15,
		time.Second*10,
	)
}
