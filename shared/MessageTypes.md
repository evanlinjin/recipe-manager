# Handshake
`"cmd": "handshake"`

**Request (Server to Client)**
string representing the key
```json
"secure_key"
```

**Response (Client to Server)**
boolean representing whether the key was accepted
```json
true
```

# New Account
`"cmd": "new_chef"`

**Request (Client to Server)**
object containing email and password
```json
{
    "email": "chef@test.com",
    "password": "secure_password"
}
```

**Response (Server to Client)**
string representing a message to display to end user
```json
"please check your email to activate your account"
```

# Login
`"cmd": "login"`

**Request (Client to Server)**
object containing email and password
```json
{
    "email": "chef@test.com",
    "password": "secure_password"
}
```

**Response (Server to Client)**
object container whether login successful and account info
```json
{
    "okay": true,
    "session": {
        "session_id": "",
        "session_key": "",
        "user_id": "",
        "user_name": "",
        "user_email": ""
    }
}
```
```json
{
    "okay": false,
    "message": "server error"
}
```
