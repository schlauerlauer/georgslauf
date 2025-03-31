package md

import (
	"bytes"
	"sync"

	"github.com/yuin/goldmark"
)

type MdService struct {
	data MdData
	lock sync.RWMutex
}

type MdData struct {
	Headline string
	Intro    []byte
}

type Input struct {
	Headline string `json:"h" schema:"headline" validate:"max=64" mod:"trim,sanitize"`
	Intro    string `json:"i" schema:"intro" validate:"max=2048" mod:"trim,sanitize"`
}

func New() *MdService {
	return &MdService{}
}

func (s *MdService) Get() MdData {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data
}

func (s *MdService) Update(data Input) (MdData, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var buffer bytes.Buffer
	if err := goldmark.Convert([]byte(data.Intro), &buffer); err != nil {
		return MdData{}, err
	}

	output := MdData{
		Intro:    buffer.Bytes(),
		Headline: data.Headline,
	}

	s.data = output

	return output, nil
}
