# DooleyOnline

Welcome to **DooleyOnline** - Emory University's marketplace platform for buying, selling, and chatting about items within the Emory community.

## 📱 Getting Started

### Creating Your Account

DooleyOnline is exclusively for Emory University students, faculty, and staff. Follow these steps to create your account:

#### Step 1: Sign Up with Your Emory Email

1. Visit [dooleyonline.net](https://dooleyonline.net)
2. Click on **"Sign Up"**
3. Enter your **@emory.edu** email address
   - ⚠️ **Important**: Only `@emory.edu` email addresses are accepted
   - Example: `yourname@emory.edu`
4. Create a secure password
5. Fill in your profile information
6. Click **"Create Account"**

#### Step 2: Verify Your Email

After signing up, you'll receive a verification email. However, **Emory's email security system may block or quarantine the verification email**.

**If you don't see the verification email:**

1. **Check your spam/junk folder** - The email might be filtered there
2. **Check Emory's Quarantine System** (Most Common Issue):
   - Go to [security.microsoft.com/quarantine](https://security.microsoft.com/quarantine)
   - Log in with your Emory credentials
   - Look for an email from DooleyOnline
   - Click on the email and select **"Release"** or **"Approve"**
   - The verification email should now appear in your inbox

3. **Open the verification link** from the email
4. You'll see a confirmation message (even if it says "invalid verification", your account is likely verified)
5. Try logging in to confirm your account is active

#### Step 3: Start Using DooleyOnline

Once verified, you can:
- **Browse Items**: Explore listings from the Emory community
- **Post Items**: List items you want to sell
- **Chat with Sellers**: Start conversations about items you're interested in
- **Manage Your Profile**: Update your information and preferences

## 🛍️ Using the Platform

### Browsing and Searching Items

- Browse all available items on the home page
- Use categories to filter items by type
- Search for specific items using keywords
- View item details, photos, and descriptions

### Posting an Item

1. Click **"Post Item"** or **"+"**
2. Upload clear photos of your item
3. Add a descriptive title and detailed description
4. Select the appropriate category
5. Set your asking price
6. Submit your listing

### Chatting with Other Users

- Click on any item to view details
- Click **"Message Seller"** to start a conversation
- All your conversations are in one place in the **"Messages"** or **"Chat"** section
- Chat rooms show the name of the person you're talking with
- You can discuss multiple items in the same conversation

### Managing Your Account

- Update your profile information in **Settings**
- View your posted items
- Track your conversations
- Manage your liked items

## ❓ Troubleshooting

### I can't receive the verification email

**Solution**: Emory's email system blocks suspicious senders. Follow these steps:
1. Visit [security.microsoft.com/quarantine](https://security.microsoft.com/quarantine)
2. Find the DooleyOnline verification email
3. Release/approve it
4. Check your inbox again

### The verification link says "invalid" but I can log in

This is expected behavior! If you can log in successfully, your account is verified. The message is a known UI issue that doesn't affect functionality.

### I forgot my password

Click **"Forgot Password"** on the login page and follow the reset instructions. Make sure to check the quarantine system if you don't receive the reset email.

### My email is not @emory.edu

Unfortunately, DooleyOnline is exclusively for the Emory community. Only `@emory.edu` email addresses are supported at this time.

## 🔒 Privacy & Safety

- Only Emory community members can access the platform
- All users are verified through their Emory email
- Report any suspicious activity or inappropriate content
- Never share sensitive personal information in chats
- Meet in safe, public locations when exchanging items

## 📞 Support

Having issues? Need help?
- Check the troubleshooting section above
- Contact the DooleyOnline team through the platform
- Report bugs or suggest features via our feedback form

---

## 👨‍💻 Technical Documentation

The following section is for developers who want to contribute to or understand the technical implementation of DooleyOnline.

### Architecture Overview

This backend service is built with Go and follows a clean architecture pattern:

- **API Layer** (`internal/api/`) - HTTP handlers and middleware
- **Service Layer** (`internal/service/`) - Business logic
- **Data Layer** (`internal/db/`) - Database operations using sqlc
- **Models** (`internal/model/`) - Domain models
- **Storage** (`internal/storage/`) - File storage management

### Core Features

- **Authentication & Authorization**
  - JWT-based authentication
  - Secure password hashing with bcrypt
  - Email verification system
  - Token management

- **Item Management**
  - CRUD operations for items
  - Category-based organization
  - Pagination support
  - Image storage integration

- **User Management**
  - User registration and profiles
  - Email verification
  - Liked items tracking
  - Viewed items history

- **Real-time Chat**
  - WebSocket-based messaging
  - One-to-one chat rooms
  - Message history
  - Participant management

- **Storage Integration**
  - AWS S3-compatible storage
  - Image processing and optimization
  - Secure file uploads

### Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Echo v4
- **Database**: PostgreSQL (via pgx/v5)
- **Database Management**: sqlc for type-safe SQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Storage**: AWS SDK v2
- **WebSocket**: Gorilla WebSocket
- **Documentation**: Swagger/OpenAPI
- **Image Processing**: disintegration/imaging
- **Task Runner**: Task (taskfile)
- **AI Integration**: Google Gen AI

### Prerequisites for Developers

- Go 1.25 or higher
- PostgreSQL database
- AWS S3 or S3-compatible storage (for file uploads)
- [Task](https://taskfile.dev/) (optional, for task automation)

### Installation & Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/dooleyonline/backend.git
   cd backend
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   Create a `.env` file in the root directory:

   ```env
   # Environment
   ENV=development
   PORT=8080

   # Database
   DATABASE_URL=postgres://username:password@localhost:5432/dooleyonline?sslmode=disable

   # Storage (AWS S3 or compatible)
   STORAGE_BUCKET=your-bucket-name
   STORAGE_URL=https://your-storage-url
   STORAGE_S3_URL=https://s3.your-region.amazonaws.com
   STORAGE_REGION=your-region
   STORAGE_ACCESS_ID=your-access-key-id
   STORAGE_ACCESS_SECRET=your-secret-access-key

   # Authentication
   AUTH_TOKEN_SECRET=your-secret-key-for-jwt
   
   # Email Service
   EMAIL_FROM=noreply@dooleyonline.net
   EMAIL_SERVICE_CONFIG=your-email-service-credentials
   ```

4. **Generate database code (if needed)**

   ```bash
   task sql
   # or
   sqlc generate
   ```

5. **Generate API documentation**
   ```bash
   task docs
   # or
   swag init --dir ./internal/api,./internal/db/item,./internal/db/chat,./internal/db/user --generalInfo api.go --parseInternal --parseDependency
   ```

### Running the Application

#### Using Task

```bash
# Run the server
task run

# Run tests
task test

# Generate documentation
task docs

# Generate database code
task sql

# Run data generator
task datagen
```

#### Using Go directly

```bash
# Run the server
go run ./cmd/server

# Run tests
go test ./... -v
```

#### Using Docker

```bash
# Build the image
docker build -t dooleyonline-backend .

# Run the container
docker run -p 8080:8080 --env-file .env dooleyonline-backend
```

### API Documentation

Once the server is running, you can access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

### Project Structure

```
.
├── cmd/
│   ├── datagen/          # Data generation utility
│   └── server/           # Main application entry point
├── internal/
│   ├── api/              # HTTP handlers
│   │   ├── auth/         # Authentication endpoints
│   │   ├── category/     # Category endpoints
│   │   ├── chat/         # Chat & WebSocket endpoints
│   │   ├── item/         # Item endpoints
│   │   ├── shared/       # Shared API utilities
│   │   ├── storage/      # File upload endpoints
│   │   └── user/         # User endpoints
│   ├── config/           # Configuration management
│   ├── db/               # Database layer (sqlc generated)
│   │   ├── chat/         # Chat-related queries
│   │   │   ├── message/  # Message operations
│   │   │   ├── participant/ # Chat participant management
│   │   │   └── room/     # Chat room operations
│   │   ├── item/         # Item-related queries
│   │   │   ├── category/ # Category operations
│   │   │   └── item/     # Item CRUD operations
│   │   └── user/         # User-related queries
│   │       ├── liked/    # Liked items tracking
│   │       ├── user/     # User management
│   │       ├── verify/   # Email verification
│   │       └── viewed/   # Viewed items tracking
│   ├── model/            # Domain models
│   ├── service/          # Business logic layer
│   │   ├── auth/         # Authentication & email service
│   │   ├── category/     # Category service
│   │   ├── chat/         # Chat & WebSocket service
│   │   ├── item/         # Item management service
│   │   ├── storage/      # File storage service
│   │   └── user/         # User service
│   └── storage/          # AWS S3 storage client
├── sqlc/                 # SQL schemas and queries
├── docs/                 # Generated API documentation
├── Dockerfile            # Container definition
├── go.mod                # Go dependencies
├── sqlc.yaml             # sqlc configuration
└── taskfile.yaml         # Task automation
```

### Testing

Run all tests:

```bash
go test ./... -v
```

Run tests with coverage:

```bash
go test ./... -cover
```

Run specific package tests:

```bash
go test ./internal/service/auth/... -v
go test ./internal/service/chat/... -v
```

### Development Workflow

#### Database Schema Changes

This project uses raw SQL managed by sqlc. Schema definitions are in `sqlc/*/schema.sql` and queries are in `sqlc/*/query.sql`.

After modifying SQL files, regenerate the Go code:

```bash
task sql
# or
sqlc generate
```

#### Adding New API Endpoints

1. Define the SQL schema in `sqlc/{domain}/{entity}/schema.sql`
2. Define queries in `sqlc/{domain}/{entity}/query.sql`
3. Generate database code: `task sql`
4. Implement service logic in `internal/service/{domain}/`
5. Create handler in `internal/api/{domain}/`
6. Register routes in `internal/api/api.go`
7. Update documentation: `task docs`

#### Email Verification System

The application uses an email verification system for new user registrations:

- Verification codes are generated and stored in the database
- Emails are sent via the configured email service
- Codes expire after a set time period
- Users must verify their `@emory.edu` email address
- Only Emory domain emails are accepted

**Note for Developers**: Emory's Microsoft Exchange security may quarantine outgoing emails. Users need to approve these emails at [security.microsoft.com/quarantine](https://security.microsoft.com/quarantine).

#### WebSocket Chat Implementation

The chat system uses WebSocket connections for real-time messaging:

- One-to-one chat rooms only
- Real-time message delivery
- Message history persistence
- Automatic room creation between participants
- Connection management and cleanup

### Deployment

The application is containerized using Docker and can be deployed to any container orchestration platform:

```bash
# Build production image
docker build -t dooleyonline-backend:latest .

# Run with production config
docker run -d \
  --name dooleyonline-backend \
  -p 8080:8080 \
  --env-file .env.production \
  dooleyonline-backend:latest
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Acknowledgments

- Echo Framework for the excellent web framework
- sqlc for type-safe SQL code generation
- Gorilla WebSocket for WebSocket implementation
- The Go community for amazing tools and libraries
- Emory University community for using and supporting DooleyOnline
