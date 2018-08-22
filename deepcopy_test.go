package deepcopy

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/ulule/deepcopier"
)

// User is a user
type User struct {
	// Basic string field
	Name string
	// Deepcopier supports https://golang.org/pkg/database/sql/driver/#Valuer
	Email *sql.NullString
}

var user = &User{
	Name: "gilles",
	Email: &sql.NullString{
		Valid:  true,
		String: "gilles@example.com",
	},
}
var clone *User

func ululeDeepCopier() *User {
	clone := &User{}

	_ = deepcopier.Copy(user).To(clone)
	return clone
}

func jinzhuCopier() *User {
	clone := &User{}

	_ = copier.Copy(clone, user)
	return clone
}

func gobCopier(a, b interface{}) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
}

func TestUluleDeepCopier(t *testing.T) {
	clone := &User{}
	err := deepcopier.Copy(user).To(clone)
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(user, clone) {
		t.Fail()
	}
}

func TestJinzhuCopier(t *testing.T) {
	clone := &User{}
	err := copier.Copy(clone, user)
	if err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(user, clone) {
		t.Fail()
	}
}

func TestGobCopier(t *testing.T) {
	clone := &User{}
	gobCopier(user, clone)
	if !reflect.DeepEqual(user, clone) {
		t.Fail()
	}
}

func BenchmarkUluleDeepCopier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		clone = ululeDeepCopier()
	}
}

func BenchmarkJinzhuCopier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		clone = jinzhuCopier()
	}
}

func BenchmarkGobCopier(b *testing.B) {
	clone := &User{}
	for i := 0; i < b.N; i++ {
		gobCopier(user, clone)
	}
}
