import logging
import os
from flask import Flask, render_template, request, redirect, url_for, flash, session, jsonify, send_from_directory
import requests
import json
import uuid
import pandas as pd
from werkzeug.utils import secure_filename


#######################################################################################
#                                  Config                                             #
#######################################################################################


logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s %(levelname)s %(name)s %(threadName)s : %(message)s')

UPLOAD_FOLDER = '/usr/src/app/templates/static/assets/img/users'
CV_UPLOAD_FOLDER = '/usr/src/app/templates/static/assets/pdf/cv'
months = [('January', '01'), ('February', '02'), ('March', '03'), ('April', '04'),
            ('May', '05'), ('June', '06'), ('July', '07'), ('August', '08'),
            ('September', '09'), ('October', '10'), ('November', '11'), ('December', '12')]
years = range(1980, 2024)


app = Flask(__name__, static_folder='templates/static')
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
app.config['CV_UPLOAD_FOLDER'] = CV_UPLOAD_FOLDER
app.secret_key = 'ed0a9528-5906-457e-a373-f288b5e42579'
backend_url = 'http://api:8080'
cv_url = 'http://cv:3000'

#######################################################################################
#                                  Helpers                                            #
#######################################################################################

def allowed_file(filename):
    ALLOWED_EXTENSIONS = {'png', 'jpg', 'jpeg', 'gif', 'mp4', 'avi'}
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

def load_majors():
    df = pd.read_csv('fields-of-study.csv')
    majors = df['Major'].dropna().str.title().unique().tolist()
    majors.sort()
    return majors

def load_country_cities():
    df = pd.read_csv('worldcities.csv')
    df_sorted = df.sort_values(by=['country', 'city'])
    country_cities = [{'city': row['city'], 'country': row['country']} for index, row in df_sorted.iterrows()]
    return country_cities

#######################################################################################
#                                  Auth Routes                                        #
#######################################################################################

@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        fname = request.form.get('fname')
        email = request.form.get('email')
        password = request.form.get('password')
        response = requests.post(f"{backend_url}/register", json={"fname": fname, "email": email, "password": password})
        
        app.logger.debug(f'Registration attempt for {email} with response {response.status_code}')
        
        if response.status_code == 201:
            flash('Registration successful. Please check your email to verify your account.', 'success')
            return render_template('registration_success.html')
        else:
            flash('Registration failed. Please try again.', 'danger')
    return render_template('register.html')

@app.route('/verify_email')
def verify_email():
    token = request.args.get('token')
    if not token:
        flash('Verification token is missing.', 'danger')
        return redirect(url_for('login'))

    response = requests.get(f"{backend_url}/verify_email?token={token}")
    if response.status_code == 200:
        flash('Your email has been successfully verified.', 'success')
    else:
        flash('Email verification failed.', 'danger')

    return redirect(url_for('login'))

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        email = request.form['email']
        password = request.form['password']
        # Sending login credentials to the backend
        response = requests.post(f"{backend_url}/login", json={"email": email, "password": password})
        
        app.logger.debug(f"Login response: {response.text}")  # Log the raw response text
        
        if response.status_code == 200:
            data = response.json()
            if data.get('success'):

                if not data.get('isVerified'):
                   return render_template('registration_success.html')

                session['logged_in'] = True
                session['user_id'] = data.get('userId')

                app.logger.debug(f"User {session['user_id']} logged in, hasSelectedRole: {data.get('hasSelectedRole')}")

                if not data.get('hasSelectedRole', True):
                    return redirect(url_for('role_selection'))
                else:
                    return redirect(url_for('main_dashboard'))
            else:
                flash('Login failed. Please check your credentials.', 'danger')
        else:
            flash('Login failed. Please check your credentials.', 'danger')
    return render_template('login.html')


@app.route('/logout')
def logout():
    session.pop('logged_in', None)
    flash('You have been logged out.', 'info')
    return redirect(url_for('login'))

