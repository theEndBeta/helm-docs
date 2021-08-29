package document

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func parseYamlValues(yamlValues string) *yaml.Node {
	var chartValues yaml.Node
	err := yaml.Unmarshal([]byte(strings.TrimSpace(yamlValues)), &chartValues)

	if err != nil {
		panic(err)
	}

	return chartValues.Content[0]
}

func TestEmptyValues(t *testing.T) {
	yamlValues := parseYamlValues(`{}`)
	valuesRows, err := getSortedValuesTableRows(yamlValues)
	assert.Nil(t, err)
	assert.Len(t, valuesRows, 0)
}

func TestSimpleValues(t *testing.T) {
	yamlValues := parseYamlValues(`
echo: 0
foxtrot: true
hello: "world"
oscar: 3.14159
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`0`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "`3.14159`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].Description)
}

func TestSimpleValuesWithDescriptions(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- echo
echo: 0
# -- foxtrot
foxtrot: true
# -- hello
hello: "world"
# -- oscar
oscar: 3.14159
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "`0`", valuesRows[0].Default)
	assert.Equal(t, "echo", valuesRows[0].Description)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "`true`", valuesRows[1].Default)
	assert.Equal(t, "foxtrot", valuesRows[1].Description)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[2].Default)
	assert.Equal(t, "hello", valuesRows[2].Description)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "`3.14159`", valuesRows[3].Default)
	assert.Equal(t, "oscar", valuesRows[3].Description)
}

func TestSimpleValuesWithDescriptionsAndDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- echo
# @default -- some
echo: 0
# -- foxtrot
# @default -- explicit
foxtrot: true
# -- hello
# @default -- default
hello: "world"
# -- oscar
# @default -- values
oscar: 3.14159
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "echo", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type, intType)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "echo", valuesRows[0].Description)

	assert.Equal(t, "foxtrot", valuesRows[1].Key)
	assert.Equal(t, boolType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "foxtrot", valuesRows[1].Description)

	assert.Equal(t, "hello", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "hello", valuesRows[2].Description)

	assert.Equal(t, "oscar", valuesRows[3].Key)
	assert.Equal(t, floatType, valuesRows[3].Type)
	assert.Equal(t, "values", valuesRows[3].Default)
	assert.Equal(t, "oscar", valuesRows[3].Description)
}

func TestNestedValues(t *testing.T) {
	yamlValues := parseYamlValues(`
recursive:
  echo: cat
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestNestedValuesWithDescriptions(t *testing.T) {
	yamlValues := parseYamlValues(`
recursive:
  # -- echo
  echo: cat
# -- oscar
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "oscar", valuesRows[0].Description)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestNestedValuesWithDescriptionsAndDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
recursive:
  # -- echo
  # @default -- custom
  echo: cat
# -- oscar
# @default -- default
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "default", valuesRows[0].Default)
	assert.Equal(t, "oscar", valuesRows[0].Description)

	assert.Equal(t, "recursive.echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "custom", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestEmptyObject(t *testing.T) {
	yamlValues := parseYamlValues(`
recursive: {}
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key, "oscar")
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{}`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestEmptyObjectWithDescription(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- an empty object
recursive: {}
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{}`", valuesRows[1].Default)
	assert.Equal(t, "an empty object", valuesRows[1].Description)
}

func TestEmptyObjectWithDescriptionAndDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- an empty object
# @default -- default
recursive: {}
oscar: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "oscar", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "recursive", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "default", valuesRows[1].Default)
	assert.Equal(t, "an empty object", valuesRows[1].Description)
}
func TestEmptyList(t *testing.T) {
	yamlValues := parseYamlValues(`
birds: []
echo: cat
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[]`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestEmptyListWithDescriptions(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- birds
birds: []
# -- echo
echo: cat
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[]`", valuesRows[0].Default)
	assert.Equal(t, "birds", valuesRows[0].Description)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestEmptyListWithDescriptionsAndDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- birds
# @default -- explicit
birds: []
# -- echo
# @default -- default value
echo: cat
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "birds", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "birds", valuesRows[0].Description)

	assert.Equal(t, "echo", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default value", valuesRows[1].Default)
	assert.Equal(t, "echo", valuesRows[1].Description)
}

func TestListOfStrings(t *testing.T) {
	yamlValues := parseYamlValues(`
cats: [echo, foxtrot]
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

}

func TestListOfStringsWithDescriptions(t *testing.T) {
	yamlValues := parseYamlValues(`
cats:
  # -- the black one
  - echo
  # -- the friendly one
  - foxtrot
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

}

func TestListOfStringsWithDescriptionsAndDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
	
cats:
  # -- the black one
  # @default -- explicit
  - echo
  # -- the friendly one
  # @default -- default value
  - foxtrot
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "cats[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "cats[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default value", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

}

func TestListOfObjects(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[3].Default)
	assert.Equal(t, "", valuesRows[3].Description)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].Description)
}

func TestListOfObjectsWithDescriptions(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  - elements: 
      # -- the black one
      - echo

      # -- the friendly one
      - foxtrot
    type: cat
  - elements: 
      # -- the sleepy one
      - oscar
    type: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"echo\"`", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[3].Default)
	assert.Equal(t, "the sleepy one", valuesRows[3].Description)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].Description)
}

func TestListOfObjectsWithDescriptionsAndDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  - elements: 
      # -- the black one
      # @default -- explicit
      - echo

      # -- the friendly one
      # @default -- default
      - foxtrot
    type: cat
  - elements: 
      # -- the sleepy one
      # @default -- value
      - oscar
    type: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 5)

	assert.Equal(t, "animals[0].elements[0]", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "explicit", valuesRows[0].Default)
	assert.Equal(t, "the black one", valuesRows[0].Description)

	assert.Equal(t, "animals[0].elements[1]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "default", valuesRows[1].Default)
	assert.Equal(t, "the friendly one", valuesRows[1].Description)

	assert.Equal(t, "animals[0].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"cat\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "value", valuesRows[3].Default)
	assert.Equal(t, "the sleepy one", valuesRows[3].Description)

	assert.Equal(t, "animals[1].type", valuesRows[4].Key)
	assert.Equal(t, stringType, valuesRows[4].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[4].Default)
	assert.Equal(t, "", valuesRows[4].Description)
}

