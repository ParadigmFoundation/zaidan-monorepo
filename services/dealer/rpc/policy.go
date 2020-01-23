package rpc

type PolicyMode int

const (
	PolicyBlackList PolicyMode = iota
	PolicyWhiteList
)
