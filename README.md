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

## Client

### Команды

#### Регистрация пользователя

```
go run ./cmd/client register --email= --password=
make client-reg email=t1@dd.kz pass=qwerty
```

Ожидаемый результат
```
Register is OK
```
