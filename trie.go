package main

import "fmt"

const alphabetSize = 26 // For lowercase English letters 'a' through 'z'

// Node represents a node in the Trie structure.
type Node struct {
	children    [alphabetSize]*Node // Changed from map[byte]*Node to fixed-size array
	isEndOfWord bool                // True if this node marks the end of a word
}

// NewNode creates and returns a new Trie Node.
func NewNode() *Node {
	return &Node{} // Arrays are zero-initialized, so children will be nil
}

// Trie represents the Trie data structure.
type Trie struct {
	root *Node // The root node of the Trie
}

// NewTrie creates and returns a new Trie.
func NewTrie() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

// charToIndex converts a lowercase English letter byte to its corresponding array index (0-25).
// It panics if the character is not a lowercase English letter.
func charToIndex(char byte) int {
	if char >= 'a' && char <= 'z' {
		return int(char - 'a')
	}
	// For production code, you might want to return an error or a special value
	// instead of panicking, or handle non-lowercase inputs upstream.
	panic("trie: character not a lowercase English letter")
}

func indexToChar(i int) byte {
	if i >= 0 && i < alphabetSize {
		return byte('a' + i)
	}

	panic("trie: char index out of range")
}

// Insert adds a word to the Trie.
// Assumes input 'word' contains only lowercase English letters.
func (t *Trie) Insert(word string) {
	currentNode := t.root
	for i := 0; i < len(word); i++ {
		idx := charToIndex(word[i])
		if currentNode.children[idx] == nil {
			currentNode.children[idx] = NewNode()
		}
		currentNode = currentNode.children[idx]
	}
	currentNode.isEndOfWord = true
}

// Search checks if a word exists in the Trie.
// Assumes input 'word' contains only lowercase English letters.
func (t *Trie) Search(word string) bool {
	currentNode := t.root
	for i := 0; i < len(word); i++ {
		idx := charToIndex(word[i])
		if currentNode.children[idx] == nil {
			return false // Character not found, word doesn't exist
		}
		currentNode = currentNode.children[idx]
	}
	return currentNode.isEndOfWord // True if it's a complete word, false otherwise (e.g., prefix)
}

// StartsWith checks if there is any word in the Trie that starts with the given prefix.
// Assumes input 'prefix' contains only lowercase English letters.
func (t *Trie) StartsWith(prefix string) bool {
	currentNode := t.root
	for i := 0; i < len(prefix); i++ {
		idx := charToIndex(prefix[i])
		if currentNode.children[idx] == nil {
			return false // Character not found, no word starts with this prefix
		}
		currentNode = currentNode.children[idx]
	}
	return true // Prefix found
}

// Delete removes a word from the Trie.
// This implementation performs a "soft" delete by just unmarking isEndOfWord.
// Assumes input 'word' contains only lowercase English letters.
func (t *Trie) Delete(word string) bool {
	currentNode := t.root
	// We need to keep track of the path for potential hard deletion later,
	// but for soft delete, just direct traversal is enough.

	for i := 0; i < len(word); i++ {
		idx := charToIndex(word[i])
		if currentNode.children[idx] == nil {
			return false // Word not found
		}
		currentNode = currentNode.children[idx]
	}

	if !currentNode.isEndOfWord {
		return false // Word exists as a prefix but not as a complete word
	}

	currentNode.isEndOfWord = false // Unmark as end of word

	// Hard delete logic (more complex):
	// To perform a hard delete, you would need to iterate backwards from the
	// currentNode through the path of nodes visited. For each node, check if:
	// 1. It's no longer 'isEndOfWord'.
	// 2. It has no other children (i.e., it's not a prefix to any other word).
	// If both conditions are met, remove the node from its parent's children array.

	return true
}

// CollectAllWordsStartingWith collects all words in the Trie that start with the given prefix.
// Assumes input 'prefix' contains only lowercase English letters.
func (t *Trie) CollectAllWordsStartingWith(prefix string) []string {
	var words []string
	currentNode := t.root

	// Traverse to the end of the prefix
	for i := 0; i < len(prefix); i++ {
		idx := charToIndex(prefix[i])
		if currentNode.children[idx] == nil {
			return []string{} // No words start with this prefix
		}
		currentNode = currentNode.children[idx]
	}

	// Now, perform a DFS from the current node to collect all words
	t.collectWordsDFS(currentNode, prefix, &words)

	return words
}

// collectWordsDFS is a helper function for CollectAllWordsStartingWith that performs a DFS.
func (t *Trie) collectWordsDFS(node *Node, currentWord string, words *[]string) {
	if node.isEndOfWord {
		*words = append(*words, currentWord)
	}

	// Iterate over the fixed-size array
	for i := 0; i < alphabetSize; i++ {
		childNode := node.children[i]
		if childNode != nil {
			// Convert index back to char and append
			char := indexToChar(i)
			t.collectWordsDFS(childNode, currentWord+string(char), words)
		}
	}
}

// Example Usage (main function to test):

func main() {
	trie := NewTrie()

	trie.Insert("cat")
	trie.Insert("car")
	trie.Insert("card")
	trie.Insert("apple")
	trie.Insert("app")
	trie.Insert("application")

	fmt.Println("Search 'cat':", trie.Search("cat"))     // true
	fmt.Println("Search 'car':", trie.Search("car"))     // true
	fmt.Println("Search 'apple':", trie.Search("apple")) // true
	fmt.Println("Search 'app':", trie.Search("app"))     // true
	fmt.Println("Search 'card':", trie.Search("card"))   // true
	fmt.Println("Search 'ca':", trie.Search("ca"))       // false (prefix only)
	fmt.Println("Search 'cow':", trie.Search("cow"))     // false

	fmt.Println("Starts with 'ca':", trie.StartsWith("ca"))   // true
	fmt.Println("Starts with 'app':", trie.StartsWith("app")) // true
	fmt.Println("Starts with 'co':", trie.StartsWith("co"))   // false

	fmt.Println("Words starting with 'a':", trie.CollectAllWordsStartingWith("a"))     // [apple app application]
	fmt.Println("Words starting with 'app':", trie.CollectAllWordsStartingWith("app")) // [apple app application]
	fmt.Println("Words starting with 'z':", trie.CollectAllWordsStartingWith("z"))     // []

	fmt.Println("Delete 'app':", trie.Delete("app"))                                    // true
	fmt.Println("Search 'app' after delete:", trie.Search("app"))                       // false
	fmt.Println("Search 'apple' after 'app' delete:", trie.Search("apple"))             // true (apple still exists)
	fmt.Println("Search 'application' after 'app' delete:", trie.Search("application")) // true

	fmt.Println("Delete 'nonexistent':", trie.Delete("nonexistent")) // false

	// Demonstrating panic for invalid input (uncomment to test):
	// trie.Insert("ApPle") // Panics because 'A' is not lowercase
}
