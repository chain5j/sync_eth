// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/4/1
package reflectutil

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

var age = 10

type Student struct {
	Name  string  `sync:"stu_name,nil"`
	Age   int
	Score float32
	sex   int
}

func (s Student) Print() {
	fmt.Println(s)
}

func (s Student) print() {
	fmt.Println(s)
}

func TestTag(t *testing.T) {
	var a Student = Student{
		Name:  "stu01",
		Age:   18,
		Score: 92.8,
	}
	TagStruct(&a)

	valueOf := reflect.ValueOf(&a)
	s := valueOf.Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, // 字段名称
			f.Type(),              // 字段类型
			f.Interface())
	}
	s.Field(0).SetString("stu")
	s.Field(1).SetInt(77)
	s.Field(2).SetFloat(77.0)
	fmt.Println("a is now", a)

	v := valueOf.MethodByName("Print")
	v.Call([]reflect.Value{})
}

func TagStruct(a interface{}) {
	typ := reflect.TypeOf(a)

	tag := typ.Elem().Field(0).Tag.Get("sync_eth")
	fmt.Printf("Tag:%s\n", tag)
}

func TestDoc(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset
	d, err := parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, f := range d {
		fmt.Println("package", k)
		for n, f := range f.Files {

			fmt.Println(fmt.Sprintf("file name: %q", n))
			for i, c := range f.Comments {
				fmt.Println(fmt.Sprintf("Comment Group %d", i))
				for _, c1 := range c.List {
					//fmt.Println(fmt.Sprintf("Comment %d: Position: %d, Text: %q", i2, c1.Slash, c1.Text))
					fmt.Println(fmt.Sprintf("Text: %q", c1.Text))
				}
			}
		}

		p := doc.New(f, "./", doc.AllDecls)

		for _, t := range p.Types {
			fmt.Println("type=", t.Name, "docs=", t.Doc)
			for _, m := range t.Methods {
				fmt.Println("type=", m.Name, "docs=", m.Doc)
			}
		}

		for _, v := range p.Vars {
			fmt.Println("type", v.Names)
			fmt.Println("docs:", v.Doc)
		}

		for _, f := range p.Funcs {
			fmt.Println("type", f.Name)
			fmt.Println("docs:", f.Doc)
		}

		for _, n := range p.Notes {
			fmt.Println("body", n[0].Body)
		}
	}
}
