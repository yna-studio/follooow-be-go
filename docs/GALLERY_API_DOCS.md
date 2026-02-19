# Gallery API Documentation

## Overview
The Gallery API provides endpoints for creating, reading, updating galleries with support for tags, images, and influencers.

## Base URL
```
http://localhost:20223
```

## Endpoints

### 1. Create Gallery (JSON)
**POST** `/galleries`

Create a new gallery with JSON payload.

#### Request Body
```json
{
  "title": "Summer Fashion 2024",
  "description": "Latest summer fashion trends and styles",
  "images": [
    {
      "is_cover": true,
      "url": "https://example.com/image1.jpg",
      "caption": "Summer dress",
      "created_on": 1640995200,
      "updated_on": 1640995200
    }
  ],
  "influencers": ["influencer_id_1", "influencer_id_2"],
  "lang": "ID",
  "tags": ["jilbab", "sport", "fashion", "summer"]
}
```

#### Curl Example
```bash
curl -X POST http://localhost:20223/galleries \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Summer Fashion 2024",
    "description": "Latest summer fashion trends and styles",
    "images": [
      {
        "is_cover": true,
        "url": "https://example.com/image1.jpg",
        "caption": "Summer dress",
        "created_on": 1640995200,
        "updated_on": 1640995200
      }
    ],
    "influencers": ["influencer_id_1", "influencer_id_2"],
    "lang": "ID",
    "tags": ["jilbab", "sport", "fashion", "summer"]
  }'
```

#### Response
```json
{
  "status": 201,
  "message": "Success create gallery",
  "data": null
}
```

---

### 2. Create Gallery with Image Upload (Multipart)
**POST** `/galleries/upload`

Create a new gallery with image file uploads.

#### Request Body (multipart/form-data)
- `title` (string, required): Gallery title
- `description` (string, optional): Gallery description
- `lang` (string, optional): Language code (default: "ID")
- `influencers` (string, optional): Comma-separated influencer IDs
- `author_id` (string, optional): Author ID
- `tags` (string, optional): Comma-separated tags
- `images` (files, required): Image files

#### Curl Example
```bash
curl -X POST http://localhost:20223/galleries/upload \
  -F "title=Summer Fashion 2024" \
  -F "description=Latest summer fashion trends" \
  -F "lang=ID" \
  -F "influencers=influencer_id_1,influencer_id_2" \
  -F "author_id=author_id_123" \
  -F "tags=jilbab,sport,fashion,summer" \
  -F "images=@image1.jpg" \
  -F "images=@image2.jpg"
```

#### Response
```json
{
  "status": 201,
  "message": "Success create gallery with images",
  "data": {
    "gallery_id": "507f1f77bcf86cd799439011"
  }
}
```

---

### 3. Update Gallery (JSON)
**PUT** `/galleries/{gallery_id}`

Update an existing gallery with JSON payload.

#### Request Body
```json
{
  "title": "Updated Summer Fashion 2024",
  "description": "Updated description",
  "images": [
    {
      "is_cover": true,
      "url": "https://example.com/updated_image.jpg",
      "caption": "Updated caption",
      "created_on": 1640995200,
      "updated_on": 1640995200
    }
  ],
  "influencers": ["new_influencer_id"],
  "lang": "EN",
  "tags": ["jilbab", "sport", "updated", "fashion"]
}
```

#### Curl Example
```bash
curl -X PUT http://localhost:20223/galleries/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Summer Fashion 2024",
    "description": "Updated description",
    "images": [
      {
        "is_cover": true,
        "url": "https://example.com/updated_image.jpg",
        "caption": "Updated caption",
        "created_on": 1640995200,
        "updated_on": 1640995200
      }
    ],
    "influencers": ["new_influencer_id"],
    "lang": "EN",
    "tags": ["jilbab", "sport", "updated", "fashion"]
  }'
```

#### Response
```json
{
  "status": 200,
  "message": "Gallery updated successfully",
  "data": {
    "gallery_id": "507f1f77bcf86cd799439011"
  }
}
```

---

### 4. Update Gallery with Image Upload (Multipart)
**PUT** `/galleries/{gallery_id}/upload`

Update an existing gallery with new image uploads.

#### Request Body (multipart/form-data)
- `title` (string, optional): Gallery title
- `description` (string, optional): Gallery description
- `lang` (string, optional): Language code
- `influencers` (string, optional): Comma-separated influencer IDs
- `tags` (string, optional): Comma-separated tags
- `images` (files, optional): New image files

#### Curl Example
```bash
curl -X PUT http://localhost:20223/galleries/507f1f77bcf86cd799439011/upload \
  -F "title=Updated Title with New Images" \
  -F "description=Updated description with new images" \
  -F "tags=jilbab,sport,new,fashion,updated" \
  -F "images=@new_image1.jpg" \
  -F "images=@new_image2.jpg"
```

