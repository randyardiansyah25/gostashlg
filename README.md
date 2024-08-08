# gostashlg

Logger integrated with Logstash base on [kpango/glg](https://github.com/kpango/glg)

## Instalasi

```go
go get github.com/randyardiansyah25/gostashlg
```

## Struktur Umum

### Struktur Log

```go
type Fields struct {
    IdentifierName string      `json:"service"`
    Timestamp      string      `json:"service_timestamp"`
    Level          Level       `json:"level"`
    Event          string      `json:"event"`
    Message        string      `json:"log_message"`
    Data           interface{} `json:"data,omitempty"`
}
```

Merupakan struktur objek log yang akan ditulis, baik ke file atau dikirim ke logstash

* Attribute ```IdentifierName``` digunakan sebagian identitas nama service untuk dikenali harus dikenali padak logstash sebagai filter pembuatan index di elasticsearch

* Attribute ```Timestamp``` digunakan sebagai tanggal dan waktu log dibuat, format **TIMESTAMP_ISO8601** atau **yyyy-MM-dd HH:mm:ss**

* Attribute ```Level``` digunakan sebagai level dari log.

* Attribute ```Event``` digunakan sebagai informasi aktivitas pada saat log dibuat. Bisa berupa path pada endpoint, modul atau function di *source code* atau aktivitasi lain
* Attribute ```Message``` digunakan sebagai informasi pesan log yang dibuat
* Attribute ```Data``` digunakan untuk informasi detail dari log, bisa berisi text atau objek

### Const

```go
type Level string

const (
    LOG   Level = "LOG"
    WARN  Level = "WARN"
    ERROR Level = "ERROR"
    INFO  Level = "INFO"
    DEBUG Level = "DEBUG"
    PRINT Level = "PRINT"
    TRACE Level = "TRACE"
    FAIL  Level = "FAIL"
)
```

Konstanta untuk digunakan sebagai value Level pada struktur ```Fields```

### Template

Template merupakan fitur yang dapat digunakan untuk mencetak log secara dinamis. Template ini memanfaatkan package ```html/template``` di golang. Untuk detail, silahkan pelajari di [dokumentasi](https://pkg.go.dev/html/template) dari package tersebut

```go
type Template struct {
 pattern map[Level]string
}

func (t *Template) Add(level Level, template string) *Template {
 t.pattern[level] = template
 return t
}
```

Struktur ini digunakan untuk mendefinisikan template log yang dicetak ke layar berdasarkan **Level**. Untuk menambahkan template berdasarkan level, bisa menggunakan fungsi ```Add()```

#### Contoh

```go
gostashlg.NewTemplate().
  Add(gostashlg.LOG, "{{.Data.Type}}, FROM:4 {{.Data.RemoteAddr}}, {{.Event}}, {{.Message}}, Data:\n{{.Data.Body}}").
  Add(gostashlg.INFO, "{{.Data.Type}}, {{.Event}}, {{.Message}}, Data:\n{{.Data.Body}}")
```

Jika template tidak didefinisikan, maka log yang dibuat akan menggunakan **Default template** yang telah disediakan. Template default sebagai berikut:

```go
{{.Event}}, {{.Message}}, {{.Data}}
```

#### Contoh saat log dicetak dilayar

```go
2024-08-08 11:20:46 [LOG]: Simpan User, Sukses menyimpan user, User-ID: 100 Name: Mr. Smith
```

Secara default, **Timestamp** dan **Level** akan selalu dicetak, sehingga place holder yang dibuat untuk **custom** template, tidak perlu menambahkan 2 field tersebut

### Method

```go
func UseDefault() (l LoggerEngine, e error)
```

Digunakan pertama kali untuk membuat objek log menggunakan template default.

```go
func UseDefine(template *Template) (l LoggerEngine, e error)
```

Digunakan pertama kali untuk membuat objek log, dengan mengirimkan parameter template custom yang telah didefinisikan

Return :

* ```LoggerEngine``` berupa objek dari yang digunakan selanjutnya untuk mencetak log.
* ```error``` bernilai error jika terjadi kesalahan pada saat method tersebut gagal saat melakukan parsing template

```go
func (l *LoggerEngine) Write(f Fields)
```

Method ini digunakan untuk melakukan perintah mencatat log dan secara otomatis, log akan dikirim ke logstash. ```f Fields``` berupa objek dari struktur ```Fields``` yang sudah didefinisikan

```go
func (l *LoggerEngine) WriteOnly(f Fields)
```

Method ini digunakan untuk melakukan perintah mencatat log tetap tidak mengirimkan ke logstash. ```f Fields``` berupa objek dari struktur ```Fields``` yang sudah didefinisikan

### Environment

Untuk dapat memanfaatkan [gostashlg](https://github.com/randyardiansyah25/gostashlg) agar secara otomatis mengirimkan log ke logstash, buat konfigurasi ```logstash.host``` di environment variable yang berisi url dari logstash. Anda juga dapat melakukan ini dengan membuat statement di *code* program dengan perintah untuk menambah environment variables tersebut. Sebagai contoh, lihat dibawah berikut:

```go
err := os.Setenv("logstash.host", "http://localhost:5044")
if err != nil {
    fmt.Println("Error setting environment variable:", err)
    return
}
```
