package doc_test

import (
	"bytes"
	"fmt"

	"github.com/umarcor/cobra"
	"github.com/umarcor/cobra/doc"
)

func er(msg interface{}) {
	cobra.Er(msg)
}

func ExampleGenManTree() {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "my test program",
	}
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	er(doc.GenManTree(cmd, header, "/tmp"))
}

func ExampleGenMan() {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "my test program",
	}
	header := &doc.GenManHeader{
		Title:   "MINE",
		Section: "3",
	}
	out := new(bytes.Buffer)
	er(doc.GenMan(cmd, header, out))
	fmt.Print(out.String())
}
