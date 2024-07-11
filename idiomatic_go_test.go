package learning

/*
From the book:
	Learning Go
	An Idiomatic Approach to Real-World Go Programming
*/

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func countTo(max int) (<-chan int, func()) {
	ch := make(chan int)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go func() {
		for i := 0; i < max; i++ {
			select {
			case <-done:
				return
			case ch <- i:
			}
		}
		close(ch)
	}()

	return ch, cancel
}

func TestCountTo(t *testing.T) {
	ch, cancel := countTo(10)
	for i := range ch {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}
	cancel()
	t.Log("Finish")
}

type pressureGauge struct {
	ch chan struct{}
}

func (pg *pressureGauge) process(f func()) error {
	select {
	case <-pg.ch:
		f()
		pg.ch <- struct{}{}
		return nil
	default:
		return errors.New("no more capacity")
	}
}

func newPG(limit int) *pressureGauge {
	ch := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		ch <- struct{}{}
	}
	return &pressureGauge{ch: ch}
}

func TestNewPG(t *testing.T) {
	pg := newPG(5)
	t.Log(len(pg.ch))
}

func doThingThatShouldBeLimited() string {
	time.Sleep(2 * time.Second)
	return "done"
}

func TestProcess(t *testing.T) {
	pg := newPG(10)
	err := pg.process(func() {
		for i := 0; i < 10; i++ {
			t.Log(doThingThatShouldBeLimited())
		}
	})
	if err != nil {
		t.Fatal(err)
	}
}

func doSomeWork() (int, error) {
	time.Sleep(1 * time.Second)
	return 7, nil
}

func timeLimit() (int, error) {
	var result int
	var err error
	done := make(chan struct{})

	go func() {
		result, err = doSomeWork()
		close(done)
	}()

	select {
	case <-done:
		return result, err
	case <-time.After(2 * time.Second):
		return 0, errors.New("work timed out")
	}
}

func TestTimeLimit(t *testing.T) {
	res, err := timeLimit()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestDefer(t *testing.T) {
	t.Log("beginning of test")
	i := rand.Intn(10)
	if i != 7 {
		defer func() {
			t.Log("call defer function")
		}()
		t.Log(i)
	}
	t.Log("end of test")
}

type data struct {
	a int
	b int
}

func getResultA(_ context.Context, a int) (int, error) {
	if a == 0 {
		return 0, errors.New("a cannot be 0")
	}
	return a * 2, nil
}

func getResultB(_ context.Context, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("b cannot be 0")
	}
	return b * 2, nil
}

func getResultC(_ context.Context, c cIn) (int, error) {
	if &c == nil {
		return 0, errors.New("c is nil")
	}
	return c.a + c.b, nil
}

type cIn struct {
	a int
	b int
}

type processor struct {
	aOut chan int
	bOut chan int
	cOut chan int
	cIn  chan cIn
	errs chan error
}

func (p *processor) launch(ctx context.Context, data data) {
	go func() {
		aOut, err := getResultA(ctx, data.a)
		if err != nil {
			p.errs <- err
			return
		}
		p.aOut <- aOut
	}()
	go func() {
		bOut, err := getResultB(ctx, data.b)
		if err != nil {
			p.errs <- err
			return
		}
		p.bOut <- bOut
	}()
	go func() {
		select {
		case <-ctx.Done():
			return
		case inputC := <-p.cIn:
			cOut, err := getResultC(ctx, inputC)
			if err != nil {
				p.errs <- err
				return
			}
			p.cOut <- cOut
		}
	}()

}

func (p *processor) waitForAB(ctx context.Context) (cIn, error) {
	var inputC cIn
	count := 0
	for count < 2 {
		select {
		case a := <-p.aOut:
			inputC.a = a
			count++
		case b := <-p.bOut:
			inputC.b = b
			count++
		case err := <-p.errs:
			return cIn{}, err
		case <-ctx.Done():
			return cIn{}, ctx.Err()
		}
	}
	return inputC, nil
}

