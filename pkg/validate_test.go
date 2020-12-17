package pkg_test

import (
	"github.com/garethjevans/jcasc-validator/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateCmd_Validate_NoArgs(t *testing.T) {
	cmd := pkg.ValidateCmd{}
	assert.Error(t, cmd.Validate(), "No args should return an error")
}

func TestValidateCmd_Validate_Template(t *testing.T) {
	cmd := pkg.ValidateCmd{
		TemplateLocation: "path.yaml",
	}
	assert.Error(t, cmd.Validate(), "Need to specify either Jenkins args or Location")
}

func TestValidateCmd_Validate_SchemaAndTemplate(t *testing.T) {
	cmd := pkg.ValidateCmd{
		TemplateLocation: "path.yaml",
		SchemaLocation: "schema.json",
	}
	assert.NoError(t, cmd.Validate())
}

func TestValidateCmd_Validate_JenkinsArgsAndTemplate(t *testing.T) {
	cmd := pkg.ValidateCmd{
		JenkinsLocation: "http://localhost",
		JenkinsUsername: "user",
		JenkinsPassword: "pass",
		TemplateLocation: "path.yaml",
	}
	assert.NoError(t, cmd.Validate())
}