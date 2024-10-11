# Golang Excel Importer with MySQL and Redis

## Overview
This application imports Excel data into a MySQL database and caches it in Redis. It also supports CRUD operations on the data.

## Features
- Excel file upload and validation
- Batch processing for MySQL inserts
- Redis caching with a 5-minute expiration
- API for paginated data retrieval
- Update and delete records

## Requirements
- Go 1.16+
- MySQL
- Redis

## Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/your-repo/my-app.git
   cd my-app
