package model

type ChunkGetResp struct {
	Ciphertext string `json:"ciphertext"`
	ChunkNum   int    `json:"chunk_number"`
}
