package model

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/citiaps/yoinformogral-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
)

// User : Usuario del sistema
type BecaUser struct {
	Beca bson.ObjectId `json:"beca" bson:"beca,omitempty"`
}

type User struct {
	ID          bson.ObjectId `json:"id"                  bson:"_id,omitempty"`
	Email       string        `json:"email"               bson:"email"`
	Phone       string        `json:"phone"               bson:"phone"`
	Name        string        `json:"name"                bson:"name"`
	Rut         string        `json:"rut"                 bson:"rut"`
	Section     Office        `json:"section"             bson:"section"`
	Rol         string        `json:"rol"                 bson:"rol"`
	Activation  bool          `json:"activation"          bson:"activation"`
	Participant string        `json:"participant"         bson:"participant"`
	Informant   string        `json:"informant"           bson:"infomant"`
	Hash        string        `json:"_hash"               bson:"_hash,omitempty"`
	Password    string        `json:"password,omitempty"  bson:"password,omitempty"`
	Genre       string        `json:"genre"               bson:"genre"`
	Work        string        `json:"work"                bson:"work"`
	Heightini   string        `json:"heightini"           bson:"heightini"`
	Weightini   string        `json:"weightini"           bson:"weightini"`
	Weightgoal  string        `json:"weightgoal"          bson:"weightgoal"`
	Birthdate   string        `json:"birthdate"           bson:"birthdate"`
	Type        []string      `json:"type"                bson:"type"`
	Becas       []BecaUser    `json:"becas,omitempty"                bson:"becas,omitempty"`
	Count       int           `json:"count,omitempty"                bson:"count,omitempty"`
	Cargo       interface{}   `json:"cargo,omitempty"                bson:"cargo,omitempty"`
	Facultad    interface{}   `json:"facultad,omitempty"                bson:"facultad,omitempty"`
	Programa    string        `json:"programa,omitempty"                bson:"programa,omitempty"`
}

// LoadFromContext : Traer usuario desde contexto
func (userModel *User) LoadFromContext(c *gin.Context) *User {
	claims := jwt.ExtractClaims(c)
	var user User
	err := mapstructure.Decode(claims["user"], &user)
	if err != nil {
		panic(err)
	}
	// Accede a la propiedad "user" y conviértela en un mapa
	userMap := claims["user"].(map[string]interface{})
	if userMap["rol"].(string) == "profesional a cargo" {
		sectionMap := userMap["section"].(map[string]interface{})
		// Accede a la propiedad "id" dentro de "section"
		sectionID := sectionMap["id"].(string)
		user.Section.ID = bson.ObjectIdHex(sectionID)
	}

	user.ID = bson.ObjectIdHex(claims["user"].(map[string]interface{})["id"].(string))
	user.Hash = ""
	return &user
}

// Create : Traer usuario desde contexto
func (userModel *User) Create(user *User) error {

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()

	err := colUser.Insert(&user)

	return err
}

// GetUser : Se obtiene el usuario
func GetUser(c *gin.Context) {
	id := c.Param("id")

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var usuario User

	if err := colUser.FindId(bson.ObjectIdHex(id)).One(&usuario); err != nil {
		c.JSON(http.StatusNotFound, util.GetError("Usuario no encontrado", err))
	} else {
		c.JSON(http.StatusCreated, usuario)
	}
}

func (userModel *User) FindOne(query bson.M) (User, error) {

	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	user := User{}

	err := col.Find(query).One(&user)
	return user, err
}
func (userModel *User) Find(query bson.M) ([]User, error) {

	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	users := []User{}

	err := col.Find(query).All(&users)
	return users, err
}

func (userModel *User) FindPaginate(query bson.M, limit int, offset int, user bool, search string) ([]User, int, error) {

	col, session := GetCollection(CollectionNameUser)
	defer session.Close()

	pageDoc := []User{}
	var err error
	// err := col.Pipe(pipeline).One(&pageDoc)
	if user {
		mgoQuery := col.Find(query)
		// Aplica el límite y el desplazamiento después de la búsqueda
		if limit > 0 {
			mgoQuery = mgoQuery.Limit(limit)
		}
		if offset > 0 {
			mgoQuery = mgoQuery.Skip(offset)
		}
		// Ejecuta la consulta con el límite y el desplazamiento aplicados
		err = mgoQuery.All(&pageDoc)
		if err != nil {
			return pageDoc, 0, err
		}
	} else {
		err = col.Find(query).Limit(limit).Skip(offset).Sort("createdAt").All(&pageDoc)
	}
	if err != nil {
		return pageDoc, 0, err
	}
	count, err := col.Find(query).Count()
	return pageDoc, count, err
}

func (userModel *User) Update(id string, userDoc *User) error {

	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	log.Printf("%v", userDoc)
	err := col.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": userDoc})
	return err
}
func (userModel *User) Get(id string) (*User, error) {
	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var userDoc User
	err := col.FindId(bson.ObjectIdHex(id)).One(&userDoc)

	return &userDoc, err
}

func (userModel *User) Delete(id string) error {

	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	err := col.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}

func (userModel *User) UpdateField(id bson.ObjectId, updatePayload bson.M) error {

	col, session := GetCollection(CollectionNameUser)
	defer session.Close()

	err := col.Update(bson.M{"_id": id}, updatePayload)
	return err
}

func (userModel *User) GetBecas(id string) ([]Beca, error) {
	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var userDoc User
	err := col.FindId(bson.ObjectIdHex(id)).One(&userDoc)
	if err != nil {
		return nil, err
	}

	colBeca, session2 := GetCollection(CollectionNameBeca)
	defer session2.Close()
	var becasResp []Beca
	/*if userDoc.Rol != "admin" {
		for _, beca := range userDoc.Becas {
			var becaAux Beca
			err := colBeca.FindId(beca.Beca).One(&becaAux)
			if err != nil {
				return nil, err
			}
			becasResp = append(becasResp, becaAux)
		}
	} else {*/
	err = colBeca.Find(nil).All(&becasResp)
	if err != nil {
		return nil, err
	}
	//}
	return becasResp, err
}

func (userModel *User) FindByIdOffice(idOffice bson.ObjectId) ([]User, error) {
	col, session := GetCollection(CollectionNameUser)
	defer session.Close()

	users := []User{}

	err := col.Find(bson.M{"section._id": idOffice}).All(&users)
	return users, err
}
func (userModel *User) UpdateSection(id bson.ObjectId, updatePayload bson.M) error {
	col, session := GetCollection(CollectionNameUser)
	defer session.Close()

	err := col.Update(bson.M{"_id": id}, bson.M{"$set": updatePayload})
	return err
}
func (userModel *User) GetByOffice(idOffice string) ([]User, error) {
	col, session := GetCollection(CollectionNameUser)
	defer session.Close()

	var users []User
	err := col.Find(bson.M{"section._id": bson.ObjectIdHex(idOffice)}).All(&users)
	return users, err
}

func (userModel *User) Search(searchQuery string, limit, offset int) ([]User, error) {
	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var users []User
	// Crear un filtro de búsqueda para buscar coincidencias en varios campos.
	query := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"email": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"rut": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"section.name": bson.M{"$regex": searchQuery, "$options": "i"}},
			{"rol": bson.M{"$regex": searchQuery, "$options": "i"}},
		},
	}

	err := col.Find(query).Limit(limit).Skip(offset).All(&users)
	return users, err
}

func (userModel *User) GetByEmail(email string) (*User, error) {
	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var userDoc User
	err := col.Find(bson.M{"email": email}).One(&userDoc)
	return &userDoc, err
}
