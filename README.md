# spotmeet    [![Build](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml/badge.svg)](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml)
Capstone - Team E

# Progress
- [ ] Pages
  - [X]  Home Page
  - [X] Login Page
  - [X] Register Page
  - [X] Matches Page
  - [ ] Messages Page
    - [ ] Message Threads
  - [ ] Settings
    - [ ] Logout Button
  - [X] Profiles
    - [ ] Edit Profiles
  - [ ] About
- [ ] Backend
  - [ ] Models
  - [ ] Migrations
  - [ ] 80%+ Testing
- [ ] Middleware
  - [ ] Spotify Auth
  - [ ] Spotify Music Capture
  - [ ] Trigger Matches

# Setup
- Inside of the myapp folder
- `cd myapp`
- Start docker
- `docker-compose up`
- Migrate data into database
- `./celeritas migrate`

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
