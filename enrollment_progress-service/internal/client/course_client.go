package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	enrollment_progress__errors "main/enrollment_progress-service/internal/enrollment_progress_errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CourseResponse struct {
	ID        uuid.UUID `json:"id"`
	TeacherID uuid.UUID `json:"teacher_id"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
}
type LessonResponse struct {
	ID       uuid.UUID `json:"id"`
	CourseID uuid.UUID `json:"teacher_id"`
	Title    string    `json:"title"`
	Position int       `json:"position"`
}
type CourseClient interface {
	GetCourseByCourseID(ctx context.Context, courseID uuid.UUID) (*CourseResponse, error)
	GetLessonsByCourseID(ctx context.Context, courseID uuid.UUID) ([]LessonResponse, error)
}
type courseClient struct {
	baseURL string
	client  *http.Client
}

func NewcourseClient(baseURL string) CourseClient {
	return &courseClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
func (c *courseClient) GetCourseByCourseID(ctx context.Context, courseID uuid.UUID) (*CourseResponse, error) {
	url := fmt.Sprintf("%s/courses/%s", c.baseURL, courseID.String())
	log.Printf("COURSE CLIENT URL: %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf(
		"COURSE CLIENT STATUS: %d",
		resp.StatusCode,
	)
	switch resp.StatusCode {
	case http.StatusOK:
		var course CourseResponse
		if err := json.NewDecoder(resp.Body).Decode(&course); err != nil {
			return nil, err
		}
		return &course, nil
	case http.StatusNotFound:
		return nil, enrollment_progress__errors.ErrCourseNotFound
	default:
		return nil, enrollment_progress__errors.ErrCourseServiceUnavailable
	}
}
func (c *courseClient) GetLessonsByCourseID(ctx context.Context, courseID uuid.UUID) ([]LessonResponse, error) {
	url := fmt.Sprintf("%s/courses/%s/lessons", c.baseURL, courseID.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		var lessons []LessonResponse
		if err := json.NewDecoder(resp.Body).Decode(&lessons); err != nil {
			return nil, err
		}
		return lessons, nil
	case http.StatusNotFound:
		return nil, enrollment_progress__errors.ErrCourseNotFound
	default:
		return nil, enrollment_progress__errors.ErrCourseServiceUnavailable
	}
}
