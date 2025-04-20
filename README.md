# Questionnaire App for Kubernetes Modernization Assessment

A Go-based questionnaire application for assessing application suitability for Kubernetes modernization, designed following the Go philosophy of simplicity, pragmatism, and efficiency.

## Features

- Questionnaire-based assessment of application modernization readiness
- Scoring system with category-specific insights
- Automatically generated recommendations and risks
- Modernization plan generation based on assessment results
- Simple file-based persistence
- RESTful API for integration

## Architecture

This application follows idiomatic Go patterns and principles:

1. **Simplicity**: Clean interfaces, direct type definitions
2. **Pragmatism**: Focused on solving real problems without unnecessary abstraction
3. **Efficiency**: Lightweight file-based storage with minimal resource usage
4. **Concurrency**: Goroutines for graceful server management
5. **Standard Library First**: Minimal external dependencies
6. **Error Handling**: Comprehensive error handling with context
7. **Platform Independence**: Cross-platform compatibility via Docker

## Project Structure

```
questionnaire-app/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── api/              # HTTP API layer
│   ├── models/           # Data models
│   ├── services/         # Business logic
│   └── storage/          # Data persistence
├── Dockerfile            # Docker build configuration
├── docker-compose.yml    # Docker Compose configuration
├── go.mod                # Go module definition
└── go.sum                # Go module dependency checksums
```

## Running the Application

### Using Docker Compose (Recommended)

1. Clone the repository:
   ```
   git clone https://github.com/RedDogSME/questionnaire-app.git
   cd questionnaire-app
   ```

2. Start with Docker Compose:
   ```
   docker-compose up -d
   ```

3. Access the API at: http://localhost:8080/api/

### Building and Running Locally

1. Clone the repository:
   ```
   git clone https://github.com/RedDogSME/questionnaire-app.git
   cd questionnaire-app
   ```

2. Build the application:
   ```
   go build -o server ./cmd/server
   ```

3. Run the application:
   ```
   ./server
   ```

4. Access the API at: http://localhost:8080/api/

## API Endpoints

- `GET /api/health` - Health check endpoint
- `GET /api/questions` - List all questions
- `POST /api/assessments` - Create a new assessment
- `GET /api/assessments/{assessmentId}` - Get an assessment
- `POST /api/assessments/{assessmentId}/answers` - Save an answer
- `POST /api/assessments/{assessmentId}/complete` - Complete assessment and generate report
- `GET /api/assessments/{assessmentId}/report` - Get assessment report

## Example Usage

### Create a New Assessment

```bash
curl -X POST http://localhost:8080/api/assessments \
  -H "Content-Type: application/json" \
  -d '{"applicationId": "app1"}'
```

Response:
```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "applicationId": "app1",
  "createdAt": "2025-04-20T12:34:56Z",
  "answers": {},
  "status": "in_progress"
}
```

### Save an Answer

```bash
curl -X POST http://localhost:8080/api/assessments/f47ac10b-58cc-4372-a567-0e02b2c3d479/answers \
  -H "Content-Type: application/json" \
  -d '{"questionId": "q1", "optionId": "q1_a1"}'
```

### Complete Assessment and Get Report

```bash
curl -X POST http://localhost:8080/api/assessments/f47ac10b-58cc-4372-a567-0e02b2c3d479/complete
```

## Persistent Storage

The application uses a simple file-based storage system by default. Data is stored in the `./data` directory with the following structure:

- `./data/applications/` - Application details
- `./data/questions/` - Assessment questions
- `./data/assessments/` - User assessments
- `./data/reports/` - Generated reports

This directory is persisted when using Docker through a volume mount.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
