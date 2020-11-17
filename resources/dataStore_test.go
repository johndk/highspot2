package resources

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddUser(t *testing.T) {
	db := NewDataStore()

	usersJSON := `[
		  {
			"id": "1",
			"name": "Albin Jaye"
		  },
		  {
			"id": "2",
			"name": "Dipika Crescentia"
		  },
		  {
			"id": "3",
			"name": "Ankit Sacnite"
		  }
		]`

	var users []*User
	err := json.Unmarshal([]byte(usersJSON), &users)
	require.NoError(t, err)

	for _, user := range users {
		err = db.AddUser(user)
		require.NoError(t, err)
	}

	i := 0
	db.ForeachUser(func(user *User) error {
		require.EqualValues(t, users[i], user)
		i++
		return nil
	})

	require.Equal(t, 3, i)
}
