package contentprovider

type ManifestContentProvider struct{}

func NewManifestContentProvider() *ManifestContentProvider {
	return &ManifestContentProvider{}
}

func (s *ManifestContentProvider) GetDefaultContent() string {
	return s.getDefaultContent()
}

func (s *ManifestContentProvider) getDefaultContent() string {
	return `# This file holds the Manifest of your module, encompassing all resources installed in the cluster once the module is activated.
# It should include the Custom Resource Definition for your module's default CustomResource, if it exists.

`
}
