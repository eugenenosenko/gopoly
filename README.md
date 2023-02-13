gopoly [![Build Status](https://github.com/eugenenosenko/gopoly/actions/workflows/build.yml/badge.svg)](https://github.com/eugenenosenko/gopoly/actions)
=================

## what is `gopoly`?

**gopoly** is a tool that generates custom `Unmarshal*` methods for the polymorphic types.

## why?

The idea came about after using [**gqlgen**](https://github.com/99designs/gqlgen) library, which
generates union types based on interfaces. Since Go doesn't support custom `Unmarshal*` methods for
interfaces, it was hard to unmarshal into GQL-generated models.

Partly inspired by [**gopolyjson**](https://github.com/polyfloyd/gopolyjson) library.

## goals
* [x] support polymorphic decoding based on two algorithms (discriminator / strict)
* [x] support decoding of multiple field types: scalar/slices/maps
* [x] support decoding of polymorphic fields
* [ ] support payload formats other than JSON

## install
```
go install -v github.com/eugenenosenko/gopoly@latest
```

## usage
1) run `gopoly init`; this will create empty config file.
2) provide interfaces, variants, marker methods to the configuration ([yaml](#yaml)|[cmd-line](#command-line))
3) run `gopoly`
4) program automatically discovers the variants and generates custom unmarshaling functions

You can now use `Unmarashl<Interface>JSON` functions in your code.

**IMPORTANT NOTE**:

Your marker methods have to comply with following requirements:
- take no arguments
- return nothing

## sample application configuration:

#### GO code

```go
// events.go
package events

import (
    "time"

    u "github.com/username/example/users"
)

type UserEvent interface{ IsUserEvent() }

type UserDeletedEvent struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    User u.User `json:"user"`
}
func (e UserDeletedEvent) IsUserEvent() {}

type UserCreatedEvent struct {
    ID   string `json:"id"`
    Type string `json:"type"`
    User u.User `json:"user"`
}
func (e UserCreatedEvent) IsUserEvent() {}

// users.go
package users

type User interface{ IsUser() }

type RegularUser struct {
    ID       string    `json:"id"`
    Type     string    `json:"kind"`
    // other properties
    Contacts []Contact `json:"contacts"`
}
func (a RegularUser) IsUser() {}

type PrivilegedUser struct {
    ID         string    `json:"id"`
    Type       string    `json:"kind"`
    // other properties
    Contacts   []Contact `json:"contacts"`
}
func (a PrivilegedUser) IsUser() {}

type Contact interface{ IsContact() }

type BusinessContact struct {
    ID           string `json:"id"`
    BusinessName string `json:"business_name"`
    Phone        string `json:"phone"`
}
func (c BusinessContact) IsContact() {}

type PrivateContact struct {
    ID       string   `json:"id"`
    FullName string   `json:"fullname"`
}
func (c PrivateContact) IsContact() {}
```

#### configuration file .gopoly.yaml

```yaml
types:
  - name: UserEvent
    variants:
      - UserDeletedEvent
      - UserCreatedEvent
    decoding_strategy: "discriminator"
    discriminator:
      field: "type"
      mapping:
        DELETED: UserDeletedEvent
        CREATED: UserCreatedEvent
    output:
      filename: "events.gen.go"
  - name: Contact # this type inherits most of its configuration from base
    package: "github.com/eugenenosenko/gopoly/tests/e2e/testdata/users"
  - name: User
    package: "github.com/eugenenosenko/gopoly/tests/e2e/testdata/users"
    variants:
      - RegularUser
      - PrivilegedUser
      - BannedUser
    discriminator:
      field: "kind"
      mapping:
        REGULAR: RegularUser
        PRIVILEGED: PrivilegedUser
        BANNED: BannedUser
    output:
      filename: "internal/models/users.gen.go"
marker_method: "Is{{ .Name }}"
decoding_strategy: "strict"
package: "github.com/eugenenosenko/gopoly/tests/e2e/testdata/events"
output:
  filename: "gopoly.gen.go"
```

## how does GOPOLY work?
`gopoly` executes following steps:
1) config processing
2) packages scanning
3) interfaces & types validation
4) collecting source information and building internal representation
5) generating unmarshaling functions using a GO template

## configuration. YAML vs command-line
Configuration can be provided via YAML file or command-line or both. In case both are provided then command-line configuration
takes precedence over YAML config. This is helpful when you can't change original config but want to override some parts of it.

Configuration itself can be divided into two types:

1) parent configuration
2) interface specific configuration

Parent configuration works sets default configuration values for all interface types, but if any interface is located in a different package
or wants to define a different marker-method, then that configuration will override the parent configuration.

### yaml

JSON-schema for the config file can be found [here](config-json-schema.json)

### command-line

| flag | short description                                             | example                               |
|:----:|---------------------------------------------------------------|---------------------------------------|
| `-c` | config filename path                                          | `-c "myconfig.yml"`                   |
| `-p` | package name                                                  | `-p "github.com/user/lib/models"`     |
| `-d` | decoder strategy `strict` or `discriminator`                  | `-d "strict"`                         |
| `-m` | marker method [marker-interfaces], string or template         | `-m "Is{{.Name}}"` or `-m "IsMyType"` |
| `-t` | variant types' information, i.e. variants, discriminator etc. | `-t "Runner variants=A,B"`            |

[marker-interfaces]: https://en.wikipedia.org/wiki/Marker_interface_pattern

#### types `-t` flag options
It's possible to provide additional variant specific configuration via `-t` flag by providing it with required options

| option                  | description                                            | example                                     |
|-------------------------|--------------------------------------------------------|---------------------------------------------|
| `variants`              | type variants that implement the i-face                | `variants=SlowRunner,FastRunner`            |
| `marker_method`         | marker method name, can be template                    | `marker_method=IsRunner`                    |
| `decoding_strategy`     | either `strict` or `discriminator`                     | `decoding_strategy=discriminator`           |
| `discriminator.field`   | field name that determines which discriminator mapping | `discriminator.field=runner_type`           |
| `discriminator.mapping` | key-value mapping of discriminator => type variant     | `discriminator.mapping=slow:Slow,fast:Fast` |

An example of such configuration would be:
```
-t "Runner variants=SlowRunner,FastRunner marker_method=IsRunner decoding_strategy=discriminator discriminator.field=type discriminator.mapping=slow:SlowRunner,fast:FastRunner
```

## glossary
| term                            | definition                                                                                         |
|---------------------------------|----------------------------------------------------------------------------------------------------|
| interface                       | a data type describing a set of method signatures                                                  |
| variant                         | a concrete type implementing a specific interface                                                  |
| marker method                   | a method on interface, implemented by types to provide run-time type information                   |
| decoding strategy               | algorithm used to decode incoming payload into polymorphic structure                               |
| strict decoding strategy        | strict decoding will try to match incoming payload against type without allowing unknown fields    |
| discriminator decoding strategy | discriminator decoding will decode payload into a variant based on the discriminator field mapping |
| discriminator                   | a field, value of which will be used to determine the concrete type payload should be decoded into |

[**about marker iface pattern**](https://en.wikipedia.org/wiki/Marker_interface_pattern)

