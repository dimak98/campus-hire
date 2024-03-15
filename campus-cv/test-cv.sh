#!/bin/bash

# Define the URL of your Flask application
FLASK_APP_URL="http://localhost:3000/generate-cv"

# Define the user details as JSON
USER_DETAILS='{"user_id":"1","user_details":{"id":"1","email":"kazindmitrey@gmail.com","fname":"Dima Kazin","role":"student","is_cv_created":"false","cv_path":"null","description":"hello","profileImage":"/usr/src/app/templates/static/assets/img/users/1-profile.jpeg","jobs":[{"title":"samurai","company":"shauli","startDate":"01, 1984","endDate":"01, 2005","description":"cutting it"}],"education":[{"school":"shauli","degree":"PhD","fieldOfStudy":"English Language And Literature","startDate":"02, 1981","endDate":"02, 1983","description":"hhh"}]}}'

# Use curl to send a POST request to the Flask application
curl -X POST "$FLASK_APP_URL" \
     -H "Content-Type: application/json" \
     -d "$USER_DETAILS" \
     --output student_cv_1.pdf

echo "CV generated and saved as student_cv_1.pdf"