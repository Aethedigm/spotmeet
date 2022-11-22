# spotmeet    [![Build](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml/badge.svg)](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml)
Capstone - Team E

# Progress
- [ ] Pages
  - [X]  Home Page
  - [X] Login Page
    - [ ] Recover Password / Forgot Login
  - [X] Register Page
  - [X] Matches Page
  - [X] Messages Page
    - [X] Message Threads
  - [X] Settings
    - [X] Logout Button
  - [X] Profiles
    - [X] Edit Profiles
  - [ ] About
- [ ] Backend
  - [X] Models
  - [X] Migrations
  - [X] 80%+ Testing
  - [ ] Matching Algorithm
- [ ] Middleware
  - [X] Spotify Auth
  - [ ] Spotify Music Capture
  - [X] Spotify Artist Capture
  - [ ] Trigger Matches

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
