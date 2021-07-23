package main

import (
	"bytes"
	"errors"
	"golang.org/x/net/html" // XXX
	"io"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/google/uuid"
)

const XPathActionInputHTMLStr = "XPathActionInputHTMLStr"
const XPathActionInputHTMLBytes = "XPathActionInputHTMLBytes"
const XPathActionOutputStr = "XPathActionOutputStr"

type XPathAction struct {
	AbstractAction
	XPath string
}

func NewXPathAction(xpath string, expectMany bool) *XPathAction {
	return &XPathAction{
		AbstractAction: AbstractAction{
			CanFail:    false,
			ExpectMany: expectMany,
			AllowedInputNames: []string{
				XPathActionInputHTMLStr,
				XPathActionInputHTMLBytes,
			},
			AllowedOutputNames: []string{
				XPathActionOutputStr,
			},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{
				XPathActionOutputStr: []*DataPipe{},
			},
			UUID: uuid.New().String(),
		},
		XPath: xpath,
	}
}

// https://stackoverflow.com/a/38855264
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	if n != nil {
		html.Render(w, n)
	}
	return buf.String()
}

func (xa *XPathAction) Run() error {
	if xa.Inputs[XPathActionInputHTMLStr] == nil && xa.Inputs[XPathActionInputHTMLBytes] == nil {
		return errors.New("Input not connected")
	}

	if xa.Outputs[XPathActionOutputStr] == nil {
		return errors.New("Output not connected")
	}

	var htmlStr string

	if xa.Inputs[XPathActionInputHTMLStr] != nil {
		htmlStr, _ = xa.Inputs[XPathActionInputHTMLStr].Remove().(string)
	} else if xa.Inputs[XPathActionInputHTMLBytes] != nil {
		htmlBytes, ok := xa.Inputs[XPathActionInputHTMLBytes].Remove().([]byte)
		if ok {
			htmlStr = string(htmlBytes)
		}
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return err
	}

	if !xa.ExpectMany {
		var n *html.Node
		n, err = htmlquery.Query(doc, xa.XPath)
		if err != nil {
			return err
		}

		result := renderNode(n)

		for _, outDP := range xa.Outputs[XPathActionOutputStr] {
			outDP.Add(result)
		}
	} else {
		var nodes []*html.Node
		nodes, err = htmlquery.QueryAll(doc, xa.XPath)
		if err != nil {
			return err
		}

		for _, n := range nodes {
			if n == nil {
				continue
			}

			result := renderNode(n)
			for _, outDP := range xa.Outputs[XPathActionOutputStr] {
				outDP.Add(result)
			}
		}
	}

	return nil
}