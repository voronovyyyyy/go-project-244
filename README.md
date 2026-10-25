# Вычислитель отличий

[![Actions Status](https://github.com/voronovyyyyy/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/voronovyyyyy/go-project-244/actions/workflows/hexlet-check.yml)
[![Go CI](https://github.com/voronovyyyyy/go-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/voronovyyyyy/go-project-244/actions/workflows/ci.yml)

Консольная программа и библиотека на Go для сравнения конфигураций JSON и YAML.
Поддерживает вложенные объекты, смешанные входные форматы и три вида отчёта:
`stylish`, `plain`, `json`. Ключи выводятся в алфавитном порядке.

## Установка и сборка

Понадобятся Go 1.23.3 или новее, Git и Make. Для проверки стиля используется
golangci-lint v2 (в CI закреплена версия 2.15.0).

```sh
git clone https://github.com/voronovyyyyy/go-project-244.git
cd go-project-244
go mod download
make build
./bin/gendiff --help
```

На Windows `make build` создаёт `bin/gendiff.exe`:

```powershell
.\bin\gendiff.exe --help
```

Без Make можно собрать и запустить программу напрямую:

```sh
go build -o bin/gendiff ./cmd/gendiff
go run ./cmd/gendiff --help
```

Собранные программы не хранятся в репозитории: каталог `bin/` добавлен в `.gitignore`.

## Использование

```text
NAME:
   gendiff - Compares two configuration files and shows a difference.

USAGE:
   gendiff [global options] <filepath1> <filepath2>

GLOBAL OPTIONS:
   --format string, -f string  output format (default: "stylish")
   --help, -h                  show help
```

```sh
# Stylish — формат по умолчанию
./bin/gendiff testdata/fixture/file1.json testdata/fixture/file2.json

# Можно смешивать JSON и YAML, использовать .yaml и .yml
./bin/gendiff testdata/fixture/file1.yml testdata/fixture/file2.json

# Текстовый отчёт
./bin/gendiff -f plain testdata/fixture/file1.yaml testdata/fixture/file2.yaml

# Структурированный JSON
./bin/gendiff --format json testdata/fixture/file1.json testdata/fixture/file2.yml
```

Флаги можно указывать перед путями или после них. Поддерживаются абсолютные пути
и относительные пути от текущего рабочего каталога. Пути с пробелами нужно
заключать в кавычки. Входной формат определяется по расширению файла без учёта регистра.

`stylish` показывает добавления с `+`, удаления с `-`, а изменения — двумя
соседними строками: сначала старое значение, затем новое. Вложенность обозначается
отступом в четыре пробела. Ключи внутри добавленных и удалённых объектов тоже сортируются.

В `plain` выводятся только изменения с полным путём свойства:

```text
Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting5' was added with value: [complex value]
```

Строки заключаются в одинарные кавычки, объекты и массивы отображаются как
`[complex value]`. В `stylish` массивы представлены списками в квадратных скобках;
они сравниваются как целые значения, без построения диффа отдельных элементов.

В `json` возвращается объект с пустым `key`, типом `root` и массивом `children`.
Узлы имеют поля `key`, `type`, `value1` (старое значение), `value2` (новое значение).
Типы узлов: `added`, `deleted`, `changed`, `unchanged`, `nested`. Вложенные узлы
содержат `children`, а исходные значения в них не дублируются. У неизменённых
узлов записывается только `value1`. Поля со значением `null` опускаются; например,
узел `changed` без `value2` означает изменение значения на `null`.
У пустого корневого объекта `children` равен `[]`.

Корень конфигурации должен быть объектом. YAML-файл должен содержать один документ,
ключи вложенных объектов должны быть строковыми. Ошибки чтения, парсинга,
неподдерживаемый формат и неверное число аргументов выводятся в stderr;
программа завершается с ненулевым кодом. Пустой отчёт `plain` ничего не печатает.

## Использование как библиотеки

Модуль называется `code` и экспортирует основную функцию:

```go
package main

import (
    "fmt"
    "log"

    "code"
)

func main() {
    result, err := code.GenDiff("first.json", "second.yaml", "plain")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result)
}
```

Сигнатура: `GenDiff(filepath1, filepath2, format string) (string, error)`.
Пустая строка формата выбирает `stylish`. Результат возвращается без завершающего
перевода строки; функция не печатает его и не изменяет исходные файлы.

## Тестирование и качество

```sh
make test
make test-coverage
make lint
make lint-fix
```

`make test-coverage` создаёт `coverage.out` и проверяет суммарное покрытие операторов
библиотеки. Порог задаётся переменной `COVERAGE_MIN` (по умолчанию **80%**):

```sh
make test-coverage COVERAGE_MIN=90
go tool cover -html=coverage.out
```

CLI и вспомогательный проверяющий скрипт тестируются командой `make test`, но не
входят в порог покрытия библиотеки. Покрытие учитывает как модульные тесты,
так и проверки полного вызова `GenDiff`. При недостаточном покрытии команда
возвращает ненулевой код, поэтому CI также падает.

GitHub Actions собирает и тестирует проект на Linux и Windows, проверяет порог
покрытия и запускает линтер. Профили покрытия доступны в артефактах сборки.
Тесты используют `testify` и вложенные фикстуры в `testdata/fixture/`, проверяют
все сочетания `.json`, `.yml`, `.yaml`, три формата вывода и ошибочные входные данные.

## Архитектура

```text
cmd/gendiff/           CLI: аргументы, справка, вывод и код завершения
gendiff.go            GenDiff: чтение → парсинг → дерево → форматирование
parsers/              Чистые функции разбора JSON/YAML и нормализации типов
internal/diff/        Построение отсортированного дерева типизированных узлов
formatters/           Выбор формата и отдельные stylish/plain/json форматеры
testdata/fixture/     Вложенные входные данные и ожидаемые отчёты
scripts/check-coverage/  Переносимая проверка порога покрытия
docs/                 Записи демонстраций
```

Дочерние узлы создаются только тогда, когда значения одного ключа в обоих файлах —
объекты. Добавленные и удалённые объекты остаются исходными значениями узла;
их отображение выполняет форматер. Парсеры, построитель дерева и форматеры
не читают файлы и не выводят данные в терминал.

## Демонстрации шагов

Записи показывают выполнение настоящих команд. Локальные копии можно открыть,
например, командой `asciinema play docs/nested-stylish.cast`.

| Шаг | Демонстрация | Локальная запись |
| --- | --- | --- |
| 2–4, 6 | [Подготовка, справка, тесты, покрытие и линтер](https://asciinema.org/a/Czzon6qNKIxHCeKA) | [setup-and-tests.cast](docs/setup-and-tests.cast) |
| 5 | [Сравнение плоских JSON](https://asciinema.org/a/bWNKu7ny9lClDY8G) | [flat-json.cast](docs/flat-json.cast) |
| 7 | [YAML и смешанные входные форматы](https://asciinema.org/a/YYZy84q56HwHMgeC) | [flat-yaml.cast](docs/flat-yaml.cast) |
| 8 | [Вложенные структуры, stylish](https://asciinema.org/a/kZ0k2l5iJrg16ASC) | [nested-stylish.cast](docs/nested-stylish.cast) |
| 9 | [Формат plain](https://asciinema.org/a/j0oR1OM87GmnzBNx) | [plain.cast](docs/plain.cast) |
| 10 | [Формат json](https://asciinema.org/a/fBVJdSjX8eGYjc1S) | [json.cast](docs/json.cast) |
