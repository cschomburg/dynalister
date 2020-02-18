package dynalist

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/xconstruct/dynalister/pkg/api"
)

type Dynalist struct {
	client *api.Client
}

func New() (*Dynalist, error) {
	token := os.Getenv("DYNALIST_TOKEN")
	client, err := api.New(token)

	return &Dynalist{client}, err
}

func (d *Dynalist) FetchFileTree() (*NodeTree, error) {
	resp, err := d.client.FileList()
	if err != nil {
		return nil, err
	}

	var files []Node
	for _, f := range resp.Files {
		files = append(files, FileNode(f))
	}

	return NewNodeTree(files), nil
}

func (d *Dynalist) PrintFiles(w io.Writer) error {
	tree, err := d.FetchFileTree()
	if err != nil {
		return err
	}

	return tree.Walk(func(parents []Node, n Node) error {
		f, _ := n.(FileNode)

		pre := strings.Repeat(" ", len(parents)*4) + "*"
		if _, err := fmt.Fprintln(w, pre, f.Title); err != nil {
			return err
		}

		return nil
	})
}

func (d *Dynalist) FetchDocumentTree(id string) (*NodeTree, error) {
	resp, err := d.client.DocRead(id)
	if err != nil {
		return nil, err
	}

	var nodes []Node
	for _, f := range resp.Nodes {
		nodes = append(nodes, DocumentNode(f))
	}

	return NewNodeTree(nodes), nil
}

func (d *Dynalist) PrintDocument(w io.Writer, id string) error {
	tree, err := d.FetchDocumentTree(id)
	if err != nil {
		return err
	}

	return tree.Walk(func(parents []Node, n Node) error {
		d, _ := n.(DocumentNode)

		pre := strings.Repeat(" ", len(parents)*4) + "*"
		if _, err := fmt.Fprintln(w, pre, d.Content); err != nil {
			return err
		}

		return nil
	})
}

func (d *Dynalist) ExportAll(outputPath string) error {
	files, err := d.FetchFileTree()
	if err != nil {
		return err
	}

	return files.Walk(func(parents []Node, n Node) error {
		f, _ := n.(FileNode)
		if f.Type != "document" {
			return nil
		}

		dir := getFullPath(outputPath, parents)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		fpath := dir + "/" + f.Title + ".md"
		w, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer w.Close()

		return d.PrintDocument(w, f.ID)
	})
}

func getFullPath(outputPath string, parents []Node) string {
	subPath := make([]string, 0, len(parents))
	for i, n := range parents {
		f, _ := n.(FileNode)

		if i == 0 && f.Title == "Untitled" {
			continue
		}

		subPath = append(subPath, f.Title)
	}

	return outputPath + "/" + strings.Join(subPath, "/")
}
