package chefs

import "testing"

func TestMakeChefsDB(t *testing.T) {
	chefs, e := MakeChefsDB()
	if e != nil {
		t.Error(e)
		return
	}
	if e := chefs.Close(); e != nil {
		t.Error(e)
		return
	}
}

func TestChefsDB_Initiate(t *testing.T) {
	chefs, e := MakeChefsDB()
	if e != nil {
		t.Error(e)
		return
	}
	defer chefs.Close()

	if e := chefs.Initiate(); e != nil {
		t.Error(e)
		return
	}
}

func TestGetRandUniqID(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Log(GetRandUniqID())
	}
}

func TestChefsDB_AddChef(t *testing.T) {
	chefs, e := MakeChefsDB()
	if e != nil {
		t.Error(e)
		return
	}
	defer chefs.Close()

	if e := chefs.Initiate(); e != nil {
		t.Error(e)
		return
	}

	e = chefs.AddChef("evan.linjin@gmail.com", "testpwd")
	if e != nil {
		t.Error(e)
	}
}
