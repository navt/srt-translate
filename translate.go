package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// Config - конфигурация утилиты, значения находятся в config.json
type Config struct {
	ApiUrl   string
	ApiKey   string
	FromTo   string // направление перевода
	TextSize int    // максимально допустимый размер текста в символах для запроса к API
	Prefix   string // название выходного файла фомируется добавлением префикса
	Source   string // название файла-источника
}

// Body - ответ API Яндекс.Переводчика
type Body struct {
	Code int
	Lang string
	Text []string
}

// Get - метод для получения значений из файла конфигурации
func (conf *Config) Get() error {
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return err
	}
	return nil
}

func (body *Body) apiQuery(conf *Config, text []string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", conf.ApiUrl, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// формируем строку запроса
	q := req.URL.Query()
	q.Add("key", conf.ApiKey)
	q.Add("lang", conf.FromTo)
	for _, line := range text {
		q.Add("text", line)
	}
	req.URL.RawQuery = q.Encode()
	// выполняем сам запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(b, &body)
	if err != nil {
		fmt.Println("Ошибка декодирования JSON : ", err)
		return
	}
	if body.Code != 200 {
		fmt.Println("Код ответа API Переводчика : ", body.Code)
	}
	return
}

func main() {
	conf := Config{}
	err := conf.Get()
	if err != nil {
		fmt.Println("Ошибка конфигурации : ", err)
		os.Exit(2)
	}
	var line string   // строка из исходного файла
	var qStrings int  // счетчик строк текста в отдельном субтитре
	var text []string // все строки текста из исходного файла 
	var body Body

	// просмотр исходного файла субтитров
	inFile, err := os.Open(conf.Source)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	scanner := bufio.NewScanner(inFile)
	// вытаскивание текстовых строк из исходного файла
	for scanner.Scan() {
		// смотрим каждую строчку
		line = scanner.Text()
		if line != "" {
			if line[0] >= 48 && line[0] <= 57 {
				// это номер субтитра
				if _, err = strconv.Atoi(line); err == nil {
					continue
				}
				// это временной интервал
				if len(line) == 29 && line[13:16] == "-->" {
					qStrings = 1
					continue
				}
			}
			if qStrings > 0 {
				// это строка с текстом субтитра
				text = append(text, line)
				qStrings++
			}
		} else {
			// пустая строка - подготовка к новому циклу
			qStrings = 0
		}
	} // end for
	inFile.Close()
	fmt.Println("В исходном файле строк текста : ", len(text))

	// В запросе может быть около 11665 символов
	// запрос к API Яндекс.Переводчика
	var begin, end int
	qSimbols := 155
	for i := 0; i < len(text); i++ {
		qSimbols = qSimbols + len(text[i]) + 6
		end = i
		if qSimbols > conf.TextSize || i == len(text)-1 {
			body.apiQuery(&conf, text[begin:end+1])
			if body.Code != 200 {
				os.Exit(4)
			}
			fmt.Println("Переведено строк текста : ", len(body.Text))
			fmt.Println("begin : ", begin, "end : ", end)
			for ii := begin; ii < end+1; ii++ {
				text[ii] = body.Text[ii-begin]
			}
			// перезаряжаем для очередной порции текста
			begin = end + 1
			qSimbols = 155
		}

	}

	// Запись выходного файла
	inFile, err = os.Open(conf.Source)
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	defer inFile.Close()
	outFile, err := os.Create(conf.Prefix + conf.Source)
	if err != nil {
		fmt.Println("Не создан выходной файл : ", err)
		os.Exit(6)
	}
	scanner = bufio.NewScanner(inFile)
	i := 0
	for scanner.Scan() {
		line = scanner.Text()
		if line != "" {
			if line[0] >= 48 && line[0] <= 57 {
				// это номер субтитра
				if _, err = strconv.Atoi(line); err == nil {
					outFile.WriteString(line + "\n")
					continue
				}
				// это временной интервал
				if len(line) == 29 && line[13:16] == "-->" {
					outFile.WriteString(line + "\n")
					qStrings = 1
					continue
				}
			}
			if qStrings > 0 {
				// это строка с текстом субтитра
				outFile.WriteString(text[i] + "\n")
				i++
				qStrings++
			}
		} else {
			// пустая строка - подготовка к новому циклу
			qStrings = 0
			outFile.WriteString("\n")
		}
	}
}

// https://tech.yandex.ru/translate/ - документация
