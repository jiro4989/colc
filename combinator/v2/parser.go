package combinator

// Combinator はコンビネータである。
type Combinator struct {
	Name      string `json:"name"`
	ArgsCount int    `json:"argsCount"`
	Format    string `json:"format"`
}
type Combinators []Combinator

type Node struct {
	Text   string
	Parent *Node
	Nodes  []*Node
}

func (n *Node) IsRoot() bool {
	return n.Parent == nil
}

func (n *Node) HasParent() bool {
	return n.Parent != nil
}

func (n *Node) IsLeaf() bool {
	return len(n.Nodes) <= 0
}

func Parse(code string) Node {
	return Node{}
}
