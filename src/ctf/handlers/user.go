package handlers 

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"
	"time"
	"log"
	"net/http"
	"ctf/utils"
	"ctf/model"
	"errors"
)


func UserRegister(w http.ResponseWriter, r *http.Request) {
	nick := r.FormValue("nick")
	password := r.FormValue("password")

	db, err := model.GetDB(w)
	if err != nil{
		return
	}

	var user model.User
	notFound := db.Where(&model.User{Nick: nick}).First(&user).RecordNotFound()
	if !notFound {
		w.WriteHeader(http.StatusConflict)
		utils.SendResponseJSON(w, utils.Message{"A user with this nick already exists."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	//err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		utils.SendResponseJSON(w, utils.InternalErrorMessage)
		log.Printf("%v\n", err)
		return
	}
	user = model.User{
		Nick: nick,
		Password: string(hashedPassword),
		IsAdmin: false,
	}

	db.Create(&user)

	w.WriteHeader(http.StatusCreated)
	utils.SendResponseJSON(w, utils.Message{"User successfully created."})
}

func UserAuthenticate(w http.ResponseWriter, r *http.Request){
	nick := r.FormValue("nick")
	password := r.FormValue("password")

	db, err := model.GetDB(w)
	if err != nil{
		return
	}

	var user model.User
	notFound := db.Where(&model.User{Nick: nick}).First(&user).RecordNotFound()
	if notFound{
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendResponseJSON(w, utils.Message{"Can't login with those credentials."})
		log.Printf("%v\n", err)
		return 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil{
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendResponseJSON(w, utils.Message{"Can't login with those credentials."})
		log.Printf("%v\n", err)
		return
	}

	currentTime := time.Now()
	user.TimeAuthenticated = currentTime
	user.Token = utils.RandString(40)

	db.Save(&user)

	w.WriteHeader(http.StatusAccepted)
	utils.SendResponseJSON(w, user.Token)
}

func IsUserAuthenticated(w http.ResponseWriter, r *http.Request) (registeredUser bool, user model.User, err error){
	token := r.Header.Get("X-CTF-AUTH")
	registeredUser = false
	
	if token == ""{
		return
	}

	db, err := model.GetDB(w)
	if err != nil {return}

	notFound := db.Where(&model.User{Token: token}).First(&user).RecordNotFound()
	if notFound {
		return 
	}

	hoursElapsed := time.Now().Sub(user.TimeAuthenticated).Hours()
	if hoursElapsed > 48 { return registeredUser, user, errors.New("Token timed out.") }
	registeredUser = true

	return 
}

func UserShow(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userID := vars["userID"]
	
	db, err := model.GetDB(w)
	if err != nil {return}

	var user model.User
	notFound := db.Where("id = ?", userID).First(&user).RecordNotFound()
	if notFound{
		w.WriteHeader(http.StatusNotFound)
		utils.SendResponseJSON(w, utils.NotFoundErrorMessage)
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.SendResponseJSON(w, user)
}

func UserShowOwn(w http.ResponseWriter, r *http.Request){
	registeredUser, user, err := IsUserAuthenticated(w, r)
	if err != nil{return}
	if !registeredUser{
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendResponseJSON(w, utils.NotLoggedInMessage)
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.SendResponseJSON(w, user)
}

func UserShowAll(w http.ResponseWriter, r *http.Request){
	db, err := model.GetDB(w)
	if err != nil {return}
	var users model.Users
	db.Find(&users)

	w.WriteHeader(http.StatusOK)
	utils.SendResponseJSON(w, users)
}

func UserShowValidatedChallenges(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userID := vars["userID"]
	
	db, err := model.GetDB(w)
	if err != nil {return}

	var validatedChallenges model.ValidatedChallenges
	db.Where(&model.ValidatedChallenge{UserID: userID}).Find(&validatedChallenges)


	w.WriteHeader(http.StatusOK)
	utils.SendResponseJSON(w, validatedChallenges)
}

func UserChangePassword(w http.ResponseWriter, r *http.Request){
	registeredUser, user, err := IsUserAuthenticated(w, r)
	if err != nil{return}
	if !registeredUser{
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendResponseJSON(w, utils.NotLoggedInMessage)
		return
	}

	password := r.FormValue("password")

	db, err := model.GetDB(w)
	if err != nil{
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		utils.SendResponseJSON(w, utils.InternalErrorMessage)
		log.Printf("%v\n", err)
		return
	}
	user.Password = string(hashedPassword)
	user.Token = ""
	user.TimeAuthenticated, _ = time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")

	db.Update(&user)

	w.WriteHeader(http.StatusAccepted)
	utils.SendResponseJSON(w, utils.Message{"Password successfully changed. Please regenerate a token."})
}


func UserDelete(w http.ResponseWriter, r *http.Request){
	registeredUser, user, err := IsUserAuthenticated(w, r)
	if err != nil{return}
	if !registeredUser{
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendResponseJSON(w, utils.NotLoggedInMessage)
		return
	}

	db, err := model.GetDB(w)
	if err != nil{
		return
	}

	user.Token = ""
	user.TimeAuthenticated, _ = time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")
	db.Update(&user)
	db.Delete(&user)

	w.WriteHeader(http.StatusAccepted)
	utils.SendResponseJSON(w, utils.Message{"User deleted. Bye !"})
}