func TestDescriptionOnList(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- all the animals of the house
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "`[{\"elements\":[\"echo\",\"foxtrot\"],\"type\":\"cat\"},{\"elements\":[\"oscar\"],\"type\":\"dog\"}]`", valuesRows[0].Default)
	assert.Equal(t, "all the animals of the house", valuesRows[0].Description)
}

func TestDescriptionAndDefaultOnList(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- all the animals of the house
# @default -- cat and dog
animals:
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "cat and dog", valuesRows[0].Default)
	assert.Equal(t, "all the animals of the house", valuesRows[0].Description)
}

func TestDescriptionAndDefaultOnObjectUnderList(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- all the cats of the house
  # @default -- only cats here
  - elements: [echo, foxtrot]
    type: cat
  - elements: [oscar]
    type: dog
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals[0]", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "only cats here", valuesRows[0].Default)
	assert.Equal(t, "all the cats of the house", valuesRows[0].Description)

	assert.Equal(t, "animals[1].elements[0]", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"oscar\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)

	assert.Equal(t, "animals[1].type", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`\"dog\"`", valuesRows[2].Default)
	assert.Equal(t, "", valuesRows[2].Description)
}

func TestDescriptionOnObjectUnderObject(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- animals listed by their various characteristics
  byTrait:
    friendly: [foxtrot, oscar]
    mean: [echo]
    sleepy: [oscar]
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.byTrait", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "`{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}`", valuesRows[0].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[0].Description)
}

func TestDescriptionAndDefaultOnObjectUnderObject(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- animals listed by their various characteristics
  # @default -- animals, you know
  byTrait:
    friendly: [foxtrot, oscar]
    mean: [echo]
    sleepy: [oscar]
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.byTrait", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "animals, you know", valuesRows[0].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[0].Description)
}

func TestDescriptionsDownChain(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- animal stuff
animals:
  # -- animals listed by their various characteristics
  byTrait:
    # -- the friendly animals of the house
    friendly:
      # -- best cat ever
      - foxtrot
      - oscar
    mean: [echo]
    sleepy: [oscar]
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "`{\"byTrait\":{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}}`", valuesRows[0].Default)
	assert.Equal(t, "animal stuff", valuesRows[0].Description)

	assert.Equal(t, "animals.byTrait", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "`{\"friendly\":[\"foxtrot\",\"oscar\"],\"mean\":[\"echo\"],\"sleepy\":[\"oscar\"]}`", valuesRows[1].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[1].Description)

	assert.Equal(t, "animals.byTrait.friendly", valuesRows[2].Key)
	assert.Equal(t, listType, valuesRows[2].Type)
	assert.Equal(t, "`[\"foxtrot\",\"oscar\"]`", valuesRows[2].Default)
	assert.Equal(t, "the friendly animals of the house", valuesRows[2].Description)

	assert.Equal(t, "animals.byTrait.friendly[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "`\"foxtrot\"`", valuesRows[3].Default)
	assert.Equal(t, "best cat ever", valuesRows[3].Description)
}

func TestDescriptionsAndDefaultsDownChain(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- animal stuff
# @default -- some
animals:
  # -- animals listed by their various characteristics
  # @default -- explicit
  byTrait:
    # -- the friendly animals of the house
    # @default -- default
    friendly:
      # -- best cat ever
      # @default -- value
      - foxtrot
      - oscar
    mean: [echo]
    sleepy: [oscar]
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 4)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "animal stuff", valuesRows[0].Description)

	assert.Equal(t, "animals.byTrait", valuesRows[1].Key)
	assert.Equal(t, objectType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "animals listed by their various characteristics", valuesRows[1].Description)

	assert.Equal(t, "animals.byTrait.friendly", valuesRows[2].Key)
	assert.Equal(t, listType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "the friendly animals of the house", valuesRows[2].Description)

	assert.Equal(t, "animals.byTrait.friendly[0]", valuesRows[3].Key)
	assert.Equal(t, stringType, valuesRows[3].Type)
	assert.Equal(t, "value", valuesRows[3].Default)
	assert.Equal(t, "best cat ever", valuesRows[3].Description)
}

