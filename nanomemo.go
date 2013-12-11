package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/burke/nanomemo/supermemo"
	"github.com/burke/ttyutils"
)

var (
	input = flag.String("input", "", "CSV fact list.")
	openQ = flag.Bool("openq", false, "Call /usr/bin/open on questions when presenting questions?")
)

func init() {
	flag.Parse()
	if *input == "" {
		fmt.Printf("Usage: %s -input=<some.csv>\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\x1b[33mSee github.com/burke/nanomemo for more details.\x1b[0m\n")
		os.Exit(1)
	}
}

func loadAllFacts(csvpath string) supermemo.FactSet {
	f, err := os.Open(csvpath)
	if err != nil {
		log.Fatalf("Couldn't open %s: %s\n", csvpath, err.Error())
	}
	defer f.Close()

	var fs supermemo.FactSet

	csvr := csv.NewReader(f)
	csvr.FieldsPerRecord = -1
	for {
		record, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Couldnt' read csv: %s\n", err.Error())
		}
		fs, err = addFact(fs, record)
		if err != nil {
			log.Fatalf("Couldnt' parse csv: %s\n", err.Error())
		}
	}

	return fs
}

func dumpFacts(csvpath string, fs supermemo.FactSet) {
	f, err := os.OpenFile(csvpath, os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal("Couldn't open CSV file to save results: %s\n", err)
	}
	csvw := csv.NewWriter(f)

	for _, f := range fs {
		q, a, ef, n, interval, intervalFrom := f.Dump()
		sef := fmt.Sprintf("%0.6f", ef)
		sn := fmt.Sprintf("%d", n)
		sinterval := fmt.Sprintf("%d", interval)
		csvw.Write([]string{q, a, sef, sn, sinterval, intervalFrom})
	}
	csvw.Flush()
	f.Close()
}

func quiz(csvpath string, fs supermemo.FactSet, allFacts supermemo.FactSet) {
	for {
		forReview := fs.ForReview()
		if len(forReview) == 0 {
			break
		}
		for _, f := range forReview {
			printQuestion(f.Question)
			getKey()
			printAnswer(f.Answer)
			q := readQuality()
			f.Assess(q)
			dumpFacts(csvpath, fs)
			fmt.Printf("\n")
		}
	}
}
func printQuestion(q string) {
	fmt.Printf("\x1b[34mQ:\x1b[0m %s\n", q)
	exec.Command("open", q).Run()
}

func printAnswer(a string) {
	fmt.Printf("\x1b[34mA:\x1b[0m %s\n", a)
}

func readQuality() int {
	for {
		fmt.Printf("\x1b[34m?:\x1b[0m ")
		os.Stdout.Sync()
		k := getKey()
		q := k - 0x30
		if q >= 0 && q < 4 {
			fmt.Printf("\x1b[31m%c\x1b[0m\n", k)
		} else if q > 3 && q < 6 {
			fmt.Printf("\x1b[32m%c\x1b[0m\n", k)
		} else {
			fmt.Printf("%c\n", k)
		}
		if q <= 5 && q >= 0 {
			return int(q)
		}
	}
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usage: nanomemo facts.csv")
	}
	csvpath := os.Args[1]

	fs := loadAllFacts(csvpath)

	for setsize := 10; ; setsize += 10 {
		if setsize > len(fs) {
			setsize = len(fs) - 1
			subfs := fs[0:setsize]
			quiz(csvpath, subfs, fs)
			break
		}
		subfs := fs[0:setsize]
		quiz(csvpath, subfs, fs)
	}
	quiz(csvpath, fs, fs)

}

func addFact(fs supermemo.FactSet, record []string) (supermemo.FactSet, error) {
	var fact *supermemo.Fact
	switch len(record) {
	case 2:
		fact = supermemo.NewFact(record[0], record[1])
	case 6:
		q := record[0]
		a := record[1]
		ef, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		n, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, err
		}
		interval, err := strconv.ParseInt(record[4], 10, 64)
		if err != nil {
			return nil, err
		}
		intervalFrom := record[5]
		fact, err = supermemo.LoadFact(q, a, ef, int(n), int(interval), intervalFrom)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Invalid record format")
	}

	fs = append(fs, fact)

	return fs, nil
}

func getKey() byte {
	termios, err := ttyutils.MakeTerminalRaw(os.Stdin.Fd())
	if err != nil {
		log.Fatal("stdin is not a terminal or something ¯\\(°_o)/¯. Don't do that.")
	}
	defer ttyutils.RestoreTerminalState(os.Stdin.Fd(), termios)

	b := make([]byte, 1)
	os.Stdin.Read(b)

	// There's probably some combination of termios flags I can set so that I
	// still receive C-c,C-z,C-\ as signals rather than characters, while also
	// handling input byte-by-byte, but it was easier to hobo this way than go
	// digging through the termios manpages again.
	switch b[0] {
	case 3:
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	case 28:
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	case 26:
		syscall.Kill(os.Getpid(), syscall.SIGTSTP)
	}
	return b[0]
}
