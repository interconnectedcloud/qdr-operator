package test

import (
	"github.com/RHsyseng/operator-utils/pkg/validation"
	"github.com/ghodss/yaml"
	qdrv1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

var crdTypeMap = map[string]interface{}{
	"interconnectedcloud_v1alpha1_qdr_crd.yaml": &qdrv1alpha1.Qdr{},
}

func TestCRDSchemas(t *testing.T) {
	for crdFileName, amqType := range crdTypeMap {
		schema := getSchema(t, crdFileName)
		missingEntries := schema.GetMissingEntries(amqType)
		for _, missing := range missingEntries {
			if strings.HasPrefix(missing.Path, "/status/conditions/transitionTime/") {
				//skill detailed properties of transition Time.
			} else {
				assert.Fail(t, "Discrepancy between CRD and Struct",
					"Missing or incorrect schema validation at %v, expected type %v  in CRD file %v", missing.Path, missing.Type, crdFileName)
			}
		}
	}
}

func TestSampleCustomResources(t *testing.T) {

	var crFileName, crdFileName string = "interconnectedcloud_v1alpha1_qdr_cr.yaml", "interconnectedcloud_v1alpha1_qdr_crd.yaml"
	assert.NotEmpty(t, crdFileName, "No matching CRD file found for CR suffixed: %s", crFileName)

	schema := getSchema(t, crdFileName)
	yamlString, err := ioutil.ReadFile("../deploy/crds/" + crFileName)
	assert.NoError(t, err, "Error reading %v CR yaml", crFileName)
	var input map[string]interface{}
	assert.NoError(t, yaml.Unmarshal([]byte(yamlString), &input))
	assert.NoError(t, schema.Validate(input), "File %v does not validate against the CRD schema", crFileName)
}

func getSchema(t *testing.T, crdFile string) validation.Schema {

	yamlString, err := ioutil.ReadFile("../deploy/crds/" + crdFile)
	assert.NoError(t, err, "Error reading CRD yaml %v", yamlString)

	schema, err := validation.New([]byte(yamlString))
	assert.NoError(t, err)

	return schema
}
