package agent

import ()

type Load interface {
	Load(file string) error
}

type Crawl interface {
}

type Save interface {
}
