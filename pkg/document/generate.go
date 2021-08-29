package document

import (
	"bytes"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func getOutputFile(dryRun bool) (*os.File, error) {
	if dryRun {
		return os.Stdout, nil
	}

	outputFile := viper.GetString("output-file")
	f, err := os.Create(outputFile)

	if err != nil {
		return nil, err
	}

	return f, err
}

func PrintDocumentation(valuesData *yaml.Node, templateFiles []string, dryRun bool, helmDocsVersion string) {
	log.Infof("Generating README Documentation")

	documentationTemplate, err := newDocumentationTemplate(templateFiles)

	if err != nil {
		log.Warnf("Error generating gotemplates: %s", err)
		return
	}

	chartTemplateDataObject, err := getChartTemplateData(valuesData, helmDocsVersion)
	if err != nil {
		log.Warnf("Error generating template data: %s", err)
		return
	}

	outputFile, err := getOutputFile(dryRun)
	if err != nil {
		log.Warnf("Could not open chart README file %s", err)
		return
	}

	if !dryRun {
		defer outputFile.Close()
	}

	var output bytes.Buffer
	err = documentationTemplate.Execute(&output, chartTemplateDataObject)
	if err != nil {
		log.Warnf("Error generating documentation: %s", err)
	}

	output = applyMarkDownFormat(output)
	_, err = output.WriteTo(outputFile)
	if err != nil {
		log.Warnf("Error generating documentation [markdown]: %s", err)
	}
}

func applyMarkDownFormat(output bytes.Buffer) bytes.Buffer {
	outputString := output.String()
	re := regexp.MustCompile(` \n`)
	outputString = re.ReplaceAllString(outputString, "\n")

	re = regexp.MustCompile(`\n{3,}`)
	outputString = re.ReplaceAllString(outputString, "\n\n")

	output.Reset()
	output.WriteString(outputString)
	return output
}
