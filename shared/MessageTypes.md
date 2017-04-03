# Handshake
`"cmd": "handshake"`

**Request (Server to Client)**<br>
string representing the key
```json
"secure_key"
```

**Response (Client to Server)**<br>
boolean representing whether the key was accepted
```json
true
```

# New Account
`"cmd": "new_chef"`

**Request (Client to Server)**<br>
object containing email and password
```json
{
    "email": "chef@test.com",
    "password": "secure_password"
}
```

**Response (Server to Client)**<br>
string representing a message to display to end user
```json
"please check your email to activate your account"
```

# Login
`"cmd": "login"`

**Request (Client to Server)**<br>
object containing email and password
```json
{
    "email": "chef@test.com",
    "password": "secure_password"
}
```

**Response (Server to Client)**<br>
object container whether login successful and account info
```json
{
    "okay": true,
    "session": {
        "session_id": "",
        "session_key": "",
        "chef_id": "",
        "chef_name": "",
        "chef_email": "",
        "teams": ["team_id_1", "team_id_2"]
    }
}
```
```json
{
    "okay": false,
    "message": "server error"
}
```
