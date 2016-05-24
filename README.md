# gossip (go search parser)
[![Circle CI](https://circleci.com/gh/shawnohare/gossip.svg?style=svg)](https://circleci.com/gh/shawnohare/gossip)
[![GoDoc](https://godoc.org/github.com/shawnohare/gossip?status.svg)](http://godoc.org/github.com/shawnohare/gossip)

The gossip package provides the ability to parse raw text queries using a
simple DSL into a tree consisting of nodes that are phrases or subqueries
together with modal verbs.  For example:
```go
package main

import (
  "log"

  "github.com/shawnohare/gossip"
)

func main() {
  queryText := `"data science", "machine learning", -hype`
  var root *gossip.Node
  root, _ = gossip.Parse(queryText)
  log.Printf("The root has %d children.", len(root.Children()))
}
```
Here, the `root` node contains three children: one for each of the phrases
"data science", "machine learning", and "hype". The "hype" term is modified
by the modal verb "must not", represented in the search query by a minus sign.
