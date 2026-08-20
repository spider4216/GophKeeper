# GophKeeper

## Server

### Endpoints

#### Создание пользователя

##### Запрос

```
POST /auth/register

{
	"email": "abc2@dd.dd",
	"password": "qwerty"
}
```

##### Ответ

```
201 OK
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