@app.route('/forgot_password', methods=['GET', 'POST'])
def forgot_password():
    if request.method == 'POST':
        email = request.form['email']
        response = requests.post(f"{backend_url}/forgot_password", json={"email": email})
        if response.status_code == 200:
            flash('Please check your email for the password reset link.', 'info')
            return render_template('reset_password_requested.html')
        else:
            flash('Error sending password reset email. Please try again.', 'danger')
    return render_template('forgot_password.html')

@app.route('/change_password', methods=['GET', 'POST'])
def change_password():
    if request.method == 'POST':
        token = request.form['token']
        newPassword = request.form['newPassword']
        response = requests.post(f"{backend_url}/change_password", json={"token": token, "newPassword": newPassword})
        if response.status_code == 200:
            flash('Your password has been changed successfully.', 'success')
            return redirect(url_for('login'))
        else:
            flash('Failed to change password. Please try again.', 'danger')
    else:
        token = request.args.get('token')
        if not token:
            flash('Reset token is missing.', 'danger')
            return redirect(url_for('forgot_password'))
    return render_template('change_password.html', token=token)


#######################################################################################
#                        Student & Company Registration                               #
#######################################################################################

@app.route('/role_selection')
def role_selection():
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))
    return render_template('role_selection.html')

@app.route('/student_registration', methods=['GET', 'POST'])
def student_registration():
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))
    if request.method == 'POST':
        user_id = session.get('user_id')
        description = request.form['description']
        profile_image = request.files['profile_image']

        # Upload handling for profile image
        if profile_image and allowed_file(profile_image.filename):
            # Rename the file to <user_id>-profile.extension
            file_ext = profile_image.filename.rsplit('.', 1)[1].lower()
            filename = f"{user_id}-profile.{file_ext}"
            filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
            profile_image.save(filepath)
            image_path = filepath
        else:
            image_path = None

        # Aggregate job information
        jobs = [
            {
                'title': title,
                'company': company,
                'startDate': f"{startMonth}, {startYear}" if startMonth and startYear else None,
                'endDate': f"{endMonth}, {endYear}" if endMonth and endYear else None,
                'description': job_description
            }
            for title, company, startMonth, startYear, endMonth, endYear, job_description in zip(
                request.form.getlist('jobs[][title]'),
                request.form.getlist('jobs[][company]'),
                request.form.getlist('jobs[][startMonth]'),
                request.form.getlist('jobs[][startYear]'),
                request.form.getlist('jobs[][endMonth]'),
                request.form.getlist('jobs[][endYear]'),
                request.form.getlist('jobs[][description]')
            )
        ]

        # Aggregate education information
        education = [
            {
                'school': school,
                'degree': degree,
                'fieldOfStudy': fieldOfStudy,
                'startDate': f"{startMonth}, {startYear}" if startMonth and startYear else None,
                'endDate': f"{endMonth}, {endYear}" if endMonth and endYear else None,
                'description': edu_description
            }
            for school, degree, fieldOfStudy, startMonth, startYear, endMonth, endYear, edu_description in zip(
                request.form.getlist('education[][school]'),
                request.form.getlist('education[][degree]'),
                request.form.getlist('education[][fieldOfStudy]'),
                request.form.getlist('education[][startMonth]'),
                request.form.getlist('education[][startYear]'),
                request.form.getlist('education[][endMonth]'),
                request.form.getlist('education[][endYear]'),
                request.form.getlist('education[][description]')
            )
        ]

        # Create the payload
        payload = {
            'userId': user_id,
            'description': description,
            'imagePath': image_path,
            'jobs': jobs,
            'education': education
        }

        headers = {'Content-Type': 'application/json'}
        response = requests.post(f'{backend_url}/student_registration', headers=headers, data=json.dumps(payload))

        if response.status_code == 200:
            flash('Student registration successful!', 'success')
            return redirect(url_for('main_dashboard'))
        else:
            flash('Failed to register student. Please try again.', 'danger')
    
    majors = load_majors()
    return render_template('student_registration.html', majors=majors, months=months, years=years)

