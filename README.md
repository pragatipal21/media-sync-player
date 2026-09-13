# Media Sync Player

A full-stack multi-window media sequencing application where multiple independent display windows continuously play their own configured playlists. The system supports dynamic playlist management and global synchronization, allowing an administrator to temporarily override all windows to display the exact same media simultaneously.

## Features

- **Multiple independent media windows**: Each window operates independently with its own state.
- **Individual playlist for each window**: Playlists are strictly assigned to specific display windows.
- **Continuous playlist playback**: Media automatically loops within the master cycle.
- **Image, video, and explicit blank media support**: Handles `.mp4`, `.jpg`, and pure HTML/CSS blank intervals seamlessly.
- **Configurable media durations**: Image and blank media progress automatically based on configured durations; videos progress using native `onEnded` events.
- **5-hour playback cycle**: Implements a unified 5-hour production playback cycle.
- **Playlist restart after the cycle**: After the master cycle ends, playback smoothly resets and restarts the configured playlist.
- **Global sync playback across all windows**: One click triggers the same media to play globally across all screens.
- **Configurable sync duration**: The sync state persists for an explicitly defined number of seconds.
- **Automatic return**: After the sync duration ends, every window seamlessly resumes its own local playlist.
- **Dynamic media addition to playlists**: Live addition of media to any window's playlist without requiring a page refresh.
- **Persistent MongoDB storage**: Playlists, windows, and media records are securely stored and persisted.
- **React frontend**: Fast, responsive UI powered by React and Vite.
- **Go backend REST APIs**: Highly performant, concurrent backend routing logic.

## How It Works

1. Each window is initialized with its own playlist containing sequential media IDs.
2. Media plays sequentially. When a media item ends (via duration or video completion), the next item in the list plays.
3. The playlist continues looping throughout a configured master cycle.
4. The production cycle length is defined as exactly 5 hours.
5. At the end of the 5-hour cycle, the configured playlist resets and starts again from the beginning.
6. A "Sync" action can temporarily override the standard cycle. During Sync, all windows immediately pause their local playlist and display the selected global sync media simultaneously.
7. After the sync duration officially ends, the global state is lifted and every window resumes its own playlist playback.
8. Administrators can dynamically append new media items to a playlist; this updates the database and persists seamlessly across browser refreshes.

## Tech Stack

**Frontend:** React, Vite  
**Backend:** Go (Standard Library `net/http`)  
**Database:** MongoDB / MongoDB Atlas (using `go.mongodb.org/mongo-driver`)  
**Deployment:** Docker backend (ready for Render)

## Project Structure

```text
media-sync-player/
├── backend/
│   ├── cmd/
│   │   ├── seed/
│   │   └── server/
│   ├── internal/
│   │   ├── config/
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── repository/
│   │   └── services/
│   ├── .env.example
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── public/
│   │   └── media/
│   ├── src/
│   ├── .env.example
│   ├── package.json
│   └── vite.config.js
└── README.md
```

## Local Setup

### Backend

1. Navigate into the backend directory:
   ```bash
   cd backend
   ```
2. Copy the example environment file and configure it:
   ```bash
   cp .env.example .env
   ```
3. Install Go dependencies:
   ```bash
   go mod download
   ```
4. Run the MongoDB seed script to populate the database:
   ```bash
   go run cmd/seed/main.go
   ```
5. Start the backend Go server:
   ```bash
   go run cmd/server/main.go
   ```

### Frontend

1. Navigate into the frontend directory:
   ```bash
   cd frontend
   ```
2. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
3. Install NPM dependencies:
   ```bash
   npm install
   ```
4. Start the Vite development server:
   ```bash
   npm run dev
   ```

## Environment Variables

> **IMPORTANT:** Never commit `.env` files to Git. They are explicitly ignored in `.gitignore`. The `.env.example` files are strictly for providing placeholder values. Do NOT commit real MongoDB passwords or credentials.

### Backend (`backend/.env`)
```env
PORT=8080
MONGODB_URI=mongodb+srv://<username>:<password>@<cluster-url>/?retryWrites=true&w=majority
MONGODB_DATABASE=media_sync
```

### Frontend (`frontend/.env`)
```env
VITE_API_URL=http://localhost:8080
VITE_CYCLE_DURATION_MS=18000000
```

## API Overview

The Go backend exposes the following clean RESTful API routes:

