package controllers

import (
	"rosella/serializers"
	"strconv"

	"github.com/astaxie/beego"
)

// SurveyController Responsible for survey.
type SurveyController struct {
	beego.Controller
}

// Answer Saves a question response.
func (s *SurveyController) WriteQuestionResponse() {
	userClaims := FakeUserClaimsFromHeader(&s.Controller)
	if userClaims == nil {
		PrepareUnauthorized(&s.Controller)
		return
	}

	qidStr := s.Ctx.Input.Param(":qid")
	qid, err := strconv.Atoi(qidStr)
	if err != nil {
		PrepareBadRequest(&s.Controller, err)
		return
	}

	answerParser, err := serializers.AnswerParserForQuestion(
		qid, userClaims.UID)
	if err != nil {
		PrepareInternalServerError(&s.Controller, err)
		return
	}

	answer, err := answerParser.Parse(s.Ctx.Input.RequestBody)
	if err != nil {
		PrepareBadRequest(&s.Controller, err)
		return
	}

	if err := answer.Save(); err != nil {
		PrepareInternalServerError(&s.Controller, err)
		return
	}

	PrepareOK(&s.Controller, map[string]string{
		"Description": "Thanks for your feedback!"})
}
