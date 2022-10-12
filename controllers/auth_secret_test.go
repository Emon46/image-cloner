package controllers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestValidateAndGetAuthSecret(t *testing.T) {
	testCases := []struct {
		name         string
		secret       *v1.Secret
		auth         string
		username     string
		password     string
		registryHost string
		errString    string
	}{
		{
			name:         "Ok",
			username:     "hremon331046",
			password:     "secret1234",
			secret:       getDummySecret("hremon331046:secret1234", ""),
			registryHost: "",
			errString:    "",
		},
		{
			name:         "Ok with RegistryHost",
			username:     "hremon331046",
			password:     "secret1234",
			secret:       getDummySecret("hremon331046:secret1234", "gcr.io"),
			registryHost: "gcr.io",
			errString:    "",
		},
		{
			name:         "Missing Auth key",
			username:     "hremon331046",
			password:     "secret1234",
			secret:       getDummySecret("", ""),
			registryHost: "",
			errString:    fmt.Sprintf("\"auth\" key field not exist inside secret"),
		},
		{
			name:         "Invalid Auth Secret",
			username:     "hremon331046",
			password:     "secret1234",
			secret:       getDummySecret("hremon331046", ""),
			registryHost: "",
			errString:    fmt.Sprintf("assigned invalid value to auth key inside %s/%s secret", "dummyNs", "dummy"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			username, password, resRegistryHost, err := validateAndGetRegistryAuthCreds(testCase.secret)
			if err != nil || testCase.errString != "" {
				require.EqualError(t, err, testCase.errString)
			} else {
				require.Equal(t, testCase.username, username)
				require.Equal(t, testCase.password, password)
				require.Equal(t, testCase.registryHost, resRegistryHost)
			}
		})
	}
}

func getDummySecret(auth string, registry string) *v1.Secret {
	data := make(map[string][]byte)
	if auth != "" {
		data["auth"] = []byte(auth)
	}
	if registry != "" {
		data[RegistryHostName] = []byte(registry)
	}
	return &v1.Secret{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "dummy",
			Namespace: "dummyNs",
		},
		Data: data,
	}
}

//
//func TestGetAuthConfig(t *testing.T) {
//
//	clientset := testclient.NewSimpleClientset()
//	testclient.NewSimpleClientset()
//	getRegistryAuthCred(context.TODO(), clientset)
//	secret, err := clientset.CoreV1().Secrets("dummyNs").Create(context.TODO(), getDummySecret("hello:secret1234", ""), meta_v1.CreateOptions{})
//	fmt.Println("1 secret: ", secret, "1 err: ", err)
//
//	secret, err = clientset.CoreV1().Secrets("dummyNs").Get(context.TODO(), "dummy", meta_v1.GetOptions{})
//	fmt.Println("secret: ", secret, "err: ", err)
//
//}
//func createDummySecret(auth string, registryHost string, clientset kubernetes.Clientset) (*v1.Secret, error) {
//	secret, err := clientset.CoreV1().Secrets("dummyNs").Create(context.TODO(), getDummySecret(auth, registryHost), meta_v1.CreateOptions{})
//	return secret, err
//}
