# GophKeeper

## Todo List

[X] CLI клиент на флагах
[X] Сервер
[X] Авторизация
[X] Регистрация
[X] Двунаправленная синхронизация
[X] Поддержка ревизий
[X] Тип данных - Пары логин/пароль
[X] HTTPS
[X] Данные о сборке
[X] Сборки под разные ОС и архитектуры
[ ] Тип данных - Данные банковских карт
[ ] Тип данных - Произвольные текстовые данные
[ ] Тип данных - Произвольные бинарные данные
[ ] Модульные тесты
[ ] Документация

## О сиситеме

GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

Первичное хранение, шифровка и извлечение данных происходит на клиенте, но реализована поддержка синхронизации, т.е. клиент может отправить данные на сервер, а затем зайти с другого устройства и запустив процесс синхронизации получить данные.

Клиент сделал по принципу простого cli инструмента. Управление данными осуществляется через вызов команд. Список всех команд можно получить из документации ниже

## Envs

### Client 

```
DB_DSN                  // Connection string для БД Postgres
LOG_LEVEL               // Уровень логирования (debug, info, etc)
ENCRYPT_KEY             // Ключ для шифрования данных (32 знака)
JWT_KEY                 // JWT ключ
SERVER_HOST             // Хост сервера
TLS_HANDSHAKE_TIMEOUT
RESPONSE_HEADER_TIMEOUT
DEALER_TIMEOUT
CLIENT_TIMEOUT
CTX_TIMEOUT
SYNC_CHANK_SIZE
SYNC_LIMIT
```

Примеры использования и примеры значений можно найти в Makefile

### Server

```
DB_DSN         // Connection string для БД
LOG_LEVEL      // Уровень логирования
MAX_BODY_SIZE
SERVER_ADDRESS // Адрес запуска HTTP-сервера
READ_TIMEOUT
WRITE_TIMEOUT
IDLE_TIMEOUT
JWT_KEY
TOKEN_EXPIRE
ENABLE_HTTPS  // Режим HTTPS
PK_PATH       // Путь до приватного ключа для режимо HTTPS
CRT_PATH      // Путь до сертификата для режимо HTTPS
```

Примеры использования и примеры значений можно найти в Makefile

## Генерация сертификатов

Для работы по HTTPS нужно сгенерировать сертификаты

```
make crt
```

После отработк команды они появятся в директории `/certs`

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
go run ./cmd/client  insert-loginpass --login= --password= --user-login=
make client-insert-loginpass login= pass= title= userid=
```

Ожидаемый результат
```
JWT token
```

#### Отправка данных на сервер (синхронизация)

```
go run ./cmd/client  sync-send --user-id=
make client-sync-send userid=
```

#### Список всех секретов пользователя

```
go run ./cmd/client list --user-id=
make client-list userid=
```

#### Расшифровка секрета

```
go run ./cmd/client view --item_id= --user-id=
make client-view: id= userid=
```

### Получение данных с сервера (синхронизация)

```
go run ./cmd/client sync-get --token=
make client-sync-get jwt=
```

### Удаление данных

```
go run ./cmd/client delete-item --id= --user-id=
make client-delete-item id= userid=
```

### Обновление логина и пароля

```
go run ./cmd/client update-loginpass --login= --password= --id= --user-id= --meta-id= --title=
make client-update-loginpass login= pass= id= userid= metaid= title=
```

### Создание произвольного текста

```
go run ./cmd/client create-text --path= --user-id= --title=
make client-create-text path= userid= title=
```

### Отображение произвольного текста

```
go run ./cmd/client get-text --user-id= --item-id=
make client-get-text userid= itemid=
```

### Создание Binary

```
go run ./cmd/client create-binary --user-id= --title= --path=
make client-create-binary title= userid= path=
```

### Загрузка Binary

```
go run ./cmd/client download-binary --user-id= --item-id= --out=
make client-download-binary userid= itemid= out=
```

## Сертификаты

Сервис может работать в режиме HTTPS. Для активации этого режима нужно на сервере установить значение для переменной окружения ENABLE_HTTPS

Также нужно сгенерировать сертификаты

```
make crt
```

Если флаг HTTPS включен, нужно не забыть установить значение хоста с протоколом https (переменная SERVER_HOST). Сделать это можно в Makefile или мануально