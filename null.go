package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/krhubert/null/v1/convert"
)

type state int

const (
	Unset state = 0
	Nil   state = 1
	Set   state = 2
)

func (s state) String() string {
	switch s {
	case Unset:
		return "unset"
	case Nil:
		return "nil"
	case Set:
		return "set"
	}
	return fmt.Sprintf("null.state(%d)", s)
}

type Null[T any] struct {
	V     T
	State state
}

func New[T any](t T) Null[T] {
	return Null[T]{V: t, State: Set}
}

func NewNull[T any]() Null[T] {
	return Null[T]{State: Nil}
}

func NewFromPtr[T any](t *T) Null[T] {
	if t == nil {
		return Null[T]{State: Nil}
	}
	return Null[T]{V: *t, State: Set}
}

func (n Null[T]) Valid() bool {
	return n.State == Set
}

func (n Null[T]) Nil() bool {
	return n.State == Nil
}

func (n Null[T]) Ptr() *T {
	if n.State == Set {
		t := n.V
		return &t
	}
	return nil
}

func (n *Null[T]) Set(t T) {
	n.V = t
	n.State = Set
}

func (n *Null[T]) SetPtr(t *T) {
	if t == nil {
		var zero T
		n.V = zero
		n.State = Nil
		return
	}

	n.V = *t
	n.State = Set
}

func (n *Null[T]) Unset(t T) {
	var zero T
	n.V = zero
	n.State = Unset
}

func (n Null[T]) MarshalJSON() ([]byte, error) {
	if n.State == Set {
		return json.Marshal(n.V)
	}
	return []byte("null"), nil
}

func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		var zero T
		n.V = zero
		n.State = Nil
		return nil
	}

	if err := json.Unmarshal(data, &n.V); err != nil {
		return err
	}
	n.State = Set
	return nil
}

func (n *Null[T]) Scan(value any) error {
	if value == nil {
		var zero T
		n.V = zero
		n.State = Nil
		return nil
	}

	if err := convert.ConvertAssign(&n.V, value); err != nil {
		return err
	}
	n.State = Set
	return nil
}

func (n Null[T]) Value() (driver.Value, error) {
	if n.State == Set {
		if valuer, ok := any(n.V).(driver.Valuer); ok {
			return valuer.Value()
		}
		return n.V, nil
	}
	return nil, nil
}
