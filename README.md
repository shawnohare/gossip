# gossip (go search parser)
[![Circle CI](https://circleci.com/gh/shawnohare/gossip.svg?style=svg)](https://circleci.com/gh/shawnohare/wemdigo)
[![GoDoc](https://godoc.org/github.com/shawnohare/gossip?status.svg)](http://godoc.org/github.com/shawnohare/wemdigo)

The gossip package provides the ability to parse raw text queries using a
simple DSL into a tree consisting of nodes that are phrases or subqueries
together with modal verbs.  For example:
```
package main

import (
  "log"

  "github.com/shawnohare/gossip"
)

func main() {
  queryText := `"data science", "machine learning", -hype`
  var root *gossip.Node
  root, _ = gossip.Parse(queryText)
  log.Println(root.Children()[2].VerbString())
  // output:
  // 
}
```
Here, the `root` node contains three children: one for each of the phrases
"data science", "machine learning", and "hype". The "hype"
