package app

func (s *cloudService) Get(cmd *PodInfoCmd) (dto PodInfoDTO, err error) {
	p, _, err := s.cloudService.CheckUserCanSubscribe(cmd.User, cmd.CloudId)
	if err != nil {
		return
	}

	cloudConf, err := s.cloudRepo.GetCloudConf(p.CloudId)
	if err != nil {
		return
	}

	err = dto.toPodInfoDTO(&p, &cloudConf)

	return
}
