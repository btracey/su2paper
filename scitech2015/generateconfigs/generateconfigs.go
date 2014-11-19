package main

import (
	"log"
	"os"
	"path/filepath"
	"su2paper/scitech2015cases/paramstudy"

	"github.com/btracey/su2tools/driver"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must specify the run json file")
	}
	runfile := os.Args[1]
	runs, err := paramstudy.LoadRuns(runfile)
	if err != nil {
		log.Fatal(err)
	}

	/* Load in the base config file */
	baseWd := filepath.Join(paramstudy.BaseDir, "baseconfig")
	baseConfigName := "turb_NACA0012.cfg"

	basedrive := &driver.Driver{
		Name:   "Base",
		Config: baseConfigName,
		Wd:     baseWd,
	}
	//base := basedrive.Options

	err = basedrive.Load()
	if err != nil {
		log.Fatal(err)
	}

	for _, run := range runs {
		// Load a new copy of the base config file
		config := &driver.Driver{
			Config: baseConfigName,
			Wd:     baseWd,
		}
		err = config.Load()
		if err != nil {
			log.Fatal(err)
		}

		// get the directory into which this config file will be saved
		dir := paramstudy.Directory(run)

		// Set the options for this config run
		err = setconfig(config, run, dir)
		if err != nil {
			log.Fatal("error setting Mesh file: " + err.Error())
		}

		// Save the config file
		os.MkdirAll(dir, 0700)
		f, err := os.Create(filepath.Join(dir, paramstudy.ConfigName))
		if err != nil {
			log.Fatal(err)
		}
		_, err = config.Options.WriteTo(f, nil)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}

// Set config sets the options from the config file. If the values in run have
// the zero value, use the defaults specified by the base file (aka leave
// the config file unchanged)
// This does not set the mesh config filename
func setconfig(config *driver.Driver, run *paramstudy.ConfigMod, dir string) error {

	// For debugging
	config.Options.ExtIter = 10

	// Most of the options are straightforward
	config.Options.Aoa = run.Aoa
	config.Options.CflNumber = run.CFL
	config.Options.Mglevel = uint16(run.Mglevel) // cast Mglevel as a uint16
	config.Options.LinearSolverIter = uint64(run.LinSolveIter)
	config.Options.LimiterCoeff = run.Limiter

	// The mesh file lives in a specific location, so we have to set the value
	// relative to the new directory.
	meshDir, err := filepath.Rel(dir, filepath.Join(paramstudy.MeshDir, run.Mesh))
	if err != nil {
		return err
	}
	config.Options.MeshFilename = meshDir
	return nil
}
