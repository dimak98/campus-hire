from flask import Flask, request, send_from_directory
import anthropic
import json
from fpdf import FPDF
import os

#######################################################################################
#                                         Config                                      #
#######################################################################################

app = Flask(__name__)

ICONS_PATH = os.environ.get('ICONS_PATH', '/campus-hire/project/campus-cv/icons')
PDF_DIR = os.environ.get('PDF_DIR', '/campus-hire/project/campus-cv/cvs')
API_KEY = os.environ.get('API_KEY', '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')

#######################################################################################
#                                         Helpers                                     #
#######################################################################################

def generate_cv(user_details):
    client = anthropic.Anthropic(api_key=API_KEY)
    prompt = f"Hi! we have an application which serves as a job search platform for students.\nHere is an information about a student in the system, please generate a professional CV written in native english and the buzzwords HR loves for the following user details, please do not exxagerate and tell lies, and only make it sound better. Another thing that I need is to limit the content into one pdf page (regulrar size) and i want you to in titles: Email, Objective, Education, Experience, Skills:\n{json.dumps(user_details, indent=2)}"
    message = client.messages.create(
        model="claude-3-opus-20240229",
        max_tokens=1000,
        temperature=0,
        messages=[{"role": "user", "content": [{"type": "text", "text": prompt}]}]
    )
    cv_content = message.content[0].text
    return cv_content

class PDF(FPDF):
    def __init__(self, user_details, icons_path):
        super().__init__()
        self.user_details = user_details
        self.icons_path = icons_path

    def header(self):
        last_education = self.user_details['education'][-1] if self.user_details['education'] else {}
        title = f"{self.user_details['fname']} - {last_education.get('degree', '')} in {last_education.get('fieldOfStudy', '')} at {last_education.get('school', '')}"
        self.set_font('Arial', 'B', 12)
        self.cell(0, 10, title, 0, 1, 'C')
        self.ln(5)

    def section_title(self, title, icon=None):
        if icon:
            self.image(f"{self.icons_path}/{icon}", x=10, y=self.get_y(), w=8)
            self.set_x(self.get_x() + 10)
        self.set_fill_color(220, 220, 220)
        self.set_font('Arial', 'B', 10)
        self.cell(0, 8, title, 0, 1, 'L', 1)
        self.ln(4)

    def section_body(self, body, is_subsection=False, is_email=False):
        if is_email:
            self.set_font('Arial', '', 9)
            self.multi_cell(0, 6, body, 0, 'C')
        elif is_subsection:
            self.set_font('Arial', 'B', 10)
            self.multi_cell(0, 6, body)
        else:
            self.set_font('Arial', '', 10)
            self.multi_cell(0, 6, body)
        self.ln(2)

def create_pdf(cv_content, file_path, user_details, icons_path):
    pdf = PDF(user_details, icons_path)
    pdf.add_page()
    pdf.set_auto_page_break(auto=True, margin=15)

    icons = {
        "Email": "icon_email.png",
        "Objective": "icon_objective.png",
        "Education": "icon_education.png",
        "Experience": "icon_experience.png",
        "Skills": "icon_skills.png"
    }

    sections = cv_content.split('\n\n')
    for section in sections:
        title, body = section.split('\n', 1) if '\n' in section else (section, '')
        icon = icons.get(title.split(':')[0], None)
        if title.startswith("Email"):
            pdf.section_title(title, icon=icons["Email"])
            pdf.section_body(body, is_email=True)
        elif title.startswith("Education") or title.startswith("Experience"):
            pdf.section_title(title, icon=icon)
            degree_details, *subsections = body.split('\n')
            pdf.section_body(degree_details, is_subsection=True)
            for subsection in subsections:
                pdf.section_body(subsection)
        else:
            pdf.section_title(title, icon=icon)
            pdf.section_body(body)

    pdf.output(file_path)

#######################################################################################
#                                         Routes                                      #
#######################################################################################

@app.route('/generate-cv', methods=['POST'])
def generate_cv_pdf():
    # Extract user_id and user_details from the POST request
    data = request.json
    user_id = data['user_id']
    user_details = data['user_details']
    
    # Define file path for the generated PDF
    file_name = f"student_cv_{user_id}.pdf"
    file_path = os.path.join(PDF_DIR, file_name)
    
    # Ensure the PDF directory exists
    os.makedirs(PDF_DIR, exist_ok=True)
    
    # Generate CV content
    cv_content = generate_cv(user_details)
    
    # Create the PDF
    create_pdf(cv_content, file_path, user_details, ICONS_PATH)
    
    return send_from_directory(PDF_DIR, file_name, as_attachment=True)