#### Low latency captcha service that provides functionality for generating captcha pictures with dynamic params, session managing and writing custom session pipelines

API:
```http
GET /
Main endpoint for checking session expiration, availability etc...
Custom checking pipelines should be described in endpoint/nginx/session.js
```

| Headers         | Type | Description                                                                                               |
|:----------------| :--- |:----------------------------------------------------------------------------------------------------------|
| `session-token` | `string` | **Required**. Your final session token that you need to get after resolving captcha and get session-token |

```http
GET /home
After successful checks described in session.js user will get data from this handler
```


```http
GET /generate_session
Generates session key for getting captcha image, saving secret session token in redis for checking and saving generated captcha image. Session key should be decrypted by AES-128 where key - text on captcha
```


```http
GET /img
Handler for getting captcha image that was generated in /generate_session
```

| Headers       | Type | Description                        |
|:--------------| :--- |:-----------------------------------|
| `captcha_key` | `string` | **Required**. equals session-token |

```http
INTERNAL
GET /redisadapter
Handler for making internal queries from another handlers
```

| Headers | Type | Description               |
|:--------| :--- |:--------------------------|
| `query` | `string` | **Required**. redis query |

