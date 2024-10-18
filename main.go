// RegexRedirectLinkForWidows project main.go
// author email:iwlb@outlook.om exgaya@gmail.com h4g@163.com
// author:Exgaya Epowsal Wlb
// country:China
// Sex: male
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func Atoi(a string) int {
	a = strings.Trim(a, " \r\n\t")
	a2 := []byte{}
	for i := 0; i < len(a); i++ {
		if a[i] >= '0' && a[i] <= '9' {
			a2 = append(a2, a[i])
		} else {
			break
		}
	}
	i, er := strconv.ParseInt(string(a2), 10, 64)
	if er != nil {
		return 0
	}
	return int(i)
}

func SplitRegexReplaceWith(repw string) []string {
	ls := []string{}
	tm := []byte{}
	for ri := 0; ri < len(repw); ri += 1 {
		if ri+3 < len(repw) && repw[ri] == '\\' && repw[ri+1] == 'x' && (repw[ri+2] >= '0' && repw[ri+2] <= '9' || repw[ri+2] >= 'a' && repw[ri+2] <= 'f' || repw[ri+2] >= 'A' && repw[ri+2] <= 'F') && (repw[ri+3] >= '0' && repw[ri+3] <= '9' || repw[ri+3] >= 'a' && repw[ri+3] <= 'f' || repw[ri+3] >= 'A' && repw[ri+3] <= 'F') {
			v64, v64e := strconv.ParseUint(repw[ri+2:ri+2+2], 16, 8)
			if v64e == nil {
				tm = append(tm, byte(v64))
			}
			ri += 2
		} else if ri+1 < len(repw) && repw[ri] == '\\' && (repw[ri+1] >= '0' && repw[ri+1] <= '9' || repw[ri+1] == '\\' || repw[ri+1] == '#' || repw[ri+1] == '%' || repw[ri+1] == 'n' || repw[ri+1] == 'r' || repw[ri+1] == 't') {
			if repw[ri+1] >= '0' && repw[ri+1] <= '9' {
				if len(tm) > 0 {
					ls = append(ls, string(tm))
					tm = tm[:0]
				}
				ls = append(ls, repw[ri+1:ri+2])
				ri += 1
			} else if repw[ri+1] == 'n' {
				if len(tm) > 0 {
					ls = append(ls, string(tm))
					tm = tm[:0]
				}
				ls = append(ls, "\n")
				ri += 1
			} else if repw[ri+1] == 'r' {
				if len(tm) > 0 {
					ls = append(ls, string(tm))
					tm = tm[:0]
				}
				ls = append(ls, "\r")
				ri += 1
			} else if repw[ri+1] == 't' {
				if len(tm) > 0 {
					ls = append(ls, string(tm))
					tm = tm[:0]
				}
				ls = append(ls, "\t")
				ri += 1
			} else if repw[ri+1] == '\\' {
				if len(tm) > 0 {
					ls = append(ls, string(tm))
					tm = tm[:0]
				}
				ls = append(ls, "\\")
				ri += 1
			} else if repw[ri+1] == '%' {
				if len(tm) > 0 {
					ls = append(ls, string(tm))
					tm = tm[:0]
				}
				ls = append(ls, "%")
				ri += 1
			} else if repw[ri+1] == '#' {
				if len(tm) > 0 {
					ls = append(ls, string(tm))
					tm = tm[:0]
				}
				ls = append(ls, "#")
				ri += 1
			} else {
				tm = append(tm, '\\')
			}

		} else if ri+2 < len(repw) && repw[ri] == '#' && repw[ri+1] >= '0' && repw[ri+1] <= '9' && repw[ri+2] >= '0' && repw[ri+2] <= '9' {
			if len(tm) > 0 {
				ls = append(ls, string(tm))
				tm = tm[:0]
			}
			matag := repw[ri : ri+3]
			ls = append(ls, matag)
			ri += 2
		} else if ri+1 < len(repw) && repw[ri] == '#' && repw[ri+1] >= '0' {
			if len(tm) > 0 {
				ls = append(ls, string(tm))
				tm = tm[:0]
			}
			matag := repw[ri : ri+2]
			ls = append(ls, matag)
			ri += 1
		} else {
			tm = append(tm, repw[ri])
		}
	}
	if len(tm) > 0 {
		ls = append(ls, string(tm))
		tm = tm[:0]
	}
	return ls
}

