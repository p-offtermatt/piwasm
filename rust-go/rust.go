package rust

import (
	"fmt"
	"strings"
)

type AST interface {
	PrettyPrint() string
}

type Type interface {
	AST
}

type TypeCons struct {
	Type

	Name   string
	Params []Type
}

func (t *TypeCons) PrettyPrint() string {
	types := make([]string, len(t.Params))
	for i, typ := range t.Params {
		types[i] = typ.PrettyPrint()
	}
	return fmt.Sprintf("%s<%s>", t.Name, strings.Join(types, ", "))
}

type UInt64Type struct {
	Type
}

func (t *UInt64Type) PrettyPrint() string {
	return "u64"
}

type StringType struct {
	Type
}

func (t *StringType) PrettyPrint() string {
	return "String"
}

type BoolType struct {
	Type
}

func (t *BoolType) PrettyPrint() string {
	return "bool"
}

type TupleType struct {
	Type

	Types []Type
}

func (t *TupleType) PrettyPrint() string {
	types := make([]string, len(t.Types))
	for i, typ := range t.Types {
		types[i] = typ.PrettyPrint()
	}
	return fmt.Sprintf("(%s)", strings.Join(types, ", "))
}

type StrType struct {
	Type
}

func (t *StrType) PrettyPrint() string {
	return "str"
}

type TypeRef struct {
	Type

	OfType  Type
	Mutable bool
}

func (t *TypeRef) PrettyPrint() string {
	mut := ""
	if t.Mutable {
		mut = "mut "
	}
	return fmt.Sprintf("&%s %s", mut, t.OfType.PrettyPrint())
}

type Import struct {
	AST

	Path string
}

func (i Import) PrettyPrint() string {
	return fmt.Sprintf("use %s;", i.Path)
}

type Program struct {
	AST

	Imports   []Import
	Structs   []StructDecl
	Functions []FunctionDecl
}

func (p *Program) PrettyPrint() string {
	var sb strings.Builder

	for _, imp := range p.Imports {
		sb.WriteString(imp.PrettyPrint())
		sb.WriteString("\n")
	}

	for _, decl := range p.Structs {
		sb.WriteString(decl.PrettyPrint())
		sb.WriteString("\n")
	}

	for _, decl := range p.Functions {
		sb.WriteString(decl.PrettyPrint())
		sb.WriteString("\n")
	}

	return sb.String()
}

type Decl interface {
	AST
}

type StructDecl struct {
	Decl

	Name   string
	Fields []Field
}

func (s *StructDecl) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString("struct ")
	sb.WriteString(s.Name)
	sb.WriteString(" {\n")

	for _, field := range s.Fields {
		sb.WriteString("    ")
		sb.WriteString(field.PrettyPrint())
		sb.WriteString(",\n")
	}

	sb.WriteString("}")

	return sb.String()
}

type Field struct {
	AST

	Name string
	Type string
}

func (f *Field) PrettyPrint() string {
	return fmt.Sprintf("%s: %s", f.Name, f.Type)
}

type FunctionDecl struct {
	Decl

	Name       string
	Params     []Param
	ReturnType string
	Body       []Stmt
}

func (f *FunctionDecl) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString("fn ")
	sb.WriteString(f.Name)
	sb.WriteString("(")

	params := make([]string, len(f.Params))
	for i, param := range f.Params {
		params[i] = param.PrettyPrint()
	}
	sb.WriteString(strings.Join(params, ", "))

	sb.WriteString(") -> ")
	sb.WriteString(f.ReturnType)
	sb.WriteString(" {\n")

	for _, stmt := range f.Body {
		sb.WriteString("    ")
		sb.WriteString(stmt.PrettyPrint())
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}

type Param struct {
	AST

	Name    string
	Type    string
	Mutable bool
}

func (p *Param) PrettyPrint() string {
	mut := ""
	if p.Mutable {
		mut = "mut "
	}
	return fmt.Sprintf("%s%s: %s", mut, p.Name, p.Type)
}

type Stmt interface {
	AST
}

type Let struct {
	Stmt

	VariableName string
	Value        Expr
}

func (l *Let) PrettyPrint() string {
	return fmt.Sprintf("let %s = %s;", l.VariableName, l.Value.PrettyPrint())
}

type Assign struct {
	Stmt

	Dest  Expr
	Value Expr
}

func (a *Assign) PrettyPrint() string {
	return fmt.Sprintf("%s = %s;", a.Dest.PrettyPrint(), a.Value.PrettyPrint())
}

type Return struct {
	Stmt

	Value Expr
}

func (r *Return) PrettyPrint() string {
	return fmt.Sprintf("return %s;", r.Value.PrettyPrint())
}

type Expr interface {
	AST
}

type Block struct {
	Expr

	Statements []Stmt
}

func (b *Block) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString("{\n")
	for _, stmt := range b.Statements {
		sb.WriteString("    ")
		sb.WriteString(stmt.PrettyPrint())
		sb.WriteString("\n")
	}
	sb.WriteString("}")

	return sb.String()
}