#### Response
```json
{
  "status": 200,
  "message": "Gallery updated successfully with images",
  "data": {
    "gallery_id": "507f1f77bcf86cd799439011"
  }
}
```

---

### 5. List Galleries
**GET** `/galleries`

Get a list of galleries with pagination and filtering.

#### Query Parameters
- `limit` (integer, optional): Number of galleries per page (default: 6)
- `page` (integer, optional): Page number
- `lang` (string, optional): Filter by language
- `influencer` (string, optional): Filter by influencer ID

#### Curl Examples
```bash
# Get all galleries with default pagination
curl http://localhost:20223/galleries

# Get galleries with custom limit and page
curl "http://localhost:20223/galleries?limit=10&page=2"

# Get galleries filtered by language
curl "http://localhost:20223/galleries?lang=ID"

# Get galleries filtered by influencer
curl "http://localhost:20223/galleries?influencer=influencer_id_123"
```

#### Response
```json
{
  "status": 200,
  "message": "Success",
  "data": {
    "galleries": [
      {
        "id": "507f1f77bcf86cd799439011",
        "title": "Summer Fashion 2024",
        "description": "Latest summer fashion trends",
        "images": [...],
        "influencers": [...],
        "influencers_data": [...],
        "lang": "ID",
        "views": 150,
        "slug": "summer-fashion-2024",
        "tags": ["jilbab", "sport", "fashion"],
        "author_id": "author_id_123",
        "author": {
          "id": "author_id_123",
          "username": "john_doe"
        },
        "created_on": 1640995200,
        "updated_on": 1640995200
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 5,
      "total_items": 25
    }
  }
}
```

---

### 6. Get Gallery Details
**GET** `/galleries/{gallery_id}`

Get detailed information about a specific gallery.

#### Curl Example
```bash
curl http://localhost:20223/galleries/507f1f77bcf86cd799439011
```

#### Response
```json
{
  "status": 200,
  "message": "Success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Summer Fashion 2024",
    "description": "Latest summer fashion trends and styles",
    "images": [
      {
        "is_cover": true,
        "url": "https://res.cloudinary.com/...",
        "caption": "Summer dress",
        "created_on": 1640995200,
        "updated_on": 1640995200
      }
    ],
    "influencers": ["influencer_id_1", "influencer_id_2"],
    "influencers_data": [
      {
        "id": "influencer_id_1",
        "name": "Jane Doe",
        "slug": "jane-doe",
        "avatar": "https://res.cloudinary.com/...",
        "label": "Fashion Influencer"
      }
    ],
    "lang": "ID",
    "views": 150,
    "slug": "summer-fashion-2024",
    "tags": ["jilbab", "sport", "fashion", "summer"],
    "author_id": "author_id_123",
    "author": {
      "id": "author_id_123",
      "username": "john_doe"
    },
    "created_on": 1640995200,
    "updated_on": 1640995200
  }
}
```

---

## Tags Field Details

The `tags` field is an array of strings that allows categorizing galleries:

- **Type**: `[]string`
- **Default**: `[]` (empty array)
- **Format**: Array of tag strings
- **Examples**: `["jilbab", "sport", "fashion", "summer"]`

### Tags in Different Endpoints

#### JSON Endpoints
```json
{
  "tags": ["jilbab", "sport", "fashion"]
}
```

#### Multipart Endpoints
```
tags=jilbab,sport,fashion
```

### Common Tag Examples
- `["jilbab"]` - Islamic modest fashion
- `["sport"]` - Sports and athletic wear
- `["fashion"]` - General fashion
- `["summer"]` - Seasonal collections
- `["casual"]` - Casual wear
- `["formal"]` - Formal attire
- `["traditional"]` - Traditional clothing
- `["modern"]` - Modern styles

## Error Responses

### Validation Errors (400)
```json
{
  "status": 400,
  "message": "Title is required",
  "data": null
}
```

### Not Found (404)
```json
{
  "status": 404,
  "message": "Gallery not found",
  "data": null
}
```

### Server Error (500)
```json
{
  "status": 500,
  "message": "Error creating gallery",
  "data": {
    "error": "detailed error message"
  }
}
```

## Notes

1. **Image Uploads**: Images are automatically uploaded to Cloudinary and stored with CDN URLs
2. **Slug Generation**: Slugs are automatically generated from titles (lowercase, hyphen-separated)
3. **Timestamps**: `created_on` and `updated_on` are automatically managed
4. **Tags**: Tags are optional and default to empty array if not provided
5. **Author**: Author information is automatically populated when `author_id` is provided
6. **Influencers**: Influencer data is automatically populated based on provided IDs
