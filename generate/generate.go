package generate

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/Newbluecake/bootstrap/clang"
)

func Generate() {
	// Put some space after C deprecation warning
	fmt.Println("")
	fmt.Println("")

	resourcesDir, err := clangResourceDir()
	if err != nil {
		panic(err)
	}

	parseArgs := []string{
		"-I", path.Join(resourcesDir, "include"),
		"-I", "./skia/skia/",
		"-x", "c++-header",
	}

	var tu clang.TranslationUnit
	index := clang.NewIndex(0, 1)
	errCode := index.ParseTranslationUnit2("./generate/skia_header_files.h", parseArgs, nil,
		clang.TranslationUnit_SkipFunctionBodies, &tu)
	if errCode != clang.Error_Success {
		panic(errCode)
	}

	tu.TranslationUnitCursor().Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
		fmt.Println(cursor.Spelling())
		fmt.Println("  ", cursor.Kind().Spelling(), cursor.DisplayName())
		return clang.ChildVisit_Continue
	})
}

func clangResourceDir() (string, error) {
	out, err := exec.Command("clang", "-print-resource-dir").Output()
	if err != nil {
		log.Fatal(err)
	}

	resDir := strings.TrimSpace(string(out))
	parts := strings.Split(resDir, "\n")
	resDir = parts[0]

	if resDir == "" {
		return "", errors.New("no output when getting clang resource dir")
	}
	if !strings.HasPrefix(resDir, "/") {
		return "", fmt.Errorf("expected clang resource dir to start with '/', but it %s", resDir)
	}

	return resDir, nil
}
