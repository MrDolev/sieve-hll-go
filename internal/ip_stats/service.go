package repo

type IPStatsServiceI interface {
	UpSert()
	Collect()
}

type IPStatsService struct {
	repo IPStatsRepoI
}

func NewIPStatsService(repo IPStatsRepoI) *IPStatsService {
	return &IPStatsService{
		repo: repo,
	}
}

func (service *IPStatsService) UpSert() {
	service.repo.UpSert()
	return
}

func (service *IPStatsService) Collect() {

}
