# REST-API
REST API golang
* SearchClient - структура с методом FindUsers, который отправляет запрос во внешнюю систему и возвращает результат, 
немного преобразуя его
* SearchServer - своего рода внешняя система. Непосредственно занимается поиском данных в файле `dataset.xml`. 
В продакшене бы запускалась в виде отдельного веб-сервиса.

специфика

Данные для работы лежат в файле `dataset.xml`
* Параметр `query` ищет по полям `Name` и `About`
* Параметр `order_field` работает по полям `Id`, `Age`, `Name`, если пустой - то возвращаем по `Name`, если что-то другое - 
SearchServer ругается ошибкой. 
`Name` - это first_name + last_name из xml.
