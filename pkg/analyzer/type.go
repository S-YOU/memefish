package analyzer

type Type interface {
	String() string
	EqualTo(t Type) bool
	CastTo(t Type) bool
}

type SimpleType string

const (
	Int64Type     SimpleType = "INT64"
	Float64Type   SimpleType = "FLOAT64"
	BoolType      SimpleType = "BOOL"
	StringType    SimpleType = "STRING"
	BytesType     SimpleType = "BYTES"
	DateType      SimpleType = "DATE"
	TimestampType SimpleType = "TIMESTAMP"
)

type ArrayType struct {
	// A nested array is not supported in Spanner, so Item never become ArrayType.
	Item Type
}

type StructType struct {
	Fields []*StructField
}

type StructField struct {
	Name string
	Type Type
}

func (s SimpleType) String() string {
	return string(s)
}

func (a *ArrayType) String() string {
	return "ARRAY<" + a.Item.String() + ">"
}

func (s *StructType) String() string {
	t := "STRUCT<"
	for i, f := range s.Fields {
		if i != 0 {
			t += ", "
		}
		if f.Name != "" {
			t += f.Name + " "
		}
		t += f.Type.String()
	}
	t += ">"
	return t
}

func (s SimpleType) EqualTo(t Type) bool {
	if t, ok := t.(SimpleType); ok {
		return s == t
	} else {
		return false
	}
}

func (a *ArrayType) EqualTo(t Type) bool {
	if t, ok := t.(*ArrayType); ok {
		return a.Item.EqualTo(t.Item)
	} else {
		return false
	}
}

func (s *StructType) EqualTo(t Type) bool {
	if t, ok := t.(*StructType); ok {
		if len(s.Fields) != len(t.Fields) {
			return false
		}
		for i, f := range s.Fields {
			if !f.Type.EqualTo(t.Fields[i].Type) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (s SimpleType) CastTo(t Type) bool {
	if t, ok := t.(SimpleType); !ok {
		return false
	} else {
		// The same types can be cast to each others of course.
		if s == t {
			return true
		}

		// See: https://cloud.google.com/spanner/docs/functions-and-operators#casting
		switch s {
		case Int64Type:
			return t == StringType || t == Float64Type || t == BoolType
		case Float64Type:
			return t == StringType || t == Int64Type
		case BoolType:
			return t == StringType || t == Int64Type
		case StringType:
			return true // StringType can cast to any types via parsing.
		case BytesType:
			return t == StringType
		case DateType:
			return t == StringType || t == TimestampType
		case TimestampType:
			return t == StringType || t == DateType
		}
	}

	panic("unreachable")
}

func (a *ArrayType) CastTo(t Type) bool {
	if t, ok := t.(*ArrayType); ok {
		return a.Item.EqualTo(t.Item)
	} else {
		return false
	}
}

func (s *StructType) CastTo(t Type) bool {
	if t, ok := t.(*StructType); ok {
		if len(s.Fields) != len(t.Fields) {
			return false
		}
		for i, f := range s.Fields {
			if !f.Type.EqualTo(t.Fields[i].Type) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}