package initialize

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Trie 短语组成的Trie树.
type Trie struct {
	Root *Node
}

// Node Trie树上的一个节点.
type Node struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node
	Failure    *Node
	Parent     *Node
	depth      int
}

// LinkList ...
type LinkList struct {
	head  *listNode
	tail  *listNode
	count int64
}

// Push appends a node
func (list *LinkList) Push(v interface{}) {
	node := &listNode{
		Value: v,
	}
	if list.head == nil {
		list.head = node
	} else {
		list.tail.Next = node

	}
	list.tail = node
	list.count++
}

// Pop returns the value of the first node
func (list *LinkList) Pop() interface{} {
	if list.Empty() {
		return nil
	}

	n := list.head
	list.head = n.Next
	list.count--
	return n.Value
}

// Empty returns true if there is none node
func (list *LinkList) Empty() bool {
	return list.count == 0
}

type listNode struct {
	Value interface{}
	Next  *listNode
}

// BuildFailureLinks 更新Aho-Corasick的失败表
func (tree *Trie) BuildFailureLinks() {
	for node := range tree.bfs() {
		pointer := node.Parent
		var link *Node
		for link == nil {
			if pointer.IsRootNode() {
				link = pointer
				break
			}
			link = pointer.Failure.Children[node.Character]
			pointer = pointer.Failure

		}
		// fmt.Printf("%s[%d] link to %s[%d] \n", string(node.Character), node.depth, string(link.Character), link.depth)
		node.Failure = link

	}
	// fmt.Println("finish build failure link")
}

// bfs Breadth First Search
func (tree *Trie) bfs() <-chan *Node {
	ch := make(chan *Node) // 创建通道
	go func() {
		queue := new(LinkList) // 初始化队列
		for _, child := range tree.Root.Children {
			queue.Push(child) // 将根节点的子节点加入队列
		}

		for !queue.Empty() { // 当队列不为空时继续遍历
			n := queue.Pop().(*Node) // 从队列中取出一个节点
			ch <- n                  // 将节点发送到通道
			for _, child := range n.Children {
				queue.Push(child) // 将当前节点的子节点加入队列
			}
		}

		close(ch) // 遍历完成后关闭通道
	}()
	return ch
}

// NewTrie 新建一棵Trie
func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(0),
	}
}

// Add 添加若干个词
func (tree *Trie) Add(words ...string) {
	for _, word := range words {
		tree.add(word)
	}
}

func (tree *Trie) add(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode(r)
			newNode.depth = current.depth + 1
			newNode.Parent = current
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.isPathEnd = true
		}
	}
}

type ac struct {
	results []string
}

func (ac *ac) fail(node *Node, c rune) *Node {
	var next *Node
	for {
		next = ac.next(node.Failure, c)
		if next == nil {
			if node.IsRootNode() {
				return node
			}
			node = node.Failure
			continue
		}

		return next
	}

}

func (ac *ac) next(node *Node, c rune) *Node {
	next, ok := node.Children[c]
	if ok {
		return next
	}
	return nil
}

func (ac *ac) output(node *Node, runes []rune, position int) {
	if node == nil || node.IsRootNode() {
		return
	}

	if node.IsPathEnd() {
		ac.results = append(ac.results, string(runes[position+1-node.depth:position+1]))
	}

	ac.output(node.Failure, runes, position)
}

func (ac *ac) firstOutput(node *Node, runes []rune, position int) string {
	if node == nil || node.IsRootNode() {
		return ""
	}

	if node.IsPathEnd() {
		return string(runes[position+1-node.depth : position+1])
	}

	return ac.firstOutput(node.Failure, runes, position)
}

func (ac *ac) replace(node *Node, runes []rune, position int, replace rune) {
	if node == nil || node.IsRootNode() {
		return
	}

	if node.IsPathEnd() {
		for i := position + 1 - node.depth; i < position+1; i++ {
			runes[i] = replace
		}
	}
	ac.replace(node.Failure, runes, position, replace)
}

// Replace 词语替换
func (tree *Trie) Replace(text string, character rune) string {
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := 0; position < len(runes); position++ {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		ac.replace(node, runes, position, character)
	}

	return string(runes)
}

