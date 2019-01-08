package main

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

// тут писать SearchServer

type UserXml struct {
	Id         int    `xml:"id"`
	FirstName  string `xml:"first_name"`
	SecondName string `xml:"last_name"`
	Age        int    `xml:"age"`
	About      string `xml:"about"`
	Gender     string `xml:"gender"`
}

type Xmlcontent struct {
	Version string    `xml:"version,attr"`
	List    []UserXml `xml:"row"`
}

//тут пошли функции и типы для сортировки структуры по разным параметрам
type byName []User

func More(i, j int) bool {

	return i > j

}

func Less(i, j int) bool {

	return i < j

}

func (a byName) Len() int      { return len(a) }
func (a byName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type byAge []User

func (a byAge) Len() int      { return len(a) }
func (a byAge) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type byId []User

func (a byId) Len() int      { return len(a) }
func (a byId) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func SearchServer(w http.ResponseWriter, r *http.Request) {

	AccesToken := r.Header.Get("AccessToken")

	if AccesToken != fmt.Sprintf("%x", md5.Sum([]byte("AccessToken"))) {
		w.WriteHeader(http.StatusUnauthorized)
		return //fmt.Errorf("Error while pars param Limit")
	}

	//устанавливаем значения из ЮРЛ
	Query := r.URL.Query().Get("query")            //поиск по полям Name,About; Если query пустой, то делаем только сортировку, т.е. возвращаем все записи
	OrderField := r.URL.Query().Get("order_field") //работает по полям `Id`, `Age`, `Name`: возвращает, только то, что мы ему сказали

	Limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return //fmt.Errorf("Error while pars param Limit")
	}
	Offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return //fmt.Errorf("Error while pars param Offset")
	}

	OrderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return //fmt.Errorf("Error while pars param order_by")
	}
	//начало работы с файлом XML
	file, err := os.Open("./dataset.xml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // fmt.Errorf("Error file not found")
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return //fmt.Errorf("Error while reading file")
	}

	var XmlUserData Xmlcontent
	err = xml.Unmarshal(fileContents, &XmlUserData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // fmt.Errorf("Error while Xml Unmarshal")
	}
	Users := make([]User, 0)
	if Query == "" {
		found_count := 0
		for idx, SingXmlUser := range XmlUserData.List {

			if found_count == Limit {
				break
			}
			if idx >= Offset {
				Users = append(Users, User{})
				Users[found_count].Id = SingXmlUser.Id
				Users[found_count].Name = SingXmlUser.FirstName + SingXmlUser.SecondName
				Users[found_count].Age = SingXmlUser.Age
				Users[found_count].About = SingXmlUser.About
				Users[found_count].Gender = SingXmlUser.Gender
				found_count++
				switch OrderField {
				case "Id":
					switch OrderBy {
					case -1:
						sort.Slice(byId(Users), More)
					case 0:

					case 1:
						sort.Slice(byId(Users), Less)
					}
				case "Age":
					switch OrderBy {
					case -1:
						sort.Slice(byAge(Users), More)
					case 0:

					case 1:
						sort.Slice(byAge(Users), Less)
					}
				case "Name":
					switch OrderBy {
					case -1:
						sort.Slice(byName(Users), More)
					case 0:

					case 1:
						sort.Slice(byName(Users), Less)
					}
				case "":
					switch OrderBy {
					case -1:
						sort.Slice(byName(Users), More)
					case 0:

					case 1:
						sort.Slice(byName(Users), Less)
					}
				default:
					w.WriteHeader(http.StatusBadRequest)
					return // fmt.Errorf("Error while get OrderField")
				}
			} else {
				break
			}
		}
	} else {

		found_count := 0
		for idx, SingXmlUser := range XmlUserData.List {

			if found_count == Limit {
				break
			}
			//здесь сделать поиск
			ok, err := regexp.MatchString(Query, SingXmlUser.FirstName+SingXmlUser.SecondName)

			ok1, err1 := regexp.MatchString(Query, SingXmlUser.About)

			if (idx >= Offset && (ok && err == nil)) || (idx >= Offset && (ok1 && err1 == nil)) {
				Users = append(Users, User{})
				Users[found_count].Id = SingXmlUser.Id
				Users[found_count].Name = SingXmlUser.FirstName + SingXmlUser.SecondName
				Users[found_count].Age = SingXmlUser.Age
				Users[found_count].About = SingXmlUser.About
				Users[found_count].Gender = SingXmlUser.Gender
				found_count++

				switch OrderField {
				case "Id":
					switch OrderBy {
					case -1:
						sort.Slice(byId(Users), More)
					case 0:

					case 1:
						sort.Slice(byId(Users), Less)
					}
				case "Age":
					switch OrderBy {
					case -1:
						sort.Slice(byAge(Users), More)
					case 0:

					case 1:
						sort.Slice(byAge(Users), Less)
					}
				case "Name":
					switch OrderBy {
					case -1:
						sort.Slice(byName(Users), More)
					case 0:

					case 1:
						sort.Slice(byName(Users), Less)
					}
				case "":
					switch OrderBy {
					case -1:
						sort.Slice(byName(Users), More)
					case 0:

					case 1:
						sort.Slice(byName(Users), Less)
					}
				default:
					w.WriteHeader(http.StatusBadRequest)
					return // fmt.Errorf("Error while get OrderField")
				}
			}
		}

	}
	w.WriteHeader(http.StatusOK)
	result, _ := json.Marshal(Users)
	_, _ = w.Write(result)
	r.Body.Close()
}

func WaitTimeoutHan(w http.ResponseWriter, r *http.Request) {

	time.Sleep(10 * time.Second)

	return
}

func FakeReqToClient(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusBadRequest)

	return
}
