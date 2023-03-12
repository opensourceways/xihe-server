package app

import "errors"

func (s *cloudService) ReleasePod(cmd *RelasePodCmd) (code string, err error) {
	// get pod
	p, err := s.podRepo.GetPodInfo(cmd.PodId)
	if err != nil {
		return
	}

	// is owner
	if !p.Pod.IsOnwer(cmd.User) {
		code = errorNoAuthorized
		err = errors.New("no authorize")

		return
	}

	// check status
	if !p.CanRelease() {
		code = errorNotRunning
		err = errors.New("pod not running")

		return
	}

	// relase
	err = s.cloudService.ReleasePod(&p.Pod)

	return
}

func (s *cloudService) Get(cmd *PodInfoCmd) (dto PodInfoDTO, err error) {
	p, err := s.podRepo.GetUserCloudIdPod(cmd.Owner, cmd.CloudId)
	if err != nil {
		return
	}

	dto.toPodInfoDTO(&p.PodInfos[0])

	return
}