func (p *processor) waitForC(ctx context.Context) (int, error) {
	select {
	case out := <-p.cOut:
		return out, nil
	case err := <-p.errs:
		return 0, err
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func TestWaitForC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	p := processor{
		aOut: make(chan int, 1),
		bOut: make(chan int, 1),
		cOut: make(chan int, 1),
		cIn:  make(chan cIn, 1),
		errs: make(chan error, 2),
	}
	p.launch(ctx, data{a: 1, b: 2})

	inputC, err := p.waitForAB(ctx)
	if err != nil {
		t.Fatal(err)
	}
	p.cIn <- inputC

	out, err := p.waitForC(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(out)
}

func TestTypeAssertion(t *testing.T) {
	a := map[string]any{
		"a": 1,
	}

	i, ok := a["a"].(string)
	if !ok {
		t.Fatal("value is not a string")
	}

	t.Log(i)
}

func countLetters(r io.Reader) (map[string]int, error) {
	buf := make([]byte, 2048)
	out := map[string]int{}
	for {
		n, err := r.Read(buf)
		for _, b := range buf[:n] {
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func TestCountLetters(t *testing.T) {
	s := "The quick brown fox jumped over the lazy dog"
	r := strings.NewReader(s)
	counts, err := countLetters(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(counts)
}

func openFile(path string) (*os.File, func(), error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return f, func() {
		f.Close()
	}, nil
}

func TestCountLettersFromText(t *testing.T) {
	f, closer, err := openFile("vim.txt")
	defer closer()

	counts, err := countLetters(f)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range counts {
		t.Logf("%s: %d", k, v)
	}
}

func fnThatReturnClosure() func() {
	i := 17
	return func() {
		fmt.Println(i)
	}
}

func TestClosure(t *testing.T) {
	f := fnThatReturnClosure()
	f()
}

func TestDuration(t *testing.T) {
	pomodoro := 25*time.Minute + 5*time.Minute
	t.Log(pomodoro.String())
	t.Logf("1 pomodoro = %.0f minutes", pomodoro.Minutes())
}

type order struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func TestJSONMarshal(t *testing.T) {
	s := `{"id": 1, "name": "Caffe Latte"}`
	var o order

	err := json.Unmarshal([]byte(s), &o)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(o)

	b, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestJSONEncoderMore(t *testing.T) {
	const data = `
		{"name": "Fred", "age": 40}
		{"name": "Mary", "age": 21}
		{"name": "Pat", "age": 30}
	`
	var tt struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	dec := json.NewDecoder(strings.NewReader(data))
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	for dec.More() {
		err := dec.Decode(&tt)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(tt)

		err = enc.Encode(tt)
		if err != nil {
			t.Fatal(err)
		}
	}

	out := b.String()
	t.Log(out)
}

func TestOctal(t *testing.T) {
	oct := 0o123
	t.Log(oct)

	hex := 0x123
	t.Log(hex)
}

func TestHTTPRequest(t *testing.T) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://jsonplaceholder.typicode.com/todos/1",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatal(fmt.Sprintf("unexpected status: got %v", res.Status))
	}
	t.Log(res.Header.Get("Content-Type"))

	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v\n", data)
}

func slowServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("slow response"))
	}))
	return s
}

func fastServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("error") == "true" {
			w.Write([]byte("error"))
			return
		}
		w.Write([]byte("ok"))
	}))
	return s
}

var client = &http.Client{}

func callBoth(ctx context.Context, errVal string, slowURL string, fastURL string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := callServer(ctx, "slow", slowURL)
		if err != nil {
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		err := callServer(ctx, "fast", fastURL+"?error="+errVal)
		if err != nil {
			cancel()
		}
	}()
	wg.Wait()

	fmt.Println("done with both")
}

func callServer(ctx context.Context, label string, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(label, "request error:", err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(label, "response error:", err)
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(label, "read error:", err)
		return err
	}
	result := string(data)
	if result != "" {
		fmt.Println(label, "result:", result)
	}
	if result == "error" {
		fmt.Println("canceling from", label)
		return errors.New("error happened")
	}
	return nil
}

