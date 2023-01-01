//go:build e2e
// +build e2e

//go:generate bash -c scripts/generate.sh

package e2e

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eugenenosenko/gopoly/tests/e2e/testdata/events"
	"github.com/eugenenosenko/gopoly/tests/e2e/testdata/users"
)

func TestE2E(t *testing.T) {
	t.Run("should correctly unmarshal incoming payload into polymorphic structures", func(t *testing.T) {
		data, err := os.ReadFile("testdata/user_event_01.json")
		require.NoError(t, err)

		event, err := events.UnmarshalUserEventJSON(data)
		require.NoError(t, err)

		require.Equal(t, &events.UserDeletedEvent{
			ID:   "12345",
			Type: "DELETED",
			User: &users.RegularUser{
				ID:      "1234",
				Type:    "REGULAR",
				Name:    "John Doe",
				Address: "Kings Road 12, London, UK",
				Contacts: []users.Contact{
					&users.BusinessContact{
						ID:           "1",
						BusinessName: "Business Ltd.",
						Phone:        "0 800 122 222",
						Email:        "john.doe@gmail.com",
					},
					&users.PrivateContact{
						ID: "2",
						FullName: users.FullName{
							Firstname: "John",
							Lastname:  "Doe",
						},
						Phone: "0 800 122 223",
						Email: "john.doe@gmail.com",
					},
				},
			},
		}, event)
	})
}
