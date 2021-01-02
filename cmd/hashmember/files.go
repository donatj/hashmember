package main

import (
	"os"

	"github.com/donatj/hashmember"
)

func loadHashmember(file string) (hashmember.Hashmember, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return hashmember.Decode(f)
}

func saveHashmember(file string, hm hashmember.Hashmember) error {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return err
	}

	return hashmember.Encode(f, hm)
}
