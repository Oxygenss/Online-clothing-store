package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var id_user int = 0

type Article struct {
	Id                int
	Name, Email, Text string
}

type Item struct {
	Id, Price, id_user, Count                                         int
	Name, Type, Material, Color, Season, Lift, Country, Img, Slide string
}

type User struct {
	Id                    int
	Name, Email, Password string
}

func Main(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/Main/main.html")
	tmpl.Execute(w, nil)
}

func AboutUs_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/AboutUs/AboutUs.html")
	tmpl.Execute(w, nil)
}

func Contacts_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/Contacts/Contacts.html")
	tmpl.Execute(w, nil)
}

func Catalog(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Catalog/Catalog.html")
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	database, err := sql.Open("sqlite3", "./data")
	if err != nil {
		panic(err)
	}
	defer database.Close()
	var items = []Item{}
	res, err := database.Query("SELECT id, Name, Type, Material, Color, Price, Img FROM Items")
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var item Item
		err = res.Scan(&item.Id, &item.Name, &item.Type, &item.Material, &item.Color, &item.Price, &item.Img)
		if err != nil {
			panic(err)
		}
		items = append(items, item)
	}
	defer res.Close()
	tmpl.ExecuteTemplate(w, "Catalog", items)
}

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
   
	tmpl, err := template.ParseFiles("templates/Catalog/Catalog.html")
	if err != nil {
	 fmt.Fprint(w, err.Error())
	 return
	}
   
	database, err := sql.Open("sqlite3", "./data")
	if err != nil {
	 fmt.Fprint(w, err.Error())
	 return
	}
	defer database.Close()
   
	var items []Item
   
	if query != "" {
	 // Выполните поиск элементов, где любое из полей Name, Type, Material или Color совпадает с поисковым запросом
	 res, err := database.Query("SELECT id, Name, Type, Material, Color, Price, Img FROM Items WHERE Name LIKE ? OR Type LIKE ? OR Material LIKE ? OR Color LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	 if err != nil {
	  fmt.Fprint(w, err.Error())
	  return
	 }
   
	 for res.Next() {
	  var item Item
	  err = res.Scan(&item.Id, &item.Name, &item.Type, &item.Material, &item.Color, &item.Price, &item.Img)
	  if err != nil {
	   fmt.Fprint(w, err.Error())
	   return
	  }
	  items = append(items, item)
	 }
	 res.Close()
	} else {
	 // Если поисковый запрос не указан, получаем все элементы
	 res, err := database.Query("SELECT id, Name, Type, Material, Color, Price, Img FROM Items")
	 if err != nil {
	  fmt.Fprint(w, err.Error())
	  return
	 }
   
	 for res.Next() {
	  var item Item
	  err = res.Scan(&item.Id, &item.Name, &item.Type, &item.Material, &item.Color, &item.Price, &item.Img)
	  if err != nil {
	   fmt.Fprint(w, err.Error())
	   return
	  }
	  items = append(items, item)
	 }
	 res.Close()
	}
   
	tmpl.ExecuteTemplate(w, "Catalog", items)
}

func Feedback(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Feedback/Feedback.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "Feedback", nil)
}

func AllFeed(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Feedback/AllFeed.html")

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	database, err := sql.Open("sqlite3", "./data")
	if err != nil {
		panic(err)
	}
	defer database.Close()

	var posts = []Article{}
	res, err := database.Query("SELECT id, Name, Email, Text FROM Articles")
	if err != nil {
		panic(err)
	}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Name, &post.Email, &post.Text)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	defer res.Close()

	tmpl.ExecuteTemplate(w, "AllFeed", posts)
}

func Login_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/Login/Login.html")
	tmpl.Execute(w, nil)
}

func Registration_page(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/Login/Registration.html")
	tmpl.Execute(w, nil)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Profile/Profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if id_user != 0 {
		database, err := sql.Open("sqlite3", "./data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer database.Close()

		row := database.QueryRow("SELECT Name, Email, Password FROM Users WHERE id = ?", id_user)

		var user User
		err = row.Scan(&user.Name, &user.Email, &user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "Profile", user)
		//http.Redirect(w, r, "/Profile", http.StatusSeeOther)
	} else {
		tmpl.ExecuteTemplate(w, "Profile", false)
	}

}

func Save_user(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	repeat := r.FormValue("repeat")
	fmt.Println(name, email, password, repeat)

	if name == "" || email == "" || password == "" || repeat == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		if password != repeat {
			http.Error(w, "Пароли не совпадают", http.StatusBadRequest)
			return
		}
		database, err := sql.Open("sqlite3", "./data")
		if err != nil {
			panic(err)
		}

		defer database.Close()

		insert, err := database.Prepare("INSERT INTO users (Name, Email, Password) VALUES (?, ?, ?)")
		insert.Exec(name, email, password)

		defer insert.Close()
		fmt.Fprintf(w, "Вы успешно зарегистрироваться")
	}

}

func ItemPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Item/Item.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := mux.Vars(r)["id"]
	//fmt.Println(id)
	database, err := sql.Open("sqlite3", "./data")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var item Item
	err = database.QueryRow("SELECT id, Name, Type, Material, Color, Price, Season, Lift, Country, Slide FROM Items WHERE id = ?", id).Scan(
		&item.Id, &item.Name, &item.Type, &item.Material, &item.Color, &item.Price, &item.Season, &item.Lift, &item.Country, &item.Slide)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "ItemPage", item)
}

func Articles(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.ParseFiles("templates/Feedback/Feedback.html")

	// if err != nil {
	//    fmt.Fprint(w, err.Error())
	// }

	if id_user == 0 {
		fmt.Fprintf(w, "Нужно авторизоваться, чтобы отправить отзыв.")
	} else {
		//name := r.FormValue("name")
		name := r.FormValue("name")
		email := r.FormValue("email")
		text := r.FormValue("text")
		//fmt.Println(name, text, email)

		if email == "" || text == "" || name == "" {
			fmt.Fprintf(w, "Не все данные заполнены")
		} else {
			database, err := sql.Open("sqlite3", "./data")
			if err != nil {
				panic(err)
			}
			defer database.Close()

			insert, err := database.Prepare("INSERT INTO articles (Name, Email, Text) VALUES (?, ?, ?)")
			if err != nil {
				panic(err)
			}
			insert.Exec(name, email, text)

			defer insert.Close()
		}
	}
	http.Redirect(w, r, "/Feedback", http.StatusSeeOther)
}

