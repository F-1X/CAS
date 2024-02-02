package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "files"
	pathname := CASPathTransformFunc(key)
	fmt.Println(pathname)

	expectedOriginalKey := "a1f13b3bc20a296e08c212be9c56c706c10abc4f"
	expectedPathName := "a1f13/b3bc2/0a296/e08c2/12be9/c56c7/06c10/abc4f"
	if pathname.Path != expectedPathName {
		t.Errorf("have %s want %s", pathname.Path, expectedPathName)
	}
	if pathname.Filename != expectedOriginalKey {
		t.Errorf("have %s want %s", pathname.Filename, expectedOriginalKey)
	}
	

}

func TestStore(t *testing.T){
	opts := StoreOpts{
		PathTransFormFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)

	key := "somekey"
	data := []byte("some bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}


	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)

	if string(b) != string(data) {
		t.Errorf("want %s have %s",string(data), string(b))
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T){
	opts := StoreOpts{
		PathTransFormFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)

	key := "somekey"

	data := []byte("some bytes for delete")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}