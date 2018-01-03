package luajit

import "testing"
import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func (s *State) printstack(t *testing.T) {
	t.Log("--- stack:")
	n := s.Gettop()
	for i := 1; i <= n; i++ {
		switch s.Type(i) {
		case Tstring:
			t.Logf("%s", s.Tostring(i))
		case Tnumber:
			t.Logf("%f", s.Tonumber(i))
		case Tboolean:
			t.Logf("%t", s.Toboolean(i))
		default:
			t.Logf("(%s)", s.Typename(s.Type(i)))
		}
	}
	t.Log("---")
}

func TestPushpop(t *testing.T) {
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate failed")
	}
	defer s.Close()
	r := rand.New(rand.NewSource(1))

	vals := make([]float64, 1000)
	var i int
	if !s.Checkstack(6 * 1000) {
		t.Fatal("not enough slots in stack")
	}
	for i = 0; i < len(vals); i++ {
		vals[i] = r.Float64() * 1000.0
		s.Pushnumber(vals[i])
		s.Pushstring(fmt.Sprintf("!!%f", vals[i]))
		s.Pushinteger(int(vals[i]))
		s.Pushnil()
		s.Pushboolean(true)
		s.Pushboolean(false)
	}
	for i--; i >= 0; i-- {
		n := vals[i]
		if s.Toboolean(-1) {
			t.Errorf("expected false, got true")
		}
		if !s.Toboolean(-2) {
			t.Errorf("expected true, got true")
		}
		s.Pop(3) // pop the nil as well
		if nn := s.Tointeger(-1); nn != int(n) {
			t.Errorf("expected int %d, got %d", int(n), nn)
		}
		ns := fmt.Sprintf("!!%f", n)
		if str := s.Tostring(-2); str != ns {
			t.Errorf("expected string %s, got %s", ns, str)
		}
		if f := s.Tonumber(-3); f != n {
			t.Errorf("expected float64 %f, got %f", n, f)
		}
		s.Pop(3)
	}
}

func TestStacktypes(t *testing.T) {
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate failed")
	}
	defer s.Close()
	r := rand.New(rand.NewSource(2))

	vals := make([]float64, 1000)
	var i int
	if !s.Checkstack(5 * 1000) {
		t.Fatal("not enough slots in stack")
	}
	for i = 0; i < len(vals); i++ {
		vals[i] = r.Float64() * 1000.0
		s.Pushnumber(vals[i])
		s.Pushstring(fmt.Sprintf("!!%f", vals[i]))
		s.Pushinteger(int(vals[i]))
		s.Pushnil()
		s.Pushboolean(true)
	}
	for i--; i >= 0; i-- {
		if !s.Isboolean(-1) {
			t.Errorf("expected boolean")
		}
		if !s.Isnil(-2) {
			t.Errorf("expected nil")
		}
		if !s.Isnumber(-3) {
			t.Errorf("expected number (from int)")
		}
		if !s.Isstring(-4) {
			t.Errorf("expected string")
		}
		if !s.Isnumber(-5) {
			t.Errorf("expected number")
		}
		s.Pop(5)
	}
}

func TestLoad(t *testing.T) {
	txt := `
		function f(x)
			return math.sqrt(x)
		end
		testx = f(400)
		testy = f(1)
		testz = f(36)
	`
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	defer s.Close()
	r := bufio.NewReader(strings.NewReader(txt))
	if r == nil {
		t.Fatal("NewReader returned nil")
	}
	s.Openlibs()
	if err := s.Load(r, "TestLoad"); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	if err := s.Pcall(0, Multret, 0); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	s.Getglobal("testz")
	s.Getglobal("testy")
	s.Getglobal("testx")
	if n := s.Tointeger(-1); n != 20 {
		t.Errorf("expected 20, got %d", n)
	}
	if n := s.Tointeger(-2); n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if n := s.Tointeger(-3); n != 6 {
		t.Errorf("expected 6, got %d", n)
	}
	s.Pop(3)
}

func TestLoadstring(t *testing.T) {
	txt := `
		function f(x)
			return math.sqrt(x)
		end
		testx = f(400)
		testy = f(1)
		testz = f(36)
	`
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	defer s.Close()
	s.Openlibs()
	if err := s.Loadstring(txt); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	if err := s.Pcall(0, Multret, 0); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	s.Getglobal("testz")
	s.Getglobal("testy")
	s.Getglobal("testx")
	if n := s.Tointeger(-1); n != 20 {
		t.Errorf("expected 20, got %d", n)
	}
	if n := s.Tointeger(-2); n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if n := s.Tointeger(-3); n != 6 {
		t.Errorf("expected 6, got %d", n)
	}
	s.Pop(3)
}