func Authorization(w http.ResponseWriter, r *http.Request) {
	var _user User
	_user.Name = r.FormValue("name")
	_user.Email = r.FormValue("email")
	_user.Password = r.FormValue("password")

	if _user.Name == "" || _user.Email == "" || _user.Password == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		database, err := sql.Open("sqlite3", "./data")
		if err != nil {
			panic(err)
		}
		defer database.Close()

		var users = []User{}

		res, err := database.Query("SELECT id, Name, Email, Password FROM Users")
		if err != nil {
			panic(err)
		}
		for res.Next() {
			var user User
			err = res.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
			if err != nil {
				panic(err)
			}
			users = append(users, user)
		}
		defer res.Close()

		for _, user := range users {
			if user.Email == _user.Email && user.Name == _user.Name && user.Password == _user.Password {
				id_user = user.Id
				fmt.Println(_user.Email, _user.Name, _user.Password, id_user)
				http.Redirect(w, r, "/Profile", http.StatusSeeOther)
			}
		}

		if id_user == 0 {
			fmt.Fprintf(w, "Данные введены неправильно")
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	id_user = 0
	http.Redirect(w, r, "/Main", http.StatusSeeOther)
}

func AddInCart(w http.ResponseWriter, r *http.Request) {

	if id_user == 0 {
		fmt.Fprintf(w, "Нужно авторизоваться, чтобы добавить товар в корзину.")
	} else {
		id := mux.Vars(r)["id"]

		database, err := sql.Open("sqlite3", "./data")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer database.Close()

		var item Item
		err = database.QueryRow("SELECT id, Name, Type, Color, Material, Price, Img FROM Items WHERE id = ?", id).Scan(
			&item.Id, &item.Name, &item.Type, &item.Color, &item.Material, &item.Price, &item.Img)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		insert, err := database.Prepare("INSERT INTO Cart (id, Name, Type, Color, Material, Price, Img, id_user, Count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err)
		}

		// Проверяем наличие товара в корзине
		var count int
		err = database.QueryRow("SELECT COUNT(*) FROM Cart WHERE id = ? AND id_user = ?", item.Id, id_user).Scan(&count)
		if err != nil {
			panic(err)
		}

		if count > 0 {
			fmt.Fprintf(w, "Товар уже находится в корзине.")
		} else {
			insert.Exec(item.Id, item.Name, item.Type, item.Color, item.Material, item.Price, item.Img, id_user, 1)
			fmt.Fprintf(w, "Товар успешно добавлен в корзину.")
		}
		defer insert.Close()
	}
}

func Cart(w http.ResponseWriter, r *http.Request) {
    if id_user == 0 {
        http.Redirect(w, r, "/Login", http.StatusSeeOther)
    } else {
        tmpl, err := template.ParseFiles("templates/Cart/Cart.html")
        if err != nil {
            fmt.Fprint(w, err.Error())
        }
    
        database, err := sql.Open("sqlite3", "./data")
        if err != nil {
            panic(err)
        }
        defer database.Close()
        var items = []Item{}
        res, err := database.Query("SELECT id, Name, Type, Color, Material, Price, Img, id_user, Count FROM Cart")
        if err != nil {
            panic(err)
        }
    
        var sum int = 0
    
        for res.Next() {
            var item Item
            err = res.Scan(&item.Id, &item.Name, &item.Type, &item.Color, &item.Material, &item.Price, &item.Img, &item.id_user, &item.Count)
            if err != nil {
                panic(err)
            }
    
            if id_user == item.id_user {
                sum += item.Price * item.Count
                items = append(items, item)
            }
    
        }
        defer res.Close()

        data := struct {
            Items []Item
            Sum   int
        }{
            Items: items,
            Sum:   sum,
        }

        tmpl.ExecuteTemplate(w, "Cart", data)
    }
}

func AddItem(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    fmt.Println(id)
    database, err := sql.Open("sqlite3", "./data")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer database.Close()

    var item Item
    err = database.QueryRow("SELECT Count FROM Cart WHERE id = ?", id).Scan(
        &item.Count)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println(item.Count)

    update, err := database.Prepare("UPDATE Cart SET Count = ? WHERE id = ?")
    if err != nil {
        panic(err)
    }
    update.Exec(item.Count+1, id)
    defer update.Close()
	http.Redirect(w, r, "/Cart", http.StatusSeeOther)
}

func DelItem(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    fmt.Println(id)
    database, err := sql.Open("sqlite3", "./data")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer database.Close()

    var item Item
    err = database.QueryRow("SELECT Count FROM Cart WHERE id = ?", id).Scan(
        &item.Count)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println(item.Count)

    if item.Count == 1 {
        DelInCart(w, r) // Передаем параметры w и r для удаления товара
    } else {
        update, err := database.Prepare("UPDATE Cart SET Count = ? WHERE id = ?")
        if err != nil {
            panic(err)
        }
        update.Exec(item.Count-1, id)
        defer update.Close()
    }

    http.Redirect(w, r, "/Cart", http.StatusSeeOther)
}

func DelInCart(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    fmt.Println(id)
    database, err := sql.Open("sqlite3", "./data")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer database.Close()

    delete, err := database.Prepare("DELETE FROM Cart WHERE id = ?")
    if err != nil {
        panic(err)
    }
    delete.Exec(id)
    defer delete.Close()

    http.Redirect(w, r, "/Cart", http.StatusSeeOther)
}



func handleRequest() {
	router := mux.NewRouter()

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	http.Handle("/", router)
	router.HandleFunc("/Main", Main)
	router.HandleFunc("/AboutUs", AboutUs_page)
	router.HandleFunc("/Contacts", Contacts_page)
	router.HandleFunc("/Catalog", Catalog)
	router.HandleFunc("/Feedback", Feedback)
	router.HandleFunc("/Login", Login_page)
	router.HandleFunc("/Registration", Registration_page)
	router.HandleFunc("/Profile", Profile)

	router.HandleFunc("/DelInCart/{id:[0-9]+}", DelInCart)
	router.HandleFunc("/Search", Search)
	router.HandleFunc("/AddItem/{id:[0-9]+}", AddItem)
	router.HandleFunc("/DelItem/{id:[0-9]+}", DelItem)
	router.HandleFunc("/Cart", Cart)
	router.HandleFunc("/AddInCart/{id:[0-9]+}", AddInCart)
	router.HandleFunc("/Logout", Logout)
	router.HandleFunc("/AllFeed", AllFeed)
	router.HandleFunc("/Item/{id:[0-9]+}", ItemPage)
	router.HandleFunc("/Articles", Articles)
	router.HandleFunc("/Authorization", Authorization)
	router.HandleFunc("/Save_user", Save_user)
	http.ListenAndServe(":8080", nil)
}

func main() {

	fmt.Println("SQL")

	handleRequest()
}
