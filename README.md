<p align="center">
  <img src="fastio-logo.svg" alt="fastio — fast buffered I/O for Go" width="420" />
</p>

<p align="center">
  <strong>fast buffered I/O for Go</strong>
</p>

# FastIO

Библиотека быстрых ввода/вывода на Go с упором на работу со стандартными потоками и файлами. Пакет `fastio` предоставляет два основных типа:

- **FastReader** — высокопроизводительное чтение из любого `io.Reader` с методами `NextInt`, `NextInt64`, `NextUint64`, `NextFloat64`, `NextWord`, `NextLine`, а также побайтовым доступом `ReadByte` и `PeekByte`.
- **FastWriter** — буферизованная запись в `io.Writer` с методами `WriteInt`, `WriteInt64`, `WriteUint64`, `WriteFloat64`, `WriteString`, `WriteLine`, `WriteByte` и общим `Write`.

Оба типа минимизируют количество аллокаций за счёт собственных буферов (по умолчанию 64 KB) и позволяют вручную управлять ошибками через `Err()` и `Flush()`.

## Установка

```bash
GOTOOLCHAIN=local go get github.com/PavelKhromykhGo/fastio/fastio
```

`GOTOOLCHAIN=local` гарантирует использование локальной версии Go без попытки загрузить другой toolchain.

## Быстрый старт

Пример суммирования чисел, поступающих на stdin (см. `examples/basic`):

```go
fr := fastio.NewReader(os.Stdin)
fw := fastio.NewWriter(os.Stdout)
defer fw.Flush()

n, _ := fr.NextInt()
sum := 0
for i := 0; i < n; i++ {
    x, _ := fr.NextInt()
    sum += x
}
_ = fw.WriteInt(sum)
_ = fw.WriteByte('\n')
```

Работа с файлами (см. `examples/fileio`):

```go
in, _ := os.Open("input.txt")
out, _ := os.Create("output.txt")
r := fastio.NewReader(in)
w := fastio.NewWriter(out)
defer w.Flush()

// ... чтение чисел и запись результата ...
```

## Тесты

В репозитории есть модульные тесты для `FastReader` и базовая проверка сборки других пакетов. Запуск:

```bash
GOTOOLCHAIN=local go test ./...
```

## Бенчмарки

Для оценки производительности доступны бенчмарки чтения и записи (сравниваются с `fmt.Fscan` и `bufio.Scanner` / `fmt.Fprintln`). Запуск:

```bash
GOTOOLCHAIN=local go test -bench . -benchmem ./fastio
```

Пример актуальных результатов на машине CI (AMD64):

```
goos: windows
goarch: amd64
pkg: github.com/PavelKhromykhGo/fastio/fastio
cpu: Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz
BenchmarkFastReader_NextInt-12              2692            449434 ns/op           65584 B/op          2 allocs/op
BenchmarkFmtFscan-12                         508           2365154 ns/op          164251 B/op      19993 allocs/op
BenchmarkBufioScanner-12                    3640            321064 ns/op            4144 B/op          2 allocs/op
BenchmarkFastWriter_WriteInt-12             5233            229743 ns/op           65600 B/op          2 allocs/op
BenchmarkFmtFprintWithBufio-12              2024            594906 ns/op          209250 B/op       9752 allocs/op

```

## Полезные файлы

- `fastio/reader.go` — реализация `FastReader` и вспомогательных методов.
- `fastio/writer.go` — реализация `FastWriter` и методов форматированной записи.
- `examples/basic` и `examples/fileio` — демонстрационные программы работы со стандартным вводом и файлами.
- `input.txt` / `output.txt` — тестовые данные для примера чтения/записи файлов.