@app.route('/company_registration', methods=['GET', 'POST'])
def company_registration():
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))
    if request.method == 'POST':
        user_id = session.get('user_id')
        # Process form fields
        form_data = {
            'userId': user_id,
            'name': request.form.get('company_name'),
            'size': request.form.get('company_size'),
            'address': request.form.get('address'),
            'description': request.form.get('description')
        }

        files = {}
        for file_field in ['company_image', 'video']:
            file = request.files.get(file_field)
            if file and allowed_file(file.filename):
                # Rename the file to <user_id>-<field>.extension
                file_ext = file.filename.rsplit('.', 1)[1].lower()
                filename = f"{user_id}-{file_field}.{file_ext}"
                file_path = os.path.join(app.config['UPLOAD_FOLDER'], filename)
                file.save(file_path)
                files[file_field] = file_path

        # Add file paths to form data if available
        if 'company_image' in files:
            form_data['image_path'] = files['company_image']
        if 'video' in files:
            form_data['video_path'] = files['video']

        # Make a request to the backend
        headers = {'Content-Type': 'application/json'}
        response = requests.post(f'{backend_url}/company_registration', json=form_data, headers=headers)

        if response.status_code == 200:
            flash('Company registration successful!', 'success')
            return redirect(url_for('main_dashboard'))
        else:
            flash('Failed to register company. Please try again.', 'danger')

    country_cities = load_country_cities()
    return render_template('company_registration.html', country_cities=country_cities)

#######################################################################################
#                                         User Profiles                               #
#######################################################################################

@app.route('/profile')
def profile():
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))
    
    user_id = session.get('user_id')
    role_response = requests.get(f'{backend_url}/user_role?userID={user_id}')
    if role_response.status_code == 200:
        role = role_response.json().get('role')
        response = requests.get(f'{backend_url}/user_details?userID={user_id}')
        if response.status_code == 200:
            user_details = response.json()
            if role == 'student':
                if 'profileImage' in user_details:
                    user_details['profileImage'] = user_details['profileImage'].split('/static/', 1)[-1]
                return render_template('student_profile.html', user_details=user_details)
            elif role == 'company':
                if 'image_path' in user_details:
                    user_details['image_path'] = user_details['image_path'].split('/static/', 1)[-1]
                if 'video_path' in user_details:
                    user_details['video_path'] = user_details['video_path'].split('/static/', 1)[-1]
                return render_template('company_profile.html', user_details=user_details)
            else:
                flash('Invalid user role.', 'danger')
                return redirect(url_for('login'))
        else:
            flash(f"Error fetching user details: {response.text}", 'danger')
            return redirect(url_for('login'))
    else:
        flash(f"Error fetching user role: {role_response.text}", 'danger')
        return redirect(url_for('login'))


#######################################################################################
#                              Main App Routes                                        #
#######################################################################################

@app.route('/')
def main_dashboard():
    jobs = None  # Initialize jobs variable
    jobs_response = requests.get(f'{backend_url}/jobs?latest=true')
    if jobs_response.status_code == 200:
        jobs = jobs_response.json()
    user_details = None
    if session.get('logged_in'):
        user_id = session.get('user_id')
        response = requests.get(f'{backend_url}/user_details?userID={user_id}')
        if response.status_code == 200:
            user_details = response.json()
            if user_details.get('role') == 'student':
                if 'profileImage' in user_details:
                    user_details['profileImage'] = user_details['profileImage'].split('/static/', 1)[-1]
            elif user_details.get('role') == 'company':
                if 'image_path' in user_details:
                    user_details['image_path'] = user_details['image_path'].split('/static/', 1)[-1]
                if 'video_path' in user_details:
                    user_details['video_path'] = user_details['video_path'].split('/static/', 1)[-1]
    else:
        return render_template('main.html', jobs=jobs)

    return render_template('index.html', user_details=user_details, jobs=jobs)

