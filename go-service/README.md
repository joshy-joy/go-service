# School Management PDF Report Service

## Overview

This project is a **Go microservice** designed to generate **student PDF reports**. The service use html and cromedp to create formatted pdf files.

---

## Workflow

1. **Authentication**
   - The Go service logs in to the Node.js backend using API credentials.
   - It retrieves authentication cookies (`accessToken`, `refreshToken`, `csrfToken`).

2. **Data Fetch**
   - For a given `student_id`, the Go service calls the backend endpoint `/api/v1/students/:id`.
   - It sends the required cookies and CSRF token to fetch the student’s profile data.

3. **Template Rendering**
   - The student data is injected into an HTML template stored in a YAML file (`reports.yaml`).
   - The template supports placeholders like `{{ .name }}`, `{{ .dob }}`, `{{ .class }}`, etc.

4. **PDF Generation**
   - The rendered HTML is converted to a PDF using a rendering engine that supports full HTML and CSS.
   - The PDF is sent back as a **downloadable file** in the HTTP response.

5. **Result**
   - The user receives a well-formatted, styled PDF containing the student’s information.

---

## Prerequisites

Make sure you have the following installed before starting:

- **Node.js** (v16 or higher)
- **PostgreSQL** (v12 or higher)
- **npm** or **yarn**
- ***google-chrome*** package for html render

---

## 1. Backend Setup

```bash
cd backend
npm install
cp .env.example .env  # Configure your environment variables
npm start
```

**Backend API URL:** `http://localhost:5007`

---

## 2. Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

---

## 3. Database Setup

You can set up the database in **two ways**:

### Option A — Manual Setup
```bash
# Create PostgreSQL database
createdb school_mgmt

# Run database migrations
psql -d school_mgmt -f seed_db/tables.sql
psql -d school_mgmt -f seed_db/seed-db.sql
```

### Option B — Docker Setup
```bash
docker compose -f docker-compose.postgres.yaml up -d
```

---

## 4. Add Sample Data

Once the database is running, connect to PostgreSQL and run:

```sql
-- 1) Ensure class & section exist
INSERT INTO classes(name) VALUES ('Grade 8') ON CONFLICT (name) DO NOTHING;
INSERT INTO sections(name) VALUES ('B')      ON CONFLICT (name) DO NOTHING;

-- (Optional) assign a class teacher
-- INSERT INTO class_teachers(teacher_id, class_name, section_name)
-- VALUES (1, 'Grade 8', 'B');

-- 2) Add the student
SELECT * FROM student_add_update('{
  "name": "Jane Smith",
  "email": "jane.smith@example.com",
  "gender": "Female",
  "phone": "0509876543",
  "dob": "2010-05-15",
  "admissionDate": "2024-09-01",
  "class": "Grade 8",
  "section": "B",
  "roll": 12,
  "currentAddress": "Dubai",
  "permanentAddress": "Sharjah",
  "fatherName": "John Smith",
  "fatherPhone": "0501112233",
  "motherName": "Mary Smith",
  "motherPhone": "0504445566",
  "systemAccess": true
}'::jsonb);
```

---

## 5. Go Service Setup

### Install Google Chrome package

```bash
# For macOS
brew install --cask google-chrome

# For Windows
winget install --id Google.Chrome

# For Debian/Ubuntu
sudo apt update
sudo apt install wget -y
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo apt install ./google-chrome-stable_current_amd64.deb
```

### Start Go service

```bash
cd go-service
go mod tidy
go run main.go
```

This will start the Go service on port **8081**.

---

## 6. Generate Student Report

Use the following curl command to fetch the PDF report:

```bash
curl --location 'http://localhost:8081/api/v1/students/3/report' \
--header 'Cookie: refreshToken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
```

If authentication is successful, this will trigger:
- Login to the backend
- Fetch student data
- Render the HTML template
- Convert HTML to PDF
- Return the PDF as a downloadable file

---