func TestNilValues(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- (list) the list of birds we have
  birds:
  # -- (int) the number of birds we have
  birdCount:
  # -- the cats that we have that are not weird
  nonWeirdCats:
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.birdCount", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type)
	assert.Equal(t, "`nil`", valuesRows[0].Default)
	assert.Equal(t, "the number of birds we have", valuesRows[0].Description)

	assert.Equal(t, "animals.birds", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, "`nil`", valuesRows[1].Default)
	assert.Equal(t, "the list of birds we have", valuesRows[1].Description)

	assert.Equal(t, "animals.nonWeirdCats", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "`nil`", valuesRows[2].Default)
	assert.Equal(t, "the cats that we have that are not weird", valuesRows[2].Description)
}

func TestNilValuesWithDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- (list) the list of birds we have
  # @default -- explicit
  birds:
  # -- (int) the number of birds we have
  # @default -- some
  birdCount:
  # -- the cats that we have that are not weird
  # @default -- default
  nonWeirdCats:
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.birdCount", valuesRows[0].Key)
	assert.Equal(t, intType, valuesRows[0].Type)
	assert.Equal(t, "some", valuesRows[0].Default)
	assert.Equal(t, "the number of birds we have", valuesRows[0].Description)

	assert.Equal(t, "animals.birds", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, "explicit", valuesRows[1].Default)
	assert.Equal(t, "the list of birds we have", valuesRows[1].Description)

	assert.Equal(t, "animals.nonWeirdCats", valuesRows[2].Key)
	assert.Equal(t, stringType, valuesRows[2].Type)
	assert.Equal(t, "default", valuesRows[2].Default)
	assert.Equal(t, "the cats that we have that are not weird", valuesRows[2].Description)
}

func TestKeysWithSpecialCharacters(t *testing.T) {
	yamlValues := parseYamlValues(`
websites:
  stupidchess.jmn23.com: defunct
fullNames:
  John Norwood: me
`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"me\"`", valuesRows[0].Default)
	assert.Equal(t, "", valuesRows[0].Description)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"defunct\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}

func TestKeysWithSpecialCharactersWithDescriptions(t *testing.T) {
	yamlValues := parseYamlValues(`
websites:
  # -- status of the stupidchess website
  stupidchess.jmn23.com: defunct
fullNames:
  # -- who am I
  John Norwood: me
`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`\"me\"`", valuesRows[0].Default)
	assert.Equal(t, "who am I", valuesRows[0].Description)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"defunct\"`", valuesRows[1].Default)
	assert.Equal(t, "status of the stupidchess website", valuesRows[1].Description)
}

func TestKeysWithSpecialCharactersWithDescriptionsAndDefaults(t *testing.T) {
	yamlValues := parseYamlValues(`
websites:
  # -- status of the stupidchess website
  # @default -- value
  stupidchess.jmn23.com: defunct
fullNames:
  # -- who am I
  # @default -- default
  John Norwood: me
`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, `fullNames."John Norwood"`, valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "default", valuesRows[0].Default)
	assert.Equal(t, "who am I", valuesRows[0].Description)

	assert.Equal(t, `websites."stupidchess.jmn23.com"`, valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "value", valuesRows[1].Default)
	assert.Equal(t, "status of the stupidchess website", valuesRows[1].Description)
}

func TestRequiredSymbols(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- is she friendly?
foxtrot: true

# doesn't show up
hello: "world"
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 2)

	assert.Equal(t, "foxtrot", valuesRows[0].Key)
	assert.Equal(t, boolType, valuesRows[0].Type)
	assert.Equal(t, "`true`", valuesRows[0].Default)
	assert.Equal(t, "is she friendly?", valuesRows[0].Description)

	assert.Equal(t, "hello", valuesRows[1].Key)
	assert.Equal(t, stringType, valuesRows[1].Type)
	assert.Equal(t, "`\"world\"`", valuesRows[1].Default)
	assert.Equal(t, "", valuesRows[1].Description)
}


