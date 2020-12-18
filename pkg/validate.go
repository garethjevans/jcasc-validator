package pkg

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	ghodss "github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
)

// ValidateCmd defines parent get.
type ValidateCmd struct {
	Cmd              *cobra.Command
	Args             []string
	JenkinsLocation  string
	SchemaLocation   string
	SoftFail bool
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

	cmd.Flags().BoolVarP(&c.SoftFail, "soft-fail", "", false,
		"Only display errors, do not fail")

	return cmd
}

// Validate.
func (c *ValidateCmd) Validate() error {
	if c.TemplateLocation == "" {
		return errors.New("--template-location must be set")
	}

	if c.SchemaLocation == "" && c.JenkinsLocation == "" {
		return errors.New("either --schema-location or --jenkins-location must be set")
	}

	return nil
}

// Run.
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

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewReferenceLoader(schemaURL(c.SchemaLocation, c.JenkinsLocation))
	logrus.Debugf("using schema %s", schemaLoader)

	filesToProcess, err := writeConfigMapsToTempDir(yamlFile, tmpDir)
	if err != nil {
		return err
	}

	valid := true
	for _, f := range filesToProcess {
		jsonFile := path.Join(tmpDir, f+".json")
		yamlFile := path.Join(tmpDir, f)

		documentLoader := gojsonschema.NewReferenceLoader("file://" + jsonFile)

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return err
		}

		if !result.Valid() {
			valid = false
			logrus.Errorf("The file %s is not valid. see errors :", f)
			for _, desc := range result.Errors() {
				logrus.Warnf(" - %s", desc.String())
			}
			yamlContents, err := ioutil.ReadFile(yamlFile)
			if err != nil {
				return err
			}
			logrus.Infof("\n%s", string(yamlContents))
		}
	}

	if !valid {
		logrus.Errorf("file %s is invalid", c.TemplateLocation)
		if !c.SoftFail {
			os.Exit(1)
		}
	}
	return nil
}

func writeConfigMapsToTempDir(yamlFile []byte, tmpDir string) ([]string, error) {
	configMaps := v1.ConfigMap{}
	dec := yaml.NewDecoder(strings.NewReader(string(yamlFile)))
	filesToProcess := []string{}
	for dec.Decode(&configMaps) == nil {
		for k, v := range configMaps.Data {
			if !contains(filesToProcess, k) {
				yamlFile := path.Join(tmpDir, k)
				err := ioutil.WriteFile(yamlFile, []byte(v), 0600)
				if err != nil {
					return nil, err
				}

				json, err := ghodss.YAMLToJSON([]byte(v))
				if err != nil {
					return nil, err
				}

				jsonFile := yamlFile + ".json"
				err = ioutil.WriteFile(jsonFile, json, 0600)
				if err != nil {
					return nil, err
				}
				filesToProcess = append(filesToProcess, k)
			}
		}
	}
	return filesToProcess, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func schemaURL(schemaLocation string, jenkinsLocation string) string {
	if schemaLocation != "" {
		return "file://" + schemaLocation
	}

	return jenkinsLocation + "/configuration-as-code/schema"
}
