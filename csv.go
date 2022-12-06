// Copyright 2014 Oleku Konko All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// This module is a Table Writer  API for the Go Programming Language.
// The protocols were written in pure Go and works on windows and unix systems

package tablewriter

import (
	"encoding/csv"
	"errors"
	"io"
)

// NewCSV builds a Table writer that reads its rows from a csv.Reader.
func NewCSV(csvReader *csv.Reader, hasHeader bool, opts ...Option) (*Table, error) {
	options := opts

	if hasHeader {
		header, err := csvReader.Read()
		if err != nil {
			return nil, err
		}

		options = append(options, WithHeader(header))
	}

	table := New(options...)

	for {
		record, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		table.Append(record)
	}

	return table, nil
}