func TestMultilineDescription(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- The best kind of animal probably, allow me to list their many varied benefits.
  # Cats are very funny, and quite friendly, in almost all cases
  # @default -- The list of cats that _I_ own
  cats:
      - echo
      - foxtrot
`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.cats", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "The list of cats that _I_ own", valuesRows[0].Default)
	assert.Equal(t, "The best kind of animal probably, allow me to list their many varied benefits. Cats are very funny, and quite friendly, in almost all cases", valuesRows[0].Description)
}

func TestMultilineDescriptionWithoutValue(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- (list) I mean, dogs are quite nice too...
  # @default -- The list of dogs that _I_ own
  dogs:
`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 1)

	assert.Equal(t, "animals.dogs", valuesRows[0].Key)
	assert.Equal(t, listType, valuesRows[0].Type)
	assert.Equal(t, "The list of dogs that _I_ own", valuesRows[0].Default)
	assert.Equal(t, "I mean, dogs are quite nice too...", valuesRows[0].Description)
}

func TestInferredTyping(t *testing.T) {
	yamlValues := parseYamlValues(`
# -- pets?
animals:
  # -- multiple cats?
  cats: 3.14159
  # -- pugs are their own species
  pugs: "Frank"
  # -- we have more?
  other: ['gerbil']
  # -- there are
  dogs: true
  # -- they keep eating each other
  fish: 24
  # -- haven't checked
  porcupines:
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 7)

	assert.Equal(t, "animals", valuesRows[0].Key)
	assert.Equal(t, objectType, valuesRows[0].Type)
	assert.Equal(t, "pets?", valuesRows[0].Description)

	assert.Equal(t, "animals.cats", valuesRows[1].Key)
	assert.Equal(t, floatType, valuesRows[1].Type)
	assert.NotEmpty(t, valuesRows[1].Description)

	assert.Equal(t, "animals.dogs", valuesRows[2].Key)
	assert.Equal(t, boolType, valuesRows[2].Type)
	assert.NotEmpty(t, valuesRows[2].Description)

	assert.Equal(t, "animals.fish", valuesRows[3].Key)
	assert.Equal(t, intType, valuesRows[3].Type)
	assert.NotEmpty(t, valuesRows[3].Description)

	assert.Equal(t, "animals.other", valuesRows[4].Key)
	assert.Equal(t, listType, valuesRows[4].Type)
	assert.NotEmpty(t, valuesRows[4].Description)

	assert.Equal(t, "animals.porcupines", valuesRows[5].Key)
	assert.Equal(t, stringType, valuesRows[5].Type)
	assert.NotEmpty(t, valuesRows[5].Description)

	assert.Equal(t, "animals.pugs", valuesRows[6].Key)
	assert.Equal(t, stringType, valuesRows[6].Type)
	assert.NotEmpty(t, valuesRows[6].Description)
}

func TestExplicitTyping(t *testing.T) {
	yamlValues := parseYamlValues(`
animals:
  # -- no type for cats
  cats:
  # -- (list) dogs should be a list
  # @default -- nil
  dogs:
  # -- (Mermen) can haz Mermen?
  fish:
	`)

	valuesRows, err := getSortedValuesTableRows(yamlValues)

	assert.Nil(t, err)
	assert.Len(t, valuesRows, 3)

	assert.Equal(t, "animals.cats", valuesRows[0].Key)
	assert.Equal(t, stringType, valuesRows[0].Type)
	assert.Equal(t, "`nil`", valuesRows[0].Default)
	assert.Equal(t, "no type for cats", valuesRows[0].Description)

	assert.Equal(t, "animals.dogs", valuesRows[1].Key)
	assert.Equal(t, listType, valuesRows[1].Type)
	assert.Equal(t, "nil", valuesRows[1].Default)
	assert.Equal(t, "dogs should be a list", valuesRows[1].Description)

	assert.Equal(t, "animals.fish", valuesRows[2].Key)
	assert.Equal(t, "Mermen", valuesRows[2].Type)
	assert.Equal(t, "`nil`", valuesRows[2].Default)
	assert.Equal(t, "can haz Mermen?", valuesRows[2].Description)
}
