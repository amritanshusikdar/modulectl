package contentprovider

import "github.com/kyma-project/modulectl/internal/scaffold/common/types"

type DefaultCRContentProvider struct{}

func NewDefaultCRContentProvider() *DefaultCRContentProvider {
	return &DefaultCRContentProvider{}
}

func (s *DefaultCRContentProvider) GetDefaultContent(_ types.KeyValueArgs) string {
	return `# This is the file that contains the defaultCR for your module, which is the Custom Resource that will be created upon module enablement.
	# Make sure this file contains *ONLY* the Custom Resource (not the Custom Resource Definition, which should be a part of your module manifest)

`
}
