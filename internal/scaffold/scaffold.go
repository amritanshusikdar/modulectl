package scaffold

import (
	"fmt"
	"path"

	"github.com/kyma-project/modulectl/internal/scaffold/common/types"
	"github.com/kyma-project/modulectl/internal/scaffold/contentprovider"
	"github.com/kyma-project/modulectl/tools/io"
)

type ModuleConfigService interface {
	FileGeneratorService
	PreventOverwrite(directory, moduleConfigFileName string, overwrite bool) error
}

type ManifestService interface {
	GenerateManifestFile(out io.Out, path string) error
}

type DefaultCRService interface {
	GenerateDefaultCRFile(out io.Out, path string) error
}

type FileGeneratorService interface {
	GenerateFile(out io.Out, path string, args types.KeyValueArgs) error
}

// TODO: @Adi: Refactor so that 'manifestService' and 'defaultCRService' to use the FileGeneratorService interface and `ReuseFileGeneratorService` implementation
// 'securityConfigService' provides an overall exmaple
// please collect all questions you got while working on this task
// please create a separate branch and open a PR targetting the current branch
// we will have a look at the PR and discuss your questions together next week
// outline:
//   - change 'manifestService' and 'defaultCRService' to be of type `FileGeneratorService`
//     adapt the code in the CreateScaffold method accordingly (i.e., use the new method signature)
//   - create new contententproviders for manifest and defaultCR
//     see contentprovider/securityconfig.go as an example
//     in this case, we don't need the yaml converter and args (just pass nil or an empty map)
//     the default content (e.g., from manifest.go#getDefaultContent) should be returned from the new content providers
//     include a simple test that assures that the returned string is the expected default content (similar to how we test flag defaults)
//   - update the composition root in cmd/modulectl/create/scaffold/scaffold.go
//     instead of 'NewManifestService' and 'NewDefaultCRService' a new 'NewReuseFileGeneratorService' should be used for each of those using the respective content provider
//     see the security config part as an example (remember, yaml converter is not needed for the new ones)
//   - refactor the related tests in 'scaffold_test'
//     the stubs for manifestService and defaultCRService can be replaced with stubs for FileGenerator
//     again, see security config as an example
//   - remove the 'ManifestService' and 'DefaultCRService' interfaces from this file
//   - remove the 'ManifestService' and 'DefaultCRService' implementations from internal/scaffold/manifest/manifest.go and internal/scaffold/defaultcr/defaultcr.go
type ScaffoldService struct {
	moduleConfigService   ModuleConfigService
	manifestService       FileGeneratorService
	defaultCRService      FileGeneratorService
	securityConfigService FileGeneratorService
}

func NewScaffoldService(moduleConfigService ModuleConfigService,
	manifestService FileGeneratorService,
	defaultCRService FileGeneratorService,
	securityConfigService FileGeneratorService) *ScaffoldService {
	return &ScaffoldService{
		moduleConfigService:   moduleConfigService,
		manifestService:       manifestService,
		defaultCRService:      defaultCRService,
		securityConfigService: securityConfigService,
	}
}

func (s *ScaffoldService) CreateScaffold(opts Options) error {
	if err := opts.validate(); err != nil {
		return err
	}

	if err := s.moduleConfigService.PreventOverwrite(opts.Directory, opts.ModuleConfigFileName, opts.ModuleConfigFileOverwrite); err != nil {
		return err
	}

	// TODO: @Adi: We use 'path.Join(opts.Directory, ...)' to create the file paths at various places here
	// it should work regardless the user provides a relative or absolute path
	// as of now, I think it only works with absolute paths for 'opts.Directory'
	// please verify if this observation is true and if so, please fix it
	// do some research on how to handle file paths in Go properly and update the code accordingly

	// DONE: `path.Join(opts.Directory, ...)` actually works for both the cases, i.e., absolute and relative
	// paths. The files were always created in the desired directory. Tested it with the following cases:
	//	1. shorthand -d
	//  2. full flag --directory
	//	3. $(pwd)/{folder_name} -> tried replacing `folder_name` with "", "LICENSES", "docs" and it worked each time
	//	4. hard coded the absolute path; works
	// 	5. hard coded relative path -> ".", "./", "../modulectl", "../../Desktop/modulectl"; works
	//	So, nothing changed in this case :D
	manifestFilePath := path.Join(opts.Directory, opts.ManifestFileName)
	if err := s.manifestService.GenerateFile(opts.Out, manifestFilePath, nil); err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.ManifestFileName, err)
	}

	defaultCRFilePath := ""
	if opts.defaultCRFileNameConfigured() {
		defaultCRFilePath = path.Join(opts.Directory, opts.DefaultCRFileName)
		if err := s.defaultCRService.GenerateFile(opts.Out, defaultCRFilePath, nil); err != nil {
			return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.DefaultCRFileName, err)
		}
	}

	securityConfigFilePath := ""
	if opts.securityConfigFileNameConfigured() {
		securityConfigFilePath = path.Join(opts.Directory, opts.SecurityConfigFileName)
		if err := s.securityConfigService.GenerateFile(
			opts.Out,
			securityConfigFilePath,
			types.KeyValueArgs{contentprovider.ArgModuleName: opts.ModuleName}); err != nil {
			return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.SecurityConfigFileName, err)
		}
	}

	moduleConfigFilePath := path.Join(opts.Directory, opts.ModuleConfigFileName)
	if err := s.moduleConfigService.GenerateFile(
		opts.Out,
		moduleConfigFilePath,
		types.KeyValueArgs{
			contentprovider.ArgModuleName:         opts.ModuleName,
			contentprovider.ArgModuleVersion:      opts.ModuleVersion,
			contentprovider.ArgModuleChannel:      opts.ModuleChannel,
			contentprovider.ArgManifestFile:       opts.ManifestFileName,
			contentprovider.ArgDefaultCRFile:      defaultCRFilePath,
			contentprovider.ArgSecurityConfigFile: securityConfigFilePath,
		}); err != nil {
		return fmt.Errorf("%w %s: %w", ErrGeneratingFile, opts.ModuleConfigFileName, err)
	}

	return nil
}
