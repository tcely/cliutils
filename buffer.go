package main

import (
	"github.com/ricochet2200/go-disk-usage/du"
	"time"
	"os"
	"bufio"
	"io"
	"flag"
)

var outOfSpace bool
var checkAgain bool = true
var appendOutput bool

var input *bufio.Reader = bufio.NewReader(os.Stdin)
var output *bufio.Writer = bufio.NewWriter(os.Stdout)

func inputFile(name string) (*os.File, error) {
	if len(name) <= 0 {
		return os.Stdin, &os.PathError{"open", name, os.ErrInvalid}
	}

	return os.Open(name)
}

func outputFile(name string) (*os.File, error) {
	if len(name) <= 0 {
		return os.Stdout, &os.PathError{"open", name, os.ErrInvalid}
	}

	if appendOutput {
		return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	} else {
		return os.Create(name)
	}
}

func flushDelay() {
	output.Flush()
	time.Sleep(10 * time.Millisecond)
}

func stillAvailable(limit uint64, filesystem string) {
	usage := du.NewDiskUsage(filesystem)
	outOfSpace = ! (usage.Available() > limit)
	checkAgain = true
}

func main() {
	var (
		looping bool = true
		rbytes, wbytes int
		rerr, werr error
	)

	inputPtr := flag.String("in", "", "read input from a file")
	outputPtr := flag.String("out", "", "write output to a file")
	appendPtr := flag.Bool("append", false, "append to the output file")
	filesystemPtr := flag.String("fs", ".", "which filesystem to check for available space")
	limitPtr := flag.Int("limit", 1024*1024*100, "how many bytes to reserve")
	flag.Parse()

	var limit uint64 = uint64(*limitPtr)
	check := func() {
		stillAvailable(limit, *filesystemPtr)
	}

	inFile, inErr := inputFile(*inputPtr)
	if nil == inErr {
		input = bufio.NewReader(inFile)
		defer inFile.Close()
	}
	appendOutput = *appendPtr
	outFile, outErr := outputFile(*outputPtr)
	if nil == outErr {
		output = bufio.NewWriter(outFile)
		defer outFile.Close()
		defer output.Flush()
	}

	defer flushDelay()
	for buf := make([]byte, input.Size()); looping; {
		if checkAgain {
			checkAgain = false
			time.AfterFunc(5 * time.Second, check)
		}

		if outOfSpace {
			time.Sleep(10 * time.Second)
			continue
		}

		rbytes, rerr = input.Read(buf)
		if rbytes > 0 && rerr == nil {
			wbytes, werr = output.Write(buf[0:rbytes])
			looping = looping && wbytes == rbytes && nil == werr
		}

		if output.Available() <= rbytes {
			flushDelay()
		}

		looping = looping && rerr != io.EOF
	}
}
