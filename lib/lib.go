package lib

import (
	"context"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/uuid"
)

type Tag struct {
	Key string
}

type User struct {
	Id   uuid.UUID
	Name string
	Tags []Tag

	multiFactorAuthInformationStore MultiFactorAuthInformationStore
}

func NewUser(
	id uuid.UUID,
	name string,
	tags []Tag,
	multiFactorAuthInformationStore MultiFactorAuthInformationStore,
) User {
	return User{
		Id:                              id,
		Name:                            name,
		Tags:                            tags,
		multiFactorAuthInformationStore: multiFactorAuthInformationStore,
	}
}

type SimulateMultiFactorAuthInformationStore struct{}

func (SimulateMultiFactorAuthInformationStore) GetMultiFactorAuthAge(ctx context.Context, userId uuid.UUID) (uint64, error) {
	return 1000, nil
}

type MultiFactorAuthInformationStore interface {
	// Returns the multi factor authentication age in seconds
	GetMultiFactorAuthAge(ctx context.Context, userId uuid.UUID) (uint64, error)

	// ... other methods
}

const MultiFactorAuthAgeMethodName = "multiFactorAuthAge"

var MultiFactorAuthAgeMethodBinding = cel.Function(MultiFactorAuthAgeMethodName, cel.MemberOverload(MultiFactorAuthAgeMethodName, []*cel.Type{
	cel.ObjectType("lib.User"), // Type of the argument (since it's a method, the first argument is the struct)
},
	cel.UintType, // Type of the return
	cel.FunctionBinding(func(values ...ref.Val) ref.Val {
		user, ok := values[0].Value().(User)
		if !ok {
			return types.NewErr("Could not convert first argument of %s to User", MultiFactorAuthAgeMethodName)
		}
		ctx := context.Background()

		multiFactorAuthAge, err := user.multiFactorAuthInformationStore.GetMultiFactorAuthAge(ctx, user.Id)
		if err != nil {
			return types.NewErr("could not get the multi factor auth age: %v", err)
		}

		return types.DefaultTypeAdapter.NativeToValue(multiFactorAuthAge)
	}),
))
