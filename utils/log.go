package utils

import "github.com/sirupsen/logrus"

func DoLog(userid, username, action, access, result string) {
	logrus.Infof("| userid: %s | username: %s | action: %s | access: %s | result: %s |",
		userid, username, action, access, result)
}
