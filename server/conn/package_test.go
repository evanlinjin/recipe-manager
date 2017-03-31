package conn

import "testing"

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMakePackage(t *testing.T) {
	pac, e := makePersonPac()
	if e != nil {
		t.Error(e)
	}
	t.Log(string(pac))
}

func TestReadPackage(t *testing.T) {
	pac, e := makePersonPac()
	if e != nil {
		t.Error(e)
	}

	var p Person
	if e = ReadPackage(pac, []byte("secret"), &p); e != nil {
		t.Error(e)
	}
	t.Log(p)
}

func TestReadPackage2(t *testing.T) {
	var p Person
	pac := []byte("eyJhZ2UiOjIxLCJuYW1lIjoiRXZhbiJ9.EagFnSRC1XqeUpXRsI_eLfLvs5CSO2M1VQVuBPnqsvw")
	if e := ReadPackage(pac, []byte("secret"), &p); e != nil {
		t.Error(e)
	}
	t.Log(p)
}

func makePersonPac() (pac []byte, e error) {
	p := Person{"Evan", 21}
	k := []byte("secret")

	pac, e = MakePackage(p, k)
	return
}
