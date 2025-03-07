package controller

import (
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/NicolasAMunozR/PAF/backend-PAF/util"
	"github.com/NicolasAMunozR/PAF/backend-PAF/mailer"
	"github.com/NicolasAMunozR/PAF/backend-PAF/middleware"
	

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// Roles en el sistema
const (
	RolAdmin = "Admin"
	RolUser  = "User"
)

// Tipo de código
const (
	ActivationCode   = 1
	RecuperationCode = 2
)

// AuthenticationController : Estructura controladora de las colecciones
type AuthenticationController struct {
}

// Routes : Define las rutas del controlador
func (authenticationController *AuthenticationController) Routes(base *gin.RouterGroup, authNormal *jwt.GinJWTMiddleware) {

	// Refresh time can be longer than token timeout
	base.GET("/refresh_token",
		middleware.SetRoles(RolAdmin, RolUser),
		authNormal.MiddlewareFunc(),
		authNormal.RefreshHandler)

	//funcion de login: recibe un objeto {email: , pass:}
	base.POST("/login", authNormal.LoginHandler)

	//creacion de usuarios
	base.POST("/user",
		CreateUser)
	base.POST("/create-professional", CreateProfessional)
	base.POST("/activate-user", authNormal.MiddlewareFunc(), ActivateUser)
	base.POST("/activation-email", authNormal.MiddlewareFunc(), SendActivation)

	base.POST("/reset-password", ResetPassword)

	base.POST("/validation-code", ValidationCode)

	base.POST("/update-password", UpdatePassword)
}

var userModel models.User

func CreateProfessional(c *gin.Context) {
	var user models.User
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo registrar al usuario", e))
		return
	}

	user.ID = bson.NewObjectId()
	query := bson.M{}
	user.Password = ""
	user.Activation = false
	query["email"] = user.Email
	users, _ := userModel.Find(query)
	if len(users) != 0 {
		c.JSON(http.StatusBadRequest, util.GetError("Este email ya ha sido registrado", e))
		return
	}
	query = bson.M{}
	query["rut"] = user.Rut
	users, _ = userModel.Find(query)
	if len(users) != 0 {
		c.JSON(http.StatusBadRequest, util.GetError("Este rut ya ha sido registrado", e))
		return
	}
	// se crea la contraseña
	/*
		arregloName := strings.Split(user.Name, " ")
		titleName := strings.Title(arregloName[0])
		code := titleName + "123"*/
	code := util.SecureRandomAlphaString(8)
	user.Hash = middleware.GeneratePassword(code)
	if err := userModel.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, util.GetError("Fallo al crear el usuario", err))
		return
	}
	panelUrl := os.Getenv("PANEL_URL")
	templateDir := "mailer/templates/professionalUserTemplate.html"
	r := mailer.NewRequest([]string{user.Email}, "Inscripción Yoinformogral")
	err := r.SendMailSkipTLS(templateDir, map[string]string{"username": user.Name, "code": code, "section.name": user.Section.Name, "panelUrl": panelUrl})
	if err != nil {
		_ = userModel.Delete(user.ID.Hex())
		c.JSON(http.StatusBadRequest, util.GetError("Email invalido", err))
		return
	}
	c.String(http.StatusCreated, "")
}

// CreateUser : Registrar usuario
func CreateUser(c *gin.Context) {
	var user models.User
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo registrar al usuario", e))
		return
	}
	user.ID = bson.NewObjectId()

	query := bson.M{}
	if user.Rol != "general" {
		user.Hash = middleware.GeneratePassword(user.Password)
		user.Password = ""
		user.Activation = false
		query["email"] = user.Email
	}

	users, _ := userModel.Find(query)
	if len(users) != 0 {
		c.JSON(http.StatusBadRequest, util.GetError("Este rut o email ya ha sido registrado", e))
		return
	}

	if err := userModel.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, util.GetError("Fallo al crear el usuario", err))
		return
	}
	if user.Rol != "general" && user.Rol != "informante" {

		safeCode := models.SafeCode{
			ID:          bson.NewObjectId(),
			User:        user.ID,
			CreatedAt:   time.Now(),
			Code:        "",
			UserPayload: user,
			Type:        ActivationCode,
		}

		err := safeCodeModel.Create(&safeCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo procesar la solicitud", err))
			return
		}
		templateDir := "mailer/templates/activationTemplate.html"
		r := mailer.NewRequest([]string{user.Email}, "Activación cuenta Yoinformogral")

		err = r.SendMailSkipTLS(templateDir, map[string]string{"username": user.Name, "code": safeCode.Code})
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("Email invalido", err))
			return
		}

	}

	c.String(http.StatusCreated, "")
}

