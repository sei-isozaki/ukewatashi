package main

import (
	"bytes"
	"fmt"
	"html/template"
	"miidas/domain/connect/enum"
	"miidas/domain/core/service/user"
	"reflect"
)

// テストで角煮ん用のStruct
type TestSampleRequest struct {
	Sample2
	Hoge   int
	HogeSt Sample
}
type Sample2 struct {
	A int
	B Sample
}
type Sample struct {
	Ho  int
	Huu enum.ApplyType
}

// テスト用のリクエストハンドラメソッド
//+ request url:'' validator:xxxx
func Hoo(i TestSampleRequest) *user.TextSearchResult {
	return nil
}

// リクエストハンドラの定義。これ以外にもあるかもね
type Endpoint struct {
	Validator interface{}
	Handler   interface{}
	Url       string
}

// リクエストハンドラを取得して、API I/Fドキュメント（？）を作る
func main() {

	// Astから取ってくるEndpointのMAP TODO
	var Mx = map[string]Endpoint{
		"/usl/hoge/user/:idxxxx": Endpoint{Handler: Hoo},
	}

	// Endpoint毎にループ
	for _, v := range Mx {

		// ハンドラの型をチェック
		handlerValue := reflect.ValueOf(v.Handler)
		if handlerValue.Kind() != reflect.Func {
			// 関数じゃないならエラー
			panic("f must be func.")
		}

		handlerType := reflect.TypeOf(v.Handler)
		p(handlerType)
		fmt.Printf("%t\n", handlerType)
		p("// 引数の一覧")

		d := StData{
			Classes: []StClass{},
		}

		for i := 0; i < handlerType.NumIn(); i++ {
			// 引数の型の取得
			p(handlerType.In(i))

			d.Classes = append(d.Classes, structTypeToSt(handlerType.In(i))...)
		}
		outputA(d)

		p("// 返り値の一覧")
		for i := 0; i < handlerType.NumOut(); i++ {
			// 返り値の型の取得

			p(handlerType.Out(i))
		}

	}
}

// これの戻り値全部。
type StData struct {
	Classes []StClass
}

// クラス一つ分
type StClass struct {
	Name   string
	Fields []StField
}

// フィールド
type StField struct {
	Name    string
	Type    string
	Comment string
}

func structTypeToSt(t reflect.Type) []StClass {

	r := StClass{
		Name:   t.Name(),
		Fields: []StField{},
	}
	ret := []StClass{r}

	fmt.Println("-------------")
	fmt.Println(t)
	fmt.Println("フィールド数")
	fmt.Println(t.NumField())

	for i := 0; i < t.NumField(); i++ {

		r.Fields = append(r.Fields,
			StField{Name: t.Field(i).Name,
				Type: string(t.Field(i).Type.Name())})
	}
	fmt.Println(ret)

	return ret
}

func outputA(t StData) {

	f := func(filename, tpl string, data StData) {
		t, e := template.New("hoge").Parse(tpl)
		if e != nil {
			panic(e)
		}

		buf := &bytes.Buffer{}

		t.Execute(buf, data)

		outputString := buf.String()
		fmt.Println(outputString)

		// Write to file.
		//e = ioutil.WriteFile(filename, []byte(outputString), 0644)
		//if e != nil {
		//	panic(e)
		//}
	}
	f("hoge"+".kt", kotlinTemplate, t)
	//f(output+".swift", swiftTemplate, pa)
	//f(output+".js", jsTemplate, pa)

}

const kotlinTemplate = `
{{range $index, $p := .Classes}}

class {{$p.Name}}{
{{range $index, $pp := $p.Fields}}
	{{$pp.Name}} {{$pp.Type}}
{{end}}
}
{{end}}
`

func p(a ...interface{}) {
	fmt.Println(a...)
}
