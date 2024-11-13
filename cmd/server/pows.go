package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type FindHashPOW struct {
	rand       randomizerIntner
	difficulty int
}

type TaskFindHashPOW struct {
	nonce      int
	difficulty int
}

func NewFindHashPOW(rand randomizerIntner, difficulty int) *FindHashPOW {
	return &FindHashPOW{rand: rand, difficulty: difficulty}
}

func (p FindHashPOW) MakeTask() POWTask {
	nonce := p.rand.Intn(1000000)

	return TaskFindHashPOW{
		nonce:      nonce,
		difficulty: p.difficulty,
	}
}

func (p TaskFindHashPOW) CheckSolve(answerStr string) bool {
	answer, err := strconv.Atoi(answerStr)
	if err != nil {
		return false
	}
	hashInput := fmt.Sprintf("%d%d", p.nonce, answer)
	hash := sha256.Sum256([]byte(hashInput))
	hashHex := hex.EncodeToString(hash[:])
	return strings.HasPrefix(hashHex, strings.Repeat("0", p.difficulty))
}

func (p TaskFindHashPOW) Present() string {
	return fmt.Sprintf("%d:%d\n", p.nonce, p.difficulty)
}
