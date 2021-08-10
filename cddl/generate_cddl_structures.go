package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"go/format"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/iancoleman/strcase"
)

type generic struct {
	block       []interface{}
	templateVar string
	UsedTypes   map[string][]interface{}
}

type composition struct {
	Name   string
	Fields []field
}

type field struct {
	Name string
	Type string
}

var generics map[string]generic
var compositions map[string]composition
var additionalTypes map[string]map[string][]interface{}
var addedCompositions map[string]int

func isIgnoredType(v interface{}) bool {
	switch name := v.(type) {
	case string:
		switch name {
		case "int64", "byte", "address":
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func preprocessName(s string) string {
	return strcase.ToCamel(strings.TrimLeft(s, "$"))
}

func getName(block interface{}) string {
	switch v := block.(type) {
	case string:
		return preprocessName(v)
	case []interface{}:
		switch v[0] {
		case "name":
			return preprocessName(v[1].(string))
		case "rep":
			return preprocessName(v[3].([]interface{})[1].([]interface{})[1].(string)) + "Map"
		case "number":
			return "Number"
		case "ary":
			switch v[1].([]interface{})[0] {
			case "rep":
				return preprocessName(v[1].([]interface{})[3].([]interface{})[2].([]interface{})[1].(string)) + "List"
			case "mem":
				return preprocessName(v[1].([]interface{})[2].([]interface{})[1].(string)) + "List"
			default:
				return "SomeArray"
			}
		case "gen":
			//return preprocessName(v[2].([]interface{})[1].(string) + "_" + v[1].(string))
			return getGeneric(v[1].(string), v[2].([]interface{}))
		case "map":
			return preprocessName(v[1].([]interface{})[3].([]interface{})[1].([]interface{})[1].(string)) + "Map"
		case "op":
			return preprocessName(v[2].([]interface{})[1].(string))
		default:
			log.Panicf("unknown name struct: %v", block)
		}
	default:
		log.Panicf("unknown name struct: %v", block)
	}
	return ""
}

func getType(block []interface{}, replaceTypeFrom string, replateTypeTo []interface{}) string {
	//fmt.Println(block)
	switch block[0] {
	case "name":
		switch block[1] {
		case replaceTypeFrom:
			return getType(replateTypeTo, "", nil)
		case "uint":
			return "uint"
		case "bool":
			return "bool"
		case "nint", "int":
			return "int"
		case "float16", "float32":
			return "float32"
		case "float64", "float":
			return "float64"
		case "bstr", "bytes":
			return "[]byte"
		case "tstr", "text":
			return "string"
		case "int64":
			return "int64"
		case "address":
			return "Address"
		default:
			return preprocessName(block[1].(string))
		}
	case "ary":
		subBlock := block[1].([]interface{})
		switch subBlock[0] {
		case "rep":
			memInfo := subBlock[3].([]interface{})
			switch memInfo[0] {
			case "mem":
				return "[]" + getType(memInfo[2].([]interface{}), replaceTypeFrom, replateTypeTo)
			default:
				panic("unknown list")
			}
		case "mem":
			switch subBlock[0] {
			case "mem":
				return "[]" + getType(subBlock[2].([]interface{}), replaceTypeFrom, replateTypeTo)
			default:
				panic("unknown list")
			}
		default:
			panic("unknown list")
		}
	case "map":
		res := "unknown map"
		subBlock := block[1].([]interface{})
		switch subBlock[0] {
		case "rep":
			memInfo := subBlock[3].([]interface{})
			switch memInfo[0] {
			case "mem":
				firstBlock := memInfo[1].([]interface{})
				secondBlock := memInfo[2].([]interface{})
				return "hash_map.HashMap //map[" + getType(firstBlock, replaceTypeFrom, replateTypeTo) + "]" + getType(secondBlock, replaceTypeFrom, replateTypeTo)
			}
		}
		return res
	case "tcho":
		switch firstBlock := block[1].(type) {
		case []interface{}:
			switch secondBlock := block[2].(type) {
			case []interface{}:
				switch firstBlock[0] {
				case "name":
					switch secondBlock[0] {
					case "name":
						switch secondBlock[1] {
						case "null":
							return "*" + getName(firstBlock[1])
						default:
							panic("unknown tcho struct")
						}
					default:
						panic("unknown type in tcho")
					}
				case "number":
					return "int"
				default:
					panic("unknown type in tcho")
				}
			default:
				panic("unknown tcho struct")

			}
		default:
			panic("unknown tcho struct")
		}
	case "gen":
		return getGeneric(block[1].(string), block[2].([]interface{}))
	case "mem":
		//log.Println(block[2])
		return getType(block[2].([]interface{}), replaceTypeFrom, replateTypeTo)
	case "op":
		return getType(block[2].([]interface{}), replaceTypeFrom, replateTypeTo)
	case "seq": //may be composition
		if isComposition(block) {
			return "composition"
		} else {
			panic("unknown seq")
		}
	case "number":
		return "int32"
	default:
		return "ERROR"
	}
}

func addGeneric(blocks []interface{}) string {
	generics[blocks[1].([]interface{})[1].(string)] = generic{
		block:       blocks[2].([]interface{}),
		templateVar: blocks[1].([]interface{})[2].(string),
		UsedTypes:   make(map[string][]interface{}),
	}
	return ""
}

func isComposition(blocks []interface{}) bool {
	switch blocks[0] {
	case "seq":
		if len(blocks) < 2 {
			return false
		}

		switch blocks[1].([]interface{})[0] {
		case "mem":
			switch blocks[1].([]interface{})[2].([]interface{})[0] {
			case "number":
				return true
			default:
				return false
			}
		default:
			return false
		}
	default:
		return false
	}
}

func addSupportType(block []interface{}, prefix string) string {
	if additionalTypes == nil {
		additionalTypes = make(map[string]map[string][]interface{})
	}
	typeName := fmt.Sprintf("%sAdditionalType%d", prefix, len(additionalTypes[prefix]))
	specialName := typeName
	switch block[0] {
	case "name":
		return getType(block, "", nil)
	}
	if len(block) > 2 {
		switch subBlock := block[2].(type) {
		case []interface{}:
			switch subBlock[0] {
			case "name":
				//typeName = getType(subBlock, "", nil)
				typeName = getName(subBlock)
				specialName = subBlock[1].(string)
			}
		}
	}
	if _, ok := additionalTypes[prefix]; !ok {
		additionalTypes[prefix] = make(map[string][]interface{})
	}
	additionalTypes[prefix][specialName] = []interface{}{
		"=",
		[]interface{}{"name", typeName},
		block,
	}
	return typeName
}

func parseComposition(blocks []interface{}) *composition {
	switch blocks[0] {
	case "seq":
		//res := composition{}
		switch v := blocks[1].(type) {
		case []interface{}:
			switch v[0] {
			case "mem":
				switch vv := v[2].(type) {
				case []interface{}:
					switch vv[0] {
					case "number":
						//res.Key = vv[1].(string)
					default:
						return nil
					}
				default:
					return nil
				}
			default:
				return nil
			}
		default:
			return nil
		}
		return nil
	}
	return nil
}

func addComposition(blocks []interface{}) string {
	if compositions == nil {
		compositions = make(map[string]composition)
	}
	name := blocks[1].([]interface{})[1].(string)
	seq := blocks[2].([]interface{})
	switch seq[0] {
	case "seq":
		comp := composition{Name: name}
		for index, rawField := range seq[1:] {
			comp.Fields = append(comp.Fields, field{
				fmt.Sprintf("V%d", index),
				getType(rawField.([]interface{}), "", nil),
			})
		}
		compositions[name] = comp
	default:
		panic("unknown composition")
	}
	return ""
}

func getGeneric(genName string, block []interface{}) string {
	if gen, ok := generics[genName]; ok {
		switch block[0] {
		case "name":
			gen.UsedTypes[block[1].(string)] = block
		default:
			panic("unknown name block")
		}
		generics[genName] = gen
		return preprocessName(genName + "_" + getName(block))
	} else {
		panic("unknown generic")
	}
}

func setNames(block interface{}, cnt int) interface{} {
	switch blocks := block.(type) {
	case []interface{}:
		switch blocks[0] {
		case "seq", "ary":
			for i, field := range blocks[1].([]interface{}) {
				blocks[1].([]interface{})[i] = setNames(field, i)
			}
			return blocks
		case "mem":
			switch blocks[1] {
			case nil:
				blocks[1] = []interface{}{"name", fmt.Sprintf("V%d", cnt)}
			}
			return blocks
		default:
			for i, block := range blocks {
				blocks[i] = setNames(block, 0)
			}
			return blocks
		}
	default:
		return blocks
	}
}

var uniqueNames map[string]bool

func preprocessCdll(path string) ([]interface{}, error) {
	jsonRaw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cddl []interface{}
	err = json.Unmarshal(jsonRaw, &cddl)
	if err != nil {
		return nil, err
	}

	if uniqueNames == nil {
		uniqueNames = make(map[string]bool)
	}
	//searching generics
	var forRemoving []int
	for i, elem := range cddl {
		if i == 0 {
			continue
		}
		lelem := elem.([]interface{})
		switch lelem[1].([]interface{})[0] {
		case "gen":
			addGeneric(lelem)
			forRemoving = append(forRemoving, i)
		case "name":
			// set names for duplicates
			cddl[i] = setNames(cddl[i], 0)

			switch name := lelem[1].([]interface{})[1].(type) {
			case string:
				if _, ok := uniqueNames[name]; ok {
					forRemoving = append(forRemoving, i)
					continue
				} else {
					uniqueNames[name] = true
				}
			}
			if isIgnoredType(lelem[1].([]interface{})[1]) {
				forRemoving = append(forRemoving, i)
				continue
			}
		}
		if isComposition(lelem[2].([]interface{})) {
			forRemoving = append(forRemoving, i)
			addComposition(lelem)
		}
	}

	for i := range forRemoving {
		index := forRemoving[len(forRemoving)-1-i]
		cddl = append(cddl[:index], cddl[index+1:]...)
	}
	return cddl, nil
}

func clearTag(blocks []interface{}) []interface{} {
	if len(blocks) < 3 {
		return blocks
	}
	switch blocks[2].([]interface{})[0] {
	case "prim":
		return []interface{}{blocks[0], blocks[1], blocks[2].([]interface{})[3]}
	default:
		return blocks
	}
}

type cddlFile struct {
	Path       string
	CddlSchema []interface{}
}

func getComposition(field []interface{}, prefix string) *composition {
	if addedCompositions == nil {
		addedCompositions = make(map[string]int)
	}
	name := fmt.Sprintf("%sComposition%d", prefix, addedCompositions[prefix])
	switch field[0] {
	case "mem":
		//addComposition([]interface{}{
		//	"=", []interface{}{"name", name}, []interface{}{
		//		"seq", field,
		//	},
		//})
		//addedCompositions[prefix] += 1
		name = getType(field[2].([]interface{}), "", nil)
		//switch name {
		//case "Number":
		//
		//}
		return &composition{
			Name:   name,
			Fields: nil,
		}
	case "seq":
		addComposition([]interface{}{
			"=", []interface{}{"name", name}, field,
		})
		addedCompositions[prefix] += 1
	default:
		panic("unknown error")
	}
	if res, ok := compositions[name]; ok {
		return &res
	}
	return nil
}

func main() {
	generics = make(map[string]generic)

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"ToCamel":          preprocessName,
		"GetType":          getType,
		"GetName":          getName,
		"AddGeneric":       addGeneric,
		"GetGeneric":       getGeneric,
		"ClearTag":         clearTag,
		"IsComposition":    isComposition,
		"AddSupportType":   addSupportType,
		"GetComposition":   getComposition,
		"AddComposition":   addComposition,
		"ParseComposition": parseComposition,
	}

	var cddl_files []cddlFile
	err := filepath.WalkDir("json-files", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() && strings.Contains(d.Name(), "json") {
			cddl_files = append(cddl_files, cddlFile{Path: s})
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	t := template.Must(template.New("shelley_cddl.tmpl").Funcs(funcMap).ParseFiles("shelley_cddl.tmpl"))
	tComposition := template.Must(template.New("composition_shelley_cddl.tmpl").Funcs(funcMap).ParseFiles("composition_shelley_cddl.tmpl"))
	commonBuf := &bytes.Buffer{}
	w := bufio.NewWriter(commonBuf)
	fe, err := os.Create("../tests/errors.txt")
	we := bufio.NewWriter(fe)

	for i, cddl_file := range cddl_files {
		res, err := preprocessCdll(cddl_file.Path)
		cddl_files[i].CddlSchema = res
		if err != nil {
			log.Fatalln(err)
		}
	}

	_, err = commonBuf.Write([]byte("package common\n"))
	if err != nil {
		panic("unexpected error")
	}

	_, err = commonBuf.Write([]byte("import (\n"))
	if err != nil {
		panic("unexpected error")
	}

	_, err = commonBuf.Write([]byte("\"github.com/fivebinaries/go-cardano-serialization/hash_map\"\n"))
	if err != nil {
		panic("unexpected error")
	}

	_, err = commonBuf.Write([]byte(")\n"))
	if err != nil {
		panic("unexpected error")
	}

	for _, cddl_file := range cddl_files {
		for _, elem := range cddl_file.CddlSchema[1:] {
			buf := &bytes.Buffer{}
			err := t.Execute(buf, elem)
			if err != nil {
				log.Println("executing template:", err)
				_, err = we.Write(buf.Bytes())
				we.Write([]byte("\n==============\n"))
			} else {
				_, err = w.Write(buf.Bytes())
				if err != nil {
					log.Println("writing error:", err)
				}
			}
		}
	}

	for _, innerMap := range additionalTypes {
		for name, elem := range innerMap {
			if len(name) > 1 {
			}
			buf := &bytes.Buffer{}
			err := t.Execute(buf, elem)
			if err != nil {
				log.Println("executing template:", err)
				_, err = we.Write(buf.Bytes())
				we.Write([]byte("\n==============\n"))
			} else {
				_, err = w.Write(buf.Bytes())
				if err != nil {
					log.Println("writing error:", err)
				}
			}
		}
	}

	for _, comp := range compositions {
		buf := &bytes.Buffer{}
		err := tComposition.Execute(buf, comp)
		if err != nil {
			log.Println("executing template:", err)
			_, err = we.Write(buf.Bytes())
			we.Write([]byte("\n==============\n"))
		} else {
			_, err = w.Write(buf.Bytes())
			if err != nil {
				log.Println("writing error:", err)
			}
		}
	}

	for genName, generic := range generics {
		for _, usedType := range generic.UsedTypes {
			buf := &bytes.Buffer{}
			genericGo := getType(generic.block, generic.templateVar, usedType)
			kek := fmt.Sprintf("type %s %s\n", preprocessName(genName+"_"+getName(usedType)), genericGo)
			_, err := w.Write([]byte(kek))
			if err != nil {
				log.Println("executing template:", err)
				_, err = we.Write(buf.Bytes())
				we.Write([]byte("\n==============\n"))
			} else {
				_, err = w.Write(buf.Bytes())
				if err != nil {
					log.Println("writing error:", err)
				}
			}
		}
	}

	err = w.Flush()
	if err != nil {
		panic(err)
	}
	f, err := os.Create(filepath.Join("..", "types", "cddl_serialization.gen.go"))
	if err != nil {
		log.Fatalln(err)
	}

	p, err := format.Source(commonBuf.Bytes())
	if err != nil {
		log.Panicf("unexpected error: %v", err)
	}
	_, err = f.Write(p)
	if err != nil {
		log.Panicf("unexpected error: %v", err)
	}
	err = we.Flush()
	if err != nil {
		log.Panicf("unexpected error: %v", err)
	}

}
