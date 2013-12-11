package main

import (
	"encoding/csv"
	"errors"
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
			fmt.Printf("\x1b[34mQ:\x1b[0m %s\n", f.Question)
			exec.Command("open", f.Question).Run()
			getKey()
			fmt.Printf("\x1b[34mA:\x1b[0m %s\n", f.Answer)
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
					f.Assess(int(q))
					dumpFacts(csvpath, fs)
					fmt.Printf("\n")
					break
				}
			}
		}
	}
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usage: nanomemo facts.csv")
	}
	csvpath := os.Args[1]

	fs := loadAllFacts(csvpath)

	for setsize := 10; ; setsize++ {
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
		log.Fatal("STDIN is not a terminal or something")
	}
	defer ttyutils.RestoreTerminalState(os.Stdin.Fd(), termios)

	b := make([]byte, 1)
	os.Stdin.Read(b)
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
