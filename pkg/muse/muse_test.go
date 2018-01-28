package muse

import (
	"reflect"
	"testing"

	"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
)

func Test001BasicTypeConversion(t *testing.T) {

	cv.Convey(`muse.Pun() should convert basic type.Types to basic reflect.Types`, t, func() {
		m := NewMuse()

		b := types.Typ[types.Bool]
		rt, err := m.Pun(b)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Bool.String())

		s := types.Typ[types.String]
		rt, err = m.Pun(s)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.String.String())

		v := types.Typ[types.Int]
		rt, err = m.Pun(v)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Int.String())

		v = types.Typ[types.Int64]
		rt, err = m.Pun(v)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Int64.String())

		v = types.Typ[types.Uint64]
		rt, err = m.Pun(v)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Uint64.String())

	})
}

func Test002StructTypeConversion(t *testing.T) {

	cv.Convey(`muse.Pun() should convert a struct type.Types `+
		`to a struct reflect.Types`, t, func() {

		m := NewMuse()

		b := types.Typ[types.Bool]
		rt, err := m.Pun(b)
		cv.So(err, cv.ShouldBeNil)
		cv.So(rt.String(), cv.ShouldResemble, reflect.Bool.String())

	})
}