func TestContextWithCancel(t *testing.T) {
	ss := slowServer()
	defer ss.Close()
	fs := fastServer()
	defer fs.Close()

	ctx := context.Background()
	callBoth(ctx, os.Getenv("ERROR"), ss.URL, fs.URL)
}

func longRunningThing(ctx context.Context, data string) (string, error) {
	time.Sleep(2 * time.Second)
	return data, nil
}

func longRunningThingManager(ctx context.Context, data string) (string, error) {
	type wrapper struct {
		result string
		err    error
	}
	ch := make(chan wrapper, 1)
	go func() {
		// do the long-running thing
		result, err := longRunningThing(ctx, data)
		ch <- wrapper{result, err}
	}()
	select {
	case data := <-ch:
		return data.result, data.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func TestLongRunningThingManager(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := longRunningThingManager(ctx, "Kore")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

var testTime time.Time

//func TestMain(m *testing.M) {
//	fmt.Println("Set up stuff for tests here")
//	testTime = time.Now()
//	exitVal := m.Run()
//	fmt.Println("Clean up stuff after tests here")
//	os.Exit(exitVal)
//}

func TestFirst(t *testing.T) {
	fmt.Println("TestFirst uses stuff set up in TestMain", testTime)
}

func TestSecond(t *testing.T) {
	fmt.Println("TestSecond also uses stuff set up in TestMain", testTime)
}

func TestReflectTypeOf(t *testing.T) {
	var x int
	xt := reflect.TypeOf(x)
	t.Log(xt.Name())

	type foo struct {
		A int    `myTag:"val1"`
		B string `myTag:"val1"`
	}

	f := foo{}
	ft := reflect.TypeOf(f)
	t.Log(ft.Name())
	for i := 0; i < ft.NumField(); i++ {
		curField := ft.Field(i)
		t.Log(curField.Name, curField.Type.Name(), curField.Tag.Get("myTag"))
	}

	xpt := reflect.TypeOf(&x)
	t.Log(xpt.Name())
	t.Log(xpt.Kind())
	t.Log(xpt.Elem().Name())
	t.Log(xpt.Elem().Kind())
}

func TestReflectValue(t *testing.T) {
	s := []string{"a", "b", "c"}
	sv := reflect.ValueOf(s)
	s2 := sv.Interface().([]string)
	t.Log(s2)

	iv := 4
	ivv := reflect.ValueOf(iv)
	iv2 := ivv.Int()
	t.Log(iv2)
}

func TestReflectSetValue(t *testing.T) {
	x := 10
	xv := reflect.ValueOf(&x)
	xe := xv.Elem()
	xe.SetInt(20)
	t.Log(x)
}

func TestMakeSlice(t *testing.T) {
	stringType := reflect.TypeOf((*string)(nil)).Elem()
	stringSliceType := reflect.TypeOf((*[]string)(nil)).Elem()
	ssv := reflect.MakeSlice(stringSliceType, 0, 10)
	sv := reflect.New(stringType).Elem()
	sv.SetString("oyen")
	ssv = reflect.Append(ssv, sv)
	sv.SetString("kore")
	ssv = reflect.Append(ssv, sv)
	ss := ssv.Interface().([]string)
	t.Log(ss)
}

// Marshal maps all structs in a slice of structs to a slice of slice of strings.
// The first row written is the header with the column names.
func Marshal(v interface{}) ([][]string, error) {
	sliceVal := reflect.ValueOf(v)
	if sliceVal.Kind() != reflect.Slice {
		return nil, errors.New("must be a slice of structs")
	}
	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return nil, errors.New("must be a slice of structs")
	}
	var out [][]string
	header := marshalHeader(structType)
	out = append(out, header)
	for i := 0; i < sliceVal.Len(); i++ {
		row, err := marshalOne(sliceVal.Index(i))
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

func marshalHeader(vt reflect.Type) []string {
	var row []string
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		if curTag, ok := field.Tag.Lookup("csv"); ok {
			row = append(row, curTag)
		}
	}
	return row
}

func marshalOne(vv reflect.Value) ([]string, error) {
	var row []string
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		fieldVal := vv.Field(i)
		if _, ok := vt.Field(i).Tag.Lookup("csv"); !ok {
			continue
		}
		switch fieldVal.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			row = append(row, strconv.FormatInt(fieldVal.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			row = append(row, strconv.FormatUint(fieldVal.Uint(), 10))
		case reflect.String:
			row = append(row, fieldVal.String())
		case reflect.Bool:
			row = append(row, strconv.FormatBool(fieldVal.Bool()))
		default:
			return nil, fmt.Errorf("cannot handle field of kind %v", fieldVal.Kind())
		}
	}
	return row, nil
}

// Unmarshal maps all the rows of data in slice of slice of strings into a slice of structs.
// The first row is assumed to be the header with the column names.
func Unmarshal(data [][]string, v interface{}) error {
	sliceValPtr := reflect.ValueOf(v)
	if sliceValPtr.Kind() != reflect.Ptr {
		return errors.New("must be a pointer to a slice of structs")
	}
	sliceVal := sliceValPtr.Elem()
	if sliceVal.Kind() != reflect.Slice {
		return errors.New("must be a pointer to a slice of structs")
	}
	structType := sliceVal.Type().Elem()
	if structType.Kind() != reflect.Struct {
		return errors.New("must be a pointer to a slice of structs")
	}

	// assume the first row is a header
	header := data[0]
	namePos := make(map[string]int, len(header))
	for k, v := range header {
		namePos[v] = k
	}

	for _, row := range data[1:] {
		newVal := reflect.New(structType).Elem()
		err := unmarshalOne(row, namePos, newVal)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, newVal))
	}
	return nil
}

func unmarshalOne(row []string, namePos map[string]int, vv reflect.Value) error {
	vt := vv.Type()
	for i := 0; i < vv.NumField(); i++ {
		typeField := vt.Field(i)
		pos, ok := namePos[typeField.Tag.Get("csv")]
		if !ok {
			continue
		}
		val := row[pos]
		field := vv.Field(i)
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(i)
		case reflect.String:
			field.SetString(val)
		case reflect.Bool:
			b, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			field.SetBool(b)
		default:
			return fmt.Errorf("cannot handle field of kind %v", field.Kind())
		}
	}
	return nil
}

