package main

import (
	"fmt"
	"strings"
)

func (t *TypeCons) PrettyPrint(level int) string {
	types := make([]string, len(t.Params))
	for i, typ := range t.Params {
		types[i] = typ.PrettyPrint(level)
	}
	return fmt.Sprintf("%s<%s>", t.Name, strings.Join(types, ", "))
}

func (t *StructType) PrettyPrint(level int) string {
	var sb strings.Builder

	sb.WriteString("struct {\n")

	for _, field := range t.Fields {
		sb.WriteString("    ")
		sb.WriteString(field.PrettyPrint(level))
		sb.WriteString("\n")
	}

	sb.WriteString("}")

	return sb.String()
}

func (t *UInt64Type) PrettyPrint(level int) string {
	return "u64"
}

func (t *BoolType) PrettyPrint(level int) string {
	return "bool"
}

func (t *SetType) PrettyPrint(level int) string {
	return "HashSet::<" + t.ElementType.PrettyPrint(level) + ">"
}

func (t *MapType) PrettyPrint(level int) string {
	sb := strings.Builder{}
	sb.WriteString("HashMap::<")
	sb.WriteString(t.ArgType.PrettyPrint(level))
	sb.WriteString(", ")
	sb.WriteString(t.ReturnType.PrettyPrint(level))
	sb.WriteString(">")
	return sb.String()
}

func (t *TupleType) PrettyPrint(level int) string {
	types := make([]string, len(t.Types))
	for i, typ := range t.Types {
		types[i] = typ.PrettyPrint(level)
	}
	return fmt.Sprintf("(%s)", strings.Join(types, ", "))
}

func (t *StrType) PrettyPrint(level int) string {
	return "String"
}

func (t *ListType) PrettyPrint(level int) string {
	return "[" + t.ElementType.PrettyPrint(level) + "]"
}

func (t *ConstType) PrettyPrint(level int) string {
	return t.Name
}

func (t *TypeRef) PrettyPrint(level int) string {
	mut := ""
	if t.Mutable {
		mut = "mut "
	}
	return fmt.Sprintf("&%s %s", mut, t.OfType.PrettyPrint(level))
}

func (i Import) PrettyPrint(level int) string {
	return fmt.Sprintf("use %s;", i.Path)
}