func ActivateUser(c *gin.Context) {

	currentUser := userModel.LoadFromContext(c)
	type CodeData struct {
		Code  string `form:"code" json:"code"`
		Email string `form:"email" json:"email" `
	}

	var codeData CodeData
	var err error
	if err = c.BindJSON(&codeData); err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo encontrar el código", err))
		return
	}

	safeCode, err := safeCodeModel.FindOne(bson.M{"code": codeData.Code, "type": ActivationCode})
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Código no válido", err))
		return
	}

	if codeData.Email != "" {
		if safeCode.UserPayload.Email != codeData.Email {
			c.JSON(http.StatusBadRequest, util.GetError("Código no válido", errors.New("Invalid Data")))
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, util.GetError("Codigo no válido", errors.New("Debe especificar un email o nombre de usuario")))
		return
	}

	// dateNow := time.Now()

	// if safeCode.Expiration.Before(dateNow) {
	// 	c.JSON(http.StatusRequestTimeout, util.GetError("El código ha caducado", errors.New("Timeout code")))
	// 	return
	// }

	err = userModel.UpdateField(currentUser.ID, bson.M{"$set": bson.M{"activation": true}})

	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Error al activar cuenta", err))
		return
	}

	safeCodeModel.Delete(safeCode.ID)

	c.String(http.StatusOK, "OK")
	return

}
func SendActivation(c *gin.Context) {

	currentUser := userModel.LoadFromContext(c)

	user, err := userModel.Get(currentUser.ID.Hex())
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Usuario no válido", err))
	}

	if user.Activation {
		c.JSON(http.StatusBadRequest, util.GetError("Usuario ya se encuentra activado", errors.New("User is Active")))
	}
	if currentSafe, err := safeCodeModel.FindOne(bson.M{"user": user.ID, "type": ActivationCode}); err == nil {
		safeCodeModel.Delete(currentSafe.ID)
	}

	// timein := time.Now().Add(time.Minute * 10)
	safeCode := model.SafeCode{
		ID:          bson.NewObjectId(),
		User:        user.ID,
		CreatedAt:   time.Now(),
		Code:        "",
		UserPayload: *user,
		Type:        ActivationCode,
	}

	err = safeCodeModel.Create(&safeCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se logró procesar la solicitud", err))
		return
	}
	templateDir := "mailer/templates/activationTemplate.html"
	r := mailer.NewRequest([]string{user.Email}, "Activación cuenta Yoinformogral")

	err = r.SendMailSkipTLS(templateDir, map[string]string{"username": user.Name, "code": safeCode.Code})
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Email invalido", err))
		return
	}

	c.String(http.StatusOK, "OK")
	return

}

