package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"s21/models"
	"s21/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MyService struct {
	baseUrl string
	token   string
	limit   int
	client  http.Client
}

func NewMySerivce() *MyService {
	return &MyService{
		baseUrl: "https://platform.21-school.ru/services/21-school/api",
		token:   utils.S21Token(),
		limit:   900,
		client:  http.Client{},
	}
}

func Common[T models.MYJSON](ms *MyService, endPoint string) (*T, error) {
	url := ms.baseUrl + endPoint
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+ms.token)
	req.Header.Set("Content-Type", "application/json")

	body, err := ms.doRequestWithRetry(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	// todo
	// if len(body) == 0 {
	// 	return nil, nil
	// }

	var respJSON T
	err = json.Unmarshal(body, &respJSON)
	if err != nil {
		return nil, fmt.Errorf("ошибка при распаковке JSON: %w", err)
	}

	return &respJSON, nil
}

func (ms *MyService) fetch(endPoint string) []string {
	var logins []string

	for offset := 0; ; offset += ms.limit {
		endPoint := fmt.Sprintf("%s?limit=%d&offset=%d", endPoint, ms.limit, offset)

		q, err := Common[models.ParticipantsResponse](ms, endPoint)
		if err != nil {
			continue
		}

		logins = append(logins, q.Participants...)

		if len(q.Participants) < ms.limit {
			break
		}
	}

	return logins
}

func (ms *MyService) Participant(login string) (*models.Participant, error) {
	endPoint := fmt.Sprintf("/v1/participants/%s", login)
	participant, err := Common[models.Participant](ms, endPoint)
	if err != nil {
		return nil, err
	}
	return participant, nil
}

func (ms *MyService) ParticipantsWorkstation(login string) (*models.ParticipantsWorkstation, error) {
	endPoint := fmt.Sprintf("/v1/participants/%s/workstation", login)
	workstation, err := Common[models.ParticipantsWorkstation](ms, endPoint)
	if err != nil {
		return nil, err
	}

	return workstation, nil
}

func (ms *MyService) Participants() ([]models.Participant, error) {
	var (
		participants = make([]models.Participant, 0, 200)
		logins       = ms.All4MskTribesLogins()
	)

	const maxConcurrentRequests = 200

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		sem     = make(chan struct{}, maxConcurrentRequests)
		errChan = make(chan error, len(logins))
	)

	for i, login := range logins {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			participant, err := ms.Participant(login)
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			participants = append(participants, *participant)
			mu.Unlock()
			time.Sleep(500 * time.Millisecond)
		}(i)
	}
	log.Println("we're in waiting")
	wg.Wait()
	log.Println("wg.Done is behind")

	close(errChan)

	var errs []string
	for err := range errChan {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("ошибки выполнения:\n%s", strings.Join(errs, "\n"))
	}

	return participants, nil
}

func (ms *MyService) TribeLogins(coalitionID int) []string {
	endPoint := fmt.Sprintf("/v1/coalitions/%d/participants", coalitionID)
	return ms.fetch(endPoint)
}

func (ms *MyService) All4MskTribesLogins() []string {
	var tribes = map[int]string{
		156: "Саламандры",
		157: "Медоеды",
		158: "Альпаки",
		159: "Капибары",
	}

	var allLogins []string
	for coalitionID := range tribes {
		logins := ms.TribeLogins(coalitionID)
		allLogins = append(allLogins, logins...)

	}
	return allLogins
}

func (ms *MyService) Usernames() string {
	participants := ms.fetch("/v1/campuses/6bfe3c56-0211-4fe1-9e59-51616caac4dd/participants")

	return fmt.Sprintf("\nВсего: %d\n", len(participants))
}

func (ms *MyService) doRequestWithRetry(req *http.Request) ([]byte, error) {
	for {
		resp, err := ms.client.Do(req)
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		switch resp.StatusCode {
		case 200:
			if err != nil {
				log.Println("Ошибка при чтении тела ответа:", err)
			}
			return body, err
		case 429:
			retryAfter := resp.Header.Get("Retry-After")
			waitSeconds := 0
			if retryAfter != "" {
				waitSeconds, _ = strconv.Atoi(retryAfter)
			} else {
				waitSeconds = 60
				waitSeconds = 1
			}
			log.Printf("Получен 429. Повтор через %d секунд...", waitSeconds)
			time.Sleep(time.Duration(waitSeconds) * time.Second)
			continue
		default:
			log.Println("Статус-код:", resp.Status)
			return nil, err
		}
	}
}
