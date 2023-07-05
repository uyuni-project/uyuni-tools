package main

type cliArg struct {
	defaultValue string
	argument     string
	description  string
	envVariable  string
}

type taggedReference struct {
	namedRepository string
	tag             string
}

func (t taggedReference) String() string {
	return t.namedRepository + ":" + t.tag
}

type certificate struct {
	useExistingCertificate string
	state                  string
	city                   string
	country                string
	org                    string
	ou                     string
	pwd                    string
	email                  string
}

type database struct {
	name     string
	host     string
	password string
	port     string
	protocol string
}

type manager struct {
	user     string
	password string
	email    string
	db       database
}

type remoteServer struct {
	user     string
	hostname string
}

type container struct {
	remoteServer
	name string
}