type MyData struct {
	Name   string `csv:"name"`
	HasPet bool   `csv:"has_pet"`
	Age    int    `csv:"age"`
}

func TestMarshaller(t *testing.T) {
	data := `name,age,has_pet
Jon,"100",true
"Fred ""The Hammer"" Smith",42,false
Martha,37,"true"
`
	r := csv.NewReader(strings.NewReader(data))
	allData, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	var entries []MyData
	Unmarshal(allData, &entries)
	fmt.Println(entries)

	//now to turn entries into output
	out, err := Marshal(entries)
	if err != nil {
		panic(err)
	}
	sb := &strings.Builder{}
	w := csv.NewWriter(sb)
	w.WriteAll(out)
	fmt.Println(sb)
}

func makeTimedFunction(f any) any {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		panic("expects a function")
	}
	vf := reflect.ValueOf(f)
	wrapperF := reflect.MakeFunc(rf, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := vf.Call(in)
		end := time.Now()
		fmt.Printf("calling %s took %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}

func timeMe(a int) int {
	time.Sleep(time.Duration(a) * time.Millisecond)
	result := a * 2
	return result
}

func TestMakeTimedFunc(t *testing.T) {
	timed := makeTimedFunction(timeMe).(func(int) int)
	t.Log(timed(2))
}

// convert float64 to uint64
func Float64bits(floatVal float64) uint64 {
	// the short way is:
	//     return *(*uint64)(unsafe.Pointer(&floatVal))

	// Take a pointer to the float64 value stored in f.
	floatPtr := &floatVal

	// Convert the *float64 to an unsafe.Pointer.
	unsafePtr := unsafe.Pointer(floatPtr)

	// Convert the unsafe.Pointer to *uint64.
	uintPtr := (*uint64)(unsafePtr)

	// Dereference the *uint64, yielding a uint64 value.
	uintVal := *uintPtr

	return uintVal
}

func TestFloat64bits(t *testing.T) {
	a := 123.4
	i := Float64bits(a)
	t.Log(i)
}

func TestUnsafe(t *testing.T) {
	i := 5
	uintPtr := uintptr(unsafe.Pointer(&i))
	t.Logf("%d, %d", uintPtr, &uintPtr)
}

type data2 struct {
	value  uint32
	label  [10]byte
	active bool
}

func dataFromBytes(b [16]byte) data2 {
	d := data2{}
	d.value = binary.BigEndian.Uint32(b[:4])
	copy(d.label[:], b[4:14])
	d.active = b[14] != 0
	return d
}

func TestDataFromByte(t *testing.T) {
	b := [16]byte{0, 132, 98, 237, 105, 80, 104, 111, 110, 101, 32, 49, 53, 80, 1, 0}
	d := dataFromBytes(b)
	t.Logf("value: %v, label: %v, active: %v", d.value, string(d.label[:]), d.active)
}

func isLE() bool {
	var x uint16 = 0xFF00
	xb := *(*[2]byte)(unsafe.Pointer(&x))
	return xb[0] == 0x00
}

func TestIsLE(t *testing.T) {
	t.Log(isLE())
}

func dataFromByteUnsafe(b [16]byte) data2 {
	d := *(*data2)(unsafe.Pointer(&b))
	if isLE() {
		d.value = bits.ReverseBytes32(d.value)
	}
	return d
}

func TestDataFromByteUnsafe(t *testing.T) {
	b := [16]byte{0, 132, 98, 237, 80, 104, 111, 110, 101, 0, 0, 0, 0, 0, 1, 0}
	d := dataFromByteUnsafe(b)
	t.Logf("value: %v, label: %v, active: %v", d.value, d.label, d.active)
}

func bytesFromData(d data2) [16]byte {
	out := [16]byte{}
	binary.BigEndian.PutUint32(out[:4], d.value)
	copy(out[4:14], d.label[:])
	if d.active {
		out[14] = 1
	}
	return out
}

func TestByteFromData(t *testing.T) {
	d := data2{
		value:  8676077,
		label:  [10]byte{80, 104, 111, 110, 101, 0, 0, 0, 0, 0},
		active: true,
	}
	b := bytesFromData(d)
	expected := [16]byte{0, 132, 98, 237, 80, 104, 111, 110, 101, 0, 0, 0, 0, 0, 1, 0}
	if b != expected {
		t.Fatalf("byte from data: %v, expected %v", b, expected)
	}
	t.Log("byte from data: ", b)
}

func bytesFromDataUnsafe(d data2) [16]byte {
	if isLE() {
		d.value = bits.ReverseBytes32(d.value)
	}
	b := *(*[16]byte)(unsafe.Pointer(&d))
	return b
}

func TestByteFromDataUnsafe(t *testing.T) {
	d := data2{
		value:  8676077,
		label:  [10]byte{80, 104, 111, 110, 101, 0, 0, 0, 0, 0},
		active: true,
	}
	b := bytesFromDataUnsafe(d)
	expected := [16]byte{0, 132, 98, 237, 80, 104, 111, 110, 101, 0, 0, 0, 0, 0, 1, 0}
	if b != expected {
		t.Fatalf("byte from data: %v, expected %v", b, expected)
	}
	t.Log("byte from data: ", b)
}

func TestUnsafeString(t *testing.T) {
	s := "oyen"
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	t.Log(sHdr.Len)

	for i := 0; i < sHdr.Len; i++ {
		bp := *(*byte)(unsafe.Pointer(sHdr.Data + uintptr(i)))
		t.Log(string(bp))
	}

	t.Log()
	runtime.KeepAlive(s)
}

func TestUnsafeInts(t *testing.T) {
	s := []int{10, 20, 30}
	sHdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	t.Log(sHdr.Len)
	t.Log(sHdr.Cap)

	intByteSize := unsafe.Sizeof(s[0])
	t.Log(intByteSize)

	for i := 0; i < sHdr.Len; i++ {
		intVal := *(*int)(unsafe.Pointer(sHdr.Data + intByteSize*uintptr(i)))
		t.Log(intVal)
	}

	runtime.KeepAlive(s)
}
