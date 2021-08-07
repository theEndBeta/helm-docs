package document

import (
	"fmt"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/theEndBeta/yaml-docs/pkg/helm"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type valueRow struct {
	Key             string
	Type            string
	AutoDefault     string
	Default         string
	AutoDescription string
	Description     string
	Column          int
	LineNumber      int
}

type chartTemplateData struct {
	helm.DocumentationInfo
	HelmDocsVersion string
	Values          []valueRow
}

func getSortedValuesTableRows(documentRoot *yaml.Node, chartValuesDescriptions map[string]helm.ValueDescription) ([]valueRow, error) {
	valuesTableRows, err := createValueRowsFromField(
		"",
		nil,
		documentRoot,
		chartValuesDescriptions,
		true,
	)

	if err != nil {
		return nil, err
	}

	sortOrder := viper.GetString("sort-values-order")
	if sortOrder == FileSortOrder {
		sort.Slice(valuesTableRows, func(i, j int) bool {
			if valuesTableRows[i].LineNumber == valuesTableRows[j].LineNumber {
				return valuesTableRows[i].Column < valuesTableRows[j].Column
			}

			return valuesTableRows[i].LineNumber < valuesTableRows[i].LineNumber
		})
	} else if sortOrder == AlphaNumSortOrder {
		sort.Slice(valuesTableRows, func(i, j int) bool {
			return valuesTableRows[i].Key < valuesTableRows[j].Key
		})
	} else {
		log.Warnf("Invalid sort order provided %s, defaulting to %s", sortOrder, AlphaNumSortOrder)
		sort.Slice(valuesTableRows, func(i, j int) bool {
			return valuesTableRows[i].Key < valuesTableRows[j].Key
		})
	}

	return valuesTableRows, nil
}


func getChartTemplateData(chartDocumentationInfo helm.DocumentationInfo, helmDocsVersion string) (chartTemplateData, error) {
	// handle empty values file case
	if chartDocumentationInfo.Values.Kind == 0 {
		return chartTemplateData{
			DocumentationInfo: chartDocumentationInfo,
			HelmDocsVersion:        helmDocsVersion,
			Values:                 make([]valueRow, 0),
		}, nil
	}

	if chartDocumentationInfo.Values.Kind != yaml.DocumentNode {
		return chartTemplateData{}, fmt.Errorf("invalid node kind supplied: %d", chartDocumentationInfo.Values.Kind)
	}
	if chartDocumentationInfo.Values.Content[0].Kind != yaml.MappingNode {
		return chartTemplateData{}, fmt.Errorf("values file must resolve to a map, not %s", strconv.Itoa(int(chartDocumentationInfo.Values.Kind)))
	}

	valuesTableRows, err := getSortedValuesTableRows(chartDocumentationInfo.Values.Content[0], chartDocumentationInfo.ValuesDescriptions)

	if err != nil {
		return chartTemplateData{}, err
	}

	return chartTemplateData{
		DocumentationInfo: chartDocumentationInfo,
		HelmDocsVersion:        helmDocsVersion,
		Values:                 valuesTableRows,
	}, nil
}
