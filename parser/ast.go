package main

type (
	Type     interface{}
	TypeCons struct {
		Type
		Name   string
		Params []Type
	}
)

type StructType struct {
	Type
	Fields []Field
}

type UInt64Type struct {
	Type
}
type StringType struct {
	Type
}
type BoolType struct {
	Type
}
type TupleType struct {
	Type
	Types []Type
}
type StrType struct {
	Type
}
type ListType struct {
	Type
	ElementType Type
}
type SetType struct {
	Type
	ElementType Type
}
type FunType struct {
	Type
	ArgType    Type
	ReturnType Type
}
type TypeRef struct {
	Type
	TypeName string
	Mutable  bool
}
type Import struct {
	Path string
}
type Program struct {
	Imports   []Import
	Structs   []StructDecl
	Functions []FunctionDecl
}
type (
	Decl       interface{}
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
		Body       Block
	}

	// declares a global constant
	ConstDecl struct {
		Decl
		Name  string
		Value Block
	}
)

type TypeDef struct {
	Name string
	Type Type
}

type Field struct {
	Name string
	Type Type
}
type Param struct {
	Name    string
	Type    Type
	Mutable bool
}
type (
	Stmt interface{}
	Let  struct {
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
type (
	Expr  interface{}
	Block struct {
		Expr
		Statements []Stmt
	}
)

type FieldValue struct {
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
