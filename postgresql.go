package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// структура таблицы UserSegments
type UserSegments struct {
	User_id int    `json:"user_id"`
	Slug    string `json:"slug"`
}

// структура ответа Json
type GetRequestSegment struct {
	Type     string   `json:"type"`
	UserID   int      `json:"user_id"`
	Segments []string `json:"segments"`
	Message  string   `json:"message"`
}

// Функция для открытия соединения с БД
func OpenConn(w http.ResponseWriter) (*sql.DB, error) {
	db, errDB := sql.Open("postgres", "user=postgres password=1 dbname=avito-tech sslmode=disable")
	if errDB != nil {
		//log.Fatal(err)
		fmt.Println("Ошибка открытия БД")
		CreateResponseJSON(w, GetRequestSegment{Type: "error", Message: "Ошибка открытия БД"})
	}
	return db, errDB
}

// Функция для записи ответа JSON
func CreateResponseJSON(w http.ResponseWriter, response GetRequestSegment) {
	json.NewEncoder(w).Encode(response)
}

// добавить сегмент
func CreateSegment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //определяем в заголовке с каким типом контента будем работать
	db, errDB := OpenConn(w)

	if errDB == nil {
		var exists bool                   //переменная для проверки существования сегмента в таблице
		slug := r.URL.Query().Get("slug") //получаем значения из URL по ключам

		w.Write([]byte("this is CreateSegment \n"))

		check, err := db.Query("SELECT EXISTS (SELECT * FROM segments WHERE slug=($1))", slug)
		//SELECT EXISTS (SELECT * FROM segments WHERE slug = 'slugname');
		//запрос чтобы проверить есть ли сегмент с таким именем в таблице
		if err != nil {
			//log.Fatal(err)
			fmt.Println("Ошибка: не удалось выполнить запрос EXISTS")
			CreateResponseJSON(w, GetRequestSegment{Type: "error", Message: "Не удалось выполнить запрос к БД на существующие сегменты"})
		}

		for check.Next() {
			err = check.Scan(&exists) //проверяем ответ
			if err != nil {
				fmt.Printf("Ошибка: не удалось считать данные из БД о существующих сегментах")
				CreateResponseJSON(w, GetRequestSegment{Type: "error", Message: "Не удалось считать данные из БД о существующих сегментах"})
			}
		}

		if !exists {
			_, err := db.Exec("INSERT INTO segments (slug) VALUES ($1)", slug)
			//пример: insert into segments (slug) values ('slugname')
			if err != nil {
				//panic(err)
				fmt.Println("Ошибка: не удалось добавить сегмент", slug)
				CreateResponseJSON(w, GetRequestSegment{Type: "error", Message: "Добавление сегмента " + slug + ":error"})
			} else {
				CreateResponseJSON(w, GetRequestSegment{Type: "success", Message: "Добавление сегмента " + slug + ":success"})
			}
		} else {
			fmt.Printf("Сегмент %s существует в таблице \n", slug)
			CreateResponseJSON(w, GetRequestSegment{Type: "error", Message: "Сегмент " + slug + " существует в таблице"})
		}
	}
}

// удалить сегмент
func DeleteSegment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, errDB := OpenConn(w)

	if errDB == nil {
		slug := r.URL.Query().Get("slug")

		w.Write([]byte("this is DeleteSegment\n"))

		_, err := db.Exec("DELETE FROM segments WHERE slug=($1)", slug)
		//DELETE FROM segments WHERE slug='add1'
		if err != nil {
			//panic(err)
			fmt.Println("Ошибка: не удалось выполнить удаление сегмента", slug)
			CreateResponseJSON(w, GetRequestSegment{Type: "error", Message: "Удаление сегмента " + slug + ":error"})
		} else {
			fmt.Println("Ошибка: удалось выполнить удаление сегмента", slug)
			CreateResponseJSON(w, GetRequestSegment{Type: "success", Message: "Удаление сегмента " + slug + ":success"})
		}
	}
}