// Filter 直接过滤掉字符串中的敏感词
func (tree *Trie) Filter(text string) string {
	var (
		parent      = tree.Root
		current     *Node
		left        = 0
		found       bool
		runes       = []rune(text)
		length      = len(runes)
		resultRunes = make([]rune, 0, length)
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found {
			resultRunes = append(resultRunes, runes[left])
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() {
			left = position + 1
		}
		parent = current
	}

	resultRunes = append(resultRunes, runes[left:]...)
	return string(resultRunes)
}

// Validate 验证字符串是否合法，如不合法则返回false和检测到
// 的第一个敏感词
func (tree *Trie) Validate(text string) (bool, string) {
	const EMPTY = ""
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := 0; position < len(runes); position++ {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		if first := ac.firstOutput(node, runes, position); len(first) > 0 {
			return false, first
		}
	}

	return true, EMPTY
}

// FindIn 判断text中是否含有词库中的词
func (tree *Trie) FindIn(text string) (bool, string) {
	validated, first := tree.Validate(text)
	return !validated, first
}

// FindAll 找有所有包含在词库中的词
func (tree *Trie) FindAll(text string) []string {
	var (
		node  = tree.Root
		next  *Node
		runes = []rune(text)
	)

	var ac = new(ac)
	for position := 0; position < len(runes); position++ {
		next = ac.next(node, runes[position])
		if next == nil {
			next = ac.fail(node, runes[position])
		}

		node = next
		ac.output(node, runes, position)
	}

	return ac.results

}

// NewNode 新建子节点
func NewNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node, 0),
	}
}

// NewRootNode 新建根节点
func NewRootNode(character rune) *Node {
	root := &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node, 0),
		depth:      0,
	}

	root.Failure = root

	return root
}

// IsLeafNode 判断是否叶子节点
func (node *Node) IsLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断是否为根节点
func (node *Node) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断是否为某个路径的结束
func (node *Node) IsPathEnd() bool {
	return node.isPathEnd
}

type Filter struct {
	trie       *Trie
	noise      *regexp.Regexp
	buildVer   int64
	updatedVer int64
}

// New 返回一个敏感词过滤器
func New() *Filter {
	return &Filter{
		trie:  NewTrie(),
		noise: regexp.MustCompile(`[\|\s&%$@*]+`),
	}
}

// UpdateNoisePattern 更新去噪模式
func (filter *Filter) UpdateNoisePattern(pattern string) {
	filter.noise = regexp.MustCompile(pattern)
}

// LoadWordDict 加载敏感词字典
func (filter *Filter) LoadWordDict(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return filter.Load(f)
}

// LoadNetWordDict 加载网络敏感词字典
func (filter *Filter) LoadNetWordDict(url string) error {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	rsp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	return filter.Load(rsp.Body)
}

// Load common method to add words
func (filter *Filter) Load(rd io.Reader) error {
	buf := bufio.NewReader(rd)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		filter.AddWord(string(line))
	}

	return nil
}

func (filter *Filter) updateFailureLink() {
	if filter.buildVer != filter.updatedVer {
		// fmt.Println("update failure link")
		filter.trie.BuildFailureLinks()
		filter.buildVer = filter.updatedVer
	}
}

// AddWord 添加敏感词
func (filter *Filter) AddWord(words ...string) {
	filter.trie.Add(words...)
	filter.updatedVer = time.Now().UnixNano()
}

// Filter 过滤敏感词
func (filter *Filter) Filter(text string) string {
	filter.updateFailureLink()
	return filter.trie.Filter(text)
}

// Replace 和谐敏感词
func (filter *Filter) Replace(text string, repl rune) string {
	filter.updateFailureLink()
	return filter.trie.Replace(text, repl)
}

// FindIn 检测敏感词
func (filter *Filter) FindIn(text string) (bool, string) {
	filter.updateFailureLink()
	text = filter.RemoveNoise(text)
	return filter.trie.FindIn(text)
}

// FindAll 找到所有匹配词
func (filter *Filter) FindAll(text string) []string {
	filter.updateFailureLink()
	return filter.trie.FindAll(text)
}

// Validate 检测字符串是否合法
func (filter *Filter) Validate(text string) (bool, string) {
	filter.updateFailureLink()
	text = filter.RemoveNoise(text)
	return filter.trie.Validate(text)
}

// RemoveNoise 去除空格等噪音
func (filter *Filter) RemoveNoise(text string) string {
	return filter.noise.ReplaceAllString(text, "")
}
