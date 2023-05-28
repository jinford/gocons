# gocons

gocons is a cli tool that generate constructor function & getter methods for structs

## Usage

```bash
NAME:
   gocons - generate constructor function & getter methods for structs

USAGE:
   gocons [global options] command [command options] [arguments...]

VERSION:
   dev

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --src value     path of file that declares structs
   --tag value     name of target struct tag (default: "cons")
   --output value  output: 'file', 'stdout' (default: "file")
   --help, -h      show help
   --version, -v   print the version

```

## Synopsis

Prepare a struct definition file with `go:generate`. It can generate a getter method for a private field by using `cons:"getter"` tag. If the field is exported, the getter method will be not generated. Also, the tag's key can be changed by option.

```go
//go:generate --src=$GOFILE
type Person struct {
    id       string         `cons:"getter"`
    name     string         `cons:"getter"`
    tags     []string       `cons:"getter"`
    desc     sql.NullString `cons:"getter"`
    *deposit `cons:"getter"`
}

type deposit struct {
    charge decimal.Decimal `cons:"getter"`
}
```

Run `go generate ./...`, then constructors and getters code for all structs in source file are generated in same directory.

```go
func NewPerson(
	id string,
	name string,
	tags []string,
	desc sql.NullString,
	deposit *deposit,
) *Person {
	return &Person{
		id:      id,
		name:    name,
		tags:    tags,
		desc:    desc,
		deposit: deposit,
	}
}

func (x *Person) Id() string {
	return x.id
}

func (x *Person) Name() string {
	return x.name
}

...

```

