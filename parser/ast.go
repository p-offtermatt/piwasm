package main

type AST interface {
	PrettyPrint(level int) string
}

// Types
type (
	Type interface {
		AST
	}
	TypeCons struct {
		Type
		Name   string
		Params []Type
	}
)

type (
	StructType struct {
		Type
		Fields []Field
	}

	UInt64Type struct {
		Type
	}
	StringType struct {
		Type
	}
	BoolType struct {
		Type
	}
	TupleType struct {
		Type
		Types []Type
	}
	StrType struct {
		Type
	}
	ListType struct {
		Type
		ElementType Type
	}
	SetType struct {
		Type
		ElementType Type
	}
	// a reference to a custom type that is already defined
	ConstType struct {
		Type
		Name string
	}
	MapType struct {
		Type
		ArgType    Type
		ReturnType Type
	}
	TypeRef struct {
		Type
		OfType  Type
		Mutable bool
	}
)

type Import struct {
	AST
	Path string
}
type Program struct {
	AST

	Imports []Import
	Decls   []Decl
}
type (
	Decl interface {
		AST
	}
	StructDecl struct {
		Decl
		Name   string
		Fields []Field
	}
	FunctionDecl struct {
		Decl
		Name       string
		Params     []Param
		ReturnType Type
		Body       []Stmt
	}

	// declares a global constant
	ConstDecl struct {
		Decl
		Name  string
		Type  Type
		Value Block
	}

	TypeDecl struct {
		Decl
		Name string
		Type Type
	}
)

type Field struct {
	Name string
	Type Type
}
type Param struct {
	Name    string
	Type    Type
	Mutable bool
}

// Statements
type (
	Stmt interface {
		AST
	}
	Let struct {
		Stmt
		VariableName string
		Value        Expr
	}
)

type Assign struct {
	Stmt
	Dest  Expr
	Value Expr
}
type Return struct {
	Stmt
	Value Expr
}

// Expressions
type (
	Expr interface {
		AST
	}
	Block struct {
		Expr
		Statements []Stmt
	}
)

type FieldValue struct {
	AST
	Name  string
	Value Expr
}
type StructCons struct {
	Expr
	StructName string
	Fields     []FieldValue
}
type EnumCons struct {
	Expr
	EnumName string
	Variant  string
	Params   []Expr
}
type Borrow struct {
	Expr
	Value Expr
}
type Tuple struct {
	Expr
	Values []Expr
}
type FunctionCall struct {
	Expr
	FunctionName string
	Arguments    []Expr
}
type StaticMethodCall struct {
	Expr
	TypeName   Type
	MethodName string
	Arguments  []Expr
}
type MethodCall struct {
	Expr
	Value      Expr
	MethodName string
	Arguments  []Expr
}
type Variable struct {
	Expr
	VariableName string
}
type FieldAccess struct {
	Expr
	Value Expr
	Field string
}
type IfElse struct {
	Expr
	Condition Expr
	Then      Expr
	Else      Expr
}
type Not struct {
	Expr
	Value Expr
}
type Add struct {
	Expr
	Left  Expr
	Right Expr
}
type Literal interface {
	Expr
}
type UInt64Literal struct {
	Literal
	Value uint64
}
type StringLiteral struct {
	Literal
	Value string
}
type BoolLiteral struct {
	Literal
	Value bool
}
