# DocuMind

DocuMind is an AI-powered contract analysis and summarization tool that helps legal professionals and businesses understand complex documents quickly. Upload PDF or DOCX files and get instant AI-powered summaries, key clause extraction, risk detection, and the ability to ask questions about your documents.

## ✨ Features

- **📄 Multi-format Support**: Upload and analyze PDF and DOCX contract files
- **🤖 AI-Powered Analysis**: Powered by Google Gemini AI for intelligent document processing
- **📊 Smart Summarization**: Get instant, comprehensive summaries of complex legal documents
- **⚠️ Risk Detection**: Automatically identify potential legal risks and important clauses
- **🔍 Key Term Extraction**: Extract and highlight the most important terms and conditions
- **💬 Document Chat**: Ask questions about your uploaded documents and get AI-powered answers
- **📱 Modern UI**: Clean, responsive React interface with beautiful animations
- **🔒 Enterprise Security**: Bank-level encryption ensures your sensitive documents stay private
- **💾 Persistent Storage**: All analyses are saved to PostgreSQL database with caching for efficiency

## 🚀 Getting Started

### Prerequisites
- **Node.js** (v18+ recommended)
- **Go** (v1.24+ recommended)
- **PostgreSQL** (v12+ recommended)
- **Google Gemini API key** ([Get one here](https://makersuite.google.com/app/apikey))

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/documind.git
cd documind
```

### 2. Backend Setup
```bash
cd backend

# Create environment file
cp configs/.env.example configs/.env

# Edit configs/.env with your credentials:
# GEMINI_API_KEY=your_gemini_api_key_here
# DATABASE_URL=host=localhost user=postgres password=yourpass dbname=documind_db port=5432 sslmode=disable
# PORT=8080

# Install Go dependencies
go mod tidy

# Run the backend server
cd cmd/api
go run main.go
```

The backend will start on `http://localhost:8080`

### 3. Frontend Setup
```bash
cd ../../documind

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:5173`

### 4. Access the Application
Open [http://localhost:5173](http://localhost:5173) in your browser to start using DocuMind.

## 🔧 Environment Variables

### Backend (.env file in `backend/configs/`)
```env
GEMINI_API_KEY=your_google_gemini_api_key
DATABASE_URL=host=localhost user=postgres password=yourpass dbname=documind_db port=5432 sslmode=disable
PORT=8080
```

## 📁 Project Structure

```
documind/
├── backend/                    # Go backend application
│   ├── cmd/api/               # Application entry point
│   │   ├── internal/          # Internal packages
│   │   │   ├── handlers/      # HTTP request handlers
│   │   │   ├── models/        # Database models
│   │   │   └── services/      # Business logic services
│   │   ├── pkg/               # Public packages
│   │   │   └── database/      # Database connection and utilities
│   │   └── configs/           # Configuration files
├── documind/                  # React frontend application
│   ├── src/
│   │   ├── components/        # React components
│   │   │   ├── ui/           # Shadcn/ui components
│   │   │   ├── DocumentUpload.tsx
│   │   │   └── DocumentChat.tsx
│   │   ├── pages/            # Page components
│   │   ├── api/              # API service functions
│   │   └── hooks/            # Custom React hooks
│   ├── public/               # Static assets
│   └── package.json          # Frontend dependencies
├── configs/                   # Shared configuration
└── README.md                 # This file
```

## 🔌 API Endpoints

### Document Analysis
- `POST /api/v1/analyze` - Upload and analyze a document
- `GET /api/v1/analyses` - Get list of all analyses
- `GET /api/v1/analyses/:id` - Get detailed analysis by ID

### Document Chat
- `POST /api/v1/contract-chat` - Ask questions about uploaded documents

### Health Check
- `GET /ping` - Health check endpoint

## 🎯 Key Features Explained

### Smart Caching
DocuMind uses file hash-based caching to avoid re-processing identical documents, saving time and API costs.

### Risk Assessment
The AI analyzes contracts for:
- Unusual terms and conditions
- Potential legal risks
- Missing important clauses
- Unfavorable terms

### Document Chat
Ask natural language questions about your uploaded documents and get instant AI-powered answers.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:
1. Check the [Issues](https://github.com/yourusername/documind/issues) page
2. Create a new issue with detailed information
3. Contact the maintainers

---

**Built with ❤️ using React, Go, and Google Gemini AI** 
