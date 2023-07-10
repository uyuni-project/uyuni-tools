package setup

type certificate struct {
	useExistingCertificate bool
	cnames                 []string
	state                  string
	city                   string
	country                string
	org                    string
	ou                     string
	password               string
	email                  string
}

type database struct {
	name          string
	host          string
	user          string
	password      string
	port          int
	protocol      string
	adminUser     string
	adminPassword string
	provider      string
}

type manager struct {
	user       string
	password   string
	email      string
	emailFrom  string
	db         database
	mirrorPath string
	issParent  string
}
