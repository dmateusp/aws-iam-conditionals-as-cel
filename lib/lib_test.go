package lib_test

import (
	"fmt"
	"log"
	"reflect"

	"github.com/dmateusp/aws-iam-conditionals-as-cel/lib"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/google/uuid"
)

func ExampleUser_getFields() {
	env, err := cel.NewEnv(
		cel.Variable("user", cel.ObjectType("lib.User")),
		ext.NativeTypes(reflect.TypeOf(lib.User{})), // Register the struct in CEL
	)

	if err != nil {
		log.Fatalf("could not instantiate CEL environment: %v", err)
	}

	user := lib.User{
		Name: "johndoe",
		Tags: []lib.Tag{
			{Key: "foo"},
			{Key: "bar"},
		},
	}

	for _, expr := range []string{
		`user.Name`,
		`user.Tags`,
		`user`,
	} {
		ast, issues := env.Compile(expr)
		if issues != nil && issues.Err() != nil {
			log.Fatalf("type-check error: %s", issues.Err())
		}

		prg, err := env.Program(ast)
		if err != nil {
			log.Fatalf("program construction error: %s", err)
		}

		// We evaluate the program with our input
		result, _, err := prg.Eval(map[string]interface{}{
			"user": user,
		})

		if err != nil {
			log.Fatalf("failed to evaluate program: %v", err)
		}

		fmt.Printf("%+v of type %T\n", result.Value(), result.Value())
	}

	// Output:
	// johndoe of type string
	// [{Key:foo} {Key:bar}] of type []lib.Tag
	// {Id:00000000-0000-0000-0000-000000000000 Name:johndoe Tags:[{Key:foo} {Key:bar}] multiFactorAuthInformationStore:<nil>} of type lib.User
}

func ExampleUser_userfulError() {
	env, err := cel.NewEnv(
		cel.Variable("user", cel.ObjectType("lib.User")),
		ext.NativeTypes(reflect.TypeOf(lib.User{})), // Register the struct in CEL
	)

	if err != nil {
		log.Fatalf("could not instantiate CEL environment: %v", err)
	}

	_, issues := env.Compile(`user.Name > 0`)
	if issues != nil && issues.Err() != nil {
		fmt.Println(issues.String())
	}

	// Output:
	// ERROR: <input>:1:11: found no matching overload for '_>_' applied to '(string, int)'
	//  | user.Name > 0
	//  | ..........^
}

func ExampleUser_multiFactorAuthAge() {
	env, err := cel.NewEnv(
		cel.Variable("user", cel.ObjectType("lib.User")),
		ext.NativeTypes(reflect.TypeOf(lib.User{})), // Register the struct in CEL
		lib.MultiFactorAuthAgeMethodBinding,
	)

	if err != nil {
		log.Fatalf("could not instantiate CEL environment: %v", err)
	}

	userId := uuid.New()
	user := lib.NewUser(
		userId,
		"johndoe",
		[]lib.Tag{
			{Key: "foo"},
			{Key: "bar"},
		},
		lib.SimulateMultiFactorAuthInformationStore{},
	)

	ast, issues := env.Compile(`user.multiFactorAuthAge()`)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}

	// We evaluate the program with our input
	result, _, err := prg.Eval(map[string]interface{}{
		"user": user,
	})

	if err != nil {
		log.Fatalf("failed to evaluate program: %v", err)
	}

	fmt.Printf("%+v of type %T\n", result.Value(), result.Value())

	// Output:
	// 1000 of type uint64
}

func ExampleUser_final() {
	env, err := cel.NewEnv(
		cel.Variable("user", cel.ObjectType("lib.User")),
		ext.NativeTypes(reflect.TypeOf(lib.User{})), // Register the struct in CEL
		lib.MultiFactorAuthAgeMethodBinding,         // Register the function binding
		ext.Strings(),
	)

	if err != nil {
		log.Fatalf("could not instantiate CEL environment: %v", err)
	}

	userId := uuid.New()
	user := lib.NewUser(
		userId,
		"johndoe",
		[]lib.Tag{
			{Key: "foo"},
			{Key: "bar"},
		},
		lib.SimulateMultiFactorAuthInformationStore{}, // Instantiate our store simulation
	)

	ast, issues := env.Compile(`
	user.Name.lowerAscii() == "johndoe" &&
  		user.multiFactorAuthAge() <= uint(3600) &&
  		user.Tags.exists(tag, tag.Key == "foo" || tag.Key == "bar")
	`)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("type-check error: %s", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("program construction error: %s", err)
	}

	// We evaluate the program with our input
	result, _, err := prg.Eval(map[string]interface{}{
		"user": user,
	})

	if err != nil {
		log.Fatalf("failed to evaluate program: %v", err)
	}

	fmt.Printf("%+v of type %T\n", result.Value(), result.Value())

	// Output:
	// true of type bool
}
