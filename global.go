package main

type GlobalContext struct {
	NavItems NavItems
}

type NavItem struct {
	Name string
	Link string
}

type NavItems []NavItem

var globalContext *GlobalContext = nil

func initGlobalContext() {
	globalContext = &GlobalContext{NavItems: NavItems{
		NavItem{"c0dart", "/c0dart"},
		NavItem{"Blog", "/blog"},
		NavItem{"Github", "https://github.com/sameer"},
	}}
}
