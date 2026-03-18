package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/Asheehan77/Bootdev_Chirpy.git/internal/auth"
	"github.com/Asheehan77/Bootdev_Chirpy.git/internal/database"
	"github.com/google/uuid"
)

func readinessHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	data := []byte("OK")
	writer.Write(data)
}

func (cfg *apiConfig) makechirpHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body    string `json:"body"`
		User_ID string `json:"user_id"`
	}
	type response struct {
		ID      string `json:"id"`
		Create  string `json:"created_at"`
		Update  string `json:"updated_at"`
		Body    string `json:"body"`
		User_ID string `json:"user_id"`
	}

	param := parameters{}
	res := response{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
	} else {
		if len(param.Body) > 140 {
			writer.WriteHeader(400)
			log.Printf("Requested Chirp too long")
			return
		} else {
			res.Body = clean_chirp_body(param.Body)

		}
	}

	userid, err := uuid.Parse(param.User_ID)
	if err != nil {
		writer.WriteHeader(500)
		log.Printf("Error creating chirp: %s\n", err)
		return
	}
	chirpparam := database.CreateChirpParams{
		Body:   param.Body,
		UserID: userid,
	}

	chirp, err := cfg.queries.CreateChirp(context.Background(), chirpparam)
	if err != nil {
		writer.WriteHeader(500)
		log.Printf("Error creating chirp: %s\n", err)
		return
	}

	res.ID = chirp.ID.String()
	res.Create = chirp.CreatedAt.String()
	res.Update = chirp.UpdatedAt.String()
	res.Body = chirp.Body
	res.User_ID = chirp.UserID.String()

	jres, err := json.Marshal(res)

	if err != nil {
		writer.WriteHeader(500)
		log.Printf("Error encoding response: %s\n", err)
		return
	}
	writer.WriteHeader(201)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(jres))
}

func (cfg *apiConfig) getchirpHandler(writer http.ResponseWriter, req *http.Request) {

	type response struct {
		ID      string `json:"id"`
		Create  string `json:"created_at"`
		Update  string `json:"updated_at"`
		Body    string `json:"body"`
		User_ID string `json:"user_id"`
	}

	res_slice := []response{}

	chirps, err := cfg.queries.GetChirps(context.Background())
	if err != nil {
		writer.WriteHeader(500)
		log.Printf("Error getting chirp: %s\n", err)
		return
	}

	for _, chirp := range chirps {
		res := response{}
		res.ID = chirp.ID.String()
		res.Create = chirp.CreatedAt.String()
		res.Update = chirp.UpdatedAt.String()
		res.Body = chirp.Body
		res.User_ID = chirp.UserID.String()
		res_slice = append(res_slice, res)
	}

	jres, err := json.Marshal(res_slice)

	if err != nil {
		writer.WriteHeader(500)
		log.Printf("Error encoding response: %s\n", err)
		return
	}
	writer.WriteHeader(200)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(jres))
}

func clean_chirp_body(body string) string {

	wordlist := []string{"kerfuffle", "sharbert", "fornax"}

	bodywords := strings.Fields(body)
	for i, word := range bodywords {
		if slices.Contains(wordlist, strings.ToLower(word)) {
			bodywords[i] = "****"
		}
	}
	return strings.Join(bodywords, " ")
}

func (cfg *apiConfig) getchirpidHandler(writer http.ResponseWriter, req *http.Request) {

	type response struct {
		ID      string `json:"id"`
		Create  string `json:"created_at"`
		Update  string `json:"updated_at"`
		Body    string `json:"body"`
		User_ID string `json:"user_id"`
	}

	res := response{}
	idstring := req.PathValue("id")
	id, err := uuid.Parse(idstring)
	if err != nil {
		writer.WriteHeader(404)
		log.Printf("Error parsing id: %s\n", err)
		return
	}
	chirp, err := cfg.queries.GetChirp(context.Background(), id)
	if err != nil {
		writer.WriteHeader(404)
		log.Printf("Error getting chirp: %s\n", err)
		return
	}

	res.ID = chirp.ID.String()
	res.Create = chirp.CreatedAt.String()
	res.Update = chirp.UpdatedAt.String()
	res.Body = chirp.Body
	res.User_ID = chirp.UserID.String()

	jres, err := json.Marshal(res)

	if err != nil {
		writer.WriteHeader(500)
		log.Printf("Error encoding response: %s\n", err)
		return
	}
	writer.WriteHeader(200)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(jres))
}

func (cfg *apiConfig) createUserHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID     string `json:"id"`
		Create string `json:"created_at"`
		Update string `json:"updated_at"`
		Email  string `json:"email"`
	}
	param := parameters{}
	res := response{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
	} else if param.Email != "" {

		newPass, err := auth.HashPassword(param.Password)
		if err != nil {
			log.Printf("Error creating user password: %s", err)
			writer.WriteHeader(500)
			return
		}
		newUser := database.CreateUserParams{
			Email:          param.Email,
			HashedPassword: newPass,
		}
		user, err := cfg.queries.CreateUser(context.Background(), newUser)
		if err != nil {
			log.Printf("Error creating user: %s", err)
			writer.WriteHeader(500)
			return
		} else {
			writer.WriteHeader(http.StatusCreated)
			res.ID = user.ID.String()
			res.Create = user.CreatedAt.String()
			res.Update = user.UpdatedAt.String()
			res.Email = user.Email
		}
	} else {
		writer.WriteHeader(400)
		return
	}

	jres, err := json.Marshal(res)

	if err != nil {
		writer.WriteHeader(500)
		log.Printf("Error encoding response: %s", err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(jres))
}

func (cfg *apiConfig) loginHandler(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID     string `json:"id"`
		Create string `json:"created_at"`
		Update string `json:"updated_at"`
		Email  string `json:"email"`
	}
	param := parameters{}
	res := response{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	user, err := cfg.queries.GetUserByEmail(context.Background(), param.Email)
	if err != nil {
		writer.WriteHeader(401)
		return
	}

	match, err := auth.CheckPasswordHash(param.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error authenticating password: %s", err)
		writer.WriteHeader(500)
		return
	}

	if match {
		res.ID = user.ID.String()
		res.Create = user.CreatedAt.String()
		res.Update = user.UpdatedAt.String()
		res.Email = user.Email

		jres, err := json.Marshal(res)

		if err != nil {
			writer.WriteHeader(500)
			log.Printf("Error encoding response: %s", err)
			return
		}
		writer.WriteHeader(200)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte(jres))

	} else {
		writer.WriteHeader(401)
		return
	}

}
