package framework

//
//func (f *Framework) CreateAuthSecret(obj *core.Secret) error {
//	_, err := f.kubeClient.CoreV1().Secrets(obj.Namespace).Create(context.TODO(), obj, meta_v1.CreateOptions{})
//	return err
//}
//
//func (f *Framework) DeleteAuthSecret(obj meta_v1.ObjectMeta) error {
//	return f.kubeClient.CoreV1().Secrets(obj.Namespace).Delete(context.TODO(), obj.Name, meta_v1.DeleteOptions{})
//}
//
//func (i *Invocation) GetDockerHubAuthSecret(secretName string) *core.Secret {
//	return &core.Secret{
//		ObjectMeta: meta_v1.ObjectMeta{
//			Name:      secretName,
//			Namespace: i.Namespace(),
//		},
//		Data: map[string][]byte{
//			"auth":         []byte("<your dockerhub username:password>"),
//			"RegistryHost": []byte(""),
//		},
//	}
//}
