package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/NicolasAMunozR/PAF/backend-PAF/util"
)

type Login struct {
	// @Desc User is usach email without @usach.cl
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type ResponseLogin struct {
	Data   map[string]interface{} `json:"data"`
	Expire string                 `json:"expire"`
	Token  string                 `json:"token"`
}

func RequestLogin(username string, password string) (*http.Response, error) {
	//Realizar una solicitud a la API de autenticacion para verificar si las credenciales del usuario son validas
	hashPassword := util.HashPassword(password)
	//print("username: ", username)
	//print("password: ", password)
	//print("hashPassword: ", hashPassword)

	body := Login{
		User:     username,
		Password: hashPassword,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Imprime url que consultara
	print("url: ", "https://usach-auth.test-citiaps.cl/auth/login","\n")
	// Imprimelas basic auth
	print("basic auth: ", GetBasicAuth(),"\n")
	request, err := http.NewRequest("POST", "https://usach-auth.test-citiaps.cl/auth/login", bytes.NewBuffer(jsonData))

	if err != nil {

		return nil, err

	}
	//Se agregan los headers a la solicitud
	request.Header.Set("Authorization", GetBasicAuth())
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err

	}

	return response, nil
}

func GetBasicAuth() string {
	username := "api.marin"
	password := "4p1m4r1n"
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

}

func DoRequestLogin(username string, password string) (*ResponseLogin, error) {
	username = strings.ToLower(username)
	username = strings.ReplaceAll(username, "@usach.cl", "")
	response, err := RequestLogin(username, password)

	if err != nil {
		return nil, err
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		var errorMessage map[string]string
		err = json.Unmarshal(responseBody, &errorMessage)
		if err != nil {
			return nil, errors.New("error al obtener el error")
		}
		// En caso de que las credenciales sean invalidas pero exista el usuario registrarlo en la base de datos
		if response.StatusCode == 401 && errorMessage["message"] == "Password usuario invalido" {
			return nil, errors.New("contraseña incorrecta")
		}
		if response.StatusCode == 401 && errorMessage["message"] == "Error LDAP : Success" {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, errors.New(errorMessage["message"])
	}

	var responseLogin ResponseLogin
	err = json.Unmarshal(responseBody, &responseLogin)
	if err != nil {
		return nil, errors.New("error al obtener la respuesta: " + err.Error())
	}
	return &responseLogin, nil

}
