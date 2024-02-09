package main

import (
	"fmt"
	"log"
	"os"
	"twistingmercury/forge/scaffolder"

	"github.com/spf13/pflag"
)

var (
	helpFlag  = pflag.BoolP("help", "h", false, "Show help for forge.")
	verFlag   = pflag.BoolP("version", "v", false, "Show the current version of forge.")
	tmpltFlag = pflag.StringP("template-path", "t", "", "Required: The full path to the template (*.zip) file.")
	pnameFlag = pflag.StringP("project-name", "p", "", "Required: The directory name for the project. It will also be the name of the project.")
	mnameFlag = pflag.StringP("module-name", "m", "", "Required: The module path, root namespace, root package name, etc., for the project.")
)

var (
	version string
	date    string
	commit  string
)

func main() {
	pflag.Parse()
	showHelp(*helpFlag)
	showVersion()
	if len(*pnameFlag) == 0 || len(*mnameFlag) == 0 {
		showHelp(true)
	}

	template, projectName, module := flagValues()

	if err := scaffolder.CreateProject(template, projectName, module); err != nil {
		wd, _ := os.Getwd()
		scaffolder.Rollback(wd)
		log.Fatal(err, "could not create project")
	}
}

func showVersion() {
	if *verFlag {
		println(logo())
		fmt.Printf("Version: %s, Build Date: %s, VCS Commit: %s\n", version, date, commit)
		os.Exit(0)
	}
}

func showHelp(show bool) {
	if show {
		println(logo())
		println("Usage: forge [options]")
		println("Example: forge -p myProject -m github.com/<gh user name>/<project name>")
		pflag.PrintDefaults()
		os.Exit(0)
	}
}

func flagValues() (projPath, modPath, template string) {
	switch {
	case len(*pnameFlag) == 0:
		log.Fatal("value for --project-name is required")
	case len(*tmpltFlag) == 0:
		log.Fatal("value for --template-path is required")
	case len(*mnameFlag) == 0:
		log.Fatal("value for --module-name is required")
	}
	return *pnameFlag, *mnameFlag, *tmpltFlag
}

func logo() string {
	return `
██████╗ ██████╗  ██████╗      ██╗███████╗ ██████╗████████╗
██╔══██╗██╔══██╗██╔═══██╗     ██║██╔════╝██╔════╝╚══██╔══╝
██████╔╝██████╔╝██║   ██║     ██║█████╗  ██║        ██║   
██╔═══╝ ██╔══██╗██║   ██║██   ██║██╔══╝  ██║        ██║   
██║     ██║  ██║╚██████╔╝╚█████╔╝███████╗╚██████╗   ██║   
╚═╝     ╚═╝  ╚═╝ ╚═════╝  ╚════╝ ╚══════╝ ╚═════╝   ╚═╝   
                                                          
███████╗ ██████╗ ██████╗  ██████╗ ███████╗                
██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██╔════╝                
█████╗  ██║   ██║██████╔╝██║  ███╗█████╗                  
██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝                  
██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗                
╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝     
`
}
