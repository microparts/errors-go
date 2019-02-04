errors
------

Пакет для работы с ошибками.

## Полезные методы (common)

### New(string) error

Создаёт новую ошибку с переданным текстом

### HasErrors(interface{}) bool

Проверяет, есть ли в переданном аргументе ошибки. Проверяет как на error тип, так и на различные 
модификации, которые могут прийти: `map[string]error`, `[]error`, `gorm.Errors`, ...

## Gin specific

### Response(*gin.Context, interface{})

Формирует `ResponseObject` для ответ ошибочным запросом. Обрабатывает различные виды ошибок 
и отдаёт их соответствующим образом.

### MakeResponse(interface{}) (int, *ErrorObject)

Формирует `ErrorObject` для предыдущего метода. Удобно пользоваться, когда ответ содержит как 
`data`, так и `error`  

### InitValidator()

Биндит валидатор к gin'у