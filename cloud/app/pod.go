package app

func (s *cloudService) Get(cmd *PodInfoCmd) (dto PodInfoDTO, err error) {
	p, _, err := s.cloudService.CheckUserCanSubsribe(cmd.User, cmd.CloudId)
	if err != nil {
		return
	}

	cloudConf, err := s.cloudRepo.GetCloudConf(p.CloudId)
	if err != nil {
		return
	}

	dto.toPodInfoDTO(&p, &cloudConf)

	return
}
