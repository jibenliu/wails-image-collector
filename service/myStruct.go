package service

import (
	"fmt"
	"github.com/wailsapp/wails"
	"log"
	"math/rand"
)

type MyStruct struct {
	Name string
	runtime *wails.Runtime
	log     *wails.CustomLogger
	store   *wails.Store
}

func NewStruct() *MyStruct {
	stu := &MyStruct{
		Name: "myStruct",
	}
	return stu
}

func (s *MyStruct) WailsInit(runtime *wails.Runtime) error {
	// Save runtime
	s.runtime = runtime
	s.log = runtime.Log.New("MyStruct")
	// Do some other initialisation
	s.store = runtime.Store.New("Counter", 0)
	return nil
}

func (s *MyStruct) WailsShutDown() {
	fmt.Println("when wails shutdown then do")
}

func (s *MyStruct) MyBoundMethod(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

func (s *MyStruct) AddUser(name string) error {
	fmt.Printf("add one user named %s\n", name)
	return nil
}

func (s *MyStruct) OpenFile() string {
	selectedFile := s.runtime.Dialog.SelectFile("打开本地文件", "*.jpg,*.png,*.pdf")
	log.Println(selectedFile)
	s.log.Infof("I'm %s with the events that are currently unfolding", selectedFile)
	return selectedFile
}

func (s *MyStruct) RandomValue() {
	s.store.Set(rand.Intn(1000))
}

func (s *MyStruct) Hello(name string) string {
	return fmt.Sprintf("Hello %s! My name is %s", name, s.Name)
}

func (s *MyStruct) Rename(name string) string {
	s.Name = name
	return fmt.Sprintf("My name is now '%s'", s.Name)
}

func (s *MyStruct) privateMethod(name string) string {
	s.Name = name
	return fmt.Sprintf("My name is now '%s'", s.Name)
}
