package controllers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIncludedNamespace(t *testing.T) {
	testCases := []struct {
		name      string
		nameSpace string
		result    bool
	}{
		{
			name:      "Ok",
			nameSpace: "demo",
			result:    true,
		},
		{
			name:      "Ok",
			nameSpace: "kube-system",
			result:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := IncludedNamespace(testCase.nameSpace)
			require.Equal(t, testCase.result, res)
		})
	}
}
