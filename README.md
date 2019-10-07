# srt-translate
Утилита перевода текста субтитров, использующая Яндекс.Переводчик.
### Использование
Склонируйте себе данный репозиторий. Предполагается, что у вас на компьютере установлен Golang. В терминале перейдите в директорию проекта и скомпилируйте исполняемый файл. Далее инструкция для Linux<br>
```$ go build translate.go```
Все изменяемые настройки вынесены в файл config.json. Вы должны получить [API-ключ](https://translate.yandex.ru/developers/keys) - это ключевой момент.<br>
Если направление перевода отличается от с английского на русский, то должно быть изменено
значение параметра FromTo - это [направление перевода](https://yandex.ru/dev/translate/doc/dg/reference/translate-docpage/) .  Префикс Prefix будет добавлен к названию файла-результата.
Source - название файла, который вы переводите.<br>
Далее: скопируйте файлы translate и config.json  в директорию с субтитрами. В терминале перейдите в директорию с субтитрами.  Правьте config.json и запускайте перевод.<br>
```$ ./transtate```