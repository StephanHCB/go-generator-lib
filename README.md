# go-generator-lib

A golang library for generating files from templates. Can be used to scaffold application code.

See [go-generator-cli](https://github.com/StephanHCB/go-generator-cli/) for end user usage instructions, 
this is just a library that provides the rendering functionality, so it can be used from both a 
CLI and (potentially) a service.

## Generators

A generator is a directory that contains one or more `generator-*.yaml` files, called 
*generator specification files*, plus a number of 
golang [text/templates](https://golang.org/pkg/text/template/). The main generator spec is
usually called `generator-main.yaml`.

Example:
```
templates:
  - source: 'src/sub/sub.go.tmpl'
    target: 'sub/sub.go'
  - source: 'src/main.go.tmpl'
    target: 'main.go'
  - source: 'src/web/controller.go.tmpl'
    target: 'web/controller/{{ .item }}.go'
    condition: '{{ if eq .item "skipped" }}false{{ end }}'
    with_items:
     - health
     - reservations
     - skipped
variables:
  serviceUrl:
    description: 'The URL of the service repository, to be used in imports etc.'
    default: 'github.com/StephanHCB/temp'
  serviceName:
    description: 'The name of the service to be rendered.'
    pattern: '^[a-z-]+$'
  helloMessage:
    description: 'A message to be inserted in the code.'
    default: 'hello {{ "world" }}'
```

This defines templates like `src/sub/sub.go.tmpl` and `src/main.go.tmpl` and what target path
they'll be rendered to in the target directory. 
It also specifies which parameter variables will be available during rendering.

  * If a variable does not have a default value, it is a required parameter.
  * default values are evaluated as templates, too, but you will not be able to refer to other variables  
  * if a variable has a pattern set, the parameter value must regex-match that pattern. Please be advised that
    you must enclose the pattern with ^...$ if you want to force the whole value to match, otherwise
    it's enough for part of the value to match the pattern.
  * variables are assumed to be string-valued by default, but the template generator actually allows any
    valid yaml structure (lists and maps, even nested) both as default values and as variable values.
    There is no type checking whatsoever, parsing templates that access missing fields or list items
    will fail, so it is not recommended to overuse this feature. Also, you should definitely provide
    a default value for any list or map typed variable, for else how will your users know what structure
    you are assuming?

The idea is that you keep your generators under version control.

Note how you can create ansible-style loops using the same template to generate multiple output files using `with_items`.
In fact, the output file name is always parsed using the same template engine as the actual templates,
so you could also use other variables in it. 

If you set `with_items`, the template is used multiple times
with the `item` variable set to the value you provided under `with_items`. These values can also be 
a whole yaml data structure, you simply access it as `{{ .item.some.field }}`. 

_At this time, it is not possible to dynamically assign the full list in with_items from a variable, 
so you can not dynamically determine the number of render runs._

Note how you can add a `condition` that will be evaluated for the template. Inside it, you can use
variables, or even `item`. If the condition evaluates to any one of `0`, `false`, `skip`, `no` the template will not be 
rendered. Note that the empty string counts as true, that means that if you do not specify a condition,
the template is rendered.

Also note how output directories are created for you on the fly if they don't exist.
  
The [golang template language](https://golang.org/pkg/text/template/#example_Template) is pretty 
versatile, vaguely similar to the .j2 templates used by ansible. Here's a very simple example
of how to include one of the parameters in your template output:

```
fmt.Println("{{ .helloMessage }}")
```

Assuming `someList` is a list variable, access the second entry as follows:

```
{{ index .someList 1 }}
```

Assuming `someMap` is a map variable with a field `message`, access it as follows:

```
{{ .someList.message }}
```

You can combine the two for structures with nested lists: `{{ (index .someList 0).someField }}`.

### Additional Template Functions

We include [Masterminds/sprig](https://github.com/Masterminds/sprig) when parsing any template,
which offers a collection of useful template functions, so you can do stuff like

```
{{ .helloMessage | upper }}
```

Read the sprig documentation, it adds much of what you would otherwise miss compared to ansible
j2 templates.

### Api for Generators

Given a generator's path, you can ask this library for the list of available generator names using
`generatorlib.FindGeneratorNames`.

Given a generator's path and one of the generator names, you can ask this library to give you the 
`api.GeneratorSpec` as a data structure read from the generator specification file (useful if
you wish to expose it as a service). Just call `generatorlib.ObtainGeneratorSpec`.

## Render Targets

A render target is a directory that contains a yaml file which records the name of the generator used
and all parameter values. We call this a *render specification file*. If you do not specify anything,
the generator expects the file to be called `generated-main.yaml`.

Example:
```
generator: main
parameters:
  helloMessage: hello world
  serviceName: 'my-service'
  serviceUrl: github.com/StephanHCB/temp
```

### Api for Rendering

Given a generator, you can ask this library to write out a render specification file with all parameters
set to their default value by calling `generatorlib.WriteRenderSpecWithDefaults`.

Given a generator and a target directory with an existing render specification file, you can call
`generatorlib.Render` to perform the rendering operation. For each template defined in the generator
specification, the corresponding target file is written.

*Note that existing target files will be overwritten by both operations. The idea is for you to have the 
target directory under source control, so you can then inspect the changes and pick what you would like to keep.*

### Example call to Render

```
func main() {
    ctx := context.TODO()
    request := &api.Request{
        SourceBaseDir: "/path/to/generator",
        TargetBaseDir: "/path/to/target",
    }
    response := generatorlib.Render(ctx, request)
}
```

The `api.Response` data structure returned by Render contains all potential `error`s, plus information about
all files rendered.

## Implementation Prerequisites

### Choose a Logging Framework Plugin

This library uses [go-autumn-logging](https://github.com/StephanHCB/go-autumn-logging)
to allow you to plug in the logging framework of your choice. You will need to include one of
the available specific wrappers among your dependencies. 

The simplest one, just using golang's standard
logger, is [go-autumn-logging-log](https://github.com/StephanHCB/go-autumn-logging-log).
We also have [go-autumn-logging-zerolog](https://github.com/StephanHCB/go-autumn-logging-zerolog).

If you do not want any logging, just call `aulogging.SetupNoLoggerForTesting` before calling any of the library 
functions. This will disable all logging, which is not really recommended:

```
import "github.com/StephanHCB/go-autumn-logging"

func init() {
    aulogging.SetupNoLoggerForTesting()
}
```
 
Or you can provide your own implementation of `auloggingapi.LoggingImplementation` and assign it to
`aulogging.Logger`.

## Build and test

This library uses go modules. If cloned outside your GOPATH, you can build and test it using
`go build ./...` and `go test ./... -coverpkg=./...`. This will also download all required dependencies.

We release automatically using 
[semantic-release](https://github.com/semantic-release/semantic-release/).
Please form your commit messages accordingly.

### Acceptance Tests (give you examples)

We have almost complete coverage with BDD-style 
[acceptance tests](https://github.com/StephanHCB/go-generator-lib/tree/master/test/acceptance). 

Running the tests and reading their code will give you lots of easy to understand examples, 
including most common error situations. Example for a happy path test:

```
=== RUN   TestRender_ShouldWriteExpectedFilesForDefault
2020/10/14 20:56:20 Given a valid generator source directory and a valid target directory
2020/10/14 20:56:20 Given a valid render spec file for generator main
2020/10/14 20:56:20 When Render is invoked
2020/10/14 20:56:20 Then the return value is as expected and the correct files are written
--- PASS: TestRender_ShouldWriteExpectedFilesForDefault (0.01s)
```

In the course of the test runs, several generator specs and templates are read from the
[test resources](https://github.com/StephanHCB/go-generator-lib/tree/master/test/resources),
and a number of render specs are written to the 
[test output directory](https://github.com/StephanHCB/go-generator-lib/tree/master/test/output).