func TestRegister(t *testing.T) {
	txt := `
		testx = mysqrt(400)
		testy = mysqrt(1)
		testz = mysqrt(36)
	`
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	defer s.Close()
	s.Openlibs()
	s.Register(func(s *State) int {
		n := s.Tonumber(-1)
		s.Pushnumber(math.Sqrt(n))
		return 1
	}, "mysqrt")
	if err := s.Loadstring(txt); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	if err := s.Pcall(0, Multret, 0); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	s.Getglobal("testz")
	s.Getglobal("testy")
	s.Getglobal("testx")
	if n := s.Tointeger(-1); n != 20 {
		t.Errorf("expected 20, got %d", n)
	}
	if n := s.Tointeger(-2); n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
	if n := s.Tointeger(-3); n != 6 {
		t.Errorf("expected 6, got %d", n)
	}
	s.Pop(3)
}

func TestXmove(t *testing.T) {
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	defer s.Close()
	s2 := s.Newthread()
	if s2 == nil {
		t.Fatal("Newthread returned nil")
	}

	if n := s2.Gettop(); n != 0 {
		t.Errorf("state 2 expected expected empty stack, found %d elems", n)
	}
	if s.Type(-1) != Tthread {
		t.Errorf("state 1 expected thread at stack top, got %s", s.Typename(s.Type(-1)))
	}

	s.Pushinteger(1)
	s.Pushinteger(2)
	s2.Xmove(s, 2)
	if n := s2.Tointeger(-1); n != 2 {
		t.Errorf("expected %d, got %d", 2, n)
	}
	if n := s2.Tointeger(-2); n != 1 {
		t.Errorf("expected %d, got %d", 1, n)
	}
	s2.Pop(2)
	s.Pop(1) // popping s2's thread closes it
	if n := s.Gettop(); n != 0 {
		t.Errorf("state 1 expected expected empty stack, found %d elems", n)
	}
}

func TestResume(t *testing.T) {
	txt := `
		function f(x)
			coroutine.yield(10, x)
		end
		function g(x)
			f(x + 1)
			return 3
		end
	`
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	defer s.Close()
	s.Openlibs()
	if err := s.Loadstring(txt); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	if err := s.Pcall(0, 0, 0); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}

	s2 := s.Newthread()
	if s2 == nil {
		t.Fatal("Newthread returned nil")
	}
	s2.Getglobal("g")
	s2.Pushinteger(20)
	if y, err := s2.Resume(1); err != nil {
		t.Errorf("resume failed: %s – %s", err.Error(), s2.Tostring(-1))
	} else if !y {
		t.Error("expected yield")
	}
	if n := s2.Gettop(); n != 2 {
		t.Errorf("expected 2 items on stack, found %d", n)
	}
	if n := s2.Tointeger(1); n != 10 {
		t.Errorf("expected 10, got %d", n)
	}
	if n := s2.Tointeger(2); n != 21 {
		t.Errorf("expected 21, got %d", n)
	}

	if y, err := s2.Resume(0); err != nil {
		t.Errorf("resume failed: %s – %s", err.Error(), s2.Tostring(-1))
	} else if y {
		t.Error("thread yielded unexpectedly")
	}
	if n := s2.Gettop(); n != 1 {
		t.Errorf("expected 1 item on stack, found %d", n)
	}
	if n := s2.Tointeger(1); n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}

func TestTogofunction(t *testing.T) {
	s := Newstate()
	if s == nil {
		t.Fatal("Newstate returned nil")
	}
	defer s.Close()
	s.Openlibs()
	s.Pushfunction(func(s *State) int {
		n := s.Tonumber(-1)
		s.Pushnumber(math.Sqrt(n))
		return 1
	})
	fn, err := s.Togofunction(-1)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	s.Pop(1)

	s.Pushclosure(fn, 0)
	s.Pushinteger(36)
	if err := s.Pcall(1, 1, 0); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	if n := s.Tointeger(-1); n != 6 {
		t.Errorf("expected 6, got %d", n)
	}
	s.Pop(1)

	s.Pushfunction(fn)
	s.Pushinteger(400)
	if err := s.Pcall(1, 1, 0); err != nil {
		t.Fatalf("%s -- %s", err.Error(), s.Tostring(-1))
	}
	if n := s.Tointeger(-1); n != 20 {
		t.Errorf("expected 20, got %d", n)
	}
	s.Pop(1)
}
