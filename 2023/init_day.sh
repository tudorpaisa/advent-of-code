INPUT=$1

if [[ $INPUT -eq "help" ]]; then
    echo "Usage: ./init_day.sh DAY_NAME

Advent of Code 2023 day initializer
-----------------------------------

The input DAY_NAME will create a folder with all the necessities for an Advent
of Code exercise in GO";
    exit 1;
fi

CURRENT_DIR=$(pwd);

mkdir $INPUT;
cd $INPUT;

go mod init advent_of_code/$INPUT
touch input1.txt
touch input2.txt

echo "package main

import (
	\"fmt\"
	\"log\"
	\"os\"
	\"strings\"
)

func readFileRaw(fname string) string {
	file, err := os.ReadFile(fname)

	if err != nil {
		msg := fmt.Sprintf(\"Encountered an error while reading '%s': %s\", fname, err)
		log.Fatal(msg)
	}

	return string(file)
}

func readLines(fname string) []string {
	var data []string = strings.Split(readFileRaw(fname), \"\n\")
	var out []string = *new([]string)

	for _, row := range data {
		if row == \"\" {
			continue
		}
		out = append(out, row)
	}
	return out
}

func main() {
	readLines(\"input1.txt\")
}
" > main.go

cd $CURRENT_DIR
