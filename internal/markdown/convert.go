package markdown

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const twitterMaxChars = 280

func ParseFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

func ToTwitter(content string) []string {
	source := []byte(content)
	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(source))

	var plainText strings.Builder
	walkNode(doc, source, &plainText, false)

	text := strings.TrimSpace(plainText.String())
	return splitThread(text)
}

func ToLinkedIn(content string) string {
	source := []byte(content)
	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(source))

	var result strings.Builder
	walkNode(doc, source, &result, true)

	return strings.TrimSpace(result.String())
}

func walkNode(node ast.Node, source []byte, buf *strings.Builder, linkedIn bool) {
	switch n := node.(type) {
	case *ast.Document:
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			walkNode(child, source, buf, linkedIn)
		}

	case *ast.Heading:
		if linkedIn {
			// Unicode bold for LinkedIn headings
			text := extractText(n, source)
			buf.WriteString(toBold(text))
			buf.WriteString("\n\n")
		} else {
			text := extractText(n, source)
			buf.WriteString(text)
			buf.WriteString("\n\n")
		}

	case *ast.Paragraph:
		text := extractText(n, source)
		buf.WriteString(text)
		buf.WriteString("\n\n")

	case *ast.List:
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			if item, ok := child.(*ast.ListItem); ok {
				text := extractText(item, source)
				if linkedIn {
					buf.WriteString("• ")
				} else {
					buf.WriteString("- ")
				}
				buf.WriteString(text)
				buf.WriteString("\n")
			}
		}
		buf.WriteString("\n")

	case *ast.FencedCodeBlock:
		var code bytes.Buffer
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			code.Write(line.Value(source))
		}
		buf.WriteString(strings.TrimSpace(code.String()))
		buf.WriteString("\n\n")

	case *ast.ThematicBreak:
		buf.WriteString("---\n\n")

	case *ast.Blockquote:
		text := extractText(n, source)
		for _, line := range strings.Split(text, "\n") {
			buf.WriteString("> ")
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		buf.WriteString("\n")

	default:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			walkNode(child, source, buf, linkedIn)
		}
	}
}

func extractText(node ast.Node, source []byte) string {
	var buf strings.Builder
	extractTextRecursive(node, source, &buf)
	return strings.TrimSpace(buf.String())
}

func extractTextRecursive(node ast.Node, source []byte, buf *strings.Builder) {
	switch n := node.(type) {
	case *ast.Text:
		buf.Write(n.Segment.Value(source))
		if n.SoftLineBreak() {
			buf.WriteString(" ")
		}
	case *ast.String:
		buf.Write(n.Value)
	case *ast.CodeSpan:
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			extractTextRecursive(child, source, buf)
		}
	case *ast.Emphasis:
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			extractTextRecursive(child, source, buf)
		}
	case *ast.Link:
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			extractTextRecursive(child, source, buf)
		}
		buf.WriteString(" (")
		buf.Write(n.Destination)
		buf.WriteString(")")
	case *ast.AutoLink:
		buf.Write(n.URL(source))
	default:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			extractTextRecursive(child, source, buf)
		}
	}
}

func splitThread(text string) []string {
	if utf8.RuneCountInString(text) <= twitterMaxChars {
		return []string{text}
	}

	paragraphs := strings.Split(text, "\n\n")
	var chunks []string
	var current strings.Builder

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if current.Len() == 0 {
			if utf8.RuneCountInString(para) > twitterMaxChars {
				// Split long paragraph by sentences
				chunks = append(chunks, splitLong(para)...)
				continue
			}
			current.WriteString(para)
			continue
		}

		candidate := current.String() + "\n\n" + para
		if utf8.RuneCountInString(candidate) <= twitterMaxChars {
			current.WriteString("\n\n")
			current.WriteString(para)
		} else {
			chunks = append(chunks, current.String())
			current.Reset()
			if utf8.RuneCountInString(para) > twitterMaxChars {
				chunks = append(chunks, splitLong(para)...)
			} else {
				current.WriteString(para)
			}
		}
	}

	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}

func splitLong(text string) []string {
	var chunks []string
	sentences := strings.SplitAfter(text, ". ")

	var current strings.Builder
	for _, s := range sentences {
		if current.Len() == 0 {
			current.WriteString(s)
			continue
		}
		candidate := current.String() + s
		if utf8.RuneCountInString(candidate) <= twitterMaxChars {
			current.WriteString(s)
		} else {
			chunks = append(chunks, strings.TrimSpace(current.String()))
			current.Reset()
			current.WriteString(s)
		}
	}
	if current.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(current.String()))
	}

	return chunks
}

var boldMap = map[rune]rune{
	'A': '𝗔', 'B': '𝗕', 'C': '𝗖', 'D': '𝗗', 'E': '𝗘', 'F': '𝗙', 'G': '𝗚',
	'H': '𝗛', 'I': '𝗜', 'J': '𝗝', 'K': '𝗞', 'L': '𝗟', 'M': '𝗠', 'N': '𝗡',
	'O': '𝗢', 'P': '𝗣', 'Q': '𝗤', 'R': '𝗥', 'S': '𝗦', 'T': '𝗧', 'U': '𝗨',
	'V': '𝗩', 'W': '𝗪', 'X': '𝗫', 'Y': '𝗬', 'Z': '𝗭',
	'a': '𝗮', 'b': '𝗯', 'c': '𝗰', 'd': '𝗱', 'e': '𝗲', 'f': '𝗳', 'g': '𝗴',
	'h': '𝗵', 'i': '𝗶', 'j': '𝗷', 'k': '𝗸', 'l': '𝗹', 'm': '𝗺', 'n': '𝗻',
	'o': '𝗼', 'p': '𝗽', 'q': '𝗾', 'r': '𝗿', 's': '𝘀', 't': '𝘁', 'u': '𝘂',
	'v': '𝘃', 'w': '𝘄', 'x': '𝘅', 'y': '𝘆', 'z': '𝘇',
	'0': '𝟬', '1': '𝟭', '2': '𝟮', '3': '𝟯', '4': '𝟰',
	'5': '𝟱', '6': '𝟲', '7': '𝟳', '8': '𝟴', '9': '𝟵',
}

func toBold(s string) string {
	var buf strings.Builder
	for _, r := range s {
		if bold, ok := boldMap[r]; ok {
			buf.WriteRune(bold)
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
