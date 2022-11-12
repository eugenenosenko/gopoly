package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestConfig_New(t *testing.T) {
	t.Run("should correctly serialize YAML config", func(t *testing.T) {
		data, err := os.ReadFile("testdata/test_config.yaml")
		require.NoError(t, err)

		var c Config
		err = yaml.Unmarshal(data, &c)
		require.NoError(t, err)

		require.Equal(t, &Config{
			Types: []*TypeInfo{
				{
					Type:         "abc/models.AdvertBase",
					Subtypes:     []string{"abc/models.RentAdvert", "abc/models.SellAdvert"},
					MarkerMethod: "IsAdvert",
					Discriminator: DiscriminatorInfo{
						Field: "type",
						Mapping: map[string]string{
							"SELL": "SellAdvert",
							"RENT": "RentAdvert",
						},
					},
					DecodingStrategy: DecodingStrategyDiscriminator,
					TemplateFile:     "template_owner.txt",
				},
				{
					Type:             "Property",
					DecodingStrategy: DecodingStrategyStrict,
				},
				{
					Type: "Owner",
					Discriminator: DiscriminatorInfo{
						Field: "kind",
						Mapping: map[string]string{
							"INDIVIDUAL": "abc/models.IndividualOwner",
							"AGENCY":     "abc/models.AgencyOwner",
							"DEVELOPER":  "abc/models.DeveloperOwner",
						},
					},
				},
			},
			DecodingStrategy: DecodingStrategyStrict,
			MarkerMethod:     "Is{{ $type.Name }}",
			TemplateFile:     "template.txt",
		}, &c)
	})
}
