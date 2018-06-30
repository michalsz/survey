package setup

import (
	"rosella/controllers"
	"rosella/dbhandler"

	"github.com/astaxie/beego"
)

func Setup() {
	if err := dbhandler.ConnectToDB(); err != nil {
		panic(err)
	}

	beego.BConfig.CopyRequestBody = true
	beego.Router("/survey/:qid", &controllers.SurveyController{},
		"post:WriteQuestionResponse")
}
