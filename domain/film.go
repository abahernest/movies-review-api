package domain

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/Kamva/mgm/v2"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"net/http"
	//"go.mongodb.org/mongo-driver/bson"
)

type Film struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string `json:"title" bson:"title"`
	CommentCount     int64  `json:"comment_count" bson:"comment_count"`
	ReleaseDate      string `json:"release_date" bson:"release_date"`
}

type PaginatedFilm struct {
	Pagination *mongopagination.PaginatedData `json:"pagination" bson:"pagination"`
	Data       []Film                         `json:"data" bson:"data"`
}

type StarwarsDataHash struct {
	mgm.DefaultModel `bson:",inline"`
	Hash             string `json:"hash" bson:"hash"`
	Count            int64  `json:"count" bson:"count"`
}

type StarwarsApiResponse struct {
	Next     string `json:"next" bson:"next"`
	Previous string `json:"previous" bson:"previous"`
	Count    int64  `json:"count" bson:"count"`
	Results  []Film `json:"results" bson:"results"`
}

func UpdateFilmFromSource(logger *zap.Logger, api string) error {
	// Make a request to the external API and retrieve the new data

	if api == "" {
		api = "https://swapi.dev/api/films/"
	}

	resp, err := http.Get(api)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the JSON response into a new struct
	var data StarwarsApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	var (
		savedDataArr []StarwarsDataHash
		savedData    StarwarsDataHash
	)

	// fetch last saved hash
	err = mgm.Coll(&StarwarsDataHash{}).SimpleFind(&savedDataArr, bson.M{})
	if err != nil {
		return err
	}

	if len(savedDataArr) > 0 {
		savedData = savedDataArr[0]
	}

	// hash results
	newHash, err := hashStarwarsApiData(data.Results)
	if err != nil {
		return err
	}

	// if no new data
	noNewData, err := compareStarwarsDataHash(&savedData, newHash)
	if err != nil {
		return err
	}

	// update db with new data
	if !noNewData {
		for _, starwarsFilm := range data.Results {
			var film []Film
			if err = mgm.Coll(&Film{}).SimpleFind(&film, bson.M{"title": starwarsFilm.Title}); err != nil {
				logger.Error("Error:", zap.Error(err))
				continue
			}
			if len(film) <= 0 {
				err = mgm.Coll(&Film{}).Create(&Film{
					Title:       starwarsFilm.Title,
					ReleaseDate: starwarsFilm.ReleaseDate,
				})
				if err != nil {
					logger.Error("Error:", zap.Error(err))
					continue
				}
			}
		}
	}

	// go to next page
	if data.Next != "" {
		if err = UpdateFilmFromSource(logger, data.Next); err != nil {
			return err
		}
	}

	return nil
}

func compareStarwarsDataHash(savedHashObj *StarwarsDataHash, newDataHash string) (bool, error) {
	if savedHashObj.Hash != newDataHash {
		//save new hash
		savedHashObj.Hash = newDataHash
		err := mgm.Coll(&StarwarsDataHash{}).UpdateWithCtx(context.Background(), savedHashObj)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func hashStarwarsApiData(data []Film) (string, error) {
	var hash string
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return hash, err
	}

	hasher := md5.New()
	hasher.Write(dataJSON)
	hash = hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}

type FilmRepository interface {
	GetById(ctx context.Context, id string) (*Film, error)
	FetchPaginatedFilms(ctx context.Context, page, limit int64) (*PaginatedFilm, error)
}

type FilmUsecase interface {
	FetchFilmsFromAllSources(ctx context.Context, page, limit int64) (*PaginatedFilm, error)
}
