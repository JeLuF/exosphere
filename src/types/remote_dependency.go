package types

import (
	"fmt"

	"github.com/Originate/exosphere/src/assets"
	yaml "gopkg.in/yaml.v2"
)

// RemoteDependency represents a production dependency
type RemoteDependency struct {
	Type           string
	TemplateConfig map[string]string `yaml:"template-config,omitempty"`
}

// remoteDependencyRequirements represents the requirements.yml file included with each remote dependency
type remoteDependencyRequirements struct {
	RequriedFields []string `yaml:"required-fields,omitempty"`
}

// ValidateFields validates that a remote config contains all required fields
func (p *RemoteDependency) ValidateFields() error {
	data, err := assets.Asset(fmt.Sprintf("remote-dependency-templates/%s/requirements.yml", p.Type))
	if err != nil {
		return err
	}
	var requirements remoteDependencyRequirements
	err = yaml.Unmarshal(data, &requirements)
	if err != nil {
		return fmt.Errorf("Failed unmarshal '%s/requirements.yml' file", p.Type)
	}
	for _, requiredField := range requirements.RequriedFields {
		if _, ok := p.TemplateConfig[requiredField]; !ok {
			return fmt.Errorf("remote dependency of type '%s' missing required field 'template-config.%s'", p.Type, requiredField)
		}
	}
	return nil
}
