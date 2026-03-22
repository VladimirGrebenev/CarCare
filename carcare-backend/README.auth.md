# CarCare Auth API

Полный auth-flow: регистрация, подтверждение email, восстановление пароля, OAuth (Яндекс, Google), JWT/сессии, rate limiting, централизованное логирование.

## Примеры curl для всех auth-сценариев

### Регистрация
```
curl -X POST http://localhost:8080/auth/register -H 'Content-Type: application/json' -d '{"email":"user@example.com","password":"secret"}'
```

### Подтверждение email
```
curl -X POST http://localhost:8080/auth/confirm -H 'Content-Type: application/json' -d '{"token":"CONFIRM_TOKEN"}'
```

### Вход (логин)
```
curl -X POST http://localhost:8080/auth/login -H 'Content-Type: application/json' -d '{"email":"user@example.com","password":"secret"}'
```

### Восстановление пароля (запрос)
```
curl -X POST http://localhost:8080/auth/forgot -H 'Content-Type: application/json' -d '{"email":"user@example.com"}'
```

### Сброс пароля
```
curl -X POST http://localhost:8080/auth/reset -H 'Content-Type: application/json' -d '{"token":"RESET_TOKEN","new_password":"newpass"}'
```

### OAuth Яндекс
```
curl -X POST http://localhost:8080/auth/oauth/yandex -H 'Content-Type: application/json' -d '{"code":"YANDEX_CODE"}'
```

### OAuth Google
```
curl -X POST http://localhost:8080/auth/oauth/google -H 'Content-Type: application/json' -d '{"code":"GOOGLE_CODE"}'
```

### Обновление JWT (refresh)
```
curl -X POST http://localhost:8080/auth/refresh -H 'Content-Type: application/json' -d '{"refresh_token":"REFRESH_TOKEN"}'
```

### Logout
```
curl -X POST http://localhost:8080/auth/logout -H 'Content-Type: application/json' -d '{"token":"ACCESS_OR_REFRESH_TOKEN"}'
```

## Документация OpenAPI
См. [carcare-backend/api/openapi.yaml](carcare-backend/api/openapi.yaml)

## Тесты
- Unit/integration тесты: [carcare-backend/test/unit/](carcare-backend/test/unit/)

## TODO
- Реализовать централизованное логирование и rate limiting в middleware.
- Покрыть все сценарии тестами.
