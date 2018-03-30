package compiler

import (
	"github.com/glycerine/zygomys/zygo"
)

func initZygo() *zygo.Zlisp {
	env := zygo.NewZlisp()
	env.StandardSetup()
	return env
}

func callZygo(env *zygo.Zlisp, s string) (interface{}, error) {
	_, err := env.EvalString(s)
	return nil, err
}
