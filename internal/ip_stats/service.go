package repo

type IPStatsServiceI interface {
	UpSert()
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