func RegexReplace(source string, gma [][]int, replacewithstr string) string {
	if len(gma) != 1 {
		return source
	}
	var newpath []byte
	rls := SplitRegexReplaceWith(replacewithstr)
	for _, rl := range rls {
		if len(rl) == 3 && rl[0] == '#' && rl[1] >= '0' && rl[1] <= '9' && rl[2] >= '0' && rl[2] <= '9' {
			maind := Atoi(rl[1:])
			if 2*maind >= len(gma[0]) {
				newpath = append(newpath, []byte(rl)...)
			} else {
				mastr := source[gma[0][2*maind]:gma[0][2*maind+1]]
				newpath = append(newpath, []byte(mastr)...)
			}
		} else if len(rl) == 2 && rl[0] == '#' && rl[1] >= '0' && rl[1] <= '9' {
			maind := Atoi(rl[1:])
			if 2*maind >= len(gma[0]) {
				newpath = append(newpath, []byte(rl)...)
			} else {
				mastr := source[gma[0][2*maind]:gma[0][2*maind+1]]
				newpath = append(newpath, []byte(mastr)...)
			}
		} else {
			newpath = append(newpath, []byte(rl)...)
		}
	}
	if len(newpath) == 0 {
		return source
	} else {
		return string(newpath)
	}
}

var p0, p1, fr0, rep0, fr1, rep1 *string

func redirect0(w http.ResponseWriter, req *http.Request) {
	oldlk := "http://" + req.Host + "" + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		oldlk += "?" + req.URL.RawQuery
	}
	mas := regexp.MustCompile(*fr0).FindAllStringIndex(oldlk, -1)
	newlk := RegexReplace(oldlk, mas, *rep0)
	log.Println("redirect", oldlk, newlk)
	http.Redirect(w, req, newlk, http.StatusTemporaryRedirect)
}
func httpServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirect0)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Http sever error: ", err)
		log.Println("Http sever error: ", err)
	}
}

func redirect1(w http.ResponseWriter, req *http.Request) {
	oldlk := "https://" + req.Host + "" + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		oldlk += "?" + req.URL.RawQuery
	}
	mas := regexp.MustCompile(*fr1).FindAllStringIndex(oldlk, -1)
	newlk := RegexReplace(oldlk, mas, *rep1)
	log.Println("redirect", oldlk, newlk)
	http.Redirect(w, req, newlk, http.StatusTemporaryRedirect)
}
func httpsServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirect1)
	err := http.ListenAndServeTLS(":"+port, `generate_certkey/server.crt`, `generate_certkey/server.key`, nil)
	if err != nil {
		fmt.Println("Https sever error: ", err)
		log.Println("Https sever error: ", err)
	}
}

func main() {
	fmt.Println(`RedirectUrlToLocalUrl 0.1 Help:
-p0= http server port
-p1= https server port
-fr0= find text regex
-rep0= replace text:#1 for match 1...
-fr1= find text regex
-rep1= replace text:#1 for match 1...

Windows modify C:\Windows\system32\driver\et\hosts can redirect any url to local server to modify return data or redirect to other link for program right running or something.
`)
	p0 = flag.String("p0", "80", "http server port")
	p1 = flag.String("p1", "443", "https server port")
	fr0 = flag.String("fr0", ".*", "http replace find text")
	rep0 = flag.String("rep0", "#0", "http replace text")
	fr1 = flag.String("fr1", ".*", "https replace find text")
	rep1 = flag.String("rep1", "#0", "https replace text")

	http.Handle("/.well-known/", http.StripPrefix("/.well-known/", http.FileServer(http.Dir(".well-known/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data/"))))

	go httpServer(*p0)
	go httpsServer(*p1)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
