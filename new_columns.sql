alter table public.survey_question_ratings  ADD COLUMN freetext text NOT NULL;
alter table public.survey_question_ratings  ADD COLUMN answer_ids integer array NOT NULL;