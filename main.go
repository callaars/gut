package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/phayes/permbits"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const Spacer = "  "
const KiB = 1024
const MiB = KiB * KiB
const GiB = MiB * MiB

// Colour definitions
var ColorModTime = color.New(color.FgBlue)
var ColorPermDir = color.New(color.FgBlue, color.Bold)
var ColorPermOther = color.New(color.FgCyan)
var ColorPermRead = color.New(color.FgYellow)
var ColorPermWrite = color.New(color.FgRed)
var ColorPermExecute = color.New(color.FgGreen)
var ColorPermNone = color.New(color.FgYellow)
var ColorFileSize = color.New(color.FgGreen, color.Bold)
var ColorOwner = color.New(color.FgYellow, color.Bold)

func main() {
	setupApp()
}

// Pads a string with whitespaces to the left with a specific size and returns a new string.
func padLeft(size int, str string) string {
	return strings.Repeat(" ", size) + str
}

// Prints a given time.
func printDate(t time.Time) {
	formattedTime := t.Format("2 Jan 15:04")

	ColorModTime.Print(padLeft(12-len(formattedTime), formattedTime) + Spacer)
}

func printPermissions(file os.FileMode) {
	permissions := permbits.FileMode(file)
	// permissions.SetUserExecute(

	if file.IsDir() {
		ColorPermDir.Print("d")
	} else if file.IsRegular() {
		ColorPermNone.Print("-")
	} else {
		// We need to do more to find out this file mode
		ColorPermOther.Print(strings.ToLower(string(file.String()[0])))
	}

	if permissions.UserRead() {
		ColorPermRead.Print("r")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.UserWrite() {
		ColorPermWrite.Print("w")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.UserExecute() {
		ColorPermExecute.Print("x")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.GroupRead() {
		ColorPermRead.Print("r")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.GroupWrite() {
		ColorPermWrite.Print("w")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.GroupExecute() {
		ColorPermExecute.Print("x")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.OtherRead() {
		ColorPermRead.Print("r")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.OtherWrite() {
		ColorPermWrite.Print("w")
	} else {
		ColorPermNone.Print("-")
	}

	if permissions.OtherExecute() {
		ColorPermExecute.Print("x")
	} else {
		ColorPermNone.Print("-")
	}

	fmt.Print(Spacer)
}

func friendlySize(size int64) string {
	if size < KiB {
		return strconv.FormatInt(size, 10)
	} else if size < MiB {
		return strconv.FormatInt(size/KiB, 10) + "Ki"
	} else if size < GiB {
		return strconv.FormatInt(size/MiB, 10) + "Mi"
	}

	return string(size)
}

func printSize(file os.FileInfo) {
	if file.IsDir() {
		ColorPermNone.Print(padLeft(4, "-") + Spacer)
	} else {
		size := friendlySize(file.Size())
		ColorFileSize.Print(padLeft(5-len(size), size), Spacer)
	}
}

func printOwner(file os.FileInfo) {
	owner, _ := user.LookupId(fmt.Sprint(file.Sys().(*syscall.Stat_t).Uid))
	group, _ := user.LookupGroupId(fmt.Sprint(file.Sys().(*syscall.Stat_t).Uid))

	ColorOwner.Print(owner.Username + " " + group.Name + Spacer)
}

func outputFiles(files []os.FileInfo) {
	boldBlue := color.New(color.FgBlue, color.Bold)

	for _, file := range files {
		printPermissions(file.Mode())
		printSize(file)
		printOwner(file)
		printDate(file.ModTime())

		if file.IsDir() {
			boldBlue.Print(file.Name())
		} else {
			fmt.Print(file.Name())
		}

		fmt.Println()
	}
}

func setupApp() {
	app := cli.NewApp()
	app.Name = "gut"
	app.Usage = "ls replacement written in go"

	app.Action = func(c *cli.Context) error {
		fmt.Println("Amount of arguments:", c.NArg())

		if c.NArg() > 0 {
			fmt.Println("Last argument", os.Args[c.NArg()])
		}

		files, err := ioutil.ReadDir(os.Args[c.NArg()] + "/")

		if err != nil {
			log.Fatal(err)
		}

		outputFiles(files)

		return nil
	}

	app.Run(os.Args)
}
