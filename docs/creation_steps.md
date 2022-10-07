```
kubebuilder init --domain image-cloner.dev --skip-go-version-check

kubebuilder edit --multigroup=true

kubebuilder create api --group apps --version v1 --kind Deployment --controller=true --resource=false

kubebuilder create api --group apps --version v1 --kind DaemonSet --controller=true --resource=false
```