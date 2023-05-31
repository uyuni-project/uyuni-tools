package main

//TODO required go version >= 18
//zypper in libbtrfs-devel libgpgme-devel device-mapper-devel gpgme libassuan >= 2.5.3
//systemct start podman

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/containers/podman/v4/pkg/specgenutil"
)

func getImage(context context.Context, imageName string) {
	fmt.Printf("%v\n", imageName)
	_, err := images.Pull(ctx, imageName, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createContainer(opt entities.ContainerCreateOptions, imageName string) {
	s := specgen.NewSpecGenerator(imageName, false)
	if err := specgenutil.FillOutSpecGen(s, &opt, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.Name = opt.Name
	if _, err := containers.CreateWithSpec(ctx, s, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Container created.")
}

func readJsonOpt(config_json_filename string) entities.ContainerCreateOptions {
  var base_option entities.ContainerCreateOptions
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile(config_json_filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	err = json.Unmarshal(content, &base_option)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return base_option
}
