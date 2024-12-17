# Dobrify Bot

Notify about new prizes appearing in the game

### Endpoints

```
curl 'https://dobrycola-promo.ru/backend/oauth/token' \
  -H 'accept: application/json' \
  -H 'content-type: application/json' \
  -H 'origin: https://dobrycola-promo.ru' \
  -H 'referer: https://dobrycola-promo.ru/?signin' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36' \
  --data-raw '{"username":"xxx","password":"xxx"}'
```

```
curl 'https://dobrycola-promo.ru/backend/private/user' \
  -H 'accept: application/json' \
  -H 'authorization: Bearer ${AUTH_TOKEN}' \
  -H 'referer: https://dobrycola-promo.ru/?signin' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'
```

```
curl 'https://dobrycola-promo.ru/backend/private/prize/shop' \
  -H 'accept: application/json' \
  -H 'authorization: Bearer ${AUTH_TOKEN}' \
  -H 'referer: https://dobrycola-promo.ru/profile' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'
```

