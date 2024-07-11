package learning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		data := "{\"id\":\"5db6e4ce-995a-4eee-b1f8-36d78af58504\",\"name\":\"Admin Key\",\"description\":\"\"},{\"id\":\"bd5fbf66-f646-4e56-9553-609d6293e5db\",\"name\":\"Admin TMN\",\"description\":\"\"},{\"id\":\"32a0b0b8-464d-4858-8971-ecad950e77a7\",\"name\":\"Admin HCIS\",\"description\":\"\"},{\"id\":\"edbbfe50-4ede-447c-a1be-cd2a90d7912a\",\"name\":\"Interviewer\",\"description\":\"\"},{\"id\":\"b056d30d-c096-4755-b6ac-92cb82cdb18b\",\"name\":\"Kandidat\",\"description\":\"\"},{\"id\":\"cbfa5412-7fe4-4fc5-aea7-906eeffa15b9\",\"name\":\"Dewa\",\"description\":\"Dewa Cinta\"},{\"id\":\"2038eff8-e092-4e95-a82d-8bf5216bf6c4\",\"name\":\"Perdana Menteri\",\"description\":\"\"}"
		buf := bytes.NewBufferString(data)
		var obj []map[string]interface{}
		err := json.Unmarshal(buf.Bytes(), &obj)
		if err != nil {
			panic(err)
		}
		fmt.Println(obj)

		s, err := strconv.Unquote(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(s)

		resp := map[string]any{
			"data": data,
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			panic(err)
		}
	})

	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic(err)
	}
}

type vehicle interface {
	getPercentFuelRemaining() float64
}

type car struct {
	fuel float64
}

func (c car) getPercentFuelRemaining() float64 {
	return c.fuel * 0.01
}

type motorcycle struct {
	fuel float64
}

func (m motorcycle) getPercentFuelRemaining() float64 {
	return m.fuel * 0.01
}

func getFuel(v vehicle) float64 {
	return v.getPercentFuelRemaining()
}

func TestVehicle(t *testing.T) {
	c := car{fuel: 30}
	m := motorcycle{fuel: 10}

	t.Log(getFuel(c))
	t.Log(getFuel(m))
}

func TestContinue(t *testing.T) {
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			t.Log("divided by 2")
			continue
		}
		t.Log(i)
	}
}
