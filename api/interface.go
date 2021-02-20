package api

import "context"

// Functionality that this library exposes.
type Api interface {
	// Obtain the list of available generator names by looking for generator-*.yaml files in sourceBaseDir
	FindGeneratorNames(ctx context.Context, sourceBaseDir string) ([]string, error)

	// Obtain a specific generator spec, read from "generator-<generatorName>.yaml" in sourceBaseDir
	ObtainGeneratorSpec(ctx context.Context, sourceBaseDir string, generatorName string) (*GeneratorSpec, error)

	// Write a fresh RenderSpec with defaults set from the GeneratorSpec for the given generator
	//
	// The name of the output file can be set in request.RenderSpecFile, but if left empty, it defaults to
	// "generated-<generatorName>.yaml".
	//
	// Warning: if the file exists, it is silently overwritten! The idea is that you keep both your
	// generators and the generator targets in source control, so you can then review the changes made.
	WriteRenderSpecWithDefaults(ctx context.Context, request *Request, generatorName string) *Response

	// Write a RenderSpec file with the provided parameter values
	//
	// The name of the output file can be set in request.RenderSpecFile, but if left empty, it defaults to
	// "generated-<generatorName>.yaml".
	//
	// Warning: if the file exists, it is silently overwritten! The idea is that you keep both your
	// generators and the generator targets in source control, so you can then review the changes made.
	WriteRenderSpecWithValues(ctx context.Context, request *Request, generatorName string, parameters map[string]string) *Response

	// Render files from templates according to RenderSpec and the GeneratorSpec it references.
	//
	// First the RenderSpec is read from <request.TargetBaseDir>/<request.RenderSpecFile>".
	// This tells the generator everything it needs to read the GeneratorSpec and execute it.
	//
	// If you leave request.RenderSpecFile empty, it defaults to "generated-main.yaml"
	//
	// Warning: existing files are silently overwritten! The idea is that you keep both your
	// generators and the generator targets in source control, so you can then review the changes made.
	Render(ctx context.Context, request *Request) *Response
}
