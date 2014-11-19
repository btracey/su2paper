package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"su2paper/scitech2015cases/paramstudy"

	"github.com/gonum/floats"
)

func main() {
	//basecase := filepath.Join(basedir, "baseconfig")
	//baseconfigfile := "turb_NACA0012.cfg"
	runloc := filepath.Join(paramstudy.BaseDir, "runs")
	runfilename := "cflsweep.txt"

	runfile := filepath.Join(runloc, runfilename)

	nCfl := 5
	minCFL := 1.0
	maxCFL := 50.0
	cfls := make([]float64, nCfl)        // allocate a new slice of float64
	floats.LogSpan(cfls, minCFL, maxCFL) // fill it with log-spaced points between
	// cfls := []float64{1, 20, 50}

	runs := make([]*paramstudy.ConfigMod, 0)

	for _, cfl := range cfls {
		c := paramstudy.DefaultConfigMod()
		c.CFL = cfl
		runs = append(runs, c)
	}
	os.MkdirAll(runloc, 0700)
	f, err := os.Create(runfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(runs, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	f.Write(b)

	log.Println("done")
}
