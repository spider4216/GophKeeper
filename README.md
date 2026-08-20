# GophKeeper

## Server

### Endpoints

#### Создание пользователя

##### Запрос

```
POST /auth/register

{
	"email": "email@example.com",
	"password": "qwerty"
}
```

##### Ответ

```
201 OK
```

#### Авторизация

##### Запрос

```
POST /auth/login

{
	"email": "email@example.com",
	"password": "qwerty"
}
```

##### Ответ

```
200 OK

{
	"status": "success",
	"message": "Login successfully",
	"token": "jwt here",
	"expired_at": 1787211264,
	"created_at": 1787207664
}
```

#### Отправка данных (синхронизация)

##### Запрос

```
POST /sync

Header
Authorization: jwt

Body
{
  "changes": [
    {
      "operation": "CREATE",
      "item": {
        "id": 42,
        "type": "login_password",
        "ciphertext": "1e33a00da21ebcbdfb49f65daa48a7a792570f0d857fe195972456c271120941"
      },
			"metadata": {
        "website": "github.com"
       }
    },
		{
      "operation": "CREATE",
      "item": {
        "id": 43,
        "type": "login_password",
        "ciphertext": "6b2913ad5ed912a3bcc50980028f153020f6d024bbf1e9878bf1274ad83625bf"
      },
			"metadata": {
        "website": "google.com",
				"note": "need change after 1 month"
       }
    }
  ]
}
```

##### Ответ

```
200 OK
```

## Client

### Команды

#### Регистрация пользователя

```
go run ./cmd/client register --email= --password=
make client-reg email= pass=
```

Ожидаемый результат
```
Register is OK
```

#### Авторизация пользователя

```
go run ./cmd/client login --email= --password=
make client-login email= pass=
```

Ожидаемый результат
```
JWT token
```

#### Создание пары логин / пароль

```
go run ./cmd/client  insert-loginpass --login= --password=
make client-insert-loginpass login= pass= title= jwt=
```

Ожидаемый результат
```
JWT token
```
