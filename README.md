# spotmeet    [![Build](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml/badge.svg)](https://github.com/Aethedigm/spotmeet/actions/workflows/go.yml)
Capstone - Team E

# Progress
- [ ] Pages
  - [X]  Home Page
  - [X] Login Page
  - [ ] Register Page
  - [ ] Matches Page
  - [ ] Messages Page
    - [ ] Message Threads
  - [ ] Settings
  - [ ] Profiles
  - [ ] About
- [ ] Backend
  - TODO
- [ ] Middleware
  - TODO

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
