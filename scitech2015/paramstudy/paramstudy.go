package paramstudy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var gopath string

var BaseDir string
var MeshDir string

var ConfigName = "turb_NACA0012.cfg"

func init() {
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("gopath not set")
	}
	BaseDir = filepath.Join(gopath, "results", "scitech2015")
	MeshDir = filepath.Join(gopath, "results", "scitech2015", "meshes")
}

type ConfigMod struct {
	Mesh         string
	Aoa          float64
	CFL          float64
	Mglevel      int
	LinSolveIter int
	Limiter      float64
}

// DefaultConfigMode returns a new ConfigMod populated with the default values
func DefaultConfigMod() *ConfigMod {
	return &ConfigMod{
		Mesh:         "mesh_NACA0012_turb_897x257.su2",
		Aoa:          0,
		CFL:          1,
		Mglevel:      3,
		LinSolveIter: 3,
		Limiter:      6,
	}
}

func LoadRuns(runfile string) ([]*ConfigMod, error) {
	f, err := os.Open(runfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	/* Load in the cases to run */
	runs := make([]*ConfigMod, 0)

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&runs)
	return runs, err
}

// Directory returns the directory to which the config file will be saved
func Directory(run *ConfigMod) string {
	fmtstr := "%4.4g"

	meshstr := run.Mesh // mesh is already a string
	meshstr = strings.Replace(meshstr, ".", "_", -1)

	aoa := run.Aoa
	aoastr := fmt.Sprintf(fmtstr, aoa) // turn it into a string
	aoastr = strings.TrimSpace(aoastr)
	aoastr = strings.Replace(aoastr, ".", "_", -1) // replace '.' with '_'
	aoastr = "aoa_" + aoastr                       // pepend a helper

	cfl := run.CFL
	cflstr := fmt.Sprintf(fmtstr, cfl)
	cflstr = strings.TrimSpace(cflstr)
	cflstr = strings.Replace(cflstr, ".", "_", -1)
	cflstr = "cfl_" + cflstr

	mglevel := run.Mglevel
	mgstr := fmt.Sprintf("%v", mglevel) // just do the default printing for the integer
	mgstr = strings.Replace(mgstr, ".", "_", -1)
	mgstr = "mglev_" + mgstr

	iter := run.LinSolveIter
	iterstr := fmt.Sprintf("%v", iter) // just do the default printing for the integer
	iterstr = strings.Replace(iterstr, ".", "_", -1)
	iterstr = "solveiter_" + iterstr

	limiter := run.Limiter
	limstr := fmt.Sprintf(fmtstr, limiter) // just do the default printing for the integer
	limstr = strings.TrimSpace(limstr)
	limstr = strings.Replace(limstr, ".", "_", -1)
	limstr = "limiter_" + limstr

	return filepath.Join(BaseDir, "configs", meshstr, aoastr, cflstr, mgstr, iterstr, limstr)
}
