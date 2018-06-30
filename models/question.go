package models

import (
	"database/sql"
	"errors"
	"rosella/dbhandler"
	"github.com/lib/pq"
)

// QuestionAnswerType Returns the answer type for a given question.
func QuestionAnswerType(qid int) (string, error) {
	var answerType string
	if err := dbhandler.GetDatabase().QueryRow(
		"SELECT answer_type FROM survey_questions WHERE id = $1",
		qid).Scan(&answerType); err == sql.ErrNoRows {
		return "", errors.New("survey: non-existent question")
	} else if err != nil {
		return "", err
	}
	return answerType, nil
}

// Answer Saveable answer to a survey question.
type Answer interface {
	Save() error
}

// RatingSIDAnswer Binary (0 or 1) or scale (0 - 5) answer for a Session.
type RatingSIDAnswer struct {
	UID        string
	QuestionID int
	SID        string
	Rating     int
	AnswerType string
}

// Save Saves the RatingSIDAnswer.
func (a *RatingSIDAnswer) Save() error {
	sqlQuery := `INSERT INTO survey_question_ratings (sid, uid, date_modified,
		 questionid, rating, answer_type) VALUES ($1, $2, CURRENT_TIMESTAMP,
		$3, $4, $5)
	ON CONFLICT (id) DO UPDATE
		SET (sid, uid, date_modified, questionid, rating, answer_type) =
			($1, $2, CURRENT_TIMESTAMP, $3, $4, $5)`

	_, err := dbhandler.GetDatabase().Exec(sqlQuery,
		a.SID, a.UID, a.QuestionID, a.Rating, a.AnswerType)
	return err
}

// RatingSIDMultiAnswer Binary (0 or 1) or scale (0 - 5) answer for a Session.
type RatingSIDMultiAnswer struct {
	UID        string
	QuestionID int
	AnswerType string
	Answers  []int
	Freetext string
	SID      string
}

// Save Saves the RatingSIDMultiAnswer.
func (a *RatingSIDMultiAnswer) Save() error {
	sqlQuery := `INSERT INTO survey_question_ratings (sid, uid, date_modified,
		 questionid, answer_type, freetext, answer_ids) VALUES ($1, $2, CURRENT_TIMESTAMP,
		$3, $4, $5, $6)`

	_, err := dbhandler.GetDatabase().Exec(sqlQuery,
		a.SID, a.UID, a.QuestionID, "multiple-choice-freetext", a.Freetext, pq.Array(a.Answers))

	//fmt.Println(a.SID)
	//fmt.Println(a.UID)
	//fmt.Println(a.QuestionID)
	//fmt.Println(a.Freetext)
	//fmt.Println(a.Answers)
	//fmt.Println(a.AnswerType)

	return err
}

func ValidAnswers(answerIDs []int) bool {
	var validAnswers = true
	var id string
	for i := 0; i < len(answerIDs); i++ {
		if err := dbhandler.GetDatabase().QueryRow(
			"SELECT * FROM multiple_choice_answers WHERE id = $1",
			answerIDs[i]).Scan(&id); err == sql.ErrNoRows {
			if err.Error() == "sql: no rows in result set" {
              validAnswers = false
			}

		}

	}
	return validAnswers
}

