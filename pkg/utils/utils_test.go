package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadTokenFromFile(t *testing.T) {
	// Create a temporary file with a test token
	tempFile, err := os.CreateTemp("", "test-token")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testToken := "test-github-token"
	_, err = tempFile.WriteString("GH_TOKEN=" + testToken)
	require.NoError(t, err)
	tempFile.Close()

	// Test reading the token
	token, err := ReadTokenFromFile(tempFile.Name())
	require.NoError(t, err)
	assert.Equal(t, testToken, token)

	// Test with non-existent file
	_, err = ReadTokenFromFile("non-existent-file")
	assert.Error(t, err)

	// Test with file without GH_TOKEN
	emptyFile, err := os.CreateTemp("", "empty-token")
	require.NoError(t, err)
	defer os.Remove(emptyFile.Name())
	emptyFile.Close()

	_, err = ReadTokenFromFile(emptyFile.Name())
	assert.Error(t, err)
}

func TestWrapError(t *testing.T) {
	originalErr := assert.AnError
	wrappedErr := WrapError(originalErr, "test message")
	assert.Error(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "test message")
	assert.Contains(t, wrappedErr.Error(), originalErr.Error())
}
