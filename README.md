# Преобразование SQL-запосов 1Cv8
[![tag](https://img.shields.io/github/v/tag/mopo3ilo/1cv8metadata?sort=semver)](https://github.com/mopo3ilo/1cv8metadata/tags)
[![go version](https://img.shields.io/github/go-mod/go-version/mopo3ilo/1cv8metadata?label=go%20version)](https://go.dev/dl)
[![go report](https://goreportcard.com/badge/github.com/mopo3ilo/1cv8metadata)](https://goreportcard.com/report/github.com/mopo3ilo/1cv8metadata)

## Зависимости
- [go-mssqldb](https://github.com/denisenkom/go-mssqldb)
- [dlgs](https://github.com/gen2brain/dlgs)
- [sql1cv8](https://github.com/mopo3ilo/sql1cv8)
- [yaml.v2](https://gopkg.in/yaml.v2)

## Сборка
Можно выполнить задачу **build** через Visual Studio Code, это будет равнозначно выполнению следующих задач в терминале:
```bash
go build -ldflags="-s" .
strip 1cv8metadata
upx -9 1cv8metadata
```

## Настройка
Для работы предварительно следует создать в директории файл **bases.yml** с параметрами подключения к базе данных. Для подключения к SQL Server используется модуль от [go-mssqldb](https://github.com/denisenkom/go-mssqldb), так что примеры строк подключения смотрите там.
Пример содержимого файла:
```yml
CRM: "sqlserver://sa:password@sqlserver/?database=crm&connection+timeout=30&encrypt=disable&app+name=1Cv8+Metadata"
ZUP: "sqlserver://sa:password@sqlserver/?database=zup&connection+timeout=30&encrypt=disable&app+name=1Cv8+Metadata"
BUH: "sqlserver://sa:password@sqlserver/?database=buh&connection+timeout=30&encrypt=disable&app+name=1Cv8+Metadata"
```

## Использование
При простом запуске исполняемого файла будет выведено три диалоговых окна:
1. Выбор базы данных из **bases.yml**
2. Выбор файла с запросом для преобразования
3. Выбор действия, сохранить преобразованный запрос, либо выполнить его в базе данных

Так же каждому диалоговому окну соответствует свой ключ командной строки:
1. -database {Имя базы данных из файла **bases.yml**}
2. -filename {Путь к оригинальному файлу с SQL-запросом}
3. -save и/или -exec, которые соответствуют выбору Сохранить и Выполнить соответственно

При сохранении создаётся новый файл с суффиксом **_new**. Например для файла script.sql будет создан файл script_new.sql
