	"bytes"
	"os"
	"path"
	"strings"
	Header string // free-form text
	File   []*File
	Verb             Verb
	Src              string // source for Verb == Copy, Verb == Rename
	Dst              string
	OldMode, NewMode int // 0 indicates not used
	Diff                 // changes to data; == NoDiff if operation does not edit file
	Add    Verb = "add"
	Copy   Verb = "copy"
	Delete Verb = "delete"
	Edit   Verb = "edit"
	Rename Verb = "rename"
	Apply(old []byte) (new []byte, err os.Error)
func (e SyntaxError) String() string { return string(e) }
	text, files := sections(text, "Index: ")
	set := &Set{string(text), make([]*File, len(files))}
		p := new(File)
		set.File[i] = p
		s, raw, _ := getLine(raw, 1)
			p.Dst = string(bytes.TrimSpace(s[7:]))
			goto HaveName
			str := string(bytes.TrimSpace(s))
			i := strings.LastIndex(str, " b/")
				p.Dst = str[i+3:]
				goto HaveName
		return nil, SyntaxError("unexpected patch header line: " + string(s))
		p.Dst = path.Clean(p.Dst)
		p.Verb = Edit
			oldraw := raw
			var l []byte
			l, raw, _ = getLine(raw, 1)
			l = bytes.TrimSpace(l)
				p.NewMode = m
				p.Verb = Add
				continue
				p.OldMode = m
				p.Verb = Delete
				p.Src = p.Dst
				p.Dst = ""
				continue
				p.OldMode = m
				continue
				p.OldMode = m
				continue
				p.NewMode = m
				continue
				p.Src = string(s)
				p.Verb = Rename
				continue
				p.Verb = Rename
				continue
				p.Src = string(s)
				p.Verb = Copy
				continue
				p.Verb = Copy
				continue
				diff, err := ParseTextDiff(oldraw)
				p.Diff = diff
				break
				diff, err := ParseGitBinary(oldraw)
				p.Diff = diff
				break
			return nil, SyntaxError("unexpected patch header line: " + string(l))
	return set, nil
	rest = data
	ok = true
		nl := bytes.Index(rest, newline)
			rest = nil
			ok = false
			break
		rest = rest[nl+1:]
	first = data[0 : len(data)-len(rest)]
	return
	n := 0
		nl := bytes.Index(b, newline)
		b = b[nl+1:]
	sect := make([][]byte, n+1)
	n = 0
			sect[n] = text[0 : len(text)-len(b)]
			n++
			text = b
		nl := bytes.Index(b, newline)
			sect[n] = text
			break
		b = b[nl+1:]
	return sect[0], sect[1:]
	return s[len(t):], true
	var i int
	return n, s[i:], true
	_, ok := skip(s, t)
	return ok
func splitLines(s []byte) [][]byte { return bytes.SplitAfter(s, newline, 0) }