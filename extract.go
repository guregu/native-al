// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/kardianos/osext"
)

var libraryPath string

func init() {
	// Initialize blob data slice.
	initBlob()

	var usr *user.User
	var err error

	// First, try same directory as the binary.
	if sameDir, err := osext.ExecutableFolder(); err == nil {
		lib := filepath.Join(sameDir, blobFileName)
		if _, err := os.Stat(lib); err == nil {
			libraryPath = lib
			log.Println("Using same-dir OpenAL:", lib)
			goto load // Sorry
		}
	}

	// Determine where the library should actually be placed on the system
	usr, err = user.Current()
	if err != nil {
		log.Fatal(err)
	}
	libraryPath = filepath.Join(usr.HomeDir, ".azul3d")
	libraryPath = filepath.Join(libraryPath, blobFileName)

	// Check if the library already exists at that location -- if it does it
	// means we've already extracted it there or the user has placed their own
	// implementation of the library there (per the LGPL restrictions).
	_, err = os.Stat(libraryPath)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(libraryPath), 0777)
		if err != nil {
			log.Fatal(err)
		}

		// There is no dynamic library at that location, we can extract our
		// copy of it then.
		err = ioutil.WriteFile(libraryPath, blob, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

load:
	// Load the library
	err = loadLibrary(libraryPath)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize pointers
	alInit()
	alcInit()
}
