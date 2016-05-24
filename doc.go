// Package gossip specifies a basic search DSL and provides parsing functions
// to convert raw text search queries into structured search query trees.
// These trees are structured to be easily converted into other search DSLs
// such as that used by Elasticsearch.
//
// Minimal search DSL
//
// A basic query consists of a set of words or phrases and modal
// verbs applied to these objects.  Phrases typically indicate
// multi-word searches where proximity is important.  These also serve
// as a way to search over otherwise reserved symbols.
//
// The search DSL does not support parsing arbitrary propsitions in first
// order logic, but does allow infinitely nested subqueries, which provides
// a great deal of flexibility when combined with the three supported
// modal verbs: Must, Should, MustNot.
//
// A typical valid search query might look like
//   "data science" +[math -hype]
// which is interpreted as a search that should (but need not) contain
// the phrase "data science" and must contain the term "math" but not
// the term "hype".
//
// The symbols ", [, ], +, - and , are are reserved and have context-dependent
// special interpretations.
//
// Phrase Literals
//
// Phrases are typically multi-word sequences where word proximity matters.
// They are contained between matching quotation marks.  For example,
//   `"data science"`
// represents a query where matching documents should contain the
// the adjacent words "data" and "science".  No symbols are parsed inside
// a phrase literal, except for the initial and terminal quotation marks
// themselves.  For example,
//   `+programming +"c++"`
// specifies a search that must contain both programming and c++.  The
// query
//   `c++`
// results in a parse error, since the + is interpreted as a modal verb.
//
//
// Modal verbs
//
// The model verbs should, must, and must not are supported.  These verbs
// can apply words, phrase, and nested queries. By default, all objects are
// assumed to be modified by the modal verb "should".  Put another way,
// a query comprised of a set of words is equivalent to a disjunction,
// whereas a query comprised entirely of words modified by "must" represents
// a conjunction.
//
// +p indicates that the object p must be matched: e.g.,
// the query
//   +moms +"video gamers"
// should be matched by documents
// that contain both moms and the literal phrase "video gamers".
//
// -p indicates that the phrase must not be matched.
//
// [] encloses a subquery. The + and - operators specify whether
// this subuery must match or must not match.  By default the subquery
// should (but need not) match. Using subqueries allows for a high level
// of search customization.
package gossip