@app.route('/jobs')
def jobs():
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))

    jobs = None  # Initialize jobs variable
    user_id = session.get('user_id')
    user_details_response = requests.get(f'{backend_url}/user_details?userID={user_id}')
    if user_details_response.status_code == 200:
        user_details = user_details_response.json()
        role = user_details.get('role')
        if role == 'student':
            if 'profileImage' in user_details:
                user_details['profileImage'] = user_details['profileImage'].split('/static/', 1)[-1]
        elif role == 'company':
            if 'image_path' in user_details:
                user_details['image_path'] = user_details['image_path'].split('/static/', 1)[-1]
            if 'video_path' in user_details:
                user_details['video_path'] = user_details['video_path'].split('/static/', 1)[-1]

    jobs_response = requests.get(f'{backend_url}/jobs?latest=true')
    if jobs_response.status_code == 200:
        jobs = jobs_response.json()

    return render_template('jobs.html', user_details=user_details, jobs=jobs)


#######################################################################################
#                              Company Specific Routes                                #
#######################################################################################

@app.route('/company/<int:userID>', methods=['GET'])
def company_view(userID):
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))

    response = requests.get(f'{backend_url}/company?userID={userID}')
    if response.status_code == 200:
        company_details = response.json()
        if 'image_path' in company_details:
            company_details['image_path'] = company_details['image_path'].split('/static/', 1)[-1]
        if 'video_path' in company_details:
            company_details['video_path'] = company_details['video_path'].split('/static/', 1)[-1]
    else:
        return redirect(url_for('main_dashboard'))
    
    session_user_id = session.get('user_id')
    session_response = requests.get(f'{backend_url}/user_details?userID={session_user_id}')
    if session_response.status_code == 200:
        user_details = session_response.json()
        role = user_details.get('role')
        if role == 'student':
            if 'profileImage' in user_details:
                user_details['profileImage'] = user_details['profileImage'].split('/static/', 1)[-1]
        elif role == 'company':
            if 'image_path' in user_details:
                user_details['image_path'] = user_details['image_path'].split('/static/', 1)[-1]
            if 'video_path' in user_details:
                user_details['video_path'] = user_details['video_path'].split('/static/', 1)[-1]
    else:
        return redirect(url_for('main_dashboard'))
    
    return render_template('profile_view_company.html', company_details=company_details, user_details=user_details)

#######################################################################################
#                                  Student Specific Routes                            #
#######################################################################################

@app.route('/student/<int:userID>', methods=['GET'])
def student_view(userID):
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))

    response = requests.get(f'{backend_url}/student?userID={userID}')
    if response.status_code == 200:
        student_details = response.json()
        if 'profileImage' in student_details:
            student_details['profileImage'] = student_details['profileImage'].split('/static/', 1)[-1]

    else:
        return redirect(url_for('main_dashboard'))
    
    session_user_id = session.get('user_id')
    session_response = requests.get(f'{backend_url}/user_details?userID={session_user_id}')
    if session_response.status_code == 200:
        user_details = session_response.json()
        role = user_details.get('role')
        if role == 'student':
            if 'profileImage' in user_details:
                user_details['profileImage'] = user_details['profileImage'].split('/static/', 1)[-1]
        elif role == 'company':
            if 'image_path' in user_details:
                user_details['image_path'] = user_details['image_path'].split('/static/', 1)[-1]
            if 'video_path' in user_details:
                user_details['video_path'] = user_details['video_path'].split('/static/', 1)[-1]
    else:
        return redirect(url_for('main_dashboard'))
    
    return render_template('profile_view_student.html', student_details=student_details, user_details=user_details)

@app.route('/generate-student-cv', methods=['POST'])
def generate_student_cv():
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))

    user_id = session.get('user_id')
    response = requests.get(f'{backend_url}/user_details?userID={user_id}')

    if response.status_code == 200:
        user_details = response.json()

        cv_response = requests.post(
            f'{cv_url}/generate-cv',
            json={'user_id': user_id, 'user_details': user_details}
        )

        if cv_response.status_code == 200:
            cv_filename = f'student_cv_{user_id}.pdf'
            cv_filepath = os.path.join(CV_UPLOAD_FOLDER, cv_filename)
            if not os.path.exists(CV_UPLOAD_FOLDER):
                os.makedirs(CV_UPLOAD_FOLDER)
            with open(cv_filepath, 'wb') as cv_file:
                cv_file.write(cv_response.content)
            flash('Your CV has been generated successfully!', 'success')
            return send_from_directory(CV_UPLOAD_FOLDER, cv_filename, as_attachment=True)
        else:
            flash('Failed to generate CV. Please try again.', 'danger')
    else:
        flash('Failed to fetch user details. Please try again.', 'danger')

    return redirect(url_for('profile'))

