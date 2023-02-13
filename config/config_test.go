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
			Types: []*TypeDefinition{
				{
					Name:         "AdvertBase",
					Variants:     []string{"RentAdvert", "SellAdvert"},
					MarkerMethod: "IsAdvert",
					Discriminator: DiscriminatorDefinition{
						Field: "type",
						Mapping: map[string]string{
							"SELL": "SellAdvert",
							"RENT": "RentAdvert",
						},
					},
					DecodingStrategy: DecodingStrategyDiscriminator,
				},
				{
					Name:             "Property",
					DecodingStrategy: DecodingStrategyStrict,
				},
				{
					Name: "Owner",
					Discriminator: DiscriminatorDefinition{
						Field: "kind",
						Mapping: map[string]string{
							"INDIVIDUAL": "IndividualOwner",
							"AGENCY":     "AgencyOwner",
							"DEVELOPER":  "DeveloperOwner",
						},
					},
				},
			},
			DecodingStrategy: DecodingStrategyStrict,
			MarkerMethod:     "Is{{ $type.Name }}",
		}, &c)
	})
}
