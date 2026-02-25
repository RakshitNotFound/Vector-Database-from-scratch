package main

import (
	"fmt"
	"math"
	"sort"
)

type Vector struct {
	ID       string
	Values   []float32
	Metadata string
}

type SearchResult struct {
	Vector     Vector
	Similarity float64
}

type VectorDB struct {
	storage []Vector
}

func NewVectorDB() *VectorDB {
	return &VectorDB{
		storage: make([]Vector, 0),
	}
}

func (db *VectorDB) Insert(v Vector) {
	db.storage = append(db.storage, v)
}

// CosineSimilarity calculates how similar two vectors are.
// Range is -1 to 1 (1 being identical).
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (db *VectorDB) Query(queryValues []float32, topK int) []SearchResult {
	results := make([]SearchResult, 0, len(db.storage))

	for _, v := range db.storage {
		similarity := CosineSimilarity(queryValues, v.Values)
		results = append(results, SearchResult{Vector: v, Similarity: similarity})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	if len(results) > topK {
		return results[:topK]
	}
	return results
}

func main() {
	db := NewVectorDB()

	// Simulating "embeddings" for three concepts
	db.Insert(Vector{"1", []float32{1.0, 0.1, 0.0}, "King"})
	db.Insert(Vector{"2", []float32{0.9, 0.2, 0.0}, "Queen"})
	db.Insert(Vector{"3", []float32{0.0, 0.8, 0.9}, "Apple"})

	// Searching for something similar to "Royalty"
	query := []float32{0.95, 0.15, 0.0}
	matches := db.Query(query, 2)

	fmt.Println("Top Search Results:")
	for _, match := range matches {
		fmt.Printf("ID: %s | Name: %s | Score: %.4f\n",
			match.Vector.ID, match.Vector.Metadata, match.Similarity)
	}
}
