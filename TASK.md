# Background
The app has a survey functionality. It currently supports `scale` (0 - 5) and `binary` (0 or 1) answer types. Both of these require a Session ID (SID)
as well.

We want to add a new type `multiple-choice-freetext`. The new type will not
require an SID, but will need either a list of answers (integer ids) or a
free-text answer.

# Steps
1. Create a new table to keep track of answers.
```
CREATE TABLE IF NOT EXISTS public.multiple_choice_answers (
    id INTEGER NOT NULL PRIMARY KEY,
    category varchar(255) NOT NULL,
    name varchar(255) NOT NULL
)
```

2. Add required fields to `survey_question_ratings` (and write me the sql query
 into a file)
    - `answer_ids` array
    - `freetext` text (not varchar)

3. Implement the `MultipleChoiceFreetextAnswerParser` parser under `serializers/question_answer.go`. It needs to accept two types of input
    1. `{"Answers": [1, 5, 3]}` where `[1, 5, 3]` are `multiple_choice_answers`
        - this should also call a function in the models package and check
          that the answers exist
    2. `{"Freetext": "some text"}` where "some text" is the text to save under
        freetext

4. Implement the model similar to how `RatingSIDAnswer` is and return the
    appropriate struct from the parser in serializers.

5. I hope I didn't leave anything out. If something doesn't make sense or you
    have some questions, don't hesitate to ask me anything.

6. Write a test. Actually this should be done as the first task, but I wanted
    you to first play with the code. Unless you want to do it as the first task,
    in that case you're more than welcome to. :)