package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/karthikmoulya/mongoapi/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongoDBURIString"
const dbName = "netflix"
const colName = "watchlist"

//MOST IMPORTANT
var collection *mongo.Collection

//Connect With mongoDB

func init(){
	//client option
	clientOptions := options.Client().ApplyURI(connectionString)

	//connect to mongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName)

	//collection instance
	fmt.Println("Collection instance is ready")
}

//MongoDb helpers - seperate file

//insert 1 record

func insertOneMovie(movie model.Netflix){
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Inserted 1 movie in db with id: ", inserted.InsertedID)
}

//update 1 record
func updateOneMovie(movieId string){
	id, err := primitive.ObjectIDFromHex(movieId)

	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil{
		log.Fatal()
	}

	fmt.Println("modified count:", result.ModifiedCount)
	}

//Delete 1 movie
func deleteOneMovie(movieId string){
	id, err := primitive.ObjectIDFromHex(movieId)

	if err != nil{
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(),filter)

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Movie got deleted with delete count: ", deleteCount)
}

//delete all records from mongodb
func deleteAllMovie() int64{
	deleteCount, err := collection.DeleteMany(context.Background(), bson.D{{}})

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Number of movies deleted: ", deleteCount.DeletedCount)
	return deleteCount.DeletedCount
}

//get all movies from database
func getAllMovies() []primitive.M{
	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil{
		log.Fatal(err)
	}

	var movies []primitive.M

	for cur.Next(context.Background()){
		var movie bson.M
		err := cur.Decode(&movie)

		if err != nil{
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}

	defer cur.Close(context.Background())
	return movies
}

//Actual Controller - file

func GetMyAllMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkasWatched(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r) 
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params)
}

func DeleteMyOneMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMyMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}