#######################################################################################
#                                  Job Specific Routes                                #
#######################################################################################

@app.route('/post_job', methods=['GET','POST'])
def post_job():
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))

    if request.method == 'POST':
        user_id = session.get('user_id')
        title = request.form.get('title')
        salary = request.form.get('salary')
        location = request.form.get('location')
        job_description = request.form.get('job_description')
        requirements = request.form.get('requirements')

        payload = {
            'user_id': user_id,
            'title': title,
            'salary': salary,
            'address': location,
            'description': job_description,
            'requirements': requirements,
            'status': 'Open'
        }

        headers = {'Content-Type': 'application/json'}
        response = requests.post(f'{backend_url}/post_job', json=payload, headers=headers)

        if response.status_code == 201:
            flash('Job posted successfully!', category='success')
            return redirect(url_for('main_dashboard'))
        else:
            flash('Failed to post job. Please try again.', category= 'error')
            return render_template("new_job.html")

    country_cities = load_country_cities()
    user_id = session.get('user_id')
    app.logger.info(f"Posting new job attempt for user with id {user_id}")
    role_response = requests.get(f'{backend_url}/user_role?userID={user_id}')
    if role_response.status_code == 200:
        role = role_response.json().get('role')
        response = requests.get(f'{backend_url}/user_details?userID={user_id}')
        if response.status_code == 200:
            user_details = response.json()
            if role == 'student':
                if 'profileImage' in user_details:
                    user_details['profileImage'] = user_details['profileImage'].split('/static/', 1)[-1]
                return render_template('new_job.html', user_details=user_details, country_cities=country_cities)
            elif role == 'company':
                if 'image_path' in user_details:
                    user_details['image_path'] = user_details['image_path'].split('/static/', 1)[-1]
                if 'video_path' in user_details:
                    user_details['video_path'] = user_details['video_path'].split('/static/', 1)[-1]
                return render_template('new_job.html', user_details=user_details, country_cities=country_cities)
            else:
                flash('Invalid user role.', 'danger')
                return redirect(url_for('login'))
        else:
            flash(f"Error fetching user details: {response.text}", 'danger')
            return redirect(url_for('login'))
    else:
        flash(f"Error fetching user role: {role_response.text}", 'danger')
        return redirect(url_for('login'))

@app.route('/job/<int:job_id>', methods=['GET'])
def job_view(job_id):
    if not session.get('logged_in'):
        flash('Please log in to access this page.', 'warning')
        return redirect(url_for('login'))

    response = requests.get(f'{backend_url}/job?jobID={job_id}')

    session_user_id = session.get('user_id')
    session_response = requests.get(f'{backend_url}/user_details?userID={session_user_id}')
    if session_response.status_code == 200:
        user_details = session_response.json()
        role = user_details.get('role')
        if role == 'student':
            if 'profileImage' in user_details:
                user_details['profileImage'] = user_details['profileImage'].split('/static/', 1)[-1]
        elif role == 'company':
            if 'image_path' in user_details:
                user_details['image_path'] = user_details['image_path'].split('/static/', 1)[-1]
            if 'video_path' in user_details:
                user_details['video_path'] = user_details['video_path'].split('/static/', 1)[-1]

    if response.status_code == 200:
        job_details = response.json()
        return render_template('job_view.html', job_details=job_details, user_details=user_details)
    else:
        flash('Failed to fetch job details. Please try again.', 'danger')
        return redirect(url_for('main_dashboard'))

    
if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
