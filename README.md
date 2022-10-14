# Image-cloner
This is a k8s controller which will watch for Deployment and Daemonset pods and store the image by re-uploading to a bakcup registry from public registry images.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### demo link
- asciinema demo link [link](https://asciinema.org/a/528475)

### Run Controller and Test it out
- Apply docker credential `secret` in `config/sample/registry-cred.yaml` file. Before applying the `secret`, change the `stringData` key `auth` to right `<docker-username>` and `<password>` in `image-cloner-cred` secret.
  - If you want to update the docker cred secret name and namespace then need to add the secret name and namespace in `ENV` of the deployment inside `./config/manager/manager.yaml` and update the  `RegistrySecretName` and `RegistrySecretNamespace`
```sh
kubectl apply -f config/sample/registry-cred.yaml
```

- Build and push your image to the location specified by `IMG`:
```sh
export REGISTRY=<registry-name>
make docker-push
```

- Deploy the controller to the cluster with the image specified by `IMG`:
```sh
export REGISTRY=<registry-name>
make deploy
```
or,
for deploying in kind
```sh
export REGISTRY=<registry-name>
make deploy-to-kind
```
- check the pods are up and running: `kubectl get pod -n image-cloner-system`  

- apply sample deployment and daemonset
```shell
kubectl apply -f config/sample/demo-deployment.yaml
kubectl apply -f config/sample/demo-daemonset.yaml
```
- now check in the sample deployment and daemonset images, they will be cloned & pushed to your backup docker registry and use new docker image in the deployment

- UnDeploy the controller to the cluster:

```sh
make undeploy
```

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

## Run e2e test
- Added e2e test for deployment and Daemonset controller
- For e2e test follow below steps:
  - run the controller and create the `docker cred secret`
  - run `ginkgo tests/e2e/`

## Run unit test
`make test`

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
