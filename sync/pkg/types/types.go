package types

import (
	"time"

	"sync/pkg/util"

	"github.com/astaxie/beego/logs"
)

type Users struct {
	CorpId     string `json:"corp_id"`    //-- 公司corps.id
	Tye        string `json:"tye"`        //-- 角色类型 0-用户; 1-高管; 2-高管后台管理员
	Pwd        string `json:"pwd"`        //-- 登录密码
	Nickname   string `json:"nickname"`   //-- 昵称
	Realname   string `json:"realname"`   //-- 真实姓名
	Mob        string `json:"mob"`        //-- 手机号码 *候选键* 可做登陆名
	Mail       string `json:"mail"`       //-- 电邮 *候选键* 可做登陆名
	Sex        int    `json:"sex"`        //-- 性别 0-未知, 1-男, 2-女
	Qq         string `json:"qq"`         //-- QQ号
	Headimgurl string `json:"headimgurl"` //-- 头像url
	IsDel      int    `json:"is_del"`     //-- 是否注销 0-启用, 1-注销
	CreateAt   string `json:"create_at"`  //-- 创建时间 yyyy-mm-dd hh:nn:ss
	ReviseAt   string `json:"revise_at"`  //-- 修订时间 yyyy-mm-dd hh:nn:ss
	Id         string `json:"user_id"`    //--ID

}

type WxUsers struct {
	//自有用户信息
	CorpId     string `json:"corp_id"`    //-- 公司corps.id
	Tye        string `json:"tye"`        //-- 角色类型 0-用户; 1-高管; 2-高管后台管理员
	Pwd        string `json:"pwd"`        //-- 登录密码
	Nickname   string `json:"nickname"`   //-- 昵称
	Realname   string `json:"realname"`   //-- 真实姓名
	Mob        string `json:"mob"`        //-- 手机号码 *候选键* 可做登陆名
	Mail       string `json:"mail"`       //-- 电邮 *候选键* 可做登陆名
	Sex        int    `json:"sex"`        //-- 性别 0-未知, 1-男, 2-女
	Qq         string `json:"qq"`         //-- QQ号
	Headimgurl string `json:"headimgurl"` //-- 头像url
	IsDel      int    `json:"is_del"`     //-- 是否注销 0-启用, 1-注销
	CreateAt   string `json:"create_at"`  //-- 创建时间 yyyy-mm-dd hh:nn:ss
	ReviseAt   string `json:"revise_at"`  //-- 修订时间 yyyy-mm-dd hh:nn:ss
	Id         string `json:"user_id"`    //--ID

	//微信用户增加信息
	WxAppId        string `json:"wxapp_id"`       //-- 公众号app_id
	City           string `json:"city"`           //-- 所在的城市
	Country        string `json:"country"`        //-- 所在的国家
	Province       string `json:"province"`       //-- 所在的省份
	Subscribe_time string `json:"subscribe_time"` //-- 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	Remark         string `json:"remark"`         //-- 公众号运营者对粉丝的备注
	IsAttn         int    `json:"is_attn"`        //-- 是否关注公众号 0-否；1-是
	OpenId         string `json:"openid"`         //-- *键* -- 微信传入的openid
}

type Helper int

const (
	USER_MOB_PREFIX  = "user_mob_"
	USER_MAIL_PREFIX = "user_mail_"
	WXUSER_PREFIX    = "wxuser_"
)

func (this *Helper) QueryUsers(syncTime time.Time) ([]*Users, error) {
	sql := "select corp_id,tye,pwd,nickname,realname,mob,mail,sex,qq,headimgurl,is_del,create_at,revise_at,id from users"
	if syncTime.Format("2006-01-02 15:04:05") != `0001-01-01 00:00:00` {
		sql += ` where revise_at > '` + syncTime.Format("2006-01-02 15:04:05") + `'`
	}
	logs.Debug("query users sql: \n %v", sql)
	rows, err := util.Db.Query(sql)
	if err != nil {
		logs.Error("query users from db err:%v", err)
		return nil, err
	}
	var (
		list []*Users

		corpId     string
		tye        string
		pwd        string
		nickname   string
		realname   string
		mob        string
		mail       string
		sex        int
		qq         string
		headimgurl string
		isDel      int
		createAt   string
		reviseAt   string
		id         string
	)

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&corpId, &tye, &pwd, &nickname, &realname, &mob, &mail, &sex, &qq, &headimgurl, &isDel, &createAt, &reviseAt, &id); err != nil {
			logs.Error("query users from db err:%v", err)
			return nil, err
		}
		user := &Users{corpId, tye, pwd, nickname, realname, mob, mail, sex, qq, headimgurl, isDel, createAt, reviseAt, id}
		list = append(list, user)
	}
	return list, nil
}

func (this *Helper) QueryWxUsers(syncTime time.Time) ([]*WxUsers, error) {
	sql := `select 
			u.corp_id,
			u.tye,
			u.pwd,
			u.nickname,
			u.realname,
			u.mob,u.mail,
			u.sex,
			u.qq,
			u.headimgurl,
			u.is_del,
			u.create_at,
			u.revise_at,
			u.id ,
			wx.wxapp_id,
			wx.city,
			wx.country,
			wx.province,
			wx.subscribe_time,
			wx.remark,
			wx.is_attn,
			wx.id as open_id 
			from wxusers wx 
			left join users u 
			on wx.user_id=u.id`
	if syncTime.Format("2006-01-02 15:04:05") != `0001-01-01 00:00:00` {
		sql += ` where u.revise_at>'` + syncTime.Format("2006-01-02 15:04:05") + `'`
	}
	logs.Debug("query wxusers sql: \n %v", sql)
	rows, err := util.Db.Query(sql)
	if err != nil {
		logs.Error("query users from db err:%v", err)
		return nil, err
	}
	var (
		list []*WxUsers

		corpId     string
		tye        string
		pwd        string
		nickname   string
		realname   string
		mob        string
		mail       string
		sex        int
		qq         string
		headimgurl string
		isDel      int
		createAt   string
		reviseAt   string
		id         string

		wxAppId        string
		city           string
		country        string
		province       string
		subscribe_time string
		remark         string
		isAttn         int
		openId         string
	)

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&corpId, &tye, &pwd, &nickname, &realname, &mob, &mail, &sex, &qq, &headimgurl, &isDel, &createAt, &reviseAt, &id, &wxAppId, &city, &country, &province, &subscribe_time, &remark, &isAttn, &openId); err != nil {
			logs.Error("query Wxusers from db err:%v", err)
			return nil, err
		}
		user := &WxUsers{corpId, tye, pwd, nickname, realname, mob, mail, sex, qq, headimgurl, isDel, createAt, reviseAt, id, wxAppId, city, country, province, subscribe_time, remark, isAttn, openId}
		list = append(list, user)
	}
	return list, nil
}

func (this *Helper) SyncDbToRedis(syncTime time.Time) error {
	list, err := this.QueryUsers(syncTime)
	wxList, err := this.QueryWxUsers(syncTime)
	if err != nil {
		return err
	}

	pipeline := util.Client.Pipeline()
	for _, user := range list {
		pipeline.Set(USER_MAIL_PREFIX+user.Mail, util.ToJson(user), 0)
		pipeline.Set(USER_MOB_PREFIX+user.Mob, util.ToJson(user), 0)
	}

	for _, wxuser := range wxList {
		pipeline.Set(WXUSER_PREFIX+wxuser.OpenId, util.ToJson(wxuser), 0)
	}

	cmds, err := pipeline.Exec()
	if len(cmds) == 0 {
		return nil
	} else {
		if err != nil {
			return err
		}
	}
	return nil
}
