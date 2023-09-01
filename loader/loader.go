package loader

import (
	"context"
	"encoding/csv"
	"items-service/model"
	"math"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type Row map[string]string

type Loader struct {
	path string
	/*
		Кол-во полей намекает что не нужно это сохранять в структуру, окей, давайте сохраним в мапу
		а мапу уже сможем парсить в json

		и таким образом сервис позволит делать поиск по документам с разной структурой, не завязываясь
		на кол-во полей и их последовательности
	*/
	buffer []Row
	model  Row
}

func NewBufferLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadCSVItems(resourcePath string) error {
	file, err := os.Open(resourcePath)
	if err != nil {
		errors.WithMessage(err, "failed to open file")
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	/*
		Я понимаю что файл может быть на 9 террабайт, но в условие не сказаны требования
		Возможно у нас очень большие требования к скорости ответа, поэтому каждый раз бегать по диску
		и искать нужные записи мы не можем. Поэтому буду читать все, текущий файл позволяет)))
	*/
	slog.Info("staring to read resource", "path", resourcePath)
	records, err := reader.ReadAll()
	if err != nil {
		return errors.WithMessage(err, "failed to read CSV file")
	}
	model := records[0]
	for x := 1; x < len(records); x++ {
		l.buffer = append(l.buffer, parseRowByModel(model, records[x]))
	}
	slog.Info("resource readed", "path", resourcePath, "rows count", len(l.buffer), "colums count in model", len(model))
	return nil
}

func (l *Loader) ItemsById(_ context.Context, ids []string) (items []map[string]string, err error) {
	for _, id := range ids {
		find := false
		for _, item := range l.buffer {
			if id == item["id"] {
				items = append(items, item)
				find = true
				break
			}
		}
		if !find {
			return nil, errors.WithMessage(model.ErrNotFound, "failed to find in buffer")
		}
	}
	return items, nil
}

func parseRowByModel(model []string, data []string) (row Row) {
	minLen := math.Min(float64(len(model)), float64(len(data)))
	row = make(Row, len(model))
	for x := 0; x < int(minLen); x++ {
		row[model[x]] = data[x]
	}
	return row
}
