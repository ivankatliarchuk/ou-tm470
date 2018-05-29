package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// Data type for the block
type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string
	Payload   string
}

var Blockchain []Block

// Hash function to take key block data and return a SHA256 hash as a string
func calculateHash(block Block) string {
  // Update this to reflect unique block has data
	record := string(string(block.Index) + block.Timestamp + block.Payload + block.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Generate a new block based on the old block and new payload
func generateBlock(oldBlock Block, newPayload string) (Block, error) {

  var newBlock Block
  t := time.Now()

  newBlock.Index = oldBlock.Index + 1
  newBlock.Timestamp = t.String()
  newBlock.Payload = newPayload
  newBlock.PrevHash = oldBlock.Hash
  newBlock.Hash = calculateHash(newBlock)

  return newBlock, nil
}

// Check that a block is valid by comparing to previous blocks attributes and
// validating stored hash is correct.  Return a boolean as to the validation
// being correct (true) or incorrect (false)
func isBlockValid(newBlock, oldBlock Block) bool {

  if oldBlock.Index + 1 != newBlock.Index {
    return false
  }

  if oldBlock.Hash != newBlock.PrevHash {
    return false
  }

  if calculateHash(newBlock) != newBlock.Hash {
    return false
  }

  return true
}

// Overwrite the blockchain is a longer chain is presented from other EncodeToString
func replaceChain(newChain []Block) {
  if len(newChain) > len(Blockchain) {
    Blockchain = newChain
  }
}

}
