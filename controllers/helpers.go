package controllers

import (
	"log"

	"github.com/astaxie/beego"
)

func PrepareUnauthorized(controller *beego.Controller) {
	controller.Data["json"] = map[string]string{"error": "unauthorized"}
	controller.Ctx.ResponseWriter.WriteHeader(403)
	controller.ServeJSON()
}

func PrepareBadRequest(controller *beego.Controller, err error) {
	controller.Data["json"] = map[string]string{"error": err.Error()}
	controller.Ctx.ResponseWriter.WriteHeader(400)
	controller.ServeJSON()

	if err != nil {
		log.Println("err:", err)
	}
}

func PrepareInternalServerError(controller *beego.Controller, err error) {
	controller.Data["json"] = map[string]string{"error": "internal error"}
	controller.Ctx.ResponseWriter.WriteHeader(500)
	controller.ServeJSON()

	if err != nil {
		log.Println("err:", err)
	}
}

func PrepareOK(controller *beego.Controller, msg interface{}) {
	controller.Data["json"] = msg
	controller.Ctx.ResponseWriter.WriteHeader(200)
	controller.ServeJSON()
}

type UserClaims struct {
	UID string
}

// FakeUserClaimsFromHeader Returns user claims using the uid from a header.
// Lacks any security whatsoever and is not used in production.
func FakeUserClaimsFromHeader(controller *beego.Controller) *UserClaims {
	userHeader := controller.Ctx.Input.Header("User")
	if userHeader == "" {
		return nil
	}

	return &UserClaims{UID: userHeader}
}
