package service

import (
	"fmt"
	"github.com/wailsapp/wails"
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

func (s *MyStruct) OpenFile() {
	//selectedFile := runtime.Dialog.SelectFile("Select your profile picture", "*.jpg,*.png")
	//s.log.Infof("I'm %s with the events that are currently unfolding", selectedFile)
}

func (s *MyStruct) RandomValue() {
	s.store.Set(rand.Intn(1000))
}
