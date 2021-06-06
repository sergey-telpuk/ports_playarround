package pkg

type PortDomainService interface {
	Run(port string) error
	Close()
}

type PortApiService interface {
	//TODO add some api
}
