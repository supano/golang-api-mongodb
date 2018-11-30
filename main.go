package main

import (
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	h := &handler{
		m: session,
	}

	e.Use(middleware.Logger())
	e.GET("/api/users", h.list)
	e.GET("/api/users/:id", h.view)
	e.POST("/api/users", h.create)
	e.PUT("/api/users/:id", h.update)
	e.PATCH("/api/users/:id", h.update)
	e.DELETE("/api/users/:id", h.delete)
	e.Logger.Fatal(e.Start(":1323"))
}

type (
	handler struct {
		m *mgo.Session
	}

	user struct {
		ID       bson.ObjectId `json:"id" bson:"_id"`
		Email    string        `json:"email" bson:"email"`
		Password string        `json:"password" bson:"password"`
		Fname    string        `json:"fname" bson:"fname"`
		Lname    string        `json:"lname" bson:"lname"`
	}

	message struct {
		Message string `json:"message"`
	}
)

func (h *handler) create(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	u.ID = bson.NewObjectId()

	collection := session.Copy().DB("golang").C("users")
	if err := collection.Insert(u); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}

func (h *handler) list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	collection := session.Copy().DB("golang").C("users")
	var userlist []user
	if err := collection.Find(nil).All(&userlist); err != nil {
		return c.JSON(http.StatusNotFound, message{Message: "error when find users"})
	}

	if userlist == nil {
		return c.JSON(http.StatusNotFound, message{Message: "found 0 user in database"})
	}

	return c.JSON(http.StatusOK, userlist)
}

func (h *handler) view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	collection := session.Copy().DB("golang").C("users")
	id := bson.ObjectIdHex(c.Param("id"))
	var u user
	if err := collection.FindId(id).One(&u); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, u)
}

func (h *handler) update(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	collection := session.Copy().DB("golang").C("users")
	id := bson.ObjectIdHex(c.Param("id"))

	n := new(user)
	c.Bind(n)

	var u user
	if err := collection.FindId(id).One(&u); err != nil {
		return err
	}

	if err := collection.UpdateId(id, n); err != nil {
		return err
	}

	var nu user
	if err := collection.FindId(id).One(&nu); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nu)

}

func (h *handler) delete(c echo.Context) error {
	sesstion := h.m.Copy()
	defer sesstion.Close()
	collection := sesstion.Copy().DB("golang").C("users")
	id := bson.ObjectIdHex(c.Param("id"))

	if err := collection.RemoveId(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, true)
}
