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
object containing session info if successful, or a string message if not
```json
{
    "session_id": "",
    "session_key": "",
    "chef_id": "",
    "chef_name": "",
    "chef_email": "",
}
```
```json
"server error"
```

# Logout
`"cmd": "logout"`

**Request (Client to Server)**<br>
empty string
```json
""
```

**Reply (Server to Client)**<br>
string message
```json
"success"
```

# Claim Session
`"cmd": "claim_session"`

**Request (Client to Server)**<br>
object with session_id and session_key
```json
{
    "chef_id": "some_id",
    "session_id": "some_other_id",
    "session_key": "some_key"
}
```

**Response (Server to Client)**<br>
object defining user info and various updated if accepted
```json
{
    "okay": true,
    "session": {
        "session_id": "",
        "session_key": "",
        "chef_id": "",
        "chef_name": "",
        "chef_email": "",
    }
}
```
```json
{
    "okay": false,
    "message": "server error"
}
```
