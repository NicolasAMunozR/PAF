package controller

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/NicolasAMunozR/PAF/backend-PAF/middleware"
	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/NicolasAMunozR/PAF/backend-PAF/util"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

type UserController struct {
}

func (userController *UserController) Routes(base *gin.RouterGroup, authNormal *jwt.GinJWTMiddleware) *gin.RouterGroup {

	userRouter := base.Group("/users")
	{
		userRouter.GET("", authNormal.MiddlewareFunc(), userController.GetAll())
		userRouter.PUT("/:id", authNormal.MiddlewareFunc(), userController.Update())
		userRouter.DELETE("/:id", authNormal.MiddlewareFunc(), userController.Delete())
		userRouter.GET("/:id", authNormal.MiddlewareFunc(), userController.One())
		userRouter.GET("/search", authNormal.MiddlewareFunc(), userController.Search())
	}
	base.GET("/usuarios/byOffice/:id", userController.GetByOffice())
	base.GET("/beca/user/:id", userController.GetBecas())
	return userRouter
}

func (userController *UserController) GetAll() func(c *gin.Context) {

	return func(c *gin.Context) {
		pagination := PaginationParams{}
		err := c.ShouldBind(&pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se puedieron encontrar los parametros limit, offset", err))
			return
		}
		query := bson.M{}
		var value []bson.M
		var user bool
		if len(pagination.Roles) > 0 {
			user = true
			for _, rol := range pagination.Roles {
				value = append(value, bson.M{"rol": rol})
				fmt.Println(rol)
			}
			query["$or"] = value
		} else {
			user = false
		}
		if pagination.Section != "" {
			query["section.name"] = pagination.Section
		}
		fmt.Println("SEARCH PAGINATON: ", pagination.Search)
		if pagination.Search != "" {
			searchQuery := bson.M{
				"$or": []bson.M{
					{"name": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"email": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"rut": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"section.name": bson.M{"$regex": pagination.Search, "$options": "i"}},
					{"rol": bson.M{"$regex": pagination.Search, "$options": "i"}},
				},
			}
			if user {
				query = bson.M{"$and": []bson.M{query, searchQuery}}
			} else {
				query = bson.M{"$or": []bson.M{query, searchQuery}}
			}
		}

		// query["rol"] = bson.M{"$ne": "admin"}

		page, total, err := userModel.FindPaginate(query, pagination.Limit, pagination.Offset, user, pagination.Search)

		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se pudo obtener la lista de usuarios", err))
		}
		c.Header("Pagination-Count", fmt.Sprintf("%d", total))
		// if len(page.Metadata) != 0 {
		// }

		c.JSON(http.StatusOK, page)
	}
}

func (userController *UserController) Update() func(c *gin.Context) {
	return func(c *gin.Context) {

		var user models.User
		err := c.Bind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo convertir collection json", err))
			return
		}
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, util.GetError("No se encuentra parametro :id", nil))
			return
		}

		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		currentUser, err := userModel.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se encontró el usuario", err))
			return
		}

		if user.Name != "" {
			currentUser.Name = user.Name
		}
		if user.Email != "" {
			currentUser.Email = user.Email
		}
		if user.Phone != "" {
			currentUser.Phone = user.Phone
		}
		if user.Section.Name != "" {
			currentUser.Section = user.Section
		}
		if user.Password != "" {
			if user.Password == currentUser.Password {
				c.JSON(http.StatusBadRequest, util.GetError("Se debe ingresar una contraseña diferente", err))
				return
			}
			currentUser.Hash = middleware.GeneratePassword(user.Password)
			currentUser.Password = ""
		}

		currentUser.Activation = user.Activation
		// Update
		err = userModel.Update(id, currentUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo actualizar usuario", err))
			return
		}
		if user.Section.ID != "" {
			err = roundRobin.UpdateOffice(user.Section.Name, id)
			if err != nil {
				c.JSON(http.StatusBadRequest, util.GetError("No se pudo actualizar round robin de usuario", err))
				return
			}
		}
		c.String(http.StatusOK, "")
	}
}

func (userController *UserController) Delete() func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		err := userModel.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo encontrar usuario", err))
			return
		}
		err = roundRobin.Update(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo eliminar usuario round robin", err))
			return
		}
		c.String(http.StatusOK, "")
	}
}

// One : Obtener user por _id
func (userController *UserController) One() func(c *gin.Context) {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		group, err := userModel.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se encontró perro", err))
			return
		}
		c.JSON(http.StatusOK, group)
	}
}
func (userController *UserController) GetByOffice() func(c *gin.Context) {
	return func(c *gin.Context) {
		idOffice := c.Param("id")
		if idOffice == "" {
			c.JSON(http.StatusNotFound, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(idOffice) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		group, err := userModel.GetByOffice(idOffice)
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se encontró usuarios", err))
			return
		}
		c.JSON(http.StatusOK, group)
	}
}
func (userController *UserController) GetBecas() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusNotFound, util.GetError("No se encuentra parametro :id", nil))
			return
		}
		if !bson.IsObjectIdHex(id) {
			c.JSON(http.StatusInternalServerError, util.GetError("El id ingresado no es válido", nil))
			return
		}
		group, err := userModel.GetBecas(id)
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se encontró perro", err))
			return
		}
		c.JSON(http.StatusOK, group)
	}
}

func (userController *UserController) Search() func(c *gin.Context) {
	return func(c *gin.Context) {
		pagination := PaginationParams{}
		err := c.ShouldBind(&pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se puedieron encontrar los parametros", err))
			return
		}
		users, err := userModel.Search(pagination.Search, pagination.Limit, pagination.Offset)
		if err != nil {
			c.JSON(http.StatusNotFound, util.GetError("No se encontraro coincidencias", err))
			return
		}
		c.JSON(http.StatusOK, users)
	}
}
