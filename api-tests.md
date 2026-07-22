# API Testing — cURL Reference

## Setup

```bash
BASE_URL="http://localhost:3000"

# Set after login:
TOKEN=""
```

## 1. Health

```bash
curl -s "$BASE_URL/api/health" | jq
```

```json
{"status":"ok","uptime":"...","timestamp":"..."}
```

## 2. Auth

### Register

```bash
curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"secret123"}' | jq
```

**Status:** 201  
**Response:** `{"id":"...","username":"test","role":"user","createdAt":"...","updatedAt":"..."}`

### Login

```bash
curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"secret123"}' | jq
```

```bash
# Capture token:
TOKEN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"secret123"}' | jq -r '.token')
```

## 3. Lists

### List all

```bash
curl -s "$BASE_URL/api/lists" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Create

```bash
curl -s -X POST "$BASE_URL/api/lists" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Watch Later"}' | jq
```

```bash
# Capture list ID:
LIST_ID=$(curl -s -X POST "$BASE_URL/api/lists" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Watch Later"}' | jq -r '.id')
```

### Get by ID

```bash
curl -s "$BASE_URL/api/lists/$LIST_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Delete

```bash
curl -s -X DELETE "$BASE_URL/api/lists/$LIST_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

## 4. Videos

### List all

```bash
curl -s "$BASE_URL/api/videos?page=1&limit=20" \
  -H "Authorization: Bearer $TOKEN" | jq
```

With filters:

```bash
curl -s "$BASE_URL/api/videos?page=1&limit=20&listId=$LIST_ID&search=example" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Upload (multipart)

```bash
curl -s -X POST "$BASE_URL/api/videos" \
  -H "Authorization: Bearer $TOKEN" \
  -F "name=My Video" \
  -F "video=@/media/titan/Workspace/Projects/Videos/000071.mp4" \
  -F "listId=$LIST_ID" | jq
```

```bash
# Capture video ID:
VIDEO_ID=$(curl -s -X POST "$BASE_URL/api/videos" \
  -H "Authorization: Bearer $TOKEN" \
  -F "name=My Video" \
  -F "video=@/path/to/video.mp4" | jq -r '.id')
```

### Get by ID

```bash
curl -s "$BASE_URL/api/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Update

```bash
curl -s -X PUT "$BASE_URL/api/videos/$VIDEO_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Renamed Video","listId":"'"$LIST_ID"'"}' | jq
```

### Delete

```bash
curl -s -X DELETE "$BASE_URL/api/videos/$VIDEO_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

## 5. Stream (HLS)

### Master playlist

```bash
curl -s "$BASE_URL/api/videos/$VIDEO_ID/stream" \
  -H "Authorization: Bearer $TOKEN"
```

### Segment

```bash
curl -s "$BASE_URL/api/videos/$VIDEO_ID/stream/segment_0.ts" \
  -H "Authorization: Bearer $TOKEN" -o segment_0.ts
```

## Docker

If running via Docker Compose:

```bash
BASE_URL="http://localhost:8080"
```

---

**Note:** Some endpoints may return 404, 400, or 409 depending on database state and data dependencies. The expected happy-path status is listed in the [AGENTS.md](./AGENTS.md) route table.
