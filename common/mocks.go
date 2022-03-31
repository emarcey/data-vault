package common

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func NewDummySecret(t *testing.T) *Secret {
	tmp := Secret{}
	err := faker.FakeData(&tmp)
	require.Nil(t, err, "Unexpected error generating dummy secret: %v", err)
	return &tmp
}

func NewDummyUser(t *testing.T) *User {
	tmp := User{}
	err := faker.FakeData(&tmp)
	require.Nil(t, err, "Unexpected error generating dummy user: %v", err)
	return &tmp
}

func NewDummyUserGroup(t *testing.T) *UserGroup {
	tmp := UserGroup{}
	err := faker.FakeData(&tmp)
	require.Nil(t, err, "Unexpected error generating dummy user group: %v", err)
	return &tmp
}
