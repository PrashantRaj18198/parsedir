# parsedir

parse golang templates directory via cli

## Installation

Download the latest release

```bash
# the below link is for linux (amd64) just replace the compressed file name with the appropriate one if you are on other system
# Latest release: https://github.com/PrashantRaj18198/parsedir/releases/latest
curl -L https://github.com/PrashantRaj18198/parsedir/releases/latest/download/parsedir_linux_amd64.tar.gz -o parsedir.tar.gz
tar -xzf parsedir.tar.gz
mv parsedir /usr/bin/
```

## How to use?

Example folder structure:

|example
|---- {{.dog.name}}
|-------- {{dog.name}}.yaml
|---- {{range .pets}}{{.name}}
|-------- detail.txt

```bash
# Run the below command to see parsedir in action
parsedir --vars-file=config.yaml --template-dir example/ --out-dir result/ #json is also supported as input
```

The filepath will be generated from the config.yaml and written to result/ dir.

## Things to keep in mind

Create a template similar to the the `example/` folder. The filenames and the contents can be golang templates which is parsed and saved to out-dir.
Things to note:

- If the path is a template, and the template fails to parse likely because of missing variable the template is skipped with a warning.
- If you want to loop the file multiple times, use "{{range .some.variable.with.list.of.values }}" before the name, parsedir will read the file
  and place the content inside the range function

```go
// path template used below
{{range .some.variable.with.list.of.values}}
// file path goes here
// parsedir separator used to slice the generated list
{{end}}
// content template used below
{{range .some.variable.with.list.of.values }}
// file content goes here
// parsedir separator used to slice the generated list
{{end}} // end block is automatically added
```
