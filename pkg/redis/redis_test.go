package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const addr = "127.0.0.1:6379"

func TestStore(t *testing.T) {
	store := NewStore(&Config{
		Addr:      addr,
		DB:        1,
		KeyPrefix: "prefix",
	})

	defer store.Close()

	ctx := context.Background()
	key := "test"
	err := store.Set(ctx, key, 0, 5)
	assert.Nil(t, err)

	b, err := store.Check(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	err = store.Delete(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)
}
