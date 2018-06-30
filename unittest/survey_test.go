package unittest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rosella/dbhandler"
	"testing"

	"github.com/astaxie/beego"
)

type surveyQuestionAnswer struct {
	Rating int
	SID    string
}

func postSurveyQuestionAnswer(uid string, questionID int,
	ans surveyQuestionAnswer) error {
	requestBytes, err := json.Marshal(ans)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST",
		fmt.Sprintf("/survey/%d", questionID),
		bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	r.Header.Add("User", uid)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
			return err
		}

		return errors.New(errorResponse.Error)
	}

	return nil
}

func clearDBA() {
	sqlQuery := `BEGIN;
			DELETE FROM public.survey_question_ratings;
			DELETE FROM public.session;
			DELETE FROM public.user;
			DELETE FROM public.survey_questions;
            DELETE FROM public.multiple_choice_answers;
			END;`
	if _, err := dbhandler.GetDatabase().Exec(sqlQuery); err != nil {
		fmt.Println(err)
	}
}

func TestPostSurveyQuestionAnswer(t *testing.T) {
	clearDBA()
	sqlQuery := `BEGIN;
		INSERT INTO public.survey_questions
			(id, question, date_added, answer_type)
			VALUES (1, 'Test Answered Question', NOW(), 'scale'),
				   (2, 'Test Unanswered Question', NOW(), 'binary');

		INSERT INTO public.user (id, uid, name, email, date_created, admin)
			VALUES (0, 'testUID', 'Test User', 'test@user.com', NOW(), FALSE);

		INSERT INTO public.session (id, uid, public, date_created)
			VALUES ('testSID', 'testUID', FALSE, NOW());

		INSERT INTO public.survey_question_ratings
			(id, sid, uid, questionid, date_modified, rating, answer_type)
			VALUES (1, 'testSID', 'testUID', 1, NOW(), 3, 'scale');
		END;`

	if _, err := dbhandler.GetDatabase().Exec(sqlQuery); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		sqlQuery := `BEGIN;
			DELETE FROM public.survey_question_ratings;
			DELETE FROM public.session;
			DELETE FROM public.user;
			DELETE FROM public.survey_questions;
			END;`
		if _, err := dbhandler.GetDatabase().Exec(sqlQuery); err != nil {
			t.Error(err)
		}
	}()

	// valid scale update
	if err := postSurveyQuestionAnswer("testUID", 1,
		surveyQuestionAnswer{
			Rating: 3,
			SID:    "testSID",
		}); err != nil {
		t.Error(err)
		return
	}

	// invalid scale update
	if err := postSurveyQuestionAnswer("testUID", 1,
		surveyQuestionAnswer{
			Rating: 10,
			SID:    "testSID",
		}); err.Error() != "survey: only numbers between 0 and 5 are valid answers to a scale question" {
		t.Error(fmt.Errorf("did not get the right error on invalid rating: %v", err))
		return
	}

	// valid binary create
	if err := postSurveyQuestionAnswer("testUID", 2,
		surveyQuestionAnswer{
			Rating: 1,
			SID:    "testSID",
		}); err != nil {
		t.Error(err)
		return
	}

	// invalid binary update
	if err := postSurveyQuestionAnswer("testUID", 2,
		surveyQuestionAnswer{
			Rating: 3,
			SID:    "testSID",
		}); err.Error() != "survey: only 0 and 1 are valid answers to a binary question" {
		t.Error(fmt.Errorf("did not get the right error on invalid rating: %v", err))
		return
	}
}

type surveyMultipleChoiceAnswer struct {
	Answers []int
	Freetext string
	SID    string
}

type surveyFreeTextAnswer struct {
	Text string
}

func postSurveyQuestionMultiAnswer(uid string, questionID int,
	ans surveyMultipleChoiceAnswer) error {
	requestBytes, err := json.Marshal(ans)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST",
		fmt.Sprintf("/survey/%d", questionID),
		bytes.NewBuffer(requestBytes))
	if err != nil {
		return err
	}

	r.Header.Add("User", uid)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 200 {
		var errorResponse struct {
			Error string `json:"error"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
			return err
		}

		return errors.New(errorResponse.Error)
	}

	return nil
}

func TestPostSurveyQuestionMultiAnswer(t *testing.T) {
	clearDBA()
	sqlQuery := `BEGIN;
		INSERT INTO public.survey_questions
			(id, question, date_added, answer_type)
			VALUES (1, 'Test Answered Question', NOW(), 'multiple-choice-freetext'),
				   (2, 'Test Unanswered Question', NOW(), 'binary');

		INSERT INTO public.user (id, uid, name, email, date_created, admin)
			VALUES (0, 'testUID', 'Test User', 'test@user.com', NOW(), FALSE);

		INSERT INTO public.session (id, uid, public, date_created)
			VALUES ('testSID', 'testUID', FALSE, NOW());

		INSERT INTO public.survey_question_ratings
			(id, sid, uid, questionid, date_modified, rating, answer_type)
			VALUES (1, 'testSID', 'testUID', 1, NOW(), 3, 'multiple-choice-freetext');

        INSERT INTO public.multiple_choice_answers (id, category, name)
			VALUES (1, 'cars', 'Best car');
        INSERT INTO public.multiple_choice_answers (id, category, name)
			VALUES (2, 'cars', 'Best super car');
		END;`

	if _, err := dbhandler.GetDatabase().Exec(sqlQuery); err != nil {
		t.Error(err)
		return
	}

	//valid save multi answers
	var choices []int
	choices = append(choices, 1, 2)
	if err := postSurveyQuestionMultiAnswer("testUID", 1,
		surveyMultipleChoiceAnswer{SID: "testSID", Answers: choices, Freetext: "Free text",}, ); err != nil {
		t.Error(err)
		return
	}

	// empty answers
	var emptyChoices []int
	if err := postSurveyQuestionMultiAnswer("testUID", 1,
		surveyMultipleChoiceAnswer{
			SID: "testSID",
			Answers: emptyChoices,
			Freetext: "Free text",
	}); err.Error() != "survey: no multiple answers" {
		t.Error(fmt.Errorf("did not get the right error on missing answers: %v", err))
		return
	}

	// invalid answers
	var invalidChoices []int
	invalidChoices = append(invalidChoices, 6, 7)
	if err := postSurveyQuestionMultiAnswer("testUID", 1,
		surveyMultipleChoiceAnswer{
			SID: "testSID",
			Answers: invalidChoices,
			Freetext: "Free text",
		}); err.Error() != "survey: invalid answers" {
		t.Error(fmt.Errorf("did not get the right error on missing answers: %v", err))
		return
	}

	//invalid - no free text
	if err := postSurveyQuestionMultiAnswer("testUID", 1,
		surveyMultipleChoiceAnswer{
			SID: "testSID",
			Answers: choices,
			Freetext: "",
		}); err.Error() != "survey: no free text" {
		t.Error(fmt.Errorf("did not get the right error on missing free text: %v", err))
		return
	}
}