// Copyright 2022 Eugene Nosenko
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cli
/*
Gopoly generates custom unmarshaling functions for provided interface.
It uses built-in go tooling to scan packages and look for implementing types as well
as polymorphic fields.

By default, it will look for a .gopoly.yaml file as a main source of the configuration but configuration can be
provided via command line. Due to complicated configuration inputs it's preferable to provide config via .gopoly.yaml
file. If both are provided, inputs from command line overwrite inputs from config file.

Usage:

	gopoly [commands] [flags]

The commands are:
	init
		Creates a .gopoly.yaml configuration file with dummy values.

The flags are:

	-c
		Provide path to the config file. Default value is .gopoly.yaml
	-p
		Scoped package where models are located.
	-d
		Decoding strategy to be used when unmarshaling functions are generated.
		Can be either 'strict' or 'discriminator'
	-o
		Output filename that will contain generated code. Please be mindful that if you provide
		a full or relative path, it will be ignored. Since unmarshaling functions need to be
		generated in the same package.
	-m
		Marker method. Marker-method is a way 'mark' types that belong to a specific interface,
		basically serving as a metadata. Default value is Is{{.Name}}. And type will be taken from
		the interface
	-t
		Types information.

# Examples:

Create a configuration file with dummy values:

	gopoly init

Generate unmarshaling functions based on the default config file:

	gopoly

Generate unmarshaling functions based on the custom named config file:

	gopoly -c .gopoly-config.yaml

Generate unmarshaling based on command input only:

	gopoly -p "github.com/username/example/models" \
		-o "out.gen.go" \
		-d "strict" \
		-m "IsRunner" \
		-t 'Runner subtypes=A,B'

Generate unmarshaling based on custom config file and command input :

	gopoly -c .gopoly-config.yaml \
		-p "github.com/username/example/models" \
		-o "out.gen.go" \
		-d "strict" \
		-m "IsRunner" \
		-t 'Runner subtypes=A,B'
*/
package cli