type FieldValue struct {
	AST

	Name  string
	Value Expr
}

func (f *FieldValue) PrettyPrint() string {
	return fmt.Sprintf("%s: %s", f.Name, f.Value.PrettyPrint())
}

type StructCons struct {
	Expr

	StructName string
	Fields     []FieldValue
}

func (s *StructCons) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString(s.StructName)
	sb.WriteString(" {\n")

	for _, field := range s.Fields {
		sb.WriteString("    ")
		sb.WriteString(field.PrettyPrint())
		sb.WriteString(",\n")
	}

	sb.WriteString("}")

	return sb.String()
}

type EnumCons struct {
	Expr

	EnumName string
	Variant  string
	Params   []Expr
}

func (e *EnumCons) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString(e.EnumName)
	sb.WriteString("::")
	sb.WriteString(e.Variant)
	sb.WriteString("(")

	params := make([]string, len(e.Params))
	for i, param := range e.Params {
		params[i] = param.PrettyPrint()
	}
	sb.WriteString(strings.Join(params, ", "))

	sb.WriteString(")")

	return sb.String()
}

type Borrow struct {
	Expr

	Value Expr
}

func (b *Borrow) PrettyPrint() string {
	return fmt.Sprintf("&%s", b.Value.PrettyPrint())
}

type Tuple struct {
	Expr

	Values []Expr
}

func (t *Tuple) PrettyPrint() string {
	values := make([]string, len(t.Values))
	for i, value := range t.Values {
		values[i] = value.PrettyPrint()
	}
	return fmt.Sprintf("(%s)", strings.Join(values, ", "))
}

type FunctionCall struct {
	Expr

	FunctionName string
	Arguments    []Expr
}

func (f *FunctionCall) PrettyPrint() string {
	args := make([]string, len(f.Arguments))
	for i, arg := range f.Arguments {
		args[i] = arg.PrettyPrint()
	}
	return fmt.Sprintf("%s(%s)", f.FunctionName, strings.Join(args, ", "))
}

type StaticMethodCall struct {
	Expr

	TypeName   Type
	MethodName string
	Arguments  []Expr
}

func (s *StaticMethodCall) PrettyPrint() string {
	args := make([]string, len(s.Arguments))
	for i, arg := range s.Arguments {
		args[i] = arg.PrettyPrint()
	}
	return fmt.Sprintf("%s::%s(%s)", s.TypeName.PrettyPrint(), s.MethodName, strings.Join(args, ", "))
}

type MethodCall struct {
	Expr

	Value      Expr
	MethodName string
	Arguments  []Expr
}

func (m *MethodCall) PrettyPrint() string {
	args := make([]string, len(m.Arguments))
	for i, arg := range m.Arguments {
		args[i] = arg.PrettyPrint()
	}
	return fmt.Sprintf("%s.%s(%s)", m.Value.PrettyPrint(), m.MethodName, strings.Join(args, ", "))
}

type Variable struct {
	Expr

	VariableName string
}

func (v *Variable) PrettyPrint() string {
	return v.VariableName
}

type FieldAccess struct {
	Expr

	Value Expr
	Field string
}

func (f *FieldAccess) PrettyPrint() string {
	return fmt.Sprintf("%s.%s", f.Value.PrettyPrint(), f.Field)
}

type IfElse struct {
	Expr

	Condition Expr
	Then      Expr
	Else      Expr
}

func (i *IfElse) PrettyPrint() string {
	var sb strings.Builder
	sb.WriteString("if ")
	sb.WriteString(i.Condition.PrettyPrint())
	sb.WriteString(" {\n")
	sb.WriteString("    ")
	sb.WriteString(i.Then.PrettyPrint())
	sb.WriteString("\n} else {\n")
	sb.WriteString("    ")
	sb.WriteString(i.Else.PrettyPrint())
	sb.WriteString("\n}")
	return sb.String()
}

type Not struct {
	Expr

	Value Expr
}

func (n *Not) PrettyPrint() string {
	return fmt.Sprintf("!%s", n.Value.PrettyPrint())
}

type Add struct {
	Expr

	Left  Expr
	Right Expr
}

func (a *Add) PrettyPrint() string {
	return fmt.Sprintf("%s + %s", a.Left.PrettyPrint(), a.Right.PrettyPrint())
}

type Literal interface {
	Expr
}

type UInt64Literal struct {
	Literal

	Value uint64
}

func (u *UInt64Literal) PrettyPrint() string {
	return fmt.Sprintf("%d_u64", u.Value)
}

type StringLiteral struct {
	Literal

	Value string
}

func (s *StringLiteral) PrettyPrint() string {
	return fmt.Sprintf("\"%s\"", strings.ReplaceAll(s.Value, "\"", "\\\""))
}

type BoolLiteral struct {
	Literal

	Value bool
}

func (b *BoolLiteral) PrettyPrint() string {
	if b.Value {
		return "true"
	}
	return "false"
}
