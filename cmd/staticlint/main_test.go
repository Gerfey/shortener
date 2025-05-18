package main

import (
	"testing"
)

func TestNoExitAnalyzer(t *testing.T) {
	if NoExitAnalyzer.Name != "noexitinmain" {
		t.Errorf("Неверное имя анализатора: %s, ожидалось: %s", NoExitAnalyzer.Name, "noexitinmain")
	}

	if NoExitAnalyzer.Doc != "Запрещает использование прямого вызова os.Exit в функции main пакета main" {
		t.Errorf("Неверное описание анализатора: %s", NoExitAnalyzer.Doc)
	}
}
