/*
 * Copyright (C) 2020, MinIO, Inc.
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

package main

import (
	"errors"

	toml "github.com/BurntSushi/toml"
)

// Input structure holds the data fetched from parsed input file
type Input struct {
	Namespace    string   `toml:"namespace"`
	Capacity     string   `toml:"capacity"`
	StorageClass string   `toml:"storageClass"`
	Hosts        []string `toml:"hosts"`
	Paths        []string `toml:"paths"`
}

func (i *Input) validate() error {
	if len(i.Hosts) == 0 {
		return errors.New("Please provide at least one host name in input file")
	}

	if len(i.Paths) == 0 {
		return errors.New("Please provide at least one path in input file")
	}

	if i.Capacity == "" {
		return errors.New("Please provide capacity in input file")
	}

	return nil
}

func parseInput(path string) (Input, error) {
	var i Input
	_, err := toml.DecodeFile(path, &i)
	if err != nil {
		return i, err
	}
	err = i.validate()
	if err != nil {
		return i, err
	}
	return i, nil
}