func ResetPassword(c *gin.Context) {
	type UserData struct {
		Email string `form:"email" json:"email" `
	}

	var credentials UserData
	var err error
	if err = c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo encontrar al usuario/a", err))
		return
	}

	var userResponse model.User

	if credentials.Email != "" {
		userResponse, err = userModel.FindOne(bson.M{"email": credentials.Email})

	} else {
		c.JSON(http.StatusBadRequest, util.GetError("Usuarion ingresado incorrecto", errors.New("Debe especificar un email o nombre de usuario")))
		return
	}

	if userResponse.Rol == "profesional a cargo" && !userResponse.Activation {
		c.JSON(http.StatusBadRequest, util.GetError("Debe acceder y modificar la contraseña asignada a su correo primero", errors.New("Debe especificar un email o nombre de usuario")))
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo encontrar al usuario/a", err))
		return
	}
	if currentSafe, err := safeCodeModel.FindOne(bson.M{"user": userResponse.ID, "type": RecuperationCode}); err == nil {
		safeCodeModel.Delete(currentSafe.ID)
	}

	timein := time.Now().Add(time.Minute * 2)
	safeCode := model.SafeCode{
		ID:          bson.NewObjectId(),
		User:        userResponse.ID,
		CreatedAt:   time.Now(),
		Expiration:  timein,
		Code:        "",
		Type:        RecuperationCode,
		UserPayload: userResponse,
	}
	err = safeCodeModel.Create(&safeCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo procesar la solicitud", err))
		return
	}
	templateDir := "mailer/templates/recuperationTemplate.html"
	r := mailer.NewRequest([]string{userResponse.Email}, "Recuperación contraseña YoInformogral")

	err = r.SendMailSkipTLS(templateDir, map[string]string{"username": userResponse.Name, "code": safeCode.Code})
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Email invalido", err))
		return
	}

	c.String(http.StatusOK, "OK")
	return
}

func ValidationCode(c *gin.Context) {

	type CodeData struct {
		Code  string `form:"code" json:"code"`
		Email string `form:"email" json:"email" `
	}

	var codeData CodeData
	var err error
	if err = c.BindJSON(&codeData); err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo encontrar el código", err))
		return
	}

	safeCode, err := safeCodeModel.FindOne(bson.M{"code": codeData.Code, "type": RecuperationCode})
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Código no válido", err))
		return
	}

	if codeData.Email != "" {
		if safeCode.UserPayload.Email != codeData.Email {
			c.JSON(http.StatusBadRequest, util.GetError("Código no válido", errors.New("Invalid Data")))
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, util.GetError("Codigo no válido", errors.New("Debe especificar un email o nombre de usuario")))
		return
	}

	dateNow := time.Now()

	if safeCode.Expiration.Before(dateNow) {
		c.JSON(http.StatusRequestTimeout, util.GetError("El código ha caducado", errors.New("Timeout code")))
		return
	}
	err = safeCodeModel.UpdateField(safeCode.ID, bson.M{"$set": bson.M{"expiration": safeCode.Expiration.Add(time.Minute * 2)}})

	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Error al verificar código", err))
		return
	}

	c.String(http.StatusOK, "OK")
	return
}

func UpdatePassword(c *gin.Context) {

	type UserData struct {
		Code        string `form:"code" json:"code"`
		NewPassword string `form:"newPassword" json:"newPassword"`
		Email       string `form:"email" json:"email" `
	}

	var userData UserData
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Datos recibidos incorrectos", err))
		return
	}

	safeCode, err := safeCodeModel.FindOne(bson.M{"code": userData.Code, "type": RecuperationCode})

	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("Información no válida", err))
	}

	if userData.Email != "" {
		if safeCode.UserPayload.Email != userData.Email {
			c.JSON(http.StatusBadRequest, util.GetError("Información no válida", errors.New("Invalid Data")))
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, util.GetError("Información no válida", errors.New("Invalid Data")))
		return
	}
	dateNow := time.Now()

	if safeCode.Expiration.Before(dateNow) {
		c.JSON(http.StatusRequestTimeout, util.GetError("El código ha caducado", errors.New("Timeout code")))
		return
	}

	newHash := middleware.GeneratePassword(userData.NewPassword)

	err = userModel.UpdateField(safeCode.User, bson.M{"$set": bson.M{"_hash": newHash}})

	if err != nil {
		c.JSON(http.StatusRequestTimeout, util.GetError("No se pudo actualizar la contraseña", errors.New("Timeout code")))
		return
	}

	safeCodeModel.Delete(safeCode.ID)

	c.String(http.StatusOK, "OK")
	return
}
