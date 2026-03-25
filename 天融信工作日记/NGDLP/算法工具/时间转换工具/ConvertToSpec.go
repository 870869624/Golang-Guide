package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var weekMap = map[string]string{
	"1": "mon",
	"2": "tue",
	"3": "wed",
	"4": "thu",
	"5": "fri",
	"6": "sat",
	"7": "sun",
}

// ConvertToSpec
// timeType: second,minute,hour,day,week,everyday,everyweek,everymonth
// timeContent:1,1,1,1,1,[10:00,11:00],[1,2,3;10:00,11:00],[11,12,13;10:00,11:00]
// s m h * * w
// nolint: funlen
func ConvertToSpec(timeType, timeContent string) ([]string, error) {
	var (
		timeInt int64
		err     error
		specs   []string
	)
	switch timeType {
	case "second", "minute", "hour", "day", "week":
		timeInt, err = strconv.ParseInt(timeContent, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("timeType: %s, cronvalue invalid: %s, error: %v", timeType, timeContent, err)
		}
	}

	switch timeType {
	//每间隔几秒钟
	case "second":
		return append(specs, fmt.Sprintf("@every %ds", timeInt)), nil
	//每间隔几分钟
	case "minute":
		return append(specs, fmt.Sprintf("@every %dm", timeInt)), nil
	//间隔几小时
	case "hour":
		return append(specs, fmt.Sprintf("@every %dh", timeInt)), nil
	//每间隔几天
	case "day":
		return append(specs, fmt.Sprintf("@every %dh", 24*timeInt)), nil
	//每间隔几周
	case "week":
		return append(specs, fmt.Sprintf("@every %dh", 7*24*timeInt)), nil
	//每天几时几分
	case "everyday":
		for _, hourAndMinute := range strings.Split(timeContent, ",") {
			vs := strings.Split(hourAndMinute, ":")
			if len(vs) != 2 {
				return nil, fmt.Errorf("timeType: %s, cronvalue invalid: %s", timeType, timeContent)
			}
			specs = append(specs, fmt.Sprintf("0 %s %s * * *", vs[1], vs[0]))
		}
		return specs, nil
	// 每周周几几时几分
	case "everyweek":
		vs := strings.Split(timeContent, ";")
		if len(vs) != 2 {
			return nil, fmt.Errorf("timeType: %s, cronvalue invalid: %s", timeType, timeContent)
		}
		var weeks []string
		for _, v := range strings.Split(vs[0], ",") {
			weeks = append(weeks, weekMap[v])
		}
		for _, hourAndMinute := range strings.Split(vs[1], ",") {
			vs := strings.Split(hourAndMinute, ":")
			if len(vs) != 2 {
				return nil, fmt.Errorf("timetype: %s, cronvalue invalid: %s", timeType, timeContent)
			}
			specs = append(specs, fmt.Sprintf("0 %s %s * * %s", vs[1], vs[0], strings.Join(weeks, ",")))
		}
		return specs, nil
	// 每月几号几时几分
	case "everymonth":
		vs := strings.Split(timeContent, ";")
		if len(vs) != 2 {
			return nil, fmt.Errorf("timeType: %s, cronvalue invalid: %s", timeType, timeContent)
		}
		month := vs[0]
		for _, hourAndMinute := range strings.Split(vs[1], ",") {
			vs := strings.Split(hourAndMinute, ":")
			if len(vs) != 2 {
				return nil, fmt.Errorf("timeType: %s, cronvalue invalid: %s", timeType, timeContent)
			}
			specs = append(specs, fmt.Sprintf("0 %s %s %s * *", vs[1], vs[0], month))
		}
		return specs, nil
	default:
		return nil, errors.New("no support timeType: " + timeType)
	}
}
