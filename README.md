# spotmeet    [![Build](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml/badge.svg)](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml)
Capstone - Team E
 (Archival Purposes Only, This Capstone Was Completed)
# Progress
- [X] Pages
  - [X]  Home Page
  - [X] Login Page
    - [X] Recover Password / Forgot Login
  - [X] Register Page
  - [X] Matches Page
  - [X] Messages Page
    - [X] Message Threads
  - [X] Settings
    - [X] Logout Button
  - [X] Profiles
    - [X] Edit Profiles
  - [X] About
- [X] Backend
  - [X] Models
  - [X] Migrations
  - [X] 80%+ Testing
  - [X] Matching Algorithm
- [X] Middleware
  - [X] Spotify Auth
  - [X] Spotify Music Capture
  - [X] Spotify Artist Capture
  - [X] Trigger Matches

# Setup
- Inside of the myapp folder
- `cd myapp`
- Start docker
- `docker-compose up`
- Migrate data into database
- `./celeritas migrate`
- Start service
- `go run .`

# Tear Down / Restart
- Inside of the myapp folder
- `cd myapp`
- Stop docker
- `docker-compose down --remove-orphans`
- Delete residual data
- `rm db-data -r`
- Restart docker, and let it recreate the folders
- `docker-compose up`
- Migrate data
- `./celeritas migrate`
- Start service
- `go run .`
