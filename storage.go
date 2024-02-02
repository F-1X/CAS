package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	fmt.Println(hash,hashStr)
	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i:=0; i<sliceLen; i++ {
		from, to := i * blockSize, (i * blockSize) + blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Path: strings.Join(paths,"/"),
		Filename: hashStr,
	}
}


type PathKey struct {
	Path string
	Filename string
}

type PathTransFormFunc func(string) PathKey

var DefaultPathTransFormFunc = func(key string) string {return key}

type StoreOpts struct {
	PathTransFormFunc PathTransFormFunc
}

type Store struct {
	StoreOpts
}
func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransFormFunc(key)

	_, err := os.Stat(pathKey.Fullpath())
	if err == os.ErrNotExist {
		return false
	}

	return true
}

func (p PathKey) Fullpath() string {
	return fmt.Sprintf("%s/%s",p.Path,p.Filename)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransFormFunc(key)

	defer func(){
		log.Printf("deleted: [%s] from disk",pathKey.Fullpath())
	}()
	
	firstFolder := strings.Split(pathKey.Path,"/")[0]
	return os.RemoveAll(firstFolder)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil,err
	}

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, f)
	if err != nil {
		return nil,err
	}
	f.Close()

	return buf, nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransFormFunc(key)

	return os.Open(pathKey.Fullpath())

}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransFormFunc(key)

	if err := os.MkdirAll(pathName.Path, os.ModePerm); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf,r)

	fullpath := pathName.Fullpath()
	
	f, err := os.Create(fullpath)
	if err != nil {
		return nil
	}

	n, err := io.Copy(f,buf)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes %s", n, fullpath)

	return nil
}