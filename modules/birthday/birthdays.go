package birthday

import (
	"strconv"
	"strings"
	"time"

	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/modules/birthday/models"
)

// monthMap maps the month name to the numerical value.
var monthMap = map[string]string{
	"January":   "01",
	"February":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

func giveBirthdayRoles(s *snorlax.Snorlax) {
	_, month, day := time.Now().Date()
	currDate := ""

	strDay := strconv.Itoa(day)
	if len(strDay) != 2 {
		strDay = "0" + strDay
	}

	currDate += monthMap[month.String()] + "/" + strDay

	birthdays, err := models.GetBirthdaysWithDate(s.DB, currDate)
	if err != nil {
		s.Log.WithError(err).Error("Error getting birthdays.")
		return
	}

	birthdayConfigs := map[string]*models.BirthdayConfig{}

	for i := 0; i < len(birthdays); i++ {
		birthday := birthdays[i]
		bdayConfig, ok := birthdayConfigs[birthday.ServerID]
		if !ok {
			bdayConfig, err = models.GetBirthdayConfig(s.DB, birthday.ServerID)
			if err != nil && err != models.ErrNoBirthdayConfigFound {
				s.Log.WithError(err).Error("Error getting birthday config.")
				return
			}

			if err == models.ErrNoBirthdayConfigFound {
				s.Log.Errorf("Could not find birthday config for %v", birthday.ServerID)
				return
			}

			birthdayConfigs[birthday.ServerID] = bdayConfig
		}

		if !bdayConfig.AssignRole || bdayConfig.BirthdayRoleID == "" {
			// Don't do anything, if a role ID hasn't been set or assign role is
			// false.
			continue
		}

		s.Session.GuildMemberRoleAdd(birthday.ServerID, birthday.UserID, bdayConfig.BirthdayRoleID)
	}
}

func removeBirthdayRoles(s *snorlax.Snorlax) {
	currBdays, err := models.GetCurrentBirthdays(s.DB)
	if err != nil {
		s.Log.WithError(err).Error("Error getting current birthdays.")
		return
	}

	now := time.Now()
	strDay := strconv.Itoa(now.Day())
	month := monthMap[now.Month().String()]
	if len(strDay) != 2 {
		strDay = "0" + strDay
	}

	for i := 0; i < len(currBdays); i++ {
		currBday := currBdays[i]

		birthdaySlice := strings.Split(currBday.Birthday, "/")
		if month != birthdaySlice[0] || strDay != birthdaySlice[1] {
			continue
		}

		err := models.DeleteCurrentBirthday(s.DB, currBday.UserID)
		if err != nil && err != models.ErrNoBirthdayConfigFound {
			s.Log.WithError(err).Error("Error deleting current birthday.")
			continue
		}

		if err == models.ErrNoBirthdayConfigFound {
			s.Log.Errorf("No current birthday was found for %v.", currBday.UserID)
		}

		if len(currBday.ServerIDs) != len(currBday.BirthdayRoleIDs) {
			s.Log.Errorf("Current birthday roleIDs and serverIDs don't match for user %v.", currBday.UserID)
		}

		for i := 0; i < len(currBday.ServerIDs); i++ {
			s.Session.GuildMemberRoleRemove(currBday.ServerIDs[i], currBday.UserID, currBday.BirthdayRoleIDs[i])
		}
	}
}
