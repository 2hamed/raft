package raft

type Options struct {
	listenPort int
	listenAddr string
}

type OptionFunc func(*Options) *Options

func WithListenPort(listenPort int) OptionFunc {
	return func(o *Options) *Options {
		o.listenPort = listenPort
		return o
	}
}

func WithListenAddr(listenAddr string) OptionFunc {
	return func(o *Options) *Options {
		o.listenAddr = listenAddr
		return o
	}
}
