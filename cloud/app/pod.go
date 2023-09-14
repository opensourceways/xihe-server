package app

func (s *cloudService) Get(cmd *PodInfoCmd) (dto PodInfoDTO, err error) {
	p, _, err := s.cloudService.CheckUserCanSubsribe(cmd.User, cmd.CloudId)
	if err != nil {
		return dto, err
	}

	dto.toPodInfoDTO(&p)

	return
}
