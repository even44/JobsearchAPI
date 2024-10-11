### Endpoints
```
GET /jobapplications - Returns a JSON list of all applications
```
```
GET /jobapplications/{id} - Returns a single JSON object containing job application with id
```
```
POST /jobapplications - Creates a job application from JSON request body
```
```
PUT /jobapplications/{id} - Updates a job application with id using JSON request body
```
```
DEL /jobapplications/{id} - Deletes job application with id
```
### Example JSON object for job application

```json
{
	"id": 3,
	"position": "Bossman",
	"company": "Bosses AS",
	"search_date": "11/10/2024",
	"deadline": "15/10/2024",
	"response": false,
	"interview": false,
	"done": false,
	"link": "https://finn.no"
}
```
