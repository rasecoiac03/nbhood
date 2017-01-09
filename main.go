package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type NbhInfo struct {
	Nbh     string `json:"name" bson:"nbh"`
	Feature string `json:"feature" bson:"feature"`
}

type Summaries []Summary

type Summary struct {
	Nbh   NbhInfo `json:"nbh" bson:"_id"`
	Count int     `json:"count" bson:"count"`
}

type Features []Feature

type Feature struct {
	Feature     string `json:"feature" bson:"feature"`
	FeatureName string `json:"featureName" bson:"featureName"`
	Address     string `json:"address" bson:"address"`
}

var nbhsCollection *mgo.Collection

func init() {
	mongoHost := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	nbhsCollection = session.DB("nbhood").C("nbhs")
}

func main() {
	middle := interpose.New()
	router := mux.NewRouter()
	middle.UseHandler(router)
	router.HandleFunc("/summary", FindSummaries).Methods("GET")
	router.HandleFunc("/summary/{sort}", FindSummaries).Methods("GET")
	router.HandleFunc("/nbh/{nbh}", FindNbh).Methods("GET")
	router.HandleFunc("/nbh/{nbh}/{feature}", FindNbh).Methods("GET")
	http.ListenAndServe(":9999", middle)
}

func FindSummaries(w http.ResponseWriter, req *http.Request) {
	sort := mux.Vars(req)["sort"]

	var summaries Summaries
	if sort != "" {
		for _, feature := range strings.Split(sort, ",") {
			summaries = append(summaries, querySummary(feature).(Summary))
		}
	} else {
		summaries = querySummary("").(Summaries)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

func FindNbh(w http.ResponseWriter, req *http.Request) {
	nbh := mux.Vars(req)["nbh"]
	feature := mux.Vars(req)["feature"]
	features := queryNbh(nbh, feature)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(features)
}

func queryNbh(nbh string, feature string) Features {
	sel := bson.M{
		"feature":     1,
		"featureName": 1,
		"address":     1,
	}
	match := bson.M{
		"nbh": nbh,
	}
	if feature != "" {
		match["feature"] = feature
	}
	var features Features
	err := nbhsCollection.Find(match).Select(sel).All(&features)
	checkError(err)
	return features
}

func querySummary(feature string) interface{} {
	var result interface{}
	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"nbh":     "$nbh",
				"feature": "$feature",
			},
			"count": bson.M{"$sum": 1},
		},
	}
	sort := bson.M{
		"$sort": bson.M{
			"_id.feature": 1,
			"count":       -1,
		},
	}
	if feature == "" {
		var summaries Summaries
		err := nbhsCollection.Pipe([]bson.M{group, sort}).All(&summaries)
		checkError(err)
		result = summaries
	} else {
		match := bson.M{
			"$match": bson.M{
				"feature": feature,
			},
		}
		var summary Summary
		err := nbhsCollection.Pipe([]bson.M{match, group, sort}).One(&summary)
		checkError(err)
		result = summary
	}
	return result
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
