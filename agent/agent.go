package agent

import ()

type Load interface {
	Load(file string) error
}

type Crawl interface {
	Crawl(loader *Loader) error
}

type Save interface {
}
