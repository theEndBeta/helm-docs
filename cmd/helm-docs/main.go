package main

import (
	"os"
	"path"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/norwoodj/helm-docs/pkg/document"
	"github.com/norwoodj/helm-docs/pkg/helm"
)

func retrieveInfoAndPrintDocumentation(valuesFilePath string, templateFiles []string, waitGroup *sync.WaitGroup, dryRun bool) {
	defer waitGroup.Done()
	valuesFileInfo, err := helm.ParseChartInformation(valuesFilePath)

	if err != nil {
		log.Warnf("Error parsing information for chart %s, skipping: %s", valuesFilePath, err)
		return
	}

	document.PrintDocumentation(valuesFileInfo, templateFiles, dryRun, version)

}

func helmDocs(cmd *cobra.Command, _ []string) {
	initializeCli()

	var valuesFiles []string
	valuesFiles = viper.GetStringSlice("values-file")

	if len(valuesFiles) == 0 {
		log.Warn("As least one `values-file` must be provided.")
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Warnf("Error getting working directory: %s", err)
		return
	}

	// fullChartSearchRoot = path.Join(cwd, chartSearchRoot)

	// chartDirs, err := helm.FindChartDirectories(fullChartSearchRoot)
	// if err != nil {
	// 	log.Errorf("Error finding chart directories: %s", err)
	// 	os.Exit(1)
	// }

	// log.Infof("Found Chart directories [%s]", strings.Join(chartDirs, ", "))

	templateFiles := viper.GetStringSlice("template-files")
	log.Debugf("Rendering from optional template files [%s]", strings.Join(templateFiles, ", "))

	dryRun := viper.GetBool("dry-run")
	waitGroup := sync.WaitGroup{}

	var fullPath string
	for _, fname := range valuesFiles {
		waitGroup.Add(1)
		fullPath = path.Join(cwd, fname)

		// On dry runs all output goes to stdout, and so as to not jumble things, generate serially
		if dryRun {
			retrieveInfoAndPrintDocumentation(fullPath, templateFiles, &waitGroup, dryRun)
		} else {
			go retrieveInfoAndPrintDocumentation(fullPath, templateFiles, &waitGroup, dryRun)
		}
	}

	waitGroup.Wait()
}

func main() {
	command, err := newHelmDocsCommand(helmDocs)
	if err != nil {
		log.Errorf("Failed to create the CLI commander: %s", err)
		os.Exit(1)
	}

	if err := command.Execute(); err != nil {
		log.Errorf("Failed to start the CLI: %s", err)
		os.Exit(1)
	}
}
