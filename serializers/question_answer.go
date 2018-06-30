package serializers

import (
	"encoding/json"
	"errors"
	"rosella/models"
)

// AnswerParserForQuestion Returns the parser for the given question.
func AnswerParserForQuestion(qid int, uid string) (AnswerParser, error) {
	answerType, err := models.QuestionAnswerType(qid)
	if err != nil {
		return nil, err
	}

	if answerType == "scale" || answerType == "binary" {
		return &RatingSIDAnswerParser{
			UID:        uid,
			QuestionID: qid,
			AnswerType: answerType,
		}, nil
	}

	if answerType == "multiple-choice-freetext" {
		return &MultipleChoiceFreetextAnswerParser{
			UID:        uid,
			QuestionID: qid,
			AnswerType: answerType,
		}, nil
	}

	return nil, errors.New("survey: unsupported answer type")
}

// AnswerParser Parses a survey question answer.
type AnswerParser interface {
	Parse(jsonBytes []byte) (models.Answer, error)
}

// RatingSIDAnswerParser Parses a RatingSID question answer.
type RatingSIDAnswerParser struct {
	UID        string
	QuestionID int
	AnswerType string
}

// Parse Parses a RatingSID question answer.
func (p *RatingSIDAnswerParser) Parse(jsonBytes []byte) (models.Answer, error) {
	var parsed struct {
		Rating int
		SID    string
	}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		return nil, err
	}

	switch p.AnswerType {
	case "scale":
		if parsed.Rating < 0 || parsed.Rating > 5 {
			return nil, errors.New(
				"survey: only numbers between 0 and 5 are valid answers to a " +
					"scale question")
		}
	case "binary":
		if parsed.Rating != 0 && parsed.Rating != 1 {
			return nil, errors.New(
				"survey: only 0 and 1 are valid answers to a binary question")
		}
	}

	if parsed.SID == "" {
		return nil, errors.New("survey: session id is required")
	}

	return &models.RatingSIDAnswer{
		UID:        p.UID,
		QuestionID: p.QuestionID,
		SID:        parsed.SID,
		Rating:     parsed.Rating,
		AnswerType: p.AnswerType,
	}, nil
}

// MultipleChoiceFreetextAnswerParser Parses a multiple-choice or a custom answer.
type MultipleChoiceFreetextAnswerParser struct {
	UID        string
	QuestionID int
	AnswerType string
	SID    string
}

// Parse Parses a multiple-choice or a custom answer.
func (p *MultipleChoiceFreetextAnswerParser) Parse(jsonBytes []byte) (models.Answer, error) {
	var parsed struct {
		Answers []int
		Freetext string
		SID string
	}

	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Answers) == 0 {
		return nil, errors.New(
			"survey: no multiple answers")
	}

	if parsed.Freetext == "" {
		return nil, errors.New("survey: no free text")
	}

	if parsed.SID == "" {
		return nil, errors.New("survey: session id is required")
	}

	if !models.ValidAnswers(parsed.Answers) {
		return nil, errors.New("survey: invalid answers")
	}

	return &models.RatingSIDMultiAnswer{
		UID:        p.UID,
		SID:        parsed.SID,
		QuestionID: p.QuestionID,
		AnswerType: p.AnswerType,
		Answers:    parsed.Answers,
		Freetext:   parsed.Freetext,
	}, nil
}