- `GET /media` — Retrieve all available media items.
- `GET /windows` — Retrieve all windows and their associated playlist IDs.
- `GET /playlists` — Retrieve all playlists and their sequential media IDs.
- `POST /playlists/{id}/media` — Append a new media item to a specific playlist dynamically.
- `GET /sync` — Fetch the current global synchronization state (polling).
- `POST /sync` — Trigger a global synchronization event with a specific media ID and duration.

## Media and Playlist Behavior

- **Representation:** Media items are explicitly defined in MongoDB with a unique ID, Type (`image`, `video`, `blank`), File Path, and Duration (for non-videos).
- **Association:** A Window maps 1:1 with a Playlist. A Playlist contains an ordered array of Media IDs.
- **Blank Media:** Blank playback is not an assumption of idle time; it is explicitly configured as a `blank` media item in the playlist, rendering a solid black screen for its specific duration.
- **Durations:** Video durations are handled dynamically by HTML5 `onEnded` events. Images and blanks enforce their database-configured durations via JavaScript `setTimeout`.
- **Dynamic Addition:** Pushing new media to a playlist securely targets the MongoDB `$push` operator, persisting state immediately.

## Synchronization

The sync flow leverages a centralized backend state and frontend polling to coordinate playback cleanly:

**Normal State:**
- Window 1 → Loop its playlist sequentially
- Window 2 → Loop its playlist sequentially
- Window 3 → Loop its playlist sequentially

**During Sync (Active Global Override):**
- Selected media → Window 1
- Selected media → Window 2
- Selected media → Window 3

**After Sync:**
- Every window natively discards the sync state and seamlessly resumes its own normal local playback sequence. The sync state is purely temporary and does not destructively overwrite the saved MongoDB playlists.

## Persistence

The core infrastructure fundamentally relies on MongoDB for state persistence. The dynamic addition of media to a playlist is a database-level update. This guarantees that refreshing the browser, opening a new tab, or restarting the frontend application entirely will accurately recover the exact configured sequences for all windows.

## Deployment

The application is architected to be deployed natively on modern cloud infrastructure:

- **Frontend:** Deployed to **Vercel** (requires configuring `VITE_API_URL` and `VITE_CYCLE_DURATION_MS` in the Vercel dashboard).
- **Backend:** Deployed to **Render** using the provided `Dockerfile` (requires configuring `MONGODB_URI` and `MONGODB_DATABASE` natively in the Render dashboard).
- **Database:** Hosted globally on **MongoDB Atlas**.

## Testing

The application supports robust end-to-end testing of the following core requirements:

1. **Independent Playback:** Verify that the three windows progress through their own individual sequences autonomously.
2. **Sync Playback:** Fire a global sync and verify that all windows immediately switch to the selected media.
3. **Sync Duration:** Verify that exactly after the inputted sync duration expires, all windows return to standard playback.
4. **Dynamic Media Addition:** Select Window 1, enter `M3`, and submit. Verify that `M3` is seamlessly appended to the sequence.
5. **Persistence:** Refresh the browser and verify `M3` is permanently part of Window 1's cycle.
6. **Master Cycle:** Test the cycle resetting logic (You can temporarily reduce `VITE_CYCLE_DURATION_MS` in `.env` from 5 hours to a lower value like `60000` to witness the cycle boundary reset).

## Assumptions / Design Decisions

- **Independent Windows:** Each window has its own explicitly configured playlist.
- **Non-Destructive Sync:** Sync temporarily overrides normal playback without rewriting the permanent playlist arrays.
- **Explicit Blanking:** Blank time is an explicit playlist entry. Playlists shorter than the cycle loop; they do not automatically "blank" unless specifically configured to do so.
- **Master Cycle Duration:** Production cycle duration is explicitly defined as 5 hours (18,000,000 milliseconds).

## Future Improvements

*Future work to extend the application's capabilities:*

- **Admin Authentication:** Implement secure login/JWT authorization for the backend APIs.
- **Drag-and-Drop Reordering:** Add a visual UI tool to physically reorganize playlist arrays natively.
- **Media Upload Management:** Expose an API to physically upload MP4/JPG files to an S3 bucket instead of serving local static paths.
- **WebSocket Synchronization:** Upgrade the HTTP polling model to real-time WebSockets for sub-millisecond global sync precision.
- **Monitoring & Logging:** Implement structured application logging.

---

**Author:** Pragati Pal  
**GitHub:** [https://github.com/pragatipal21/media-sync-player](https://github.com/pragatipal21/media-sync-player)