func (p *Program) PrettyPrint(level int) string {
	var sb strings.Builder

	for _, imp := range p.Imports {
		sb.WriteString(imp.PrettyPrint(level))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	for _, decl := range p.Decls {
		sb.WriteString(decl.PrettyPrint(level))
		sb.WriteString("\n\n")
	}

	return sb.String()
}

func (s *StructDecl) PrettyPrint(level int) string {
	var sb strings.Builder

	sb.WriteString("pub struct ")
	sb.WriteString(s.Name)
	sb.WriteString(" {\n")

	for _, field := range s.Fields {
		sb.WriteString("    ")
		sb.WriteString(field.PrettyPrint(level))
		sb.WriteString(",\n")
	}

	sb.WriteString("}")

	return sb.String()
}

func (f *Field) PrettyPrint(level int) string {
	return fmt.Sprintf("pub %s: %s", f.Name, f.Type.PrettyPrint(level))
}

func (t *TypeDecl) PrettyPrint(level int) string {
	var sb strings.Builder

	sb.WriteString("type ")
	sb.WriteString(t.Name)
	sb.WriteString(" = ")
	sb.WriteString(t.Type.PrettyPrint(level))
	sb.WriteString(";")
	return sb.String()
}

func (f *FunctionDecl) PrettyPrint(level int) string {
	var sb strings.Builder

	sb.WriteString("pub fn ")
	sb.WriteString(f.Name)
	sb.WriteString("(")

	params := make([]string, len(f.Params))
	for i, param := range f.Params {
		params[i] = param.PrettyPrint(level)
	}
	sb.WriteString(strings.Join(params, ", "))

	sb.WriteString(") -> ")
	sb.WriteString(f.ReturnType.PrettyPrint(level))
	sb.WriteString(" {\n")

	for _, stmt := range f.Body {
		if f.Body != nil {
			sb.WriteString(stmt.PrettyPrint(level + 1))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("}")

	return sb.String()
}

func (f *ConstDecl) PrettyPrint(level int) string {
	return fmt.Sprintf("pub const %s: %s = %s;", f.Name, f.Type.PrettyPrint(level), f.Value.PrettyPrint(level))
}

func (p *Param) PrettyPrint(level int) string {
	mut := ""
	if p.Mutable {
		mut = "mut "
	}
	return fmt.Sprintf("%s%s: %s", mut, p.Name, p.Type.PrettyPrint(level))
}

func (l *Let) PrettyPrint(level int) string {
	indent := strings.Repeat("    ", level)
	return fmt.Sprintf("%slet %s = %s;", indent, l.VariableName, l.Value.PrettyPrint(0))
}

func (a *Assign) PrettyPrint(level int) string {
	indent := strings.Repeat("    ", level)
	return fmt.Sprintf("%s%s = %s;", indent, a.Dest.PrettyPrint(0), a.Value.PrettyPrint(0))
}

func (r *Return) PrettyPrint(level int) string {
	indent := strings.Repeat("    ", level)
	return fmt.Sprintf("%s %s", indent, r.Value.PrettyPrint(level))
}

func (b *Block) PrettyPrint(level int) string {
	var sb strings.Builder

	indent := strings.Repeat("    ", level)

	sb.WriteString(indent)
	sb.WriteString("{\n")
	for _, stmt := range b.Statements {
		sb.WriteString(stmt.PrettyPrint(level + 1))
		sb.WriteString("\n")
	}
	sb.WriteString(indent)
	sb.WriteString("}\n")

	return sb.String()
}

func (f *FieldValue) PrettyPrint(level int) string {
	indent := strings.Repeat("    ", level)
	return fmt.Sprintf("%s%s: %s", indent, f.Name, f.Value.PrettyPrint(level))
}

func (s *StructCons) PrettyPrint(level int) string {
	var sb strings.Builder

	indent := strings.Repeat("    ", level)

	sb.WriteString(s.StructName)
	sb.WriteString(" {\n")

	for _, field := range s.Fields {
		sb.WriteString(field.PrettyPrint(level + 1))
		sb.WriteString(",\n")
	}

	sb.WriteString(indent)
	sb.WriteString("}")

	return sb.String()
}

func (e *EnumCons) PrettyPrint(level int) string {
	var sb strings.Builder

	sb.WriteString(e.EnumName)
	sb.WriteString("::")
	sb.WriteString(e.Variant)
	sb.WriteString("(")

	params := make([]string, len(e.Params))
	for i, param := range e.Params {
		params[i] = param.PrettyPrint(0)
	}
	sb.WriteString(strings.Join(params, ", "))

	sb.WriteString(")")

	return sb.String()
}

func (b *Borrow) PrettyPrint(level int) string {
	return fmt.Sprintf("&%s", b.Value.PrettyPrint(0))
}

func (t *Tuple) PrettyPrint(level int) string {
	values := make([]string, len(t.Values))
	for i, value := range t.Values {
		values[i] = value.PrettyPrint(0)
	}
	return fmt.Sprintf("(%s)", strings.Join(values, ", "))
}

func (f *FunctionCall) PrettyPrint(level int) string {
	args := make([]string, len(f.Arguments))
	for i, arg := range f.Arguments {
		args[i] = arg.PrettyPrint(0)
	}
	return fmt.Sprintf("%s(%s)", f.FunctionName, strings.Join(args, ", "))
}

func (s *StaticMethodCall) PrettyPrint(level int) string {
	args := make([]string, len(s.Arguments))
	for i, arg := range s.Arguments {
		args[i] = arg.PrettyPrint(0)
	}
	return fmt.Sprintf("%s::%s(%s)", s.TypeName.PrettyPrint(level), s.MethodName, strings.Join(args, ", "))
}

func (m *MethodCall) PrettyPrint(level int) string {
	args := make([]string, len(m.Arguments))
	for i, arg := range m.Arguments {
		args[i] = arg.PrettyPrint(0)
	}
	return fmt.Sprintf("%s.%s(%s)", m.Value.PrettyPrint(level), m.MethodName, strings.Join(args, ", "))
}

func (v *Variable) PrettyPrint(level int) string {
	return v.VariableName
}

func (f *FieldAccess) PrettyPrint(level int) string {
	return fmt.Sprintf("%s.%s", f.Value.PrettyPrint(0), f.Field)
}

func (i *IfElse) PrettyPrint(level int) string {
	var sb strings.Builder

	indent := strings.Repeat("    ", level)

	sb.WriteString(indent)
	sb.WriteString("if ")
	sb.WriteString(i.Condition.PrettyPrint(0))
	sb.WriteString(" {\n")
	sb.WriteString(i.Then.PrettyPrint(level + 1))
	sb.WriteString(indent)
	sb.WriteString("\n} else {\n")
	sb.WriteString(i.Else.PrettyPrint(level + 1))
	sb.WriteString(indent)
	sb.WriteString("\n}")

	return sb.String()
}

func (n *Not) PrettyPrint(level int) string {
	return fmt.Sprintf("!%s", n.Value.PrettyPrint(0))
}

func (a *Add) PrettyPrint(level int) string {
	return fmt.Sprintf("%s + %s", a.Left.PrettyPrint(0), a.Right.PrettyPrint(0))
}

func (u *UInt64Literal) PrettyPrint(level int) string {
	return fmt.Sprintf("%d_u64", u.Value)
}

func (s *StringLiteral) PrettyPrint(level int) string {
	str := fmt.Sprintf("\"%s\"", strings.ReplaceAll(s.Value, "\"", "\\\""))
	return str + ".to_string()"
}

func (b *BoolLiteral) PrettyPrint(level int) string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (t *Todo) PrettyPrint(level int) string {
	indent := strings.Repeat("    ", level)
	return fmt.Sprintf("%stodo!()", indent)
}
