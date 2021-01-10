package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type CepValidade struct {
	Cep    string
	Valido bool
}

func main() {
	b, err := ioutil.ReadFile("./ceps.txt")
	if err != nil {
		log.Fatal(err)
	}
	ceps := strings.Split(string(b), "\r\n")
	_, err = ValidaCeps(ceps)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 10)
}

func ValidaCeps(ceps []string) ([]CepValidade, error) {
	t0 := time.Now()
	for _, cep := range ceps {
		//USO DE SEMÁFOROS
		go func(cep string) {
			valido, err := validaCep(cep)
			fmt.Println(valido, err, time.Since(t0))
		}(cep)
	}
	fmt.Println("acabou de rodar o loop")
	return nil, nil
}

func validaCep(cep string) (bool, error) {
	var v struct {
		Erro bool `json:"erro"`
	}
	res, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%v/json", cep))
	if err != nil {
		return false, fmt.Errorf("validacep: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return false, fmt.Errorf("validacep: %w", err)
	}
	return !v.Erro, nil
}
