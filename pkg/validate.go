package pkg

import (
	"errors"
	ghodss "github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"os"
	"path"
	"strings"
)

// ValidateCmd defines parent get.
type ValidateCmd struct {
	Cmd  *cobra.Command
	Args []string
	JenkinsLocation string
	SchemaLocation string
	TemplateLocation string
}

// NewValidateCmd creates get cmd.
func NewValidateCmd() *cobra.Command {
	c := &ValidateCmd{}
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "",
		Long:    "",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				logrus.WithError(err).Fatal("unable to run command")
			}
		},
	}

	cmd.Flags().StringVarP(&c.JenkinsLocation, "jenkins-location", "j", "",
		"Location of the Jenkins server")

	cmd.Flags().StringVarP(&c.SchemaLocation, "schema-location", "s", "",
		"Path to the schema")

	cmd.Flags().StringVarP(&c.TemplateLocation, "template-location", "t", "",
		"Path to the generated jcasc-config.yaml")

	return cmd
}

// Validate
func (c *ValidateCmd) Validate() error {
	if c.TemplateLocation == "" {
		return errors.New("--template-location must be set")
	}

	if c.SchemaLocation == "" && c.JenkinsLocation == "" {
		return errors.New("either --schema-location or --jenkins-location must be set")
	}

return nil
}

// Run
func (c *ValidateCmd) Run() error {
	err := c.Validate()
	if err != nil {
		return err
	}

	file := c.TemplateLocation

	tmpDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		return err
	}
	defer os.Remove(tmpDir)

	logrus.Debugf("created temp dir %s", tmpDir)

	configMaps := v1.ConfigMap{}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewReferenceLoader(schemaUrl(c.SchemaLocation, c.JenkinsLocation))
	logrus.Debugf("using schema %s", schemaLoader)

	filesToProcess := []string{}

	dec := yaml.NewDecoder(strings.NewReader(string(yamlFile)))
	for dec.Decode(&configMaps) == nil {
		for k, v := range configMaps.Data {
			if !contains(filesToProcess, k) {
				yamlFile := path.Join(tmpDir, k)
				err = ioutil.WriteFile(yamlFile, []byte(v), 0644)
				if err != nil {
					return err
				}

				json, err := ghodss.YAMLToJSON([]byte(v))
				if err != nil {
					return err
				}

				jsonFile := yamlFile + ".json"
				err = ioutil.WriteFile(jsonFile, json, 0644)
				if err != nil {
					return err
				}
				filesToProcess = append(filesToProcess, k)
			}
		}
	}

	for _, f := range filesToProcess {
		jsonFile := path.Join(tmpDir, f + ".json")
		yamlFile := path.Join(tmpDir, f)

		documentLoader := gojsonschema.NewReferenceLoader("file://" + jsonFile)

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return err
		}

		if !result.Valid() {
			logrus.Errorf("The file %s is not valid. see errors :", f)
			for _, desc := range result.Errors() {
				logrus.Warnf("- %s", desc)
			}
			yamlContents, err := ioutil.ReadFile(yamlFile)
			if err != nil {
				return err
			}
			logrus.Infof("\n%s", string(yamlContents))
		}
	}
return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func schemaUrl(schemaLocation string, jenkinsLocation string) string {
	if schemaLocation != "" {
		return "file://" + schemaLocation
	}

	return jenkinsLocation + "/configuration-as-code/schema"
}

