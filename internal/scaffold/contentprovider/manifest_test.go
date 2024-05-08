package contentprovider

import (
	"fmt"
	"strings"
	"testing"
)

func Test_ManifestContentProvider_SetsDefaultCorrectly(t *testing.T) {
	manifestContentProvider := NewManifestContentProvider()

	manifestGeneratedDefaultContent, _ := manifestContentProvider.GetDefaultContent(nil)
	t.Parallel()
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:  "Manifest",
			value: manifestGeneratedDefaultContent,
			expected: `# This file holds the Manifest of your module, encompassing all resources installed in the cluster once the module is activated.
# It should include the Custom Resource Definition for your module's default CustomResource, if it exists.

`,
		},
	}

	for _, testcase := range tests {
		testcase := testcase
		testName := fmt.Sprintf("TestCorrectContentProviderFor_%s", testcase.name)

		testcase.value = strings.TrimSpace(testcase.value)
		testcase.expected = strings.TrimSpace(testcase.expected)

		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			if testcase.value != testcase.expected {
				t.Errorf("ContentProvider for '%s' did not return correct default: expected = '%s', but got = '%s'",
					testcase.name, testcase.expected, testcase.value)
			}
		})
	}
}
