package api

// Parameters you will need to provide for a render run. All the rest is read from parameters
type Request struct {
	// Directory where to find e.g. 'main.yaml' describing the generator. Required.
	SourceBaseDir string `yaml:"sourcedir"`

	// Directory where to find 'generator-main.yaml' specifying values and the generator to use. Required.
	TargetBaseDir string `yaml:"targetdir"`

	// yaml-file to read for RenderSpec, if not set, defaults to "generated-main.yaml".
	RenderSpecFile string `yaml:"renderspec"`
}

// Information about the results of a render run
type Response struct {
	Success       bool
	RenderedFiles []FileResult
	Errors        []error
}

type FileResult struct {
	Success          bool
	RelativeFilePath string
	Errors           []error
}
