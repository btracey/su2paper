package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"su2paper/scitech2015cases/paramstudy"
	"sync"

	"github.com/btracey/su2tools/driver"
)

// This code runs a list of cases. It assumes the config files have already been
// generated using generate configs.
//
// This function takes two arguments. The first is the list of json cases. The second
// is an integer specifying the maximum number of simulations at the same time.
func main() {
	if len(os.Args) < 3 {
		log.Fatal("Must specify the run json file and the number of cores")
	}
	runfile := os.Args[1]
	runs, err := paramstudy.LoadRuns(runfile)
	if err != nil {
		log.Fatal(err)
	}

	nParallel, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("Second argument is the number of parallel runs. It must be an integer. %v found", os.Args[2])
	}

	runchan := make(chan *paramstudy.ConfigMod)

	errCases := make([]*paramstudy.ConfigMod, 0)
	errMux := &sync.Mutex{}

	var wg sync.WaitGroup
	// Launch all of the workers
	for i := 0; i < nParallel; i++ {
		wg.Add(1)
		go func(runchan <-chan *paramstudy.ConfigMod) {
			defer wg.Done()
			for run := range runchan {
				err := RunSU2(run)
				if err != nil {
					errMux.Lock()
					errCases = append(errCases, run)
					errMux.Unlock()
				}
			}
		}(runchan)
	}

	// Send all of the cases to be run
	for i, run := range runs {
		log.Print("Sending case ", i, " of ", len(runs))
		runchan <- run
	}
	close(runchan)
	wg.Wait()

	if len(errCases) == 0 {
		log.Print("All cases finished successfully")
	} else {
		log.Print("Cases finished with err. The cases that erred were:")
		for _, run := range runs {
			fmt.Printf(paramstudy.Directory(run))
		}
	}
}

func RunSU2(run *paramstudy.ConfigMod) error {
	// Need to set log file location
	dir := paramstudy.Directory(run)
	reldir, err := filepath.Rel(paramstudy.BaseDir, dir)
	if err != nil {
		log.Fatal(err)
	}
	configName := paramstudy.ConfigName
	drive := &driver.Driver{
		Name:   reldir,
		Config: configName,
		Wd:     dir,
		Stdout: "log.txt",
		Stderr: "su2err.txt",
	}
	err = drive.Load()
	if err != nil {
		return err
	}
	if drive.IsComputed(drive.Status()) {
		log.Print(dir + ": already computed")
		return nil
	}

	log.Print(reldir + ": starting computation")
	err = drive.Run(driver.Serial{})

	if err != nil {
		log.Print(reldir + ": stopped running with error")
		return err
	}
	log.Print(reldir + ": finished successfully")
	// Copy the restart to the solution
	err = drive.CopyRestartToSolution()
	return err
}
