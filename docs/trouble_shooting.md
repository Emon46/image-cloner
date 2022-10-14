## For make deploy issue:
`$(pwd)/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
bash: /home/emon/go/src/github.com/Emon46/image-cloner/bin/controller-gen: No such file or directory
make: *** [Makefile:48: generate] Error 127`

Solution:
```
GOBIN=$(pwd)/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.8.0
```