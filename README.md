# DocuMind

DocuMind is a contract analysis and summarization tool powered by AI. It allows users to upload PDF, DOCX, or TXT files, extracts the text, analyzes the content using Google Gemini AI, and returns a structured summary, key clauses, and potential risks. All results are stored in a PostgreSQL database for future reference.

## Features
- Upload and analyze PDF, DOCX, or TXT contract files
- AI-powered summary, clause extraction, and risk detection
- Clean, modern React UI with Tailwind CSS
- Results displayed in organized tabs for easy navigation
- All analyses are saved to a PostgreSQL database

## Tech Stack
- **Frontend:** React, Tailwind CSS, Headless UI, Heroicons
- **Backend:** Go (Gin, GORM, Google Gemini API)
- **Database:** PostgreSQL

## Getting Started

### Prerequisites
- Node.js (v18+ recommended)
- Go (v1.20+ recommended)
- PostgreSQL

### 1. Clone the repository
```sh
git clone https://github.com/yourusername/documind.git
cd documind
```

### 2. Backend Setup
```sh
cd backend
cp configs/.env.example configs/.env  # Edit with your Gemini API key and DB credentials
# Install Go dependencies
go mod tidy
# Run the backend server
cd cmd/api
go run main.go
```

### 3. Frontend Setup
```sh
cd ../../documind
npm install
npm run dev
```

### 4. Access the App
Open [http://localhost:5173](http://localhost:5173) in your browser.

## Environment Variables
- `GEMINI_API_KEY`: Your Google Gemini API key
- `DATABASE_URL`: PostgreSQL connection string (e.g. `host=localhost user=postgres password=yourpass dbname=documind_db port=5432 sslmode=disable`)

## Folder Structure
```
backend/         # Go backend (API, services, models)
documind/        # React frontend (UI, components)
```

## License
MIT 