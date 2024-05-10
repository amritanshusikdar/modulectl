package contentprovider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
)

func Test_DefaultCRContentProvider_SetsDefaultCorrectly(t *testing.T) {
	defaultCRContentProvider := NewDefaultCRContentProvider()

	expectedDefault := `# This is the file that contains the defaultCR for your module, which is the Custom Resource that will be created upon module enablement.
# Make sure this file contains *ONLY* the Custom Resource (not the Custom Resource Definition, which should be a part of your module manifest)

`

	defaultCRGeneratedDefaultContentWithNil, _ := defaultCRContentProvider.GetDefaultContent(nil)
	defaultCRGeneratedDefaultContentWithEmptyMap, _ := defaultCRContentProvider.GetDefaultContent(make(types.KeyValueArgs))

	t.Parallel()
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "DefaultCR Default Content with Nil",
			value:    defaultCRGeneratedDefaultContentWithNil,
			expected: expectedDefault,
		},
		{
			name:     "DefaultCR Default Content with Empty Map",
			value:    defaultCRGeneratedDefaultContentWithEmptyMap,
			expected: expectedDefault,
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