// добавить/убрать пользователя из сегментов
func ChangeUserSegments(w http.ResponseWriter, r *http.Request) {
	db, errDB := OpenConn(w)

	if errDB == nil {
		str_add_seg := r.URL.Query().Get("add_seg")
		str_del_seg := r.URL.Query().Get("del_seg")
		user_id := r.URL.Query().Get("user_id")

		str_user_id, err := strconv.Atoi(user_id)

		if err != nil {
			//panic(err)
			fmt.Printf("Ошибка: некорректный id пользователя \"%s\" \n", user_id)
			CreateResponseJSON(w, GetRequestSegment{Type: "error", UserID: str_user_id, Message: "Некорректный id=" + user_id})
		} else {
			//формируем строковый массив сегментов, которые нужно добавить,
			// и сегментов, которые нужно удалить
			add_seg := strings.FieldsFunc(str_add_seg, Split)
			del_seg := strings.FieldsFunc(str_del_seg, Split)

			var messageArray []string

			w.Write([]byte("this is ChangeUserSegments\n"))

			//fmt.Println(add_seg, del_seg, user_id)

			//формируем запросы на добавление пользователя в сегменты
			for i := range add_seg {
				_, err := db.Exec("INSERT INTO usersegments VALUES ($1,$2)", user_id, add_seg[i])
				if err != nil {
					//panic(err)
					fmt.Printf("Ошибка: не удалось добавить пользователя %s в сегмент %s \n", user_id, add_seg[i])
					messageArray = append(messageArray, "Добавление в пользователя в сегмент "+add_seg[i]+": error")
				} else {
					messageArray = append(messageArray, "Добавление в пользователя в сегмент "+add_seg[i]+": success")
				}
			}

			//формируем запросы на удаление пользователя из сегментов
			for i := range del_seg {
				_, err := db.Exec("DELETE FROM usersegments WHERE (user_id=$1 and slug=$2)", user_id, del_seg[i])
				if err != nil {
					//panic(err)
					fmt.Printf("Ошибка: не удалось удалить пользователя %s из сегмента %s \n", user_id, del_seg[i])
					messageArray = append(messageArray, "Удаление пользователя из сегмента "+del_seg[i]+": error")
				} else {
					messageArray = append(messageArray, "Удаление пользователя из сегмента "+del_seg[i]+": success")
				}
			}
			CreateResponseJSON(w, GetRequestSegment{Type: "success", UserID: str_user_id, Message: strings.Join(messageArray, ", ")})
		}
	}
}

// получить сегменты, в которых состоит пользователь
func GetUserSegments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := OpenConn(w)

	if err == nil {
		var listOfSeg []string

		user_id := r.URL.Query().Get("user_id")
		str_user_id, err := strconv.Atoi(user_id)

		if err != nil {
			//panic(err)
			fmt.Printf("Ошибка: некорректный id пользователя \"%s\" \n", user_id)
			CreateResponseJSON(w, GetRequestSegment{Type: "error", UserID: str_user_id, Message: "Некорректный id=" + user_id})
		} else {
			//выполняем запрос если нет ошибок; может стоило сделать через МЕТКИ этот участок кода для удобства чтения?
			rows, err := db.Query("SELECT * FROM usersegments WHERE user_id=($1)", user_id)
			if err != nil {
				//panic(err)
				fmt.Printf("Ошибка: не удалось получить список сегментов пользователя \"%s\" из БД", user_id)
				CreateResponseJSON(w, GetRequestSegment{Type: "error", UserID: str_user_id, Message: "Не удалось выполнить запрос к БД"})
			} else {

				for rows.Next() {
					var user_id int
					var slug string

					err = rows.Scan(&user_id, &slug) //считываем значения строк
					if err != nil {
						fmt.Printf("Ошибка: не удалось считать сегменты пользователя %v", user_id)
						CreateResponseJSON(w, GetRequestSegment{Type: "error", UserID: str_user_id, Message: "Не удалось считать сегменты из БД"})
					}
					listOfSeg = append(listOfSeg, slug)
				}
				//после успешного считывания строк можно записать, что всё OK
				CreateResponseJSON(w, GetRequestSegment{Type: "success", UserID: str_user_id, Segments: listOfSeg})
			}
		}
	}
}

// функция для раздления строки
func Split(r rune) bool {
	return r == ',' //|| r == '.'
}
