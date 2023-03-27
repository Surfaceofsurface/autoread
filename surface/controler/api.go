package oapi

import (
	"encoding/json"
	"fmt"
	"io"
	"private/autoread/reactor/controler"
	"private/autoread/reactor/model"
	"private/autoread/reactor/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func X(cm *controler.ControlMan) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// rg := r.Group("/api")
	r.GET("/api/alluser", func(c *gin.Context) {
		s := map[string]*model.TaskMaker{}
		q := cm.Tmk
		for uid, v := range q {
			t, _ := v.(*service.TaskMakerS)
			s[uid] = t.TMaker
		}
		c.JSON(200, s)
	})
	r.POST("/api/user", func(c *gin.Context) {
		defer c.Request.Body.Close()
		b, e := io.ReadAll(c.Request.Body)
		if e != nil {
			c.JSON(500, e)
		}
		userdto := &LoginDTO{}
		e = json.Unmarshal(b, &userdto)
		if e != nil {
			c.JSON(500, e)
			return
		}
		if e != nil {
			c.JSON(500, e)
			return
		}
		u := &service.TaskMakerS{
			TMaker: &model.TaskMaker{
				Password: userdto.Psw,
				Username: userdto.Uid,
				SchoolID: "1cdceffd0000020bce",
			},
		}

		cm.RegisterUser(u, userdto.Conc, userdto.Gap)
		xx := []string{}
		for _, b := range u.TMaker.PendingBook {
			xx = append(xx, b.BookID)
		}
		cm.AddTask(userdto.Uid, xx...)
	})
	r.GET("/api/alltask", func(c *gin.Context) {
		s := map[string][]model.TProcess{}
		q := cm.GetAllProcess()
		for uid, v := range q {
			s[uid] = []model.TProcess{}
			for _, ea := range v {
				eax, _ := ea.(*service.TasksS)
				s[uid] = append(s[uid], *eax.TP)
			}

		}
		c.JSON(200, s)
	})
	r.Static("/h", "./dist")
	r.GET("/api/per", func(c *gin.Context) {

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		for {
			err := conn.WriteJSON(<-cm.Poster)
			if err != nil {
				conn.Close()
				return
			}

		}
	})
	r.Run(":44551")
}

type LoginDTO struct {
	Conc int    `json:"conc"`
	Psw  string `json:"psw"`
	Uid  string `json:"uid"`
	Gap  int    `json:"gap"`
}
type PassUser struct {
}
