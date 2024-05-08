package contentprovider

import (
	"fmt"
	"strings"
	"testing"
)

func Test_DefaultCRContentProvider_SetsDefaultCorrectly(t *testing.T) {
	defaultCRContentProvider := NewDefaultCRContentProvider()

	defaultCRGeneratedDefaultContent, _ := defaultCRContentProvider.GetDefaultContent(nil)

	t.Parallel()
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:  "DefaultCR",
			value: defaultCRGeneratedDefaultContent,
			expected: `# This is the file that contains the defaultCR for your module, which is the Custom Resource that will be created upon module enablement.
# Make sure this file contains *ONLY* the Custom Resource (not the Custom Resource Definition, which should be a part of your module manifest)

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
