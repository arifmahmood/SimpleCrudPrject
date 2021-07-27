package api

import (
	"github.com/go-chi/chi"
	"net/http"
	"simple-crud-project/errors"
	"simple-crud-project/model"
	"strings"
)

type createUrlRequest struct {
	Url string `json:"url"`
	Description string `json:"description"`
}

func (c *createUrlRequest) validate() *validationError {
	c.Url = strings.TrimSpace(c.Url)

	errV := validationError{}
	if c.Url == "" {
		errV.add("Url", "is required")
	}

	if len(errV) > 0 {
		return &errV
	}

	return nil
}


func (rt *Router) CreateNewUrl(w http.ResponseWriter, r *http.Request) {
	body := createUrlRequest{}
	if err := parseBody(r, &body); err != nil {
		handleAPIError(w, newAPIError("Unable to parse body", errBadRequest, err))
		return
	}

	if err := body.validate(); err != nil {
		handleAPIError(w, newAPIError("Invalid data", errInvalidData, err))
		return
	}

	url := &model.Url{
		Url:     body.Url,
		Description: body.Description,
		Status:  "A",
	}

	if err := rt.urlRepo.Create(url); err != nil {
		if err == errors.ErrDuplicateKey {
			vErr := validationError{}
			vErr.add("Url", "is not unique")
			handleAPIError(w, newAPIError("Invalid data", errEntityNotUnique, &vErr))
			return
		}

		panic(newAPIError("Internal Server Error", errInternalServer, err))
	}

	resp := response{
		code: http.StatusOK,
		Data: url,
	}

	resp.serveJSON(w)
}



func (rt *Router) GetUrl(writer http.ResponseWriter, request *http.Request) {
	urlName := strings.TrimSpace(chi.URLParam(request, "urlName"))

	url, err := rt.urlRepo.Fetch(urlName)
	if err != nil {
		panic(newAPIError("DB failed", errInternalServer, err))
	}
	if url == nil {
		handleAPIError(writer, newAPIError("Url not found", errUserNotFound, nil))
	}

	resp := response{
		code: http.StatusOK,
		Data: url,
	}

	resp.serveJSON(writer)

}

