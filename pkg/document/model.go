package document

import (
	"fmt"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type valueRow struct {
	Key             string
	Type            string
	Default         string
	Description     string
	Column          int
	LineNumber      int
}

type chartTemplateData struct {
	YamlDocsVersion string
	Values          []valueRow
}

func getSortedValuesTableRows(documentRoot *yaml.Node) ([]valueRow, error) {
	valuesTableRows, err := createValueRowsFromField(
		"",
		nil,
		documentRoot,
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
	} else { // Default to AlphaNumSortOrder
		if sortOrder == "" {
			log.Debugf("No sort order provided, defaulting to %s", AlphaNumSortOrder)
		} else if sortOrder != AlphaNumSortOrder {
			log.Infof("Invalid sort order `%s`, defaulting to %s", sortOrder, AlphaNumSortOrder)
		}

		sort.Slice(valuesTableRows, func(i, j int) bool {
			return valuesTableRows[i].Key < valuesTableRows[j].Key
		})
	}

	return valuesTableRows, nil
}


func getChartTemplateData(valuesData *yaml.Node, yamlDocsVersion string) (chartTemplateData, error) {
	// handle empty values file case
	if valuesData.Kind == 0 {
		return chartTemplateData{
			YamlDocsVersion:        yamlDocsVersion,
			Values:                 make([]valueRow, 0),
		}, nil
	}

	if valuesData.Kind != yaml.DocumentNode {
		return chartTemplateData{}, fmt.Errorf("invalid node kind supplied: %d", valuesData.Kind)
	}
	if valuesData.Content[0].Kind != yaml.MappingNode {
		return chartTemplateData{}, fmt.Errorf("values file must resolve to a map, not %s", strconv.Itoa(int(valuesData.Kind)))
	}

	valuesTableRows, err := getSortedValuesTableRows(valuesData.Content[0])

	if err != nil {
		return chartTemplateData{}, err
	}

	return chartTemplateData{
		YamlDocsVersion:        yamlDocsVersion,
		Values:                 valuesTableRows,
	}, nil
}
