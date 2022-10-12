package controllers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsBackupRegistry(t *testing.T) {
	testCases := []struct {
		name         string
		RegistryName string
		ImageName    string
		result       bool
	}{
		{
			name:         "Ok",
			RegistryName: "hremon331046",
			ImageName:    "hremon331046/demo:latest",
			result:       true,
		},
		{
			name:         "Ok with Docker Url",
			RegistryName: "hremon331046",
			ImageName:    "gcr.io/hremon331046/demo:latest",
			result:       true,
		},
		{
			name:         "Empty Registry",
			RegistryName: "hremon331046",
			ImageName:    "postgres:latest",
			result:       false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ok := isBackupRegistry(testCase.ImageName, testCase.RegistryName)
			require.Equal(t, testCase.result, ok)
		})
	}
}

func TestGetRegistry(t *testing.T) {
	testCases := []struct {
		name         string
		RegistryName string
		ImageName    string
	}{
		{
			name:         "Ok",
			RegistryName: "hremon331046",
			ImageName:    "hremon331046/demo:latest",
		},
		{
			name:         "Ok with Docker Url",
			RegistryName: "hremon-331046",
			ImageName:    "gcr.io/hremon-331046/demo:latest",
		},
		{
			name:         "Empty Registry",
			RegistryName: "",
			ImageName:    "postgres:latest",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resRegistry := getRegistryNameFromImage(testCase.ImageName)
			require.Equal(t, testCase.RegistryName, resRegistry)
		})
	}
}
