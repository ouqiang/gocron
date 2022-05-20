package service

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
)

type LdapService struct{}

func (ldapService *LdapService) Match(username, password string) bool {
	settings := models.Setting{}
	l, err := ldap.DialURL(settings.Get(models.LdapCode, models.LdapKeyUrl).Value)
	if err != nil {
		logger.Debug(err)
		return false
	}
	defer l.Close()

	err = l.Bind(settings.Get(models.LdapCode, models.LdapKeyBindDn).Value, settings.Get(models.LdapCode, models.LdapKeyBindPassword).Value)
	if err != nil {
		logger.Debug(err)
		return false
	}

	req := ldap.NewSearchRequest(
		settings.Get(models.LdapCode, models.LdapKeyBaseDn).Value,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf(settings.Get(models.LdapCode, models.LdapKeyFilterRule).Value, username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(req)
	if err != nil {
		logger.Debug(err)
		return false
	}

	// 如果没有数据返回或者超过1条数据返回,这对于用户认证而言都是不允许的.
	// 前这意味着没有查到用户,后者意味着存在重复数据
	if len(sr.Entries) != 1 {
		logger.Debug("User does not exist or too many entries returned")
		return false
	}

	err = l.Bind(sr.Entries[0].DN, password)
	if err != nil {
		logger.Debug(err)
		return false
	}
	return true
}
