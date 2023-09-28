package ssl

import (
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type CaChain struct {
	Root         string
	Intermediate []string
}

type SslPair struct {
	Cert string
	Key  string
}

func CheckPaths(chain *CaChain, serverPair *SslPair) {
	mandatoryFile(chain.Root, "root CA")
	for _, ca := range chain.Intermediate {
		optionalFile(ca)
	}
	mandatoryFile(serverPair.Cert, "server certificate")
	mandatoryFile(serverPair.Key, "server key")
}

func mandatoryFile(file string, msg string) {
	if file == "" {
		log.Fatal().Msgf("%s is required", msg)
	}
	optionalFile(file)
}

func optionalFile(file string) {
	if file != "" && !utils.FileExists(file) {
		log.Fatal().Msgf("%s file is not accessible", file)
	}
}
