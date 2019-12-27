package types_test

import (
	"testing"

	"github.com/dapperlabs/flow-go/language/runtime/ast"
	"github.com/dapperlabs/flow-go/language/runtime/common"
	"github.com/dapperlabs/flow-go/language/runtime/sema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk/language/types"
)

func TestConvert(t *testing.T) {

	t.Run("structs", func(t *testing.T) {
		position := ast.Position{
			Offset: 1, Line: 2, Column: 3,
		}
		identifier := "my_structure"

		ty := &sema.CompositeType{
			Location:     nil,
			Identifier:   identifier,
			Kind:         common.CompositeKindStructure,
			Conformances: nil,
			Members: map[string]*sema.Member{
				"fieldA": {
					ContainerType: nil,
					Access:        0,
					Identifier: ast.Identifier{
						Identifier: "fieldA",
						Pos:        position,
					},
					TypeAnnotation:  &sema.TypeAnnotation{Type: &sema.IntType{}},
					DeclarationKind: 0,
					VariableKind:    ast.VariableKindVariable,
					ArgumentLabels:  nil,
				},
			},
			ConstructorParameterTypeAnnotations: []*sema.TypeAnnotation{
				{
					IsResource: false,
					Type:       &sema.Int8Type{},
				},
			},
		}

		program := &ast.Program{
			Declarations: []ast.Declaration{
				&ast.CompositeDeclaration{
					Identifier: ast.Identifier{
						Identifier: identifier, Pos: position,
					},
					Members: &ast.Members{
						SpecialFunctions: []*ast.SpecialFunctionDeclaration{
							{
								DeclarationKind: common.DeclarationKindInitializer,
								FunctionDeclaration: &ast.FunctionDeclaration{
									Identifier: ast.Identifier{},
									ParameterList: &ast.ParameterList{
										Parameters: []*ast.Parameter{
											{
												Label: "labelA",
												Identifier: ast.Identifier{
													Identifier: "fieldA",
													Pos:        ast.Position{},
												},
											},
										},
									},
								},
							},
						},
					},
					Range: ast.Range{},
				},
			},
		}

		variable := &sema.Variable{
			Identifier:      identifier,
			DeclarationKind: common.DeclarationKindStructure,
			Pos:             &position,
		}

		ex, err := types.Convert(ty, program, variable)
		assert.NoError(t, err)

		assert.IsType(t, types.Struct{}, ex)
		s := ex.(types.Struct)

		assert.Equal(t, identifier, s.Identifier)
		require.Len(t, s.Fields, 1)

		assert.Equal(t, "fieldA", s.Fields[0].Identifier)
		assert.IsType(t, types.Int{}, s.Fields[0].Type)
	})

	t.Run("string", func(t *testing.T) {
		ty := &sema.StringType{}

		ex, err := types.Convert(ty, nil, nil)
		assert.NoError(t, err)

		assert.IsType(t, types.String{}, ex)
	})

	t.Run("events", func(t *testing.T) {
		position := ast.Position{
			Offset: 2, Line: 1, Column: 37,
		}

		ty := &sema.EventType{
			Identifier: "MagicEvent",
			Fields: []sema.EventFieldType{
				{
					Identifier: "who",
					Type:       &sema.StringType{},
				},
				{
					Identifier: "where",
					Type:       &sema.IntType{},
				},
			},
		}

		program := &ast.Program{
			Declarations: []ast.Declaration{
				&ast.EventDeclaration{
					Identifier: ast.Identifier{
						Identifier: "MagicEvent",
						Pos:        position,
					},
					ParameterList: &ast.ParameterList{
						Parameters: []*ast.Parameter{
							{
								Label: "magic_caster",
								Identifier: ast.Identifier{
									Identifier: "who",
								},
							},
							{
								Label: "magic_place",
								Identifier: ast.Identifier{
									Identifier: "where",
								},
							},
						},
					},
				},
			},
		}

		variable := &sema.Variable{
			Identifier: "MagicEvent",
			Pos:        &position,
		}

		ex, err := types.Convert(ty, program, variable)
		assert.NoError(t, err)

		assert.IsType(t, types.Event{}, ex)

		event := ex.(types.Event)

		require.Len(t, event.Fields, 2)
		assert.Equal(t, "who", event.Fields[0].Identifier)
		assert.IsType(t, types.String{}, event.Fields[0].Type)

		assert.Equal(t, "where", event.Fields[1].Identifier)
		assert.IsType(t, types.Int{}, event.Fields[1].Type)

		require.Len(t, event.Initializer, 2)
		assert.Equal(t, "magic_caster", event.Initializer[0].Label)
		assert.Equal(t, "magic_place", event.Initializer[1].Label)
	})
}
