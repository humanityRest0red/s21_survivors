package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"s21/service"
	"s21/utils"
	"time"
)

func PartisipantSeatInfo(ms service.MyService, login string) (string, error) {
	workstation, err := ms.ParticipantsWorkstation(login)
	if err != nil {
		return fmt.Sprintf("Ошибка: %v", err), err
	}

	msg := fmt.Sprintf("%s: %s-%s%d\n", login, workstation.ClusterName, workstation.Row, workstation.Number)

	return msg, nil
}

func Survivors(ms *service.MyService) (string, error) {
	var (
		etalonClassName        = "23_10_MSK"
		msg             string = etalonClassName + "\n\n"
		count                  = 0
		participants, _        = ms.Participants()
		levels                 = make(map[int32]uint)
	)

	for _, participant := range participants {
		if participant.ClassName == etalonClassName {
			msg += fmt.Sprintf("%s: %v lvl\n", participant.Login, participant.Level)
			levels[participant.Level]++
			count++
		}
	}

	msg += fmt.Sprintf("\n%v students survived from %s\n", count, etalonClassName)

	for level, cnt := range levels {
		msg += fmt.Sprintf("\n%v students have level %v", cnt, level)
	}
	return msg, nil
}

func All4MskTribesLogins(ms *service.MyService) (string, []string) {
	var tribes = map[int]string{
		156: "Саламандры",
		157: "Медоеды",
		158: "Альпаки",
		159: "Капибары",
	}

	var (
		allLogins []string
		msg       string
		sum       = 0
	)

	for coalitionID, coalitionName := range tribes {
		logins := ms.TribeLogins(coalitionID)
		allLogins = append(allLogins, logins...)

		count := len(logins)
		sum += count
		msg += fmt.Sprintf("%s: %d\n", coalitionName, count)
	}

	msg += fmt.Sprintf("\nВсего: %d\n", sum)
	// log.Printf("all: %v\n", len(allParticipants))
	data, err := json.Marshal(map[string][]string{"participants": allLogins})
	if err != nil {
		log.Println(err)
		return msg, nil
	}

	oldLogins, _ := utils.ParticipantsFromJSON()
	log.Println(len(oldLogins))
	diff := utils.SlicesDiff(oldLogins, allLogins)
	log.Println(len(diff))
	for _, login := range diff {
		// p, err := ms.Participant(login)
		// if err != nil {
		// 	continue
		// }
		// log.Println(*p)
		// msg += "\n" + p.Login + " " + p.ClassName
		msg += "\n" + login
		time.Sleep(500 * time.Millisecond)
	}

	os.WriteFile(filepath.Join("jsons", "participants_new.json"), data, 0644)

	return msg, allLogins
}
