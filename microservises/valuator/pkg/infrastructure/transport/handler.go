package transport

import (
	"context"
	"html/template"
	"net/http"
	"server/pkg/app"
)

type Handler interface {
	Index(w http.ResponseWriter, r *http.Request)
	Summary(w http.ResponseWriter, r *http.Request)
}

type Response struct {
	Rank       float64 `json:"rank"`
	Similarity int     `json:"similarity"`
}

type SummaryData struct {
	Text       string
	Rank       float64
	Similarity int
}

// Структура, инкапсулирующая все зависимости
type handler struct {
	ctx             context.Context
	valuatorService app.ValuatorService
}

func NewHandler(
	ctx context.Context,
	valuatorService app.ValuatorService,
) *handler {
	// Возвращаем новый экземпляр приложения
	return &handler{
		valuatorService: valuatorService,
	}
}

// Обработчик главной страницы (форма для ввода текста)
func (a *handler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Valuator - Оценка текста</title>
	</head>
	<body>
		<h1>Введите текст для оценки:</h1>
		<form action="/summary" method="POST">
			<textarea name="text" rows="4" cols="50" placeholder="Введите текст"></textarea><br><br>
			<input type="submit" value="Отправить">
		</form>
	</body>
	</html>
	`

	// Рендерим форму с помощью Go-шаблона
	tmplParsed, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplParsed.Execute(w, nil)
}

// Обработчик для отображения результата обработки
func (a *handler) Summary(w http.ResponseWriter, r *http.Request) {
	// Получаем текст из формы
	text := r.FormValue("text")

	ctx := context.Background()

	// Вычисляем rank и similarity с использованием valuatorService
	rank := a.valuatorService.CalculateRank(text)
	similarity := a.valuatorService.AddText(ctx, text)

	// Формируем данные для страницы summary
	data := SummaryData{
		Text:       text,
		Rank:       rank,
		Similarity: similarity,
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Результат обработки</title>
	</head>
	<body>
		<h1>Результат обработки текста</h1>
		<p><strong>Текст:</strong> {{.Text}}</p>
		<p><strong>Рейтинг (rank):</strong> {{.Rank}}</p>
		<p><strong>Похожесть (similarity):</strong> {{.Similarity}}</p>
		<br>
		<a href="/">Вернуться на главную</a>
	</body>
	</html>
	`

	// Рендерим результат с помощью Go-шаблона
	tmplParsed, err := template.New("summary").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmplParsed.Execute(w, data)
}
