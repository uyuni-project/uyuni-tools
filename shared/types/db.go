// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package types

import "strconv"

// DBFlags can store all values required to connect to a database.
type DBFlags struct {
	Host     string
	Name     string
	Port     int
	User     string
	Password string
	Provider string
	Admin    struct {
		User     string
		Password string
	}
}

// IsLocal indicates if the database is a local or a third party one.
func (flags *DBFlags) IsLocal() bool {
	return flags.Host == "" || flags.Host == "db" || flags.Host == "reportdb"
}

func (flags *DBFlags) GetPort() string {
	port := "5432"
	if flags.Port != 0 {
		port = strconv.Itoa(flags.Port)
	}
	return port
}
