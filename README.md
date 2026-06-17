### Database design:

```
┌──────────────────────┐
│        users         │
├──────────────────────┤
│ id PK                │
│ username UNIQUE      │
│ password_hash        │
│ mfa_enabled          │
│ totp_secret          │
│ failed_attempts      │
│ locked_until         │
│ created_at           │
│ last_login_at        │
└──────────┬───────────┘
           │1
           │
           │
           │1
┌──────────▼───────────┐
│      sessions        │
├──────────────────────┤
│ id PK                │
│ session_token        │
│ user_id FK           │
│ created_at           │
│ expires_at           │
│ is_active            │
└──────────────────────┘
```

- One user can has exactly one session
- Each session has exactly